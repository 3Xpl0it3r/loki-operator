package crd_test

import (
	"context"
	"fmt"
	"github.com/l0calh0st/loki-operator/pkg/crd"
	"github.com/l0calh0st/loki-operator/pkg/crd/loki"
	"github.com/l0calh0st/loki-operator/pkg/crd/promtail"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	fakeextclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)



var _ = Describe("Register CRD", func() {
	var (
		fakeExtClient extclientset.Interface
		lokicrd = loki.NewCustomResourceDefine()
		promtailcrd = promtail.NewCustomResourceDefine()
	)
	BeforeSuite(func() {
		fakeExtClient = fakeextclientset.NewSimpleClientset()
	})

	Describe("Loki CRD Register" , func() {
		BeforeEach(func() {
			err := crd.RegisterCRDWithObj(fakeExtClient, lokicrd)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		AfterEach(func() {
			err := crd.UnRegisterCRD(fakeExtClient, lokicrd.GetName())
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		It("get loki crd", func() {
			tcrd,err := fakeExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), lokicrd.GetName(), metav1.GetOptions{})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(tcrd.GetName()).Should(gomega.Equal(lokicrd.GetName()))
		})
	})

	Describe("Promtail CRD Register" , func() {
		BeforeEach(func() {
			err := crd.RegisterCRDWithObj(fakeExtClient, promtailcrd)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		AfterEach(func() {
			err := crd.UnRegisterCRD(fakeExtClient, promtailcrd.GetName())
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		It("get promtail crd", func() {
			pcrd,err := fakeExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), promtailcrd.GetName(), metav1.GetOptions{})
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(pcrd.GetName()).Should(gomega.Equal(promtailcrd.GetName()))
		})
	})

})


type crdMatcher struct {
	crd *apiextensionsv1.CustomResourceDefinition
}

func (c crdMatcher) Match(actual interface{}) (success bool, err error) {
	actualCrd,ok := actual.(*apiextensionsv1.CustomResourceDefinition)
	if !ok {
		return false, fmt.Errorf("Except apiextensionsv1.CustomResourceDefinition, but got %T ", actual)
	}
	if c.crd.GetName() != actualCrd.GetName(){
		return false, fmt.Errorf("name is not same")
	}
	// ... todo
	return true, nil
}

func (c crdMatcher) FailureMessage(actual interface{}) (message string) {
	return ""
}

func (c crdMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return ""
}

func CustomResourceDefineMatch(crd *apiextensionsv1.CustomResourceDefinition)*crdMatcher{
	return &crdMatcher{crd: crd}
}
