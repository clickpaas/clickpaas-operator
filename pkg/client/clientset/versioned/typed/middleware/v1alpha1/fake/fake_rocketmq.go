/*
Copyright The clickpaas-controller Authors.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

// FakeRocketmqs implements RocketmqInterface
type FakeRocketmqs struct {
	Fake *FakeMiddlewareV1alpha1
	ns   string
}

var rocketmqsResource = schema.GroupVersionResource{Group: "middleware.clickpaas.cn", Version: "v1alpha1", Resource: "rocketmqs"}

var rocketmqsKind = schema.GroupVersionKind{Group: "middleware.clickpaas.cn", Version: "v1alpha1", Kind: "Rocketmq"}

// Get takes name of the rocketmq, and returns the corresponding rocketmq object, and an error if there is any.
func (c *FakeRocketmqs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Rocketmq, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rocketmqsResource, c.ns, name), &v1alpha1.Rocketmq{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rocketmq), err
}

// List takes label and field selectors, and returns the list of Rocketmqs that match those selectors.
func (c *FakeRocketmqs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RocketmqList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rocketmqsResource, rocketmqsKind, c.ns, opts), &v1alpha1.RocketmqList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RocketmqList{ListMeta: obj.(*v1alpha1.RocketmqList).ListMeta}
	for _, item := range obj.(*v1alpha1.RocketmqList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rocketmqs.
func (c *FakeRocketmqs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rocketmqsResource, c.ns, opts))

}

// Create takes the representation of a rocketmq and creates it.  Returns the server's representation of the rocketmq, and an error, if there is any.
func (c *FakeRocketmqs) Create(ctx context.Context, rocketmq *v1alpha1.Rocketmq, opts v1.CreateOptions) (result *v1alpha1.Rocketmq, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rocketmqsResource, c.ns, rocketmq), &v1alpha1.Rocketmq{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rocketmq), err
}

// Update takes the representation of a rocketmq and updates it. Returns the server's representation of the rocketmq, and an error, if there is any.
func (c *FakeRocketmqs) Update(ctx context.Context, rocketmq *v1alpha1.Rocketmq, opts v1.UpdateOptions) (result *v1alpha1.Rocketmq, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rocketmqsResource, c.ns, rocketmq), &v1alpha1.Rocketmq{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rocketmq), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRocketmqs) UpdateStatus(ctx context.Context, rocketmq *v1alpha1.Rocketmq, opts v1.UpdateOptions) (*v1alpha1.Rocketmq, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(rocketmqsResource, "status", c.ns, rocketmq), &v1alpha1.Rocketmq{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rocketmq), err
}

// Delete takes name of the rocketmq and deletes it. Returns an error if one occurs.
func (c *FakeRocketmqs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rocketmqsResource, c.ns, name), &v1alpha1.Rocketmq{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRocketmqs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rocketmqsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.RocketmqList{})
	return err
}

// Patch applies the patch and returns the patched rocketmq.
func (c *FakeRocketmqs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Rocketmq, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rocketmqsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Rocketmq{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rocketmq), err
}
