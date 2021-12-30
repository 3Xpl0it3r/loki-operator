package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getLokiAppName return the name of app according loki with mod target
func getLokiAppName(loki *crapiv1alpha1.Loki, mod, target string) string {
	return loki.GetName() + "-" + "loki" + "-" + mod + "-" + target
}

func getLokiConfigMapName(loki *crapiv1alpha1.Loki, modKind string) string {
	return loki.GetName() + "-loki-" + modKind
}

// LokiGatewayService is used to exposed to promtail,
func getLokiGatewayServiceName(loki *crapiv1alpha1.Loki) string {
	return loki.GetName() + "-loki-gateway"
}

// getResourceLabels generate labels according crResource object
func getResourceLabels(loki *crapiv1alpha1.Loki, mod, target string) map[string]string {
	labels := map[string]string{
		"app":        loki.GetName(),
		"controller": loki.Kind,
		"deployMode": mod,
		"target":     target,
	}
	return labels
}

// getResourceAnnotations generate annotations according crResource object
func getResourceAnnotations(loki *crapiv1alpha1.Loki) map[string]string {
	annotations := map[string]string{}
	return annotations
}

// getResourceOwnerReference generate OwnerReference according crResource object
func getResourceOwnerReference(loki *crapiv1alpha1.Loki) []metav1.OwnerReference {
	ownerReference := []metav1.OwnerReference{}
	ownerReference = append(ownerReference, *metav1.NewControllerRef(loki, crapiv1alpha1.SchemeGroupVersion.WithKind("Loki")))
	return ownerReference
}
