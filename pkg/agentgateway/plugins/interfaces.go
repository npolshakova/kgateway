package plugins

import (
	"github.com/agentgateway/agentgateway/go/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PolicyPlugin defines the base interface for all policy plugins
type PolicyPlugin interface {
	// GroupKind returns the GroupKind of the policy this plugin handles
	GroupKind() schema.GroupKind

	// Name returns the name of the plugin
	Name() string

	// ApplyPolicies generates agentgateway policies for the given common collections
	ApplyPolicies() []ADPPolicy
}

// ADPPolicy wraps an ADP policy for collection handling
type ADPPolicy struct {
	Policy *api.Policy
	errors []error
}
