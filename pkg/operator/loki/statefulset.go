package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// newLokiStatefulset build statefulset for loki, according target with configmap
func NewStatefulSet(loki *crapiv1alpha1.Loki, modKind crapiv1alpha1.ModeKind, target crapiv1alpha1.LokiTargetKind, cm *apicorev1.ConfigMap) *apiappsv1.StatefulSet {
	resourceLabels := getResourceLabels(loki, string(modKind), string(target))
	sts := &apiappsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            getLokiAppName(loki, string(modKind), string(target)),
			Namespace:       loki.GetNamespace(),
			OwnerReferences: getResourceOwnerReference(loki),
		},
		Spec: apiappsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: resourceLabels},
			Template: apicorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: resourceLabels},
				Spec: apicorev1.PodSpec{
					Containers: []apicorev1.Container{
						{
							Name:            "loki",
							Image:           loki.Spec.Image,
							Args:            []string{"-config.file=/etc/loki/config/loki.yaml", "-target=" + string(target)},
							ImagePullPolicy: apicorev1.PullIfNotPresent,
							Ports: []apicorev1.ContainerPort{
								{Name: "http", ContainerPort: 3100, Protocol: "TCP"},
								{Name: "grpc", ContainerPort: 9095, Protocol: "TCP"},
							},
							VolumeMounts: []apicorev1.VolumeMount{
								{Name: "config", MountPath: "/etc/loki/config/loki.yaml", SubPath: "loki.yaml"},
								{Name: "storage", MountPath: "/data"},
							},
							LivenessProbe: &apicorev1.Probe{FailureThreshold: 3, ProbeHandler: apicorev1.ProbeHandler{HTTPGet: &apicorev1.HTTPGetAction{
								Path:   "/ready",
								Port:   intstr.IntOrString{Type: intstr.String, StrVal: "http"},
								Scheme: "HTTP",
							}}, InitialDelaySeconds: 45, PeriodSeconds: 10, SuccessThreshold: 1, TimeoutSeconds: 5},
							ReadinessProbe: &apicorev1.Probe{FailureThreshold: 3, ProbeHandler: apicorev1.ProbeHandler{HTTPGet: &apicorev1.HTTPGetAction{
								Path:   "/ready",
								Port:   intstr.IntOrString{Type: intstr.String, StrVal: "http"},
								Scheme: "HTTP",
							}}, InitialDelaySeconds: 45, PeriodSeconds: 10, SuccessThreshold: 1, TimeoutSeconds: 5},
						},
					},
					Volumes: []apicorev1.Volume{
						{Name: "config", VolumeSource: apicorev1.VolumeSource{ConfigMap: &apicorev1.ConfigMapVolumeSource{LocalObjectReference: apicorev1.LocalObjectReference{Name: cm.GetName()}}}},
						{Name: "storage", VolumeSource: apicorev1.VolumeSource{EmptyDir: &apicorev1.EmptyDirVolumeSource{}}},
					},
				},
			},
		},
	}
	return sts
}
