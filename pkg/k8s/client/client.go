package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type Client interface {
	Get(obj runtime.Object) error
	Update(obj runtime.Object) error
	Delete(obj runtime.Object)
}

type client struct {
	kubeClient kubernetes.Interface
}
