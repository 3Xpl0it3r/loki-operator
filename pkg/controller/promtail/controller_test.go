package promtail_test

import (
	"context"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	"github.com/l0calh0st/loki-operator/pkg/controller"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("Controller", func() {
	var (
		crdObj     *crapiv1alpha1.Promtail
		fakeCtrl   *fakeController
		eventsHook controller.EventsHook
		event      controller.Event

	)

	BeforeEach(func() {
		eventsHook = controller.NewEventsHook(10)
		crdObj = newPromtail()
		fakeCtrl = newFakeController()

	})
	JustBeforeEach(func() {
		crapiv1alpha1.WithDefaultsPromtail(crdObj)
		gomega.Ω(fakeCtrl.controller.AddHook(eventsHook)).ShouldNot(gomega.HaveOccurred())

	})
	

	Context("Create Promtail", func() {
		It("should receive addEvent from eventsHooks", func() {
			fakeCtrl.crClient.LokioperatorV1alpha1().Promtails(apicorev1.NamespaceDefault).Create(context.TODO(), crdObj, metav1.CreateOptions{})
			gomega.Eventually(eventsHook.GetEventsChan()).Should(gomega.Receive(&event), 2*time.Second)
			gomega.Ω(event.Type).To(gomega.Equal(controller.EventAdded))
			gomega.Ω(event.Object).To(gomega.Equal(crdObj))
		})
	})

})

func newPromtail() *crapiv1alpha1.Promtail {
	return &crapiv1alpha1.Promtail{
		TypeMeta: metav1.TypeMeta{
			APIVersion: crapiv1alpha1.SchemeGroupVersion.String(),
			Kind:       "Promtail",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: apicorev1.NamespaceDefault,
		},
		Spec:   crapiv1alpha1.PromtailSpec{},
		Status: crapiv1alpha1.PromtailStatus{},
	}
}
