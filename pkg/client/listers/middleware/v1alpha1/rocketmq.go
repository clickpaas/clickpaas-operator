/*
Copyright The clickpaas-controller Authors.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

// RocketmqLister helps list Rocketmqs.
// All objects returned here must be treated as read-only.
type RocketmqLister interface {
	// List lists all Rocketmqs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Rocketmq, err error)
	// Rocketmqs returns an object that can list and get Rocketmqs.
	Rocketmqs(namespace string) RocketmqNamespaceLister
	RocketmqListerExpansion
}

// rocketmqLister implements the RocketmqLister interface.
type rocketmqLister struct {
	indexer cache.Indexer
}

// NewRocketmqLister returns a new RocketmqLister.
func NewRocketmqLister(indexer cache.Indexer) RocketmqLister {
	return &rocketmqLister{indexer: indexer}
}

// List lists all Rocketmqs in the indexer.
func (s *rocketmqLister) List(selector labels.Selector) (ret []*v1alpha1.Rocketmq, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Rocketmq))
	})
	return ret, err
}

// Rocketmqs returns an object that can list and get Rocketmqs.
func (s *rocketmqLister) Rocketmqs(namespace string) RocketmqNamespaceLister {
	return rocketmqNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RocketmqNamespaceLister helps list and get Rocketmqs.
// All objects returned here must be treated as read-only.
type RocketmqNamespaceLister interface {
	// List lists all Rocketmqs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Rocketmq, err error)
	// Get retrieves the Rocketmq from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Rocketmq, error)
	RocketmqNamespaceListerExpansion
}

// rocketmqNamespaceLister implements the RocketmqNamespaceLister
// interface.
type rocketmqNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Rocketmqs in the indexer for a given namespace.
func (s rocketmqNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Rocketmq, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Rocketmq))
	})
	return ret, err
}

// Get retrieves the Rocketmq from the indexer for a given namespace and name.
func (s rocketmqNamespaceLister) Get(name string) (*v1alpha1.Rocketmq, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("rocketmq"), name)
	}
	return obj.(*v1alpha1.Rocketmq), nil
}
