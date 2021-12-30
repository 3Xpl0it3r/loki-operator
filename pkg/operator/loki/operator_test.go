package loki_test

import (
	"fmt"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	crfakeclients "github.com/l0calh0st/loki-operator/pkg/client/clientset/versioned/fake"
	"github.com/l0calh0st/loki-operator/pkg/k8s/testing"
	"github.com/l0calh0st/loki-operator/pkg/operator/loki"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

var _ = Describe("Test Loki Operator", func() {

	Context("Run Monolithic", func() {
		var (
			crObj       *crapiv1alpha1.Loki
			configMap   *apicorev1.ConfigMap
			service     *apicorev1.Service
			statefulSet *apiappsv1.StatefulSet

			fakeCtrl   *fakeController
			kubeClient *k8sfake.Clientset
			crClient   *crfakeclients.Clientset
		)
		BeforeEach(func() {
			crObj = newLoki("test")
			statefulSet = nil
			service = nil
		})
		JustBeforeEach(func() {
			crapiv1alpha1.WithDefaultsLoki(crObj)
			crObj.Spec.DeployMode = map[crapiv1alpha1.ModeKind]map[crapiv1alpha1.LokiTargetKind]int32{crapiv1alpha1.ModeKindMonolithic: {crapiv1alpha1.TargetKindAllInOne: 1}}
			service = loki.NewLokiGatewayService(crObj, string(crapiv1alpha1.ModeKindMonolithic))
			kubeClient = k8sfake.NewSimpleClientset()
			crClient = crfakeclients.NewSimpleClientset()
			fakeCtrl = newFakeController(kubeClient, crClient)
		})
		It("Create Loki, External ConfigMap", func() {
			configMap = fakeExternalConfigMap()
			crObj.Spec.ConfigMap = configMap.GetName()
			statefulSet = loki.NewStatefulSet(crObj, crapiv1alpha1.ModeKindMonolithic, crapiv1alpha1.TargetKindAllInOne, configMap)
			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.fixture.AddConfigMapLister(configMap))

			fakeCtrl.fixture.RecordKubeActions(testing.ExpectCreateStatefulSetAction(statefulSet), testing.ExpectCreateServiceAction(service))

			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())

			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())

		})
		It("Create Loki, Internal Default ConfigMap", func() {
			var err error
			configMap, err = loki.NewLokiConfigMap(crObj, string(crapiv1alpha1.ModeKindMonolithic))
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			statefulSet = loki.NewStatefulSet(crObj, crapiv1alpha1.ModeKindMonolithic, crapiv1alpha1.TargetKindAllInOne, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			fakeCtrl.fixture.RecordKubeActions(testing.ExpectCreateConfigMapAction(configMap), testing.ExpectCreateStatefulSetAction(statefulSet),
				testing.ExpectCreateServiceAction(service))

			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())

		})
		It("Create Loki, Internal Custom ConfigMap", func() {
			var err error
			configMap, err = loki.NewLokiConfigMap(crObj, string(crapiv1alpha1.ModeKindMonolithic))
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			statefulSet = loki.NewStatefulSet(crObj, crapiv1alpha1.ModeKindMonolithic, crapiv1alpha1.TargetKindAllInOne, configMap)

			gomega.Ω(fakeCtrl.fixture.AddCustomResourceLister(crObj)).ShouldNot(gomega.HaveOccurred())
			fakeCtrl.fixture.RecordKubeActions(testing.ExpectCreateConfigMapAction(configMap), testing.ExpectCreateStatefulSetAction(statefulSet),
				testing.ExpectCreateServiceAction(service))

			key, err := getKey(crObj)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(fakeCtrl.runController(key, true)).ShouldNot(gomega.HaveOccurred())

		})
		It("Update Loki, Modify Replicas", func() {})
		It("Update Loki, Modify Do Nothing", func() {})
	})

	Context("Reconcile", func() {
	})

	Context("Reconcile", func() {})

})

func newLoki(lokiName string) *crapiv1alpha1.Loki {
	return &crapiv1alpha1.Loki{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Loki",
			APIVersion: crapiv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: apicorev1.NamespaceDefault,
			Name:      lokiName,
		},
		Spec:   crapiv1alpha1.LokiSpec{},
		Status: crapiv1alpha1.LokiStatus{},
	}
}

func fakeExternalConfigMap() *apicorev1.ConfigMap {
	return &apicorev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: apicorev1.NamespaceDefault,
		},
		Data:       map[string]string{"test": "test"},
		BinaryData: nil,
	}
}

func getKey(obj *crapiv1alpha1.Loki) (string, error) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return "", fmt.Errorf("unexpected error getting key for foo %v: %v", obj.Name, err)
	}
	return key, nil
}
