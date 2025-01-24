// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// AIUpstreamApplyConfiguration represents a declarative configuration of the AIUpstream type for use
// with apply.
type AIUpstreamApplyConfiguration struct {
	CustomHost *HostApplyConfiguration       `json:"customHost,omitempty"`
	LLM        *LLMBackendApplyConfiguration `json:"llm,omitempty"`
}

// AIUpstreamApplyConfiguration constructs a declarative configuration of the AIUpstream type for use with
// apply.
func AIUpstream() *AIUpstreamApplyConfiguration {
	return &AIUpstreamApplyConfiguration{}
}

// WithCustomHost sets the CustomHost field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CustomHost field is set to the value of the last call.
func (b *AIUpstreamApplyConfiguration) WithCustomHost(value *HostApplyConfiguration) *AIUpstreamApplyConfiguration {
	b.CustomHost = value
	return b
}

// WithLLM sets the LLM field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LLM field is set to the value of the last call.
func (b *AIUpstreamApplyConfiguration) WithLLM(value *LLMBackendApplyConfiguration) *AIUpstreamApplyConfiguration {
	b.LLM = value
	return b
}
