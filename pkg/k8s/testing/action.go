package testing

import (
	"fmt"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubediff "k8s.io/apimachinery/pkg/util/diff"
	k8stest "k8s.io/client-go/testing"
	"reflect"
)

func ValidateActions(expectedActions, actualActions []k8stest.Action)error{
	for i, action := range actualActions {
		if len(expectedActions) < i+1 {
			return fmt.Errorf("%d unexpected actions: %+v ", len(actualActions)-len(expectedActions), actualActions[i:])
		}
		if err := checkAction(expectedActions[i], action);err != nil{
			return err
		}
	}
	if len(expectedActions) > len(actualActions) {
		return fmt.Errorf("%d additional expected cr actions: %+v", len(expectedActions)-len(actualActions), expectedActions[len(actualActions):])
	}
	return nil
}

// validate validate action
func checkAction(expected, actual k8stest.Action) error {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		return fmt.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Errorf("Actions has wrong type. Expected : %t . Got %t ", expected, actual)
	}

	switch actAction := actual.(type) {
	case k8stest.CreateActionImpl:
		expAction, _ := expected.(k8stest.CreateActionImpl)
		expObject := expAction.GetObject()
		actObject := actAction.GetObject()
		if !reflect.DeepEqual(expObject, actObject) {
			return fmt.Errorf("Action %s %s has wrong object \n Diff :\n %s", actAction.GetVerb(), actAction.GetResource().Resource, kubediff.ObjectGoPrintSideBySide(expObject, actObject))
		}
	case k8stest.UpdateActionImpl:
		expAction, _ := expected.(k8stest.UpdateActionImpl)
		expObject := expAction.GetObject()
		actObject := actAction.GetObject()
		if !reflect.DeepEqual(expObject, actObject) {
			return fmt.Errorf("Action %s %s has wrong object \n Diff :\n %s", actAction.GetVerb(), actAction.GetResource().Resource, kubediff.ObjectGoPrintSideBySide(expObject, actObject))
		}
	case k8stest.PatchActionImpl:
		expAction, _ := expected.(k8stest.PatchActionImpl)
		expPath := expAction.GetPatch()
		actPatch := actAction.GetPatch()
		if !reflect.DeepEqual(expPath, actPatch) {
			return fmt.Errorf("Action %s %s has wrong object \n Diff :\n %s", actAction.GetVerb(), actAction.GetResource().Resource, kubediff.ObjectGoPrintSideBySide(expPath, actPatch))
		}
	default:
		return fmt.Errorf("Uncaptured action %s %s, you should explicity add a case to capture it ", actual.GetVerb(), actual.GetResource().Resource)
	}
	return nil
}


// Pod Action Creator

// ExpectCreatePodAction return pod's CreateAction
func ExpectCreatePodAction(pod *apicorev1.Pod) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "pods"}, pod.GetNamespace(), pod)
}

// ExpectUpdatePodAction return pod's UpdateAction
func ExpectUpdatePodAction(pod *apicorev1.Pod) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "pods"}, pod.GetNamespace(), pod)
}

// ExpectGetPodAction return pod's GetAction
func ExpectGetPodAction(pod *apicorev1.Pod) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "pods"}, pod.GetNamespace(), pod.GetName())
}

// Deployment Action Creator

// ExpectCreateDeploymentAction return deployment's CreateAction
func ExpectCreateDeploymentAction(dpl *apiappsv1.Deployment) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "deployments"}, dpl.GetNamespace(), dpl)
}

// ExpectUpdateDeploymentAction return deployment's UpdateAction
func ExpectUpdateDeploymentAction(dpl *apiappsv1.Deployment) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "deployments"}, dpl.GetNamespace(), dpl)
}

// ExpectGetDeploymentAction return deployment's GetAction
func ExpectGetDeploymentAction(dpl *apiappsv1.Deployment) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "deployments"}, dpl.GetNamespace(), dpl.GetName())
}

// DaemonSet Action Creator

// ExpectCreateDaemonSetAction return daemonSet's CreateAction
func ExpectCreateDaemonSetAction(ds *apiappsv1.DaemonSet) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "daemonsets"}, ds.GetNamespace(), ds)
}

// ExpectUpdateDaemonSetAction return daemonSet's UpdateAction
func ExpectUpdateDaemonSetAction(ds *apiappsv1.DaemonSet) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "daemonsets"}, ds.GetNamespace(), ds)
}

// ExpectGetDaemonSetAction return daemonSet's GetAction
func ExpectGetDaemonSetAction(ds *apiappsv1.DaemonSet) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "daemonsets"}, ds.GetNamespace(), ds.GetName())
}

// StatefulSet Action Creator

// ExpectCreateStatefulSetAction return statefulSet's CreateAction
func ExpectCreateStatefulSetAction(sts *apiappsv1.StatefulSet) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "statefulsets"}, sts.GetNamespace(), sts)
}

// ExpectUpdateStatefulSetAction return statefulSet's UpdateAction
func ExpectUpdateStatefulSetAction(sts *apiappsv1.StatefulSet) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "statefulsets"}, sts.GetNamespace(), sts)
}

// ExpectGetStatefulSetAction return statefulSet's GetAction
func ExpectGetStatefulSetAction(sts *apiappsv1.StatefulSet) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "statefulsets"}, sts.GetNamespace(), sts.GetName())
}

// Service Action Creator

// ExpectCreateServiceAction return service's CreateAction
func ExpectCreateServiceAction(svc *apicorev1.Service) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "services"}, svc.GetNamespace(), svc)
}

// ExpectUpdateServiceAction return service's UpdateAction
func ExpectUpdateServiceAction(svc *apicorev1.Service) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "services"}, svc.GetNamespace(), svc)
}

// ExpectGetServiceAction return service's GetAction
func ExpectGetServiceAction(svc *apicorev1.Service) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "services"}, svc.GetNamespace(), svc.GetName())
}

// ConfigMap Action Creator

// ExpectCreateConfigMapAction return create configMap action
func ExpectCreateConfigMapAction(cm *apicorev1.ConfigMap) k8stest.Action {
	return k8stest.NewCreateAction(schema.GroupVersionResource{Resource: "configmaps"}, cm.GetNamespace(), cm)
}

// ExpectUpdateConfigMapAction return update configMap action
func ExpectUpdateConfigMapAction(cm *apicorev1.ConfigMap) k8stest.Action {
	return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "configmaps"}, cm.GetNamespace(), cm)
}

// ExpectGetConfigMapAction return get configMap action
func ExpectGetConfigMapAction(cm *apicorev1.ConfigMap) k8stest.Action {
	return k8stest.NewGetAction(schema.GroupVersionResource{Resource: "configmaps"}, cm.GetNamespace(), cm.GetName())
}

// CustomResource Action Creator

// for custom resource actions
func ExpectUpdateCustomResourceAction(cr runtime.Object) k8stest.Action {
	switch cr.GetObjectKind().GroupVersionKind().Kind {
	case "Loki":
		return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "lokis"}, cr.(*crapiv1alpha1.Loki).GetNamespace(), cr)
	case "Promtail":
		return k8stest.NewUpdateAction(schema.GroupVersionResource{Resource: "promtails"}, cr.(*crapiv1alpha1.Loki).GetNamespace(), cr)
	}
	return nil
}

func ExpectUpdateCustomResourceStatusAction(cr runtime.Object) k8stest.Action {
	switch cr.GetObjectKind().GroupVersionKind().Kind {
	case "Loki":
		return k8stest.NewUpdateSubresourceAction(schema.GroupVersionResource{Resource: "lokis"}, "status", cr.(*crapiv1alpha1.Loki).GetNamespace(), cr)
	case "Promtail":
		return k8stest.NewUpdateSubresourceAction(schema.GroupVersionResource{Resource: "promtails"}, "status", cr.(*crapiv1alpha1.Promtail).GetNamespace(), cr)
	}
	return nil
}



