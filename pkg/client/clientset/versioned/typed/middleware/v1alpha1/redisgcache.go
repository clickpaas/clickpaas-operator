/*
Copyright The clickpaas-controller Authors.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	scheme "l0calh0st.cn/clickpaas-operator/pkg/client/clientset/versioned/scheme"
)

// RedisGCachesGetter has a method to return a RedisGCacheInterface.
// A group's client should implement this interface.
type RedisGCachesGetter interface {
	RedisGCaches(namespace string) RedisGCacheInterface
}

// RedisGCacheInterface has methods to work with RedisGCache resources.
type RedisGCacheInterface interface {
	Create(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.CreateOptions) (*v1alpha1.RedisGCache, error)
	Update(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.UpdateOptions) (*v1alpha1.RedisGCache, error)
	UpdateStatus(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.UpdateOptions) (*v1alpha1.RedisGCache, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.RedisGCache, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RedisGCacheList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.RedisGCache, err error)
	RedisGCacheExpansion
}

// redisGCaches implements RedisGCacheInterface
type redisGCaches struct {
	client rest.Interface
	ns     string
}

// newRedisGCaches returns a RedisGCaches
func newRedisGCaches(c *MiddlewareV1alpha1Client, namespace string) *redisGCaches {
	return &redisGCaches{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the redisGCache, and returns the corresponding redisGCache object, and an error if there is any.
func (c *redisGCaches) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.RedisGCache, err error) {
	result = &v1alpha1.RedisGCache{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisgcaches").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RedisGCaches that match those selectors.
func (c *redisGCaches) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RedisGCacheList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.RedisGCacheList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisgcaches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested redisGCaches.
func (c *redisGCaches) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("redisgcaches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a redisGCache and creates it.  Returns the server's representation of the redisGCache, and an error, if there is any.
func (c *redisGCaches) Create(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.CreateOptions) (result *v1alpha1.RedisGCache, err error) {
	result = &v1alpha1.RedisGCache{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("redisgcaches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(redisGCache).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a redisGCache and updates it. Returns the server's representation of the redisGCache, and an error, if there is any.
func (c *redisGCaches) Update(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.UpdateOptions) (result *v1alpha1.RedisGCache, err error) {
	result = &v1alpha1.RedisGCache{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redisgcaches").
		Name(redisGCache.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(redisGCache).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *redisGCaches) UpdateStatus(ctx context.Context, redisGCache *v1alpha1.RedisGCache, opts v1.UpdateOptions) (result *v1alpha1.RedisGCache, err error) {
	result = &v1alpha1.RedisGCache{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redisgcaches").
		Name(redisGCache.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(redisGCache).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the redisGCache and deletes it. Returns an error if one occurs.
func (c *redisGCaches) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisgcaches").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *redisGCaches) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisgcaches").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched redisGCache.
func (c *redisGCaches) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.RedisGCache, err error) {
	result = &v1alpha1.RedisGCache{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("redisgcaches").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
