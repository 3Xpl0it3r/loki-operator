package promtail

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
)

// promtailEventHandler
type promtailEventHandler struct {
	c *controller
}

func (h *promtailEventHandler) OnAdd(obj interface{}) {
	cr, ok := obj.(*crapiv1alpha1.Promtail)
	if !ok {
		return
	}
	crapiv1alpha1.WithDefaultsPromtail(cr)
	for _,hook := range h.c.GetHooks(){
		hook.OnAdd(cr)
	}
}

func (h *promtailEventHandler) OnUpdate(oldObj, newObj interface{}) {
	oldCr, ok := oldObj.(*crapiv1alpha1.Promtail)
	if !ok {
		return
	}
	newCr, ok := newObj.(*crapiv1alpha1.Promtail)
	if !ok {
		return
	}
	if oldCr.ResourceVersion == newCr.ResourceVersion {
		return
	}
	h.c.enqueueFunc(newCr)
	for _, hook := range h.c.GetHooks(){
		hook.OnUpdate(newCr)
	}
}

func (h *promtailEventHandler) OnDelete(obj interface{}) {
	cr, ok := obj.(*crapiv1alpha1.Promtail)
	if !ok {
		return
	}
	h.c.enqueueFunc(cr)
	for _, hook := range h.c.GetHooks(){
		hook.OnDelete(cr)
	}
}

func newPromtailEventHandler(c *controller) *promtailEventHandler {
	return &promtailEventHandler{c: c}
}
