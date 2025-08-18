package plugins

import (
	"slices"

	"github.com/agentgateway/agentgateway/go/api"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PolicyWrapper struct {
	// A reference to the original policy object
	ir.ObjectSource `json:",inline"`
	// The policy object itself. TODO: we can probably remove this
	Policy metav1.Object

	// Errors processing it for status.
	// note: these errors are based on policy itself, regardless of whether it's attached to a resource.
	// Errors should be formatted for users, so do not include internal lib errors.
	// Instead use a well-defined error such as ErrInvalidConfig
	Errors []error

	// The IR of the policy objects. ideally with structural errors removed.
	// Opaque to us other than metadata.
	PolicyIR ir.PolicyIR
}

func (c PolicyWrapper) ResourceName() string {
	return c.ObjectSource.ResourceName()
}

func versionEquals(a, b metav1.Object) bool {
	var versionEquals bool
	if a.GetGeneration() != 0 && b.GetGeneration() != 0 {
		versionEquals = a.GetGeneration() == b.GetGeneration()
	} else {
		versionEquals = a.GetResourceVersion() == b.GetResourceVersion()
	}
	return versionEquals && a.GetUID() == b.GetUID()
}

func (c PolicyWrapper) Equals(in PolicyWrapper) bool {
	if c.ObjectSource != in.ObjectSource {
		return false
	}

	if !slices.EqualFunc(c.Errors, in.Errors, func(e1, e2 error) bool {
		if e1 == nil && e2 != nil {
			return false
		}
		if e1 != nil && e2 == nil {
			return false
		}
		if (e1 != nil && e2 != nil) && e1.Error() != e2.Error() {
			return false
		}

		return true
	}) {
		return false
	}

	return versionEquals(c.Policy, in.Policy) && c.PolicyIR.Equals(in.PolicyIR)
}

type PolicyPlugin struct {
	Policies krt.Collection[PolicyWrapper]
}

// ApplyPolicies extracts all policies from the collection
func (p *PolicyPlugin) ApplyPolicies() []ADPPolicy {
	var allPolicies []ADPPolicy
	for _, wrapper := range p.Policies.List() {
		if policyPass, ok := wrapper.PolicyIR.(PolicyPluginPass); ok {
			policies := policyPass.ApplyPolicies()
			allPolicies = append(allPolicies, policies...)
		}
	}
	return allPolicies
}

// PolicyPluginPass represents a single translation pass for translating agentgateway policies.
type PolicyPluginPass interface {
	ApplyPolicies() []ADPPolicy
}

// ADPPolicy wraps an ADP policy for collection handling
type ADPPolicy struct {
	Policy *api.Policy
}

func (p *ADPPolicy) Equals(in *ADPPolicy) bool {
	if p == nil && in == nil {
		return true
	}
	if p.Policy == nil || in.Policy == nil {
		return false
	}
	return proto.Equal(p.Policy, in.Policy)
}
