package controller_test

import (
	"github.com/l0calh0st/loki-operator/pkg/controller"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("Base", func() {
	var (
		baseController controller.Base
		hook           = controller.NewEventsHook(1)
		anotherHook    = controller.NewEventsHook(1)
	)
	BeforeEach(func() {
		baseController = controller.NewControllerBase()
	})

	Context("Add Hook", func() {
		It("hook is not register into controller", func() {
			gomega.Ω(baseController.AddHook(hook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(baseController.AddHook(anotherHook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(len(baseController.GetHooks())).To(gomega.Equal(2))
		})

		It("hook is already register into controller", func() {
			gomega.Ω(baseController.AddHook(hook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(baseController.AddHook(hook)).Should(gomega.HaveOccurred())
			gomega.Ω(len(baseController.GetHooks())).To(gomega.Equal(1))
		})
	})

	Context("Remove Hook", func() {
		It("hook is already register to controller, delete success", func() {
			gomega.Ω(baseController.AddHook(hook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(len(baseController.GetHooks())).To(gomega.Equal(1))
			gomega.Ω(baseController.RemoveHook(hook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(len(baseController.GetHooks())).To(gomega.Equal(0))
		})
		It("hook is not register to controller, delete failed", func() {
			gomega.Ω(baseController.AddHook(hook)).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(baseController.RemoveHook(anotherHook)).Should(gomega.HaveOccurred())
			gomega.Ω(len(baseController.GetHooks())).To(gomega.Equal(1))
		})
	})
})
