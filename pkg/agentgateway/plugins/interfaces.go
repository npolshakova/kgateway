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

// PolicyInputsInterface defines a generic interface for accessing plugin inputs
type PolicyInputsInterface[T any] interface {
	// GetServices returns the services collection
	GetServices() krt.Collection[*corev1.Service]

	// GetInferencePools returns the inference pools collection
	GetInferencePools() krt.Collection[*inf.InferencePool]

	// GetTrafficPolicies returns the traffic policies collection
	GetTrafficPolicies() krt.Collection[*v1alpha1.TrafficPolicy]

	// GetGatewayExtensions returns the gateway extensions collection
	GetGatewayExtensions() krt.Collection[*v1alpha1.GatewayExtension]

	// GetDomainSuffix returns the domain suffix
	GetDomainSuffix() string

	// GetExtensions returns plugin-specific extensions of type T
	GetExtensions() T
}

// PolicyPlugin defines the base interface for all policy plugins
type PolicyPlugin[T any] interface {
	// PolicyType returns the type of policy this plugin handles
	PolicyType() PolicyType

	// Name returns the name of the plugin
	Name() string

	// GeneratePolicies generates ADP policies for the given inputs
	GeneratePolicies(ctx krt.HandlerContext, inputs PolicyInputsInterface[T]) ([]ADPPolicy, error)
}

// PolicyInputs contains all the input collections needed by policy plugins
type PolicyInputs struct {
	Services          krt.Collection[*corev1.Service]
	InferencePools    krt.Collection[*inf.InferencePool]
	TrafficPolicies   krt.Collection[*v1alpha1.TrafficPolicy]
	GatewayExtensions krt.Collection[*v1alpha1.GatewayExtension]
	DomainSuffix      string
}

// GetServices implements PolicyInputsInterface
func (p *PolicyInputs) GetServices() krt.Collection[*corev1.Service] {
	return p.Services
}

// GetInferencePools implements PolicyInputsInterface
func (p *PolicyInputs) GetInferencePools() krt.Collection[*inf.InferencePool] {
	return p.InferencePools
}

// GetTrafficPolicies implements PolicyInputsInterface
func (p *PolicyInputs) GetTrafficPolicies() krt.Collection[*v1alpha1.TrafficPolicy] {
	return p.TrafficPolicies
}

// GetGatewayExtensions implements PolicyInputsInterface
func (p *PolicyInputs) GetGatewayExtensions() krt.Collection[*v1alpha1.GatewayExtension] {
	return p.GatewayExtensions
}

// GetDomainSuffix implements PolicyInputsInterface
func (p *PolicyInputs) GetDomainSuffix() string {
	return p.DomainSuffix
}

// GetExtensions implements PolicyInputsInterface - returns nil for the base PolicyInputs
func (p *PolicyInputs) GetExtensions() interface{} {
	return nil
}

// ADPPolicy wraps an ADP policy for collection handling
type ADPPolicy struct {
	Policy *api.Policy
}

// PolicyManager coordinates all policy plugins
type PolicyManager interface {
	// RegisterPlugin registers a policy plugin
	RegisterPlugin(plugin PolicyPlugin[any]) error

	// GetPluginsByType returns all plugins of a specific type
	GetPluginsByType(policyType PolicyType) []PolicyPlugin[any]

	// GenerateAllPolicies generates policies from all registered ADP plugins
	GenerateAllPolicies(ctx krt.HandlerContext, inputs PolicyInputsInterface[any]) ([]ADPPolicy, error)
}
