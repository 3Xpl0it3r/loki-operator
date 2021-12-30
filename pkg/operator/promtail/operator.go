package promtail

import (
	"context"
	"fmt"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	crclientset "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned"
	crlisterv1alpha1 "github.com/l0calh0st/loki-operator/pkg/client/listers/lokioperator.l0calh0st.cn/v1alpha1"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeclientset "k8s.io/client-go/kubernetes"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type promtailOperator struct {
	kubeClientSet kubeclientset.Interface
	crClientSet   crclientset.Interface

	promtailLister  crlisterv1alpha1.PromtailLister
	configMapLister listercorev1.ConfigMapLister
	serviceLister   listercorev1.ServiceLister
	daemonSetLister listerappsv1.DaemonSetLister
	recorder        record.EventRecorder
}

func NewOperator(kubeClientSet kubeclientset.Interface, crClientSet crclientset.Interface, promtailLister crlisterv1alpha1.PromtailLister,
	configMapLister listercorev1.ConfigMapLister, serviceLister listercorev1.ServiceLister, daemonSetLister listerappsv1.DaemonSetLister,
	recorder record.EventRecorder) *promtailOperator {
	return &promtailOperator{
		kubeClientSet:   kubeClientSet,
		crClientSet:     crClientSet,
		promtailLister:  promtailLister,
		configMapLister: configMapLister,
		serviceLister:   serviceLister,
		daemonSetLister: daemonSetLister,
		recorder:        recorder,
	}
}

func (op *promtailOperator) Reconcile(obj interface{}) error {
	key, ok := obj.(string)
	if !ok {
		return fmt.Errorf("except go string, but got %t", obj)
	}
	// convert namespace/name string into distinct namespace and name occurred an error
	// for it's may a invalid object, in this case we should drop this key, not continue
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(err)
		return nil
	}

	promtail, err := op.promtailLister.Promtails(namespace).Get(name)
	if k8serror.IsNotFound(err) {
		// promTail resource is not existed,in this case we should stop continue reconcile
		utilruntime.HandleError(err)
		return nil
	}
	if err != nil {
		return reconcileErrorHandler(err)
	}
	// if promTail exist, check correspond configmap
	var cm *apicorev1.ConfigMap
	if promtail.Spec.ConfigMap != "" {
		// use external configmap that user provided
		cm, err = op.configMapLister.ConfigMaps(namespace).Get(promtail.Spec.ConfigMap)
		if err != nil {
			return fmt.Errorf("promtail use external configmap failed: %v", err)
		}
	} else {
		cm, err = NewConfigMap(promtail)
		if err != nil {
			return fmt.Errorf("promtail use internal configmap, template failed: %v", err)
		}
		_, err = op.configMapLister.ConfigMaps(cm.GetNamespace()).Get(cm.GetName())
		if k8serror.IsNotFound(err) {
			_, err = op.kubeClientSet.CoreV1().ConfigMaps(cm.GetNamespace()).Create(context.TODO(), cm, metav1.CreateOptions{})
		}
		if err != nil {
			return fmt.Errorf("promtail use internal configmap, create failed: %v", err)
		}
	}

	// deploy daemonSet for promTail
	ds, err := op.daemonSetLister.DaemonSets(namespace).Get(getPromtailAppName(promtail))
	if k8serror.IsNotFound(err) {
		ds, err = op.kubeClientSet.AppsV1().DaemonSets(namespace).Create(context.TODO(), NewDaemonSet(promtail, cm), metav1.CreateOptions{})
	}
	if err != nil {
		return fmt.Errorf("promtail deploy daemonset failed: %v", err)
	}

	if !metav1.IsControlledBy(ds, promtail) {
		return fmt.Errorf("promtail daemonset exists, but is not owns it: %s", promtail.GetName())
	}
	// create service
	svc, err := op.serviceLister.Services(namespace).Get(getPromtailAppName(promtail))
	if k8serror.IsNotFound(err) {
		svc, err = op.kubeClientSet.CoreV1().Services(namespace).Create(context.TODO(), NewService(promtail), metav1.CreateOptions{})
	}
	if err != nil {
		return fmt.Errorf("promtail expose gateway failed: %v", err)
	}
	if !metav1.IsControlledBy(svc, promtail) {
		return fmt.Errorf("promtail gateway existed, but is not owns it: %s", svc.GetName())
	}

	return nil
}

// return error means logical need rehandler object
func reconcileErrorHandler(err error) error {
	if k8serror.IsNotFound(err) {
		utilruntime.HandleError(err)
		return nil
	}
	return err
}

func NewDaemonSet(promtail *crapiv1alpha1.Promtail, cm *apicorev1.ConfigMap) *apiappsv1.DaemonSet {
	dm := &apiappsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getPromtailAppName(promtail),
			Namespace:       promtail.Namespace,
			OwnerReferences: getResourceOwnerReference(promtail),
			Labels:          getResourceLabels(promtail),
		},
		Spec: apiappsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: getResourceLabels(promtail)},
			Template: apicorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getResourceLabels(promtail),
				},
				Spec: apicorev1.PodSpec{
					Containers: []apicorev1.Container{
						{
							Name:  "promtail",
							Image: promtail.Spec.Image,
							VolumeMounts: []apicorev1.VolumeMount{
								{
									Name:      "config",
									ReadOnly:  true,
									MountPath: "/etc/promtail/promtail.yaml",
									SubPath:   "promtail.yaml",
								}, {
									Name:      "run",
									MountPath: "/run/promtail",
								}, {
									Name:      "pods",
									ReadOnly:  true,
									MountPath: "/var/logs/pods",
								}, {
									Name:      "containers",
									ReadOnly:  true,
									MountPath: "/var/lib/docker/containers",
								},
							},
							Args: []string{"-config.file", "/etc/promtail/promtail.yaml"},
							Env: []apicorev1.EnvVar{
								{Name: "HOSTNAME", ValueFrom: &apicorev1.EnvVarSource{FieldRef: &apicorev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}}},
							},
						},
					},
					Volumes: []apicorev1.Volume{
						{
							Name:         "config",
							VolumeSource: apicorev1.VolumeSource{ConfigMap: &apicorev1.ConfigMapVolumeSource{LocalObjectReference: apicorev1.LocalObjectReference{Name: cm.GetName()}}},
						}, {
							Name:         "run",
							VolumeSource: apicorev1.VolumeSource{HostPath: &apicorev1.HostPathVolumeSource{Path: "/run/promtail"}},
						}, {
							Name:         "pods",
							VolumeSource: apicorev1.VolumeSource{HostPath: &apicorev1.HostPathVolumeSource{Path: "/var/logs/pods"}},
						}, {
							Name:         "containers",
							VolumeSource: apicorev1.VolumeSource{HostPath: &apicorev1.HostPathVolumeSource{Path: "/var/lib/docker/containers"}},
						},
					},
				},
			},
		},
	}
	return dm
}

func NewService(promtail *crapiv1alpha1.Promtail) *apicorev1.Service {
	svcGw := &apicorev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getPromtailAppName(promtail) + "-gateway",
			Namespace:       promtail.GetNamespace(),
			OwnerReferences: getResourceOwnerReference(promtail),
		},
		Spec: apicorev1.ServiceSpec{
			Selector: getResourceLabels(promtail),
			Ports: []apicorev1.ServicePort{
				{Name: "http", Port: 9080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80}},
			},
			Type: apicorev1.ServiceTypeClusterIP,
		},
		Status: apicorev1.ServiceStatus{},
	}
	return svcGw
}
