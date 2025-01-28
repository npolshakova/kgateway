// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/kgateway-dev/kgateway/projects/gateway2/pkg/client/clientset/versioned/typed/api/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeGatewayV1alpha1 struct {
	*testing.Fake
}

func (c *FakeGatewayV1alpha1) DirectResponses(namespace string) v1alpha1.DirectResponseInterface {
	return &FakeDirectResponses{c, namespace}
}

func (c *FakeGatewayV1alpha1) GatewayParameterses(namespace string) v1alpha1.GatewayParametersInterface {
	return &FakeGatewayParameterses{c, namespace}
}

func (c *FakeGatewayV1alpha1) HttpListenerPolicies(namespace string) v1alpha1.HttpListenerPolicyInterface {
	return &FakeHttpListenerPolicies{c, namespace}
}

func (c *FakeGatewayV1alpha1) ListenerPolicies(namespace string) v1alpha1.ListenerPolicyInterface {
	return &FakeListenerPolicies{c, namespace}
}

func (c *FakeGatewayV1alpha1) RoutePolicies(namespace string) v1alpha1.RoutePolicyInterface {
	return &FakeRoutePolicies{c, namespace}
}

func (c *FakeGatewayV1alpha1) Upstreams(namespace string) v1alpha1.UpstreamInterface {
	return &FakeUpstreams{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeGatewayV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
