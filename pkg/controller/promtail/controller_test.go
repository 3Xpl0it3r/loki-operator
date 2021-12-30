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
		fakeCtrl    *fakeController
		eventsHook controller.EventsHook
		event controller.Event
		stopCh chan struct{}

		ctx context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		eventsHook = controller.NewEventsHook(10)
		crdObj = newPromtail()
		fakeCtrl = newFakeController()
		stopCh = make(chan struct{})
		ctx,cancel = context.WithCancel(context.TODO())
	})
	JustBeforeEach(func() {
		crapiv1alpha1.WithDefaultsPromtail(crdObj)
		gomega.Ω(fakeCtrl.controller.AddHook(eventsHook)).ShouldNot(gomega.HaveOccurred())
		fakeCtrl.crInformerFactory.Start(stopCh)
		gomega.Ω(fakeCtrl.controller.Start(ctx)).ShouldNot(gomega.HaveOccurred())
	})
	JustAfterEach(func() {
		// should stop controller first before informer
		cancel()
		close(stopCh)
	})

	Context("Create Promtail", func() {
		It("should receive addEvent from eventsHooks", func() {
			fakeCtrl.crClient.LokioperatorV1alpha1().Promtails(apicorev1.NamespaceDefault).Create(context.TODO(), crdObj, metav1.CreateOptions{})
			gomega.Eventually(eventsHook.GetEventsChan()).Should(gomega.Receive(&event),2 * time.Second)
			gomega.Ω(event.Type).To(gomega.Equal(controller.EventAdded))
			gomega.Ω(event.Object).To(gomega.Equal(crdObj))
		})
	})



})

func newPromtail()*crapiv1alpha1.Promtail{
	return &crapiv1alpha1.Promtail{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: crapiv1alpha1.SchemeGroupVersion.String(),
			Kind: "Promtail",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Namespace: apicorev1.NamespaceDefault,
		},
		Spec:       crapiv1alpha1.PromtailSpec{},
		Status:     crapiv1alpha1.PromtailStatus{},
	}
}
