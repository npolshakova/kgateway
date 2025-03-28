package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=routepolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=routepolicies/status,verbs=get;update;patch

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=kgateway,app.kubernetes.io/name=kgateway}
// +kubebuilder:resource:categories=kgateway
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="gateway.networking.k8s.io/policy=Direct"
type RoutePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoutePolicySpec `json:"spec,omitempty"`
	Status SimpleStatus    `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type RoutePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoutePolicy `json:"items"`
}

type RoutePolicySpec struct {
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	TargetRefs []LocalPolicyTargetReference `json:"targetRefs,omitempty"`

	AI *AIRoutePolicy `json:"ai,omitempty"`

	Transformation TransformationPolicy `json:"transformation,omitempty"`

	// RateLimit specifies the rate limiting configuration for the policy.
	// This controls the rate at which requests are allowed to be processed.
	// +optional
	RateLimit *RateLimit `json:"rateLimit,omitempty"`
}

// TransformationPolicy config is used to modify envoy behavior at a route level.
// These modifications can be performed on the request and response paths.
type TransformationPolicy struct {
	// +optional
	Request *Transform `json:"request,omitempty"`
	// +optional
	Response *Transform `json:"response,omitempty"`
}

// Transform defines the operations to be performed by the transformation.
// These operations may include changing the actual request/response but may also cause side effects.
// Side effects may include setting info that can be used in future steps (e.g. dynamic metadata) and can cause envoy to buffer.
type Transform struct {
	// Set is a list of headers and the value they should be set to.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=16
	Set []HeaderTransformation `json:"set,omitempty"`

	// Add is a list of headers to add to the request and what that value should be set to.
	// If there is already a header with these values then append the value as an extra entry.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=16
	Add []HeaderTransformation `json:"add,omitempty"`

	// Remove is a list of header names to remove from the request/response.
	// +optional
	// +listType=set
	// +kubebuilder:validation:MaxItems=16
	Remove []string `json:"remove,omitempty"`

	// Body controls both how to parse the body and if needed how to set.
	// +optional
	//
	// If empty, body will not be buffered.
	Body *BodyTransformation `json:"body,omitempty"`
}

type InjaTemplate string

// EnvoyHeaderName is the name of a header or pseudo header
// Based on gateway api v1.Headername but allows a singular : at the start
//
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=256
// +kubebuilder:validation:Pattern=`^:?[A-Za-z0-9!#$%&'*+\-.^_\x60|~]+$`
// +k8s:deepcopy-gen=false
type HeaderName string
type HeaderTransformation struct {
	// Name is the name of the header to interact with.
	// +required
	Name HeaderName `json:"name,omitempty"`
	// Value is the template to apply to generate the output value for the header.
	Value InjaTemplate `json:"value,omitempty"`
}

// BodyparseBehavior defines how the body should be parsed
// If set to json and the body is not json then the filter will not perform the transformation.
// +kubebuilder:validation:Enum=AsString;AsJson
type BodyParseBehavior string

const (
	BodyParseBehaviorAsString BodyParseBehavior = "AsString"
	BodyParseBehaviorAsJSON   BodyParseBehavior = "AsJson"
)

// BodyTransformation controls how the body should be parsed and transformed.
type BodyTransformation struct {
	// ParseAs defines what auto formatting should be applied to the body.
	// This can make interacting with keys within a json body much easier if AsJson is selected.
	// +kubebuilder:default=AsString
	ParseAs BodyParseBehavior `json:"parseAs"`
	// Value is the template to apply to generate the output value for the body.
	Value *InjaTemplate `json:"value,omitempty"`
}

// SimpleStatus defines the observed state of the policy.
type SimpleStatus struct {
	// Conditions is the list of conditions for the policy.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// RateLimit defines a rate limiting policy.
type RateLimit struct {
	// Local defines a local rate limiting policy.
	// +required
	Local *LocalRateLimitPolicy `json:"local,omitempty"`
}

// LocalRateLimitPolicy represents a policy for local rate limiting.
// It defines the configuration for rate limiting using a token bucket mechanism.
type LocalRateLimitPolicy struct {
	// TokenBucket represents the configuration for a token bucket local rate-limiting mechanism.
	// It defines the parameters for controlling the rate at which requests are allowed.
	// +required
	TokenBucket *TokenBucket `json:"tokenBucket"`
}

// TokenBucket defines the configuration for a token bucket rate-limiting mechanism.
// It controls the rate at which tokens are generated and consumed for a specific operation.
type TokenBucket struct {
	// MaxTokens specifies the maximum number of tokens that the bucket can hold.
	// This value must be greater than or equal to 1.
	// It determines the burst capacity of the rate limiter.
	// +required
	// +kubebuilder:validation:Minimum=1
	MaxTokens uint32 `json:"maxTokens"`

	// TokensPerFill specifies the number of tokens added to the bucket during each fill interval.
	// If not specified, it defaults to 1.
	// This controls the steady-state rate of token generation.
	// +optional
	// kubebuilder:default:=1
	TokensPerFill *uint32 `json:"tokensPerFill,omitempty"`

	// FillInterval defines the time duration between consecutive token fills.
	// This value must be a valid duration string (e.g., "1s", "500ms").
	// It determines the frequency of token replenishment.
	// +required
	// +kubebuilder:validation:Format=duration
	FillInterval string `json:"fillInterval"`
}
