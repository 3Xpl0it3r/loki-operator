package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// buildLokiGatewayService generate service used to exposed an entrypoint to receive logs from client like promtail
// buildLokiGatewayService 生成一个service，用于提供一个entry point用来接受来自客户端的日志
func NewLokiGatewayService(loki *crapiv1alpha1.Loki, mod string)*apicorev1.Service{
	svc := &apicorev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: getLokiGatewayServiceName(loki),
			Namespace: loki.GetNamespace(),
			OwnerReferences: getResourceOwnerReference(loki),
		},
		Spec:       apicorev1.ServiceSpec{
			Ports: []apicorev1.ServicePort{
				{Name: "http", Protocol: "TCP", Port: 3100, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 3100}},
			},
		},
	}
	return svc
}


