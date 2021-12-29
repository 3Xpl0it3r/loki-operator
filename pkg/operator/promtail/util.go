package promtail

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getPromtailAppName(promtail *crapiv1alpha1.Promtail )string{
	return  promtail.GetName() + "-promtail"
}



// getResourceLabels generate labels according crResource object
func getResourceLabels(promtail *crapiv1alpha1.Promtail)map[string]string{
	labels := map[string]string{
		"app": promtail.GetName(),
		"controller": promtail.Kind,
	}
	return labels
}

// getResourceAnnotations generate annotations according crResource object
//func getResourceAnnotations(promtail *crapiv1alpha1.Promtail)map[string]string{
//	annotations := map[string]string{
//	}
//	return annotations
//}

// getResourceOwnerReference generate OwnerReference according crResource object
func getResourceOwnerReference(promtail *crapiv1alpha1.Promtail)[]metav1.OwnerReference{
	ownerReference := []metav1.OwnerReference{}
	ownerReference = append(ownerReference, *metav1.NewControllerRef(promtail, crapiv1alpha1.SchemeGroupVersion.WithKind("Promtail")))
	return ownerReference
}
