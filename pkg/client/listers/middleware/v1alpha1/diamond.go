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

// DiamondLister helps list Diamonds.
// All objects returned here must be treated as read-only.
type DiamondLister interface {
	// List lists all Diamonds in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Diamond, err error)
	// Diamonds returns an object that can list and get Diamonds.
	Diamonds(namespace string) DiamondNamespaceLister
	DiamondListerExpansion
}

// diamondLister implements the DiamondLister interface.
type diamondLister struct {
	indexer cache.Indexer
}

// NewDiamondLister returns a new DiamondLister.
func NewDiamondLister(indexer cache.Indexer) DiamondLister {
	return &diamondLister{indexer: indexer}
}

// List lists all Diamonds in the indexer.
func (s *diamondLister) List(selector labels.Selector) (ret []*v1alpha1.Diamond, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Diamond))
	})
	return ret, err
}

// Diamonds returns an object that can list and get Diamonds.
func (s *diamondLister) Diamonds(namespace string) DiamondNamespaceLister {
	return diamondNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DiamondNamespaceLister helps list and get Diamonds.
// All objects returned here must be treated as read-only.
type DiamondNamespaceLister interface {
	// List lists all Diamonds in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Diamond, err error)
	// Get retrieves the Diamond from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Diamond, error)
	DiamondNamespaceListerExpansion
}

// diamondNamespaceLister implements the DiamondNamespaceLister
// interface.
type diamondNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Diamonds in the indexer for a given namespace.
func (s diamondNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Diamond, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Diamond))
	})
	return ret, err
}

// Get retrieves the Diamond from the indexer for a given namespace and name.
func (s diamondNamespaceLister) Get(name string) (*v1alpha1.Diamond, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("diamond"), name)
	}
	return obj.(*v1alpha1.Diamond), nil
}
