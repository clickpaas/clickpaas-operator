/*
Copyright The clickpaas-controller Authors.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	middlewarev1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	versioned "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned"
	internalinterfaces "l0calh0st.cn/clickpaas-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/client/listers/middleware/v1alpha1"
)

// LtsJobTrackerInformer provides access to a shared informer and lister for
// LtsJobTrackers.
type LtsJobTrackerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.LtsJobTrackerLister
}

type ltsJobTrackerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewLtsJobTrackerInformer constructs a new informer for LtsJobTracker type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewLtsJobTrackerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredLtsJobTrackerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredLtsJobTrackerInformer constructs a new informer for LtsJobTracker type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredLtsJobTrackerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MiddlewareV1alpha1().LtsJobTrackers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MiddlewareV1alpha1().LtsJobTrackers(namespace).Watch(context.TODO(), options)
			},
		},
		&middlewarev1alpha1.LtsJobTracker{},
		resyncPeriod,
		indexers,
	)
}

func (f *ltsJobTrackerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredLtsJobTrackerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *ltsJobTrackerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&middlewarev1alpha1.LtsJobTracker{}, f.defaultInformer)
}

func (f *ltsJobTrackerInformer) Lister() v1alpha1.LtsJobTrackerLister {
	return v1alpha1.NewLtsJobTrackerLister(f.Informer().GetIndexer())
}
