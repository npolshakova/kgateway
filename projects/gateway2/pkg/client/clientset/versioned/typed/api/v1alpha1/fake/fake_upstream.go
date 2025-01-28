// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	apiv1alpha1 "github.com/kgateway-dev/kgateway/projects/gateway2/api/applyconfiguration/api/v1alpha1"
	v1alpha1 "github.com/kgateway-dev/kgateway/projects/gateway2/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeUpstreams implements UpstreamInterface
type FakeUpstreams struct {
	Fake *FakeGatewayV1alpha1
	ns   string
}

var upstreamsResource = v1alpha1.SchemeGroupVersion.WithResource("upstreams")

var upstreamsKind = v1alpha1.SchemeGroupVersion.WithKind("Upstream")

// Get takes name of the upstream, and returns the corresponding upstream object, and an error if there is any.
func (c *FakeUpstreams) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Upstream, err error) {
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(upstreamsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// List takes label and field selectors, and returns the list of Upstreams that match those selectors.
func (c *FakeUpstreams) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.UpstreamList, err error) {
	emptyResult := &v1alpha1.UpstreamList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(upstreamsResource, upstreamsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.UpstreamList{ListMeta: obj.(*v1alpha1.UpstreamList).ListMeta}
	for _, item := range obj.(*v1alpha1.UpstreamList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested upstreams.
func (c *FakeUpstreams) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(upstreamsResource, c.ns, opts))

}

// Create takes the representation of a upstream and creates it.  Returns the server's representation of the upstream, and an error, if there is any.
func (c *FakeUpstreams) Create(ctx context.Context, upstream *v1alpha1.Upstream, opts v1.CreateOptions) (result *v1alpha1.Upstream, err error) {
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(upstreamsResource, c.ns, upstream, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// Update takes the representation of a upstream and updates it. Returns the server's representation of the upstream, and an error, if there is any.
func (c *FakeUpstreams) Update(ctx context.Context, upstream *v1alpha1.Upstream, opts v1.UpdateOptions) (result *v1alpha1.Upstream, err error) {
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(upstreamsResource, c.ns, upstream, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeUpstreams) UpdateStatus(ctx context.Context, upstream *v1alpha1.Upstream, opts v1.UpdateOptions) (result *v1alpha1.Upstream, err error) {
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(upstreamsResource, "status", c.ns, upstream, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// Delete takes name of the upstream and deletes it. Returns an error if one occurs.
func (c *FakeUpstreams) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(upstreamsResource, c.ns, name, opts), &v1alpha1.Upstream{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeUpstreams) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(upstreamsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.UpstreamList{})
	return err
}

// Patch applies the patch and returns the patched upstream.
func (c *FakeUpstreams) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Upstream, err error) {
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(upstreamsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied upstream.
func (c *FakeUpstreams) Apply(ctx context.Context, upstream *apiv1alpha1.UpstreamApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Upstream, err error) {
	if upstream == nil {
		return nil, fmt.Errorf("upstream provided to Apply must not be nil")
	}
	data, err := json.Marshal(upstream)
	if err != nil {
		return nil, err
	}
	name := upstream.Name
	if name == nil {
		return nil, fmt.Errorf("upstream.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(upstreamsResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeUpstreams) ApplyStatus(ctx context.Context, upstream *apiv1alpha1.UpstreamApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Upstream, err error) {
	if upstream == nil {
		return nil, fmt.Errorf("upstream provided to Apply must not be nil")
	}
	data, err := json.Marshal(upstream)
	if err != nil {
		return nil, err
	}
	name := upstream.Name
	if name == nil {
		return nil, fmt.Errorf("upstream.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Upstream{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(upstreamsResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Upstream), err
}
