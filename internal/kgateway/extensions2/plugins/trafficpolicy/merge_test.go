package trafficpolicy

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugins/trafficpolicy/agentgateway"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestMergePoliciesPreservesErrors(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	gk := schema.GroupKind{Group: "test", Kind: "TrafficPolicy"}

	p1 := ir.PolicyAtt{
		GroupKind: gk,
		PolicyRef: &ir.AttachedPolicyRef{Name: "p1"},
		PolicyIr:  &TrafficPolicy{ct: time.Now()},
		Errors:    []error{err1},
	}
	p2 := ir.PolicyAtt{
		GroupKind: gk,
		PolicyRef: &ir.AttachedPolicyRef{Name: "p2"},
		PolicyIr:  &TrafficPolicy{ct: time.Now().Add(time.Minute)},
		Errors:    []error{err2},
	}

	merged := policy.MergePolicies([]ir.PolicyAtt{p1, p2}, MergeTrafficPolicies)
	require.Len(t, merged.Errors, 2)
	assert.Contains(t, merged.Errors, err1)
	assert.Contains(t, merged.Errors, err2)
}

func TestMergeCompositeTrafficPolicies(t *testing.T) {
	gk := schema.GroupKind{Group: "test", Kind: "TrafficPolicy"}

	envoyPolicy1 := &TrafficPolicy{ct: time.Now()}
	envoyPolicy2 := &TrafficPolicy{ct: time.Now().Add(time.Minute)}

	// Import agentgateway package to access TrafficPolicy
	agwPolicy1 := &agentgateway.TrafficPolicy{}
	agwPolicy2 := &agentgateway.TrafficPolicy{}

	// Create composite policies
	composite1 := &CompositeTrafficPolicy{
		ct:                 time.Now(),
		EnvoyPolicy:        envoyPolicy1,
		AgentGatewayPolicy: agwPolicy1,
	}

	composite2 := &CompositeTrafficPolicy{
		ct:                 time.Now().Add(time.Minute),
		EnvoyPolicy:        envoyPolicy2,
		AgentGatewayPolicy: agwPolicy2,
	}

	p1 := ir.PolicyAtt{
		GroupKind: gk,
		PolicyRef: &ir.AttachedPolicyRef{Name: "composite1"},
		PolicyIr:  composite1,
	}
	p2 := ir.PolicyAtt{
		GroupKind: gk,
		PolicyRef: &ir.AttachedPolicyRef{Name: "composite2"},
		PolicyIr:  composite2,
	}

	merged := policy.MergePolicies([]ir.PolicyAtt{p1, p2}, MergeCompositeTrafficPolicies)

	// Verify merged result is a composite policy
	mergedComposite, ok := merged.PolicyIr.(*CompositeTrafficPolicy)
	require.True(t, ok, "Merged policy should be a CompositeTrafficPolicy")

	// Verify both envoy and agentgateway policies are present
	assert.NotNil(t, mergedComposite.EnvoyPolicy, "Merged composite should have an envoy policy")
	assert.NotNil(t, mergedComposite.AgentGatewayPolicy, "Merged composite should have an agentgateway policy")
}
