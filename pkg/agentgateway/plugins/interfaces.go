package plugins

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	agwir "github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
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

	// Priority returns the priority of the plugin (lower numbers have higher priority)
	Priority() int
}

// ADPPolicyPlugin defines the interface for plugins that generate ADP (Agent Data Plane) policies
type ADPPolicyPlugin interface {
	PolicyPlugin

	// GeneratePolicies generates ADP policies for the given inputs
	GeneratePolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error)
}

// AgentGatewayPolicyPlugin defines the interface for plugins that participate in agent gateway translation
type AgentGatewayPolicyPlugin interface {
	PolicyPlugin

	// NewTranslationPass creates a new translation pass for this policy
	NewTranslationPass(reporter reports.Reporter) agwir.AgentGatewayTranslationPass
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

// A2APolicyPlugin handles agent-to-agent communication policies
type A2APolicyPlugin interface {
	ADPPolicyPlugin

	// GenerateA2APolicies generates A2A policies for services with a2a protocol
	GenerateA2APolicies(ctx krt.HandlerContext, services krt.Collection[*corev1.Service]) ([]ADPPolicy, error)
}

// InferencePoolPolicyPlugin handles inference pool policies
type InferencePoolPolicyPlugin interface {
	ADPPolicyPlugin

	// GenerateInferencePoolPolicies generates policies for inference pools
	GenerateInferencePoolPolicies(ctx krt.HandlerContext, inferencePools krt.Collection[*inf.InferencePool], domainSuffix string) ([]ADPPolicy, error)
}

// TrafficPolicyPlugin handles traffic policies
type TrafficPolicyPlugin interface {
	ADPPolicyPlugin

	// GenerateTrafficPolicies generates policies for traffic policies
	GenerateTrafficPolicies(ctx krt.HandlerContext, trafficPolicies krt.Collection[*v1alpha1.TrafficPolicy], gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension]) ([]ADPPolicy, error)
}

// PolicyManager coordinates all policy plugins
type PolicyManager interface {
	// RegisterPlugin registers a policy plugin
	RegisterPlugin(plugin PolicyPlugin) error

	// GetPluginsByType returns all plugins of a specific type
	GetPluginsByType(policyType PolicyType) []PolicyPlugin

	// GetADPPlugins returns all ADP policy plugins
	GetADPPlugins() []ADPPolicyPlugin

	// GetAgentGatewayPlugins returns all agent gateway policy plugins
	GetAgentGatewayPlugins() []AgentGatewayPolicyPlugin

	// GenerateAllPolicies generates policies from all registered ADP plugins
	GenerateAllPolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error)
}
