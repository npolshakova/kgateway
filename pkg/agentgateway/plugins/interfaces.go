package plugins

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
)

// PolicyInput represents a collection input that can be used by policy plugins
// This allows plugins to accept different types of krt.Collection without type conflicts
// Plugins will use type assertion to convert to their specific collection type
type PolicyInput interface{}

// PolicyType defines the different types of policies that can be handled
type PolicyType string

const (
	PolicyTypeA2A           PolicyType = "a2a"
	PolicyTypeInferencePool PolicyType = "inference-pool"
	PolicyTypeTraffic       PolicyType = "traffic"
)

// PolicyPlugin defines the base interface for all policy plugins
type PolicyPlugin interface {
	// PolicyType returns the type of policy this plugin handles
	PolicyType() PolicyType

	// Name returns the name of the plugin
	Name() string

	// GeneratePolicies generates ADP policies for the given common collections
	GeneratePolicies(ctx krt.HandlerContext, agentgatewayCol *AgwCollections, policyInput PolicyInput) ([]ADPPolicy, error)
}

// ADPPolicy wraps an ADP policy for collection handling
type ADPPolicy struct {
	Policy *api.Policy
}

// PolicyManager coordinates all policy plugins
type PolicyManager interface {
	// RegisterPlugin registers a policy plugin
	RegisterPlugin(plugin PolicyPlugin) error

	// GetPluginsByType returns all plugins of a specific type
	GetPluginsByType(policyType PolicyType) []PolicyPlugin

	// GenerateAllPolicies generates policies from all registered ADP plugins
	GenerateAllPolicies(ctx krt.HandlerContext, agw *AgwCollections, policyInput PolicyInput) ([]ADPPolicy, error)
}
