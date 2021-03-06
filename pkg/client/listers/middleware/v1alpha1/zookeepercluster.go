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

// ZookeeperClusterLister helps list ZookeeperClusters.
// All objects returned here must be treated as read-only.
type ZookeeperClusterLister interface {
	// List lists all ZookeeperClusters in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ZookeeperCluster, err error)
	// ZookeeperClusters returns an object that can list and get ZookeeperClusters.
	ZookeeperClusters(namespace string) ZookeeperClusterNamespaceLister
	ZookeeperClusterListerExpansion
}

// zookeeperClusterLister implements the ZookeeperClusterLister interface.
type zookeeperClusterLister struct {
	indexer cache.Indexer
}

// NewZookeeperClusterLister returns a new ZookeeperClusterLister.
func NewZookeeperClusterLister(indexer cache.Indexer) ZookeeperClusterLister {
	return &zookeeperClusterLister{indexer: indexer}
}

// List lists all ZookeeperClusters in the indexer.
func (s *zookeeperClusterLister) List(selector labels.Selector) (ret []*v1alpha1.ZookeeperCluster, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ZookeeperCluster))
	})
	return ret, err
}

// ZookeeperClusters returns an object that can list and get ZookeeperClusters.
func (s *zookeeperClusterLister) ZookeeperClusters(namespace string) ZookeeperClusterNamespaceLister {
	return zookeeperClusterNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ZookeeperClusterNamespaceLister helps list and get ZookeeperClusters.
// All objects returned here must be treated as read-only.
type ZookeeperClusterNamespaceLister interface {
	// List lists all ZookeeperClusters in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ZookeeperCluster, err error)
	// Get retrieves the ZookeeperCluster from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ZookeeperCluster, error)
	ZookeeperClusterNamespaceListerExpansion
}

// zookeeperClusterNamespaceLister implements the ZookeeperClusterNamespaceLister
// interface.
type zookeeperClusterNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ZookeeperClusters in the indexer for a given namespace.
func (s zookeeperClusterNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ZookeeperCluster, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ZookeeperCluster))
	})
	return ret, err
}

// Get retrieves the ZookeeperCluster from the indexer for a given namespace and name.
func (s zookeeperClusterNamespaceLister) Get(name string) (*v1alpha1.ZookeeperCluster, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("zookeepercluster"), name)
	}
	return obj.(*v1alpha1.ZookeeperCluster), nil
}
