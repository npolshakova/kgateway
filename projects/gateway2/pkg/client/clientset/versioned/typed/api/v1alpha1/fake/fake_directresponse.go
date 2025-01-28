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

// FakeDirectResponses implements DirectResponseInterface
type FakeDirectResponses struct {
	Fake *FakeGatewayV1alpha1
	ns   string
}

var directresponsesResource = v1alpha1.SchemeGroupVersion.WithResource("directresponses")

var directresponsesKind = v1alpha1.SchemeGroupVersion.WithKind("DirectResponse")

// Get takes name of the directResponse, and returns the corresponding directResponse object, and an error if there is any.
func (c *FakeDirectResponses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.DirectResponse, err error) {
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(directresponsesResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// List takes label and field selectors, and returns the list of DirectResponses that match those selectors.
func (c *FakeDirectResponses) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DirectResponseList, err error) {
	emptyResult := &v1alpha1.DirectResponseList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(directresponsesResource, directresponsesKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DirectResponseList{ListMeta: obj.(*v1alpha1.DirectResponseList).ListMeta}
	for _, item := range obj.(*v1alpha1.DirectResponseList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested directResponses.
func (c *FakeDirectResponses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(directresponsesResource, c.ns, opts))

}

// Create takes the representation of a directResponse and creates it.  Returns the server's representation of the directResponse, and an error, if there is any.
func (c *FakeDirectResponses) Create(ctx context.Context, directResponse *v1alpha1.DirectResponse, opts v1.CreateOptions) (result *v1alpha1.DirectResponse, err error) {
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(directresponsesResource, c.ns, directResponse, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// Update takes the representation of a directResponse and updates it. Returns the server's representation of the directResponse, and an error, if there is any.
func (c *FakeDirectResponses) Update(ctx context.Context, directResponse *v1alpha1.DirectResponse, opts v1.UpdateOptions) (result *v1alpha1.DirectResponse, err error) {
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(directresponsesResource, c.ns, directResponse, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDirectResponses) UpdateStatus(ctx context.Context, directResponse *v1alpha1.DirectResponse, opts v1.UpdateOptions) (result *v1alpha1.DirectResponse, err error) {
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(directresponsesResource, "status", c.ns, directResponse, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// Delete takes name of the directResponse and deletes it. Returns an error if one occurs.
func (c *FakeDirectResponses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(directresponsesResource, c.ns, name, opts), &v1alpha1.DirectResponse{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDirectResponses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(directresponsesResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DirectResponseList{})
	return err
}

// Patch applies the patch and returns the patched directResponse.
func (c *FakeDirectResponses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.DirectResponse, err error) {
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(directresponsesResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied directResponse.
func (c *FakeDirectResponses) Apply(ctx context.Context, directResponse *apiv1alpha1.DirectResponseApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.DirectResponse, err error) {
	if directResponse == nil {
		return nil, fmt.Errorf("directResponse provided to Apply must not be nil")
	}
	data, err := json.Marshal(directResponse)
	if err != nil {
		return nil, err
	}
	name := directResponse.Name
	if name == nil {
		return nil, fmt.Errorf("directResponse.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(directresponsesResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeDirectResponses) ApplyStatus(ctx context.Context, directResponse *apiv1alpha1.DirectResponseApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.DirectResponse, err error) {
	if directResponse == nil {
		return nil, fmt.Errorf("directResponse provided to Apply must not be nil")
	}
	data, err := json.Marshal(directResponse)
	if err != nil {
		return nil, err
	}
	name := directResponse.Name
	if name == nil {
		return nil, fmt.Errorf("directResponse.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.DirectResponse{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(directresponsesResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.DirectResponse), err
}
