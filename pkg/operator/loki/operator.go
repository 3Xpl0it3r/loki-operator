package loki

import (
	"context"
	"fmt"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	crclientset "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned"
	crlisterv1alpha1 "github.com/l0calh0st/loki-operator/pkg/client/listers/lokioperator.l0calh0st.cn/v1alpha1"
	apicorev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeclientset "k8s.io/client-go/kubernetes"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type lokiOperator struct {
	kubeClientSet kubeclientset.Interface
	crClientSet   crclientset.Interface

	lokiLister        crlisterv1alpha1.LokiLister
	configMapLister   listercorev1.ConfigMapLister
	serviceLister     listercorev1.ServiceLister
	statefulSetLister listerappsv1.StatefulSetLister
	deploymentLister  listerappsv1.DeploymentLister
	recorder          record.EventRecorder
}

func NewOperator(kubeClientSet kubeclientset.Interface, crClientSet crclientset.Interface, lokiLister crlisterv1alpha1.LokiLister,
	configMapLister listercorev1.ConfigMapLister, serviceLister listercorev1.ServiceLister, statefulSetListers listerappsv1.StatefulSetLister,
	deploymentLister listerappsv1.DeploymentLister, recorder record.EventRecorder) *lokiOperator {
	return &lokiOperator{
		kubeClientSet:     kubeClientSet,
		crClientSet:       crClientSet,
		lokiLister:        lokiLister,
		configMapLister:   configMapLister,
		serviceLister:     serviceLister,
		statefulSetLister: statefulSetListers,
		deploymentLister:  deploymentLister,
		recorder:          recorder,
	}
}

// Reconcile reads that state of the cluster for a Loki object and makes changes based on the state read and what is in the Loki.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.
func (op *lokiOperator) Reconcile(obj interface{}) error {
	key, ok := obj.(string)
	if !ok {
		return fmt.Errorf("except go string, but got %t", obj)
	}
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(err)
		return nil
	}
	loki, err := op.lokiLister.Lokis(namespace).Get(name)
	if k8serror.IsNotFound(err) {
		utilruntime.HandleError(err)
		return nil
	}
	crapiv1alpha1.WithDefaultsLoki(loki)

	for mod, _ := range loki.Spec.DeployMode {
		switch mod {
		case crapiv1alpha1.ModeKinMicroservice:
			// 以microservice 方式运行
			return op.syncLokiStackInMicroServiceMode()
		case crapiv1alpha1.ModeKindSampleScalable:
			// 以sample方式运行
			return op.syncLokiStackInSimpleMode()
		case crapiv1alpha1.ModeKindMonolithic:
			// 以monolithic 方式运行, 最简单的模式
			return op.syncLokiStackInMonolithicMode(loki)
		}
	}
	return nil
}

// syncLokiStackInMonolithicMode sync loki stack
func (op *lokiOperator) syncLokiStackInMonolithicMode(loki *crapiv1alpha1.Loki) error {
	// 以单体方式运行，所有的component都运行在一个进程里面
	modKind := crapiv1alpha1.ModeKindMonolithic
	target := crapiv1alpha1.TargetKindAllInOne
	var (
		cm  *apicorev1.ConfigMap
		err error
	)
	if loki.Spec.ConfigMap != "" {
		// 使用外置configmap
		cm, err = op.configMapLister.ConfigMaps(loki.GetNamespace()).Get(loki.Spec.ConfigMap)
		if err != nil {
			return fmt.Errorf("loki use external configmap failed, %v", err)
		}
	} else {
		// 使用内置默认的configmap
		cm, err = op.configMapLister.ConfigMaps(loki.GetNamespace()).Get(getLokiConfigMapName(loki, string(modKind)))
		if k8serror.IsNotFound(err) {
			// conigmap is not existed, then create new one
			cm, err = NewLokiConfigMap(loki, string(modKind))
			if err != nil {
				return fmt.Errorf("loki use internal configmap, build  failed, %v", err)
			}
			cm, err = op.kubeClientSet.CoreV1().ConfigMaps(cm.GetNamespace()).Create(context.TODO(), cm, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("loki use internal configmap, create failed, %v", err)
			}
		}
	}

	ss, err := op.statefulSetLister.StatefulSets(loki.GetNamespace()).Get(getLokiAppName(loki, string(modKind), string(target)))
	if k8serror.IsNotFound(err) {
		ss, err = op.kubeClientSet.AppsV1().StatefulSets(loki.GetNamespace()).Create(context.TODO(), NewStatefulSet(loki, modKind, target, cm), metav1.CreateOptions{})
	}
	if err != nil {
		return fmt.Errorf("loki deploy statefulset failed %v", err)
	}

	if !metav1.IsControlledBy(ss, loki) {
		return fmt.Errorf("loki statefulset existed, but is not controller by %v", loki.GetName())
	}

	// 创建gateway
	_, err = op.serviceLister.Services(loki.GetName()).Get(getLokiGatewayServiceName(loki))
	if k8serror.IsNotFound(err) {
		_, err = op.kubeClientSet.CoreV1().Services(loki.GetNamespace()).Create(context.TODO(), NewLokiGatewayService(loki, string(modKind)), metav1.CreateOptions{})
	}
	if err != nil {
		return fmt.Errorf("loki expose gateway failed, %v", err)
	}
	return nil
}

// syncLokiStackInSimpleMode sync loki stack
func (op *lokiOperator) syncLokiStackInSimpleMode() error {
	// simple模式运行
	return nil
}

// syncLokiStackInMicroServiceMode sync loki stack
func (op *lokiOperator) syncLokiStackInMicroServiceMode() error {
	// Loki 以微服务方式运行

	return nil
}

func reconcileErrorHandler(err error) error {
	if k8serror.IsNotFound(err) {
		utilruntime.HandleError(err)
		return nil
	}
	return err
}
