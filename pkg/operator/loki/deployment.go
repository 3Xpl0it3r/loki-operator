package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newLokiDeployment create an new loki instance according given target
func newLokiDeployment(loki *crapiv1alpha1.Loki, modKind crapiv1alpha1.ModeKind, target crapiv1alpha1.LokiTargetKind, cm *apicorev1.ConfigMap) *apiappsv1.Deployment {
	var (
		resourceLabels = getResourceLabels(loki, string(modKind), string(target))
		replicas       = loki.Spec.DeployMode[modKind][target]
	)

	dpl := &apiappsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getLokiAppName(loki, string(modKind), string(target)),
			Namespace:       loki.GetNamespace(),
			OwnerReferences: getResourceOwnerReference(loki),
		},
		Spec: apiappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: resourceLabels},
			Template: apicorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: resourceLabels,
				},
				Spec: apicorev1.PodSpec{
					Containers: []apicorev1.Container{
						{
							Name:  "loki",
							Image: loki.Spec.Image,
							Args:  []string{"-config.file=/etc/loki/config/config.yaml", "-target=" + string(target)},
							VolumeMounts: []apicorev1.VolumeMount{
								{
									Name:      "conf",
									ReadOnly:  true,
									MountPath: "/etc/loki/config/config.yaml",
									SubPath:   "config.yaml",
								},
							},
						},
					},
					Volumes: []apicorev1.Volume{
						{
							Name: "config",
							VolumeSource: apicorev1.VolumeSource{
								ConfigMap: &apicorev1.ConfigMapVolumeSource{
									LocalObjectReference: apicorev1.LocalObjectReference{Name: cm.GetName()},
								},
							},
						},
					},
				},
			},
		},
	}
	return dpl
}
