package plugins

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

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

	// GeneratePolicies generates ADP policies for the given inputs
	GeneratePolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error)
}

// PolicyInputs contains all the input collections needed by policy plugins
type PolicyInputs struct {
	Services          krt.Collection[*corev1.Service]
	InferencePools    krt.Collection[*inf.InferencePool]
	TrafficPolicies   krt.Collection[*v1alpha1.TrafficPolicy]
	GatewayExtensions krt.Collection[*v1alpha1.GatewayExtension]
	DomainSuffix      string
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
	GenerateAllPolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error)
}
