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

// FakeMongoClusters implements MongoClusterInterface
type FakeMongoClusters struct {
	Fake *FakeMiddlewareV1alpha1
	ns   string
}

var mongoclustersResource = schema.GroupVersionResource{Group: "middleware.l0calh0st.cn", Version: "v1alpha1", Resource: "mongoclusters"}

var mongoclustersKind = schema.GroupVersionKind{Group: "middleware.l0calh0st.cn", Version: "v1alpha1", Kind: "MongoCluster"}

// Get takes name of the mongoCluster, and returns the corresponding mongoCluster object, and an error if there is any.
func (c *FakeMongoClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MongoCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mongoclustersResource, c.ns, name), &v1alpha1.MongoCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoCluster), err
}

// List takes label and field selectors, and returns the list of MongoClusters that match those selectors.
func (c *FakeMongoClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MongoClusterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mongoclustersResource, mongoclustersKind, c.ns, opts), &v1alpha1.MongoClusterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MongoClusterList{ListMeta: obj.(*v1alpha1.MongoClusterList).ListMeta}
	for _, item := range obj.(*v1alpha1.MongoClusterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mongoClusters.
func (c *FakeMongoClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mongoclustersResource, c.ns, opts))

}

// Create takes the representation of a mongoCluster and creates it.  Returns the server's representation of the mongoCluster, and an error, if there is any.
func (c *FakeMongoClusters) Create(ctx context.Context, mongoCluster *v1alpha1.MongoCluster, opts v1.CreateOptions) (result *v1alpha1.MongoCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mongoclustersResource, c.ns, mongoCluster), &v1alpha1.MongoCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoCluster), err
}

// Update takes the representation of a mongoCluster and updates it. Returns the server's representation of the mongoCluster, and an error, if there is any.
func (c *FakeMongoClusters) Update(ctx context.Context, mongoCluster *v1alpha1.MongoCluster, opts v1.UpdateOptions) (result *v1alpha1.MongoCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mongoclustersResource, c.ns, mongoCluster), &v1alpha1.MongoCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoCluster), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMongoClusters) UpdateStatus(ctx context.Context, mongoCluster *v1alpha1.MongoCluster, opts v1.UpdateOptions) (*v1alpha1.MongoCluster, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mongoclustersResource, "status", c.ns, mongoCluster), &v1alpha1.MongoCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoCluster), err
}

// Delete takes name of the mongoCluster and deletes it. Returns an error if one occurs.
func (c *FakeMongoClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mongoclustersResource, c.ns, name), &v1alpha1.MongoCluster{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMongoClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mongoclustersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.MongoClusterList{})
	return err
}

// Patch applies the patch and returns the patched mongoCluster.
func (c *FakeMongoClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MongoCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mongoclustersResource, c.ns, name, pt, data, subresources...), &v1alpha1.MongoCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoCluster), err
}