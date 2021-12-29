package loki

import (
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
)

type lokiEventHandler struct {
	c *controller
}


func (h *lokiEventHandler) OnAdd(obj interface{}) {
	loki, ok := obj.(*crapiv1alpha1.Loki)
	if !ok {
		return
	}
	crapiv1alpha1.WithDefaultsLoki(loki)
	h.c.enqueueFunc(loki)
	for _,hook := range h.c.GetHooks(){
		hook.OnAdd(loki)
	}
}

func (h *lokiEventHandler) OnUpdate(oldObj, newObj interface{}) {
	oldLoki, ok := oldObj.(*crapiv1alpha1.Loki)
	if !ok {
		return
	}
	newLoki, ok := newObj.(*crapiv1alpha1.Loki)
	if !ok {
		return
	}
	if oldLoki.ResourceVersion == newLoki.ResourceVersion{
		return
	}
	h.c.enqueueFunc(newLoki)
	for _,hook := range h.c.GetHooks(){
		hook.OnAdd(newLoki)
	}
}

func (h *lokiEventHandler) OnDelete(obj interface{}) {
	loki,ok := obj.(*crapiv1alpha1.Loki)
	if !ok {
		return
	}
	h.c.enqueueFunc(loki)
	for _,hook := range h.c.GetHooks(){
		hook.OnAdd(loki)
	}
}

func newLokiEventHandler(c *controller)*lokiEventHandler{
	return &lokiEventHandler{c: c}
}

