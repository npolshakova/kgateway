package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=trafficpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=trafficpolicies/status,verbs=get;update;patch

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=kgateway,app.kubernetes.io/name=kgateway}
// +kubebuilder:resource:categories=kgateway
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="gateway.networking.k8s.io/policy=Direct"
type TrafficPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TrafficPolicySpec `json:"spec,omitempty"`

	Status gwv1alpha2.PolicyStatus `json:"status,omitempty"`
	// TODO: embed this into a typed Status field when
	// https://github.com/kubernetes/kubernetes/issues/131533 is resolved
}

// +kubebuilder:object:root=true
type TrafficPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrafficPolicy `json:"items"`
}

// TrafficPolicySpec defines the desired state of a traffic policy.
// +kubebuilder:validation:XValidation:rule="!has(self.hashPolicies) || ((has(self.targetRefs) && self.targetRefs.all(r, r.kind == 'HTTPRoute')) || (has(self.targetSelectors) && self.targetSelectors.all(r, r.kind == 'HTTPRoute')))",message="hash policies can only be used when targeting HTTPRoute resources"
// +kubebuilder:validation:XValidation:rule="!has(self.autoHostRewrite) || ((has(self.targetRefs) && self.targetRefs.all(r, r.kind == 'HTTPRoute')) || (has(self.targetSelectors) && self.targetSelectors.all(r, r.kind == 'HTTPRoute')))",message="autoHostRewrite can only be used when targeting HTTPRoute resources"
type TrafficPolicySpec struct {
	// TargetRefs specifies the target resources by reference to attach the policy to.
	// +optional
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:XValidation:rule="self.all(r, (r.kind == 'Gateway' || r.kind == 'HTTPRoute' || (r.kind == 'XListenerSet' && r.group == 'gateway.networking.x-k8s.io')) && (!has(r.group) || r.group == 'gateway.networking.k8s.io' || r.group == 'gateway.networking.x-k8s.io'))",message="targetRefs may only reference Gateway, HTTPRoute, or XListenerSet resources"
	TargetRefs []LocalPolicyTargetReferenceWithSectionName `json:"targetRefs,omitempty"`

	// TargetSelectors specifies the target selectors to select resources to attach the policy to.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self.all(r, (r.kind == 'Gateway' || r.kind == 'HTTPRoute' || (r.kind == 'XListenerSet' && r.group == 'gateway.networking.x-k8s.io')) && (!has(r.group) || r.group == 'gateway.networking.k8s.io' || r.group == 'gateway.networking.x-k8s.io'))",message="targetSelectors may only reference Gateway, HTTPRoute, or XListenerSet resources"
	TargetSelectors []LocalPolicyTargetSelector `json:"targetSelectors,omitempty"`

	// AI is used to configure AI-based policies for the policy.
	// +optional
	AI *AIPolicy `json:"ai,omitempty"`

	// Transformation is used to mutate and transform requests and responses
	// before forwarding them to the destination.
	// +optional
	Transformation *TransformationPolicy `json:"transformation,omitempty"`

	// ExtProc specifies the external processing configuration for the policy.
	// +optional
	ExtProc *ExtProcPolicy `json:"extProc,omitempty"`

	// ExtAuth specifies the external authentication configuration for the policy.
	// This controls what external server to send requests to for authentication.
	// +optional
	ExtAuth *ExtAuthPolicy `json:"extAuth,omitempty"`

	// RateLimit specifies the rate limiting configuration for the policy.
	// This controls the rate at which requests are allowed to be processed.
	// +optional
	RateLimit *RateLimit `json:"rateLimit,omitempty"`

	// Cors specifies the CORS configuration for the policy.
	// +optional
	Cors *CorsPolicy `json:"cors,omitempty"`

	// Csrf specifies the Cross-Site Request Forgery (CSRF) policy for this traffic policy.
	// +optional
	Csrf *CSRFPolicy `json:"csrf,omitempty"`

	// HashPolicies specifies the hash policies for hashing load balancers (RingHash, Maglev).
	// Should be used in conjunction with Load Balancer on the BackendConfigPolicy.
	// Note: can only be used when targeting routes.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	HashPolicies []*HashPolicy `json:"hashPolicies,omitempty"`

	// AutoHostRewrite rewrites the Host header to the DNS name of the selected upstream.
	// NOTE: This field is only honoured for HTTPRoute targets.
	// NOTE: If `autoHostRewrite` is set on a route that also has a [URLRewrite filter](https://gateway-api.sigs.k8s.io/reference/spec/#httpurlrewritefilter)
	// configured to override the `hostname`, the `hostname` value will be used and `autoHostRewrite` will be ignored.
	// +optional
	AutoHostRewrite *bool `json:"autoHostRewrite,omitempty"`

	// Buffer can be used to set the maximum request size that will be buffered.
	// Requests exceeding this size will return a 413 response.
	// +optional
	Buffer *Buffer `json:"buffer,omitempty"`
}

// TransformationPolicy config is used to modify envoy behavior at a route level.
// These modifications can be performed on the request and response paths.
type TransformationPolicy struct {
	// Request is used to modify the request path.
	// +optional
	Request *Transform `json:"request,omitempty"`

	// Response is used to modify the response path.
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
	// If empty, body will not be buffered.
	// +optional
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
type (
	HeaderName           string
	HeaderTransformation struct {
		// Name is the name of the header to interact with.
		// +required
		Name HeaderName `json:"name,omitempty"`
		// Value is the template to apply to generate the output value for the header.
		Value InjaTemplate `json:"value,omitempty"`
	}
)

// BodyparseBehavior defines how the body should be parsed
// If set to json and the body is not json then the filter will not perform the transformation.
// +kubebuilder:validation:Enum=AsString;AsJson
type BodyParseBehavior string

const (
	// BodyParseBehaviorAsString will parse the body as a string.
	BodyParseBehaviorAsString BodyParseBehavior = "AsString"
	// BodyParseBehaviorAsJSON will parse the body as a json object.
	BodyParseBehaviorAsJSON BodyParseBehavior = "AsJson"
)

// BodyTransformation controls how the body should be parsed and transformed.
type BodyTransformation struct {
	// ParseAs defines what auto formatting should be applied to the body.
	// This can make interacting with keys within a json body much easier if AsJson is selected.
	// +kubebuilder:default=AsString
	ParseAs BodyParseBehavior `json:"parseAs"`

	// Value is the template to apply to generate the output value for the body.
	// +optional
	Value *InjaTemplate `json:"value,omitempty"`
}

// ExtAuthEnabled determines the enabled state of the ExtAuth filter.
// +kubebuilder:validation:Enum=DisableAll
type ExtAuthEnabled string

// When we add a new field here we have to be specific around which extensions are enabled/disabled
// and how these can be overridden by other policies.
const (
	// ExtAuthDisableAll disables all instances of the ExtAuth filter for this route.
	// This is to enable a global disable such as for a health check route.
	ExtAuthDisableAll ExtAuthEnabled = "DisableAll"
)

// ExtAuthPolicy configures external authentication for a route.
// This policy will determine the ext auth server to use and how to  talk to it.
// Note that most of these fields are passed along as is to Envoy.
// For more details on particular fields please see the Envoy ExtAuth documentation.
// https://raw.githubusercontent.com/envoyproxy/envoy/f910f4abea24904aff04ec33a00147184ea7cffa/api/envoy/extensions/filters/http/ext_authz/v3/ext_authz.proto
//
// +kubebuilder:validation:ExactlyOneOf=extensionRef;enablement
type ExtAuthPolicy struct {
	// ExtensionRef references the ExternalExtension that should be used for authentication.
	// +optional
	ExtensionRef *corev1.LocalObjectReference `json:"extensionRef,omitempty"`

	// Enablement determines the enabled state of the ExtAuth filter.
	// When set to "DisableAll", the filter is disabled for this route.
	// When empty, the filter is enabled as long as it is not disabled by another policy.
	// +optional
	Enablement *ExtAuthEnabled `json:"enablement,omitempty"`

	// WithRequestBody allows the request body to be buffered and sent to the authorization service.
	// Warning buffering has implications for streaming and therefore performance.
	// +optional
	WithRequestBody *BufferSettings `json:"withRequestBody,omitempty"`

	// Additional context for the authorization service.
	// +optional
	ContextExtensions map[string]string `json:"contextExtensions,omitempty"`
}

// BufferSettings configures how the request body should be buffered.
type BufferSettings struct {
	// MaxRequestBytes sets the maximum size of a message body to buffer.
	// Requests exceeding this size will receive HTTP 413 and not be sent to the authorization service.
	// +required
	// +kubebuilder:validation:Minimum=1
	MaxRequestBytes uint32 `json:"maxRequestBytes"`

	// AllowPartialMessage determines if partial messages should be allowed.
	// When true, requests will be sent to the authorization service even if they exceed maxRequestBytes.
	// When unset, the default behavior is false.
	// +optional
	AllowPartialMessage *bool `json:"allowPartialMessage,omitempty"`

	// PackAsBytes determines if the body should be sent as raw bytes.
	// When true, the body is sent as raw bytes in the raw_body field.
	// When false, the body is sent as UTF-8 string in the body field.
	// When unset, the default behavior is false.
	// +optional
	PackAsBytes *bool `json:"packAsBytes,omitempty"`
}

// RateLimit defines a rate limiting policy.
type RateLimit struct {
	// Local defines a local rate limiting policy.
	// +optional
	Local *LocalRateLimitPolicy `json:"local,omitempty"`

	// Global defines a global rate limiting policy using an external service.
	// +optional
	Global *RateLimitPolicy `json:"global,omitempty"`
}

// LocalRateLimitPolicy represents a policy for local rate limiting.
// It defines the configuration for rate limiting using a token bucket mechanism.
type LocalRateLimitPolicy struct {
	// TokenBucket represents the configuration for a token bucket local rate-limiting mechanism.
	// It defines the parameters for controlling the rate at which requests are allowed.
	// +optional
	TokenBucket *TokenBucket `json:"tokenBucket,omitempty"`
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
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	TokensPerFill *uint32 `json:"tokensPerFill,omitempty"`

	// FillInterval defines the time duration between consecutive token fills.
	// This value must be a valid duration string (e.g., "1s", "500ms").
	// It determines the frequency of token replenishment.
	// +required
	// +kubebuilder:validation:XValidation:rule="matches(self, '^([0-9]{1,5}(h|m|s|ms)){1,4}$')",message="invalid duration value"
	FillInterval metav1.Duration `json:"fillInterval"`
}

// RateLimitPolicy defines a global rate limiting policy using an external service.
type RateLimitPolicy struct {
	// Descriptors define the dimensions for rate limiting.
	// These values are passed to the rate limit service which applies configured limits based on them.
	// Each descriptor represents a single rate limit rule with one or more entries.
	// +required
	// +kubebuilder:validation:MinItems=1
	Descriptors []RateLimitDescriptor `json:"descriptors"`

	// ExtensionRef references a GatewayExtension that provides the global rate limit service.
	// +required
	ExtensionRef *corev1.LocalObjectReference `json:"extensionRef"`
}

// RateLimitDescriptor defines a descriptor for rate limiting.
// A descriptor is a group of entries that form a single rate limit rule.
type RateLimitDescriptor struct {
	// Entries are the individual components that make up this descriptor.
	// When translated to Envoy, these entries combine to form a single descriptor.
	// +required
	// +kubebuilder:validation:MinItems=1
	Entries []RateLimitDescriptorEntry `json:"entries"`
}

// RateLimitDescriptorEntryType defines the type of a rate limit descriptor entry.
// +kubebuilder:validation:Enum=Generic;Header;RemoteAddress;Path
type RateLimitDescriptorEntryType string

const (
	// RateLimitDescriptorEntryTypeGeneric represents a generic key-value descriptor entry.
	RateLimitDescriptorEntryTypeGeneric RateLimitDescriptorEntryType = "Generic"

	// RateLimitDescriptorEntryTypeHeader represents a descriptor entry that extracts its value from a request header.
	RateLimitDescriptorEntryTypeHeader RateLimitDescriptorEntryType = "Header"

	// RateLimitDescriptorEntryTypeRemoteAddress represents a descriptor entry that uses the client's IP address as its value.
	RateLimitDescriptorEntryTypeRemoteAddress RateLimitDescriptorEntryType = "RemoteAddress"

	// RateLimitDescriptorEntryTypePath represents a descriptor entry that uses the request path as its value.
	RateLimitDescriptorEntryTypePath RateLimitDescriptorEntryType = "Path"
)

// RateLimitDescriptorEntry defines a single entry in a rate limit descriptor.
// Only one entry type may be specified.
// +kubebuilder:validation:XValidation:message="exactly one entry type must be specified",rule="(has(self.type) && (self.type == 'Generic' && has(self.generic) && !has(self.header)) || (self.type == 'Header' && has(self.header) && !has(self.generic)) || (self.type == 'RemoteAddress' && !has(self.generic) && !has(self.header)) || (self.type == 'Path' && !has(self.generic) && !has(self.header)))"
type RateLimitDescriptorEntry struct {
	// Type specifies what kind of rate limit descriptor entry this is.
	// +required
	Type RateLimitDescriptorEntryType `json:"type"`

	// Generic contains the configuration for a generic key-value descriptor entry.
	// This field must be specified when Type is Generic.
	// +optional
	Generic *RateLimitDescriptorEntryGeneric `json:"generic,omitempty"`

	// Header specifies a request header to extract the descriptor value from.
	// This field must be specified when Type is Header.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Header *string `json:"header,omitempty"`
}

// RateLimitDescriptorEntryGeneric defines a generic key-value descriptor entry.
type RateLimitDescriptorEntryGeneric struct {
	// Key is the name of this descriptor entry.
	// +required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`

	// Value is the static value for this descriptor entry.
	// +required
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value"`
}

type CorsPolicy struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	*gwv1.HTTPCORSFilter `json:",inline"`
}

// CSRFPolicy can be used to set percent of requests for which the CSRF filter is enabled,
// enable shadow-only mode where policies will be evaluated and tracked, but not enforced and
// add additional source origins that will be allowed in addition to the destination origin.
//
// +kubebuilder:validation:AtMostOneOf=percentageEnabled;percentageShadowed
type CSRFPolicy struct {
	// Specifies the percentage of requests for which the CSRF filter is enabled.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	PercentageEnabled *uint32 `json:"percentageEnabled,omitempty"`

	// Specifies that CSRF policies will be evaluated and tracked, but not enforced.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	PercentageShadowed *uint32 `json:"percentageShadowed,omitempty"`

	// Specifies additional source origins that will be allowed in addition to the destination origin.
	// +optional
	// +kubebuilder:validation:MaxItems=16
	AdditionalOrigins []StringMatcher `json:"additionalOrigins,omitempty"`
}

// +kubebuilder:validation:ExactlyOneOf=header;cookie;sourceIP
type HashPolicy struct {
	// Header specifies a header's value as a component of the hash key.
	// +optional
	Header *Header `json:"header,omitempty"`

	// Cookie specifies a given cookie as a component of the hash key.
	// +optional
	Cookie *Cookie `json:"cookie,omitempty"`

	// SourceIP specifies whether to use the request's source IP address as a component of the hash key.
	// +optional
	SourceIP *SourceIP `json:"sourceIP,omitempty"`

	// Terminal, if set, and a hash key is available after evaluating this policy, will cause Envoy to skip the subsequent policies and
	// use the key as it is.
	// This is useful for defining "fallback" policies and limiting the time Envoy spends generating hash keys.
	// +optional
	Terminal *bool `json:"terminal,omitempty"`
}

type Header struct {
	// Name is the name of the header to use as a component of the hash key.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

type Cookie struct {
	// Name of the cookie.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Path is the name of the path for the cookie.
	// +optional
	Path *string `json:"path,omitempty"`

	// TTL specifies the time to live of the cookie.
	// If specified, a cookie with the TTL will be generated if the cookie is not present.
	// If the TTL is present and zero, the generated cookie will be a session cookie.
	// +optional
	// +kubebuilder:validation:XValidation:rule="matches(self, '^([0-9]{1,5}(h|m|s|ms)){1,4}$')",message="invalid duration value"
	TTL *metav1.Duration `json:"ttl,omitempty"`

	// Attributes are additional attributes for the cookie.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	Attributes map[string]string `json:"attributes,omitempty"`
}

type SourceIP struct{}

type Buffer struct {
	// MaxRequestSize sets the maximum size in bytes of a message body to buffer.
	// Requests exceeding this size will receive HTTP 413.
	// Example format: "1Mi", "512Ki", "1Gi"
	// +required
	// +kubebuilder:validation:XValidation:message="maxRequestSize must be greater than 0 and less than 4Gi",rule="quantity(self).isGreaterThan(quantity('0')) && quantity(self).isLessThan(quantity('4Gi'))"
	MaxRequestSize *resource.Quantity `json:"maxRequestSize"`
}
