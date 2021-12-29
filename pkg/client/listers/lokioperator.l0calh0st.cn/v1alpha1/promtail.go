/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PromtailLister helps list Promtails.
// All objects returned here must be treated as read-only.
type PromtailLister interface {
	// List lists all Promtails in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Promtail, err error)
	// Promtails returns an object that can list and get Promtails.
	Promtails(namespace string) PromtailNamespaceLister
	PromtailListerExpansion
}

// promtailLister implements the PromtailLister interface.
type promtailLister struct {
	indexer cache.Indexer
}

// NewPromtailLister returns a new PromtailLister.
func NewPromtailLister(indexer cache.Indexer) PromtailLister {
	return &promtailLister{indexer: indexer}
}

// List lists all Promtails in the indexer.
func (s *promtailLister) List(selector labels.Selector) (ret []*v1alpha1.Promtail, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Promtail))
	})
	return ret, err
}

// Promtails returns an object that can list and get Promtails.
func (s *promtailLister) Promtails(namespace string) PromtailNamespaceLister {
	return promtailNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PromtailNamespaceLister helps list and get Promtails.
// All objects returned here must be treated as read-only.
type PromtailNamespaceLister interface {
	// List lists all Promtails in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Promtail, err error)
	// Get retrieves the Promtail from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Promtail, error)
	PromtailNamespaceListerExpansion
}

// promtailNamespaceLister implements the PromtailNamespaceLister
// interface.
type promtailNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Promtails in the indexer for a given namespace.
func (s promtailNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Promtail, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Promtail))
	})
	return ret, err
}

// Get retrieves the Promtail from the indexer for a given namespace and name.
func (s promtailNamespaceLister) Get(name string) (*v1alpha1.Promtail, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("promtail"), name)
	}
	return obj.(*v1alpha1.Promtail), nil
}
