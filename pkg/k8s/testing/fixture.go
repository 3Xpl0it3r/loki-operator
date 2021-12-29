package testing

import (
	"fmt"

	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	k8stest "k8s.io/client-go/testing"

	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	crinformers "github.com/l0calh0st/loki-operator/pkg/client/informers/externalversions"
)

type Fixture struct {

	// todo write your code here
	kubeInformers informers.SharedInformerFactory
	crInformers   crinformers.SharedInformerFactory

	// Object to put in the store
	podLister        []*apicorev1.Pod
	deploymentLister []*apiappsv1.Deployment
	statefulSet      []*apiappsv1.StatefulSet
	daemonSet        []*apiappsv1.DaemonSet
	serviceLister    []*apicorev1.Service
	configMapLister  []*apicorev1.ConfigMap
	// for custom resource
	customListers []*runtime.Object

	// Actions expected to happen on client.
	kubeActions           []k8stest.Action
	customResourceActions []k8stest.Action

	// Objects from here preloaded into NewSimpleFake
	kubeObjects           []runtime.Object
	customResourceObjects []runtime.Object
}

func NewFixture(kubeInformer informers.SharedInformerFactory, crInformer crinformers.SharedInformerFactory) *Fixture {
	return &Fixture{
		kubeInformers:         kubeInformer,
		crInformers:           crInformer,
		podLister:             make([]*apicorev1.Pod, 0),
		deploymentLister:      make([]*apiappsv1.Deployment, 0),
		statefulSet:           make([]*apiappsv1.StatefulSet, 0),
		daemonSet:             make([]*apiappsv1.DaemonSet, 0),
		serviceLister:         make([]*apicorev1.Service, 0),
		configMapLister:       make([]*apicorev1.ConfigMap, 0),
		customListers:         make([]*runtime.Object, 0),
		kubeActions:           make([]k8stest.Action, 0),
		customResourceActions: make([]k8stest.Action, 0),
	}
}

// AddPodLister add pod objects into informers
func (f *Fixture) AddPodLister(pods ...*apicorev1.Pod) error {
	for _, pod := range pods {
		if err := f.kubeInformers.Core().V1().Pods().Informer().GetIndexer().Add(pod); err != nil {
			return err
		}
	}
	return nil
}

// AddDeploymentLister add deployment objects into informers
func (f *Fixture) AddDeploymentLister(dpls ...*apiappsv1.Deployment) error {
	for _, dpl := range dpls {
		if err := f.kubeInformers.Apps().V1().Deployments().Informer().GetIndexer().Add(dpl); err != nil {
			return err
		}
	}
	return nil
}

// AddStatefulSetLister add statefulSet objects into informers
func (f *Fixture) AddStatefulSetLister(sts ...*apiappsv1.StatefulSet) error {
	for _, st := range sts {
		if err := f.kubeInformers.Apps().V1().StatefulSets().Informer().GetIndexer().Add(st); err != nil {
			return err
		}
	}
	return nil
}

// AddDaemonSetLister add daemonSet objects into informers
func (f *Fixture) AddDaemonSetLister(dss ...*apiappsv1.DaemonSet) error {
	for _, ds := range dss {
		if err := f.kubeInformers.Apps().V1().DaemonSets().Informer().GetIndexer().Add(ds); err != nil {
			return err
		}
	}
	return nil
}

// AddServiceLister add Service objects into informers
func (f *Fixture) AddServiceLister(svs ...*apicorev1.Service) error {
	for _, sv := range svs {
		if err := f.kubeInformers.Core().V1().Services().Informer().GetIndexer().Add(sv); err != nil {
			return err
		}
	}
	return nil
}

// AddConfigMapLister add configMap objects into informers
func (f *Fixture) AddConfigMapLister(cms ...*apicorev1.ConfigMap) error {
	for _, cm := range cms {
		if err := f.kubeInformers.Core().V1().ConfigMaps().Informer().GetIndexer().Add(cm); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fixture) AddCustomResourceLister(cr runtime.Object) error {
	f.customResourceObjects = append(f.customResourceObjects, cr)
	switch cr.GetObjectKind().GroupVersionKind().Kind {
	case "Loki":
		if err := f.crInformers.Lokioperator().V1alpha1().Lokis().Informer().GetIndexer().Add(cr.(*crapiv1alpha1.Loki)); err != nil {
			return err
		}
	case "Promtail":
		if err := f.crInformers.Lokioperator().V1alpha1().Promtails().Informer().GetIndexer().Add(cr.(*crapiv1alpha1.Promtail)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unexpect Custom Resource Type %s %s ", cr.GetObjectKind().GroupVersionKind().Kind, cr.GetObjectKind().GroupVersionKind().GroupVersion())
	}
	return nil
}

// add expect actions
func (f *Fixture) RecordKubeActions(kubeActions ...k8stest.Action) {
	f.kubeActions = append(f.kubeActions, kubeActions...)
}

func (f *Fixture) RecordCustomResourceActions(crActions ...k8stest.Action) {
	f.customResourceActions = append(f.customResourceActions, crActions...)
}

// GetKubeActions return all saved actions
func (f *Fixture) GetKubeActions() []k8stest.Action {
	return f.kubeActions
}

// GetCustomResourceActions return saved actions
func (f *Fixture) GetCustomResourceActions() []k8stest.Action {
	return f.customResourceActions
}
