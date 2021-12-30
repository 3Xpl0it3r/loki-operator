package promtail_test

import (
	"fmt"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	crfakeclients "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned/fake"
	croperatortesting "github.com/l0calh0st/loki-operator/pkg/k8s/testing"
	"github.com/l0calh0st/loki-operator/pkg/operator/promtail"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

var _ = Describe("Test Promtail Operator", func() {
	var (
		crObj     *crapiv1alpha1.Promtail
		configMap *apicorev1.ConfigMap
		service   *apicorev1.Service
		daemonSet *apiappsv1.DaemonSet

		fakeCtrl   *fakeController
		kubeClient *k8sfake.Clientset
		crClient   *crfakeclients.Clientset
	)

	BeforeEach(func() {
		crObj = newPromtail("test")
	})

	JustBeforeEach(func() {
		crapiv1alpha1.WithDefaultsPromtail(crObj)
		service = promtail.NewService(crObj)

		kubeClient = k8sfake.NewSimpleClientset()
		crClient = crfakeclients.NewSimpleClientset()
		fakeCtrl = newFakeController(kubeClient, crClient)
	})

	JustAfterEach(func() {
		crObj = nil
		configMap = nil
		service = nil
		daemonSet = nil
		kubeClient = nil
		crClient = nil
		fakeCtrl = nil
	})

	Context("Reconcile", func() {
		It("Create Promtail, External Configmap", func() {
			configMap = fakeExternalConfigmap()
			crObj.Spec.ConfigMap = configMap.GetName()
			daemonSet = promtail.NewDaemonSet(crObj, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.fixture.AddConfigMapLister(configMap)).ShouldNot(gomega.HaveOccurred())

			fakeCtrl.fixture.RecordKubeActions(croperatortesting.ExpectCreateDaemonSetAction(daemonSet),
				croperatortesting.ExpectCreateServiceAction(service))
			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())
		})

		It("Create Promtail, Internal Config, Has Spec/Config Value", func() {
			var err error
			crObj.Spec.Config.Clients.URL = "http://fakeurl"
			crObj.Spec.Config.Server.HttpListenAddress = "127.0.0.1:9090"
			configMap, err = promtail.NewConfigMap(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			daemonSet = promtail.NewDaemonSet(crObj, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			fakeCtrl.fixture.RecordKubeActions(croperatortesting.ExpectCreateConfigMapAction(configMap),
				croperatortesting.ExpectCreateDaemonSetAction(daemonSet), croperatortesting.ExpectCreateServiceAction(service))

			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())
		})

		It("Create Promtail, Internal ConfigMap, No Spec/Conf Value", func() {
			var err error
			configMap, err = promtail.NewConfigMap(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			daemonSet = promtail.NewDaemonSet(crObj, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			fakeCtrl.fixture.RecordKubeActions(croperatortesting.ExpectCreateConfigMapAction(configMap),
				croperatortesting.ExpectCreateDaemonSetAction(daemonSet), croperatortesting.ExpectCreateServiceAction(service))

			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())
		})

		It("Create Promtail Failed, Reconcile key not string", func() {
			var err error
			crapiv1alpha1.WithDefaultsPromtail(crObj)
			configMap, err = promtail.NewConfigMap(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			daemonSet = promtail.NewDaemonSet(crObj, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			fakeCtrl.fixture.RecordKubeActions(croperatortesting.ExpectCreateConfigMapAction(configMap),
				croperatortesting.ExpectCreateDaemonSetAction(daemonSet), croperatortesting.ExpectCreateServiceAction(service))

			gomega.Ω(fakeCtrl.runController(crObj, true)).Should(gomega.HaveOccurred())
		})

		It("Create Promtail, Resource Is out of controller", func() {})

	})

})

func newPromtail(name string) *crapiv1alpha1.Promtail {
	return &crapiv1alpha1.Promtail{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Promtail",
			APIVersion: crapiv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: apicorev1.NamespaceDefault,
		},
		Spec:   crapiv1alpha1.PromtailSpec{},
		Status: crapiv1alpha1.PromtailStatus{},
	}
}

// newExternalConfig return external configmap
func fakeExternalConfigmap() *apicorev1.ConfigMap {
	return &apicorev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: apicorev1.NamespaceDefault,
		},
		Data:       map[string]string{"test": "test"},
		BinaryData: nil,
	}
}

func getKey(obj *crapiv1alpha1.Promtail) (string, error) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return "", fmt.Errorf("unexpected error getting key for foo %v: %v", obj.Name, err)
	}
	return key, nil
}
