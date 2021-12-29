package controller_test

import (
	"github.com/l0calh0st/loki-operator/pkg/controller"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("Event", func() {

	var (
		eventsHook controller.EventsHook
		fakeObj    string
	)

	BeforeEach(func() {
		fakeObj = "test"
		eventsHook = controller.NewEventsHook(10)
	})

	Context("Test EventHooks", func() {
		It("Add Object", func() {
			eventsHook.OnAdd(fakeObj)
			gomega.Ω(<-eventsHook.GetEventsChan()).To(gomega.Equal(controller.Event{
				Type:   controller.EventAdded,
				Object: fakeObj,
			}))
		})

		It("Update Object", func() {
			eventsHook.OnUpdate(fakeObj)
			gomega.Ω(<-eventsHook.GetEventsChan()).To(gomega.Equal(controller.Event{
				Type:   controller.EventUpdated,
				Object: fakeObj,
			}))
		})

		It("Delete Object", func() {
			eventsHook.OnDelete(fakeObj)
			gomega.Ω(<-eventsHook.GetEventsChan()).To(gomega.Equal(controller.Event{
				Type:   controller.EventDeleted,
				Object: fakeObj,
			}))
		})
	})
})
