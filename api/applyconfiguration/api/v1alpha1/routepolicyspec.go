// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// RoutePolicySpecApplyConfiguration represents a declarative configuration of the RoutePolicySpec type for use
// with apply.
type RoutePolicySpecApplyConfiguration struct {
	TargetRef *LocalPolicyTargetReferenceApplyConfiguration `json:"targetRef,omitempty"`
	Timeout   *int                                          `json:"timeout,omitempty"`
	AI        *AIApplyConfiguration                         `json:"ai,omitempty"`
}

// RoutePolicySpecApplyConfiguration constructs a declarative configuration of the RoutePolicySpec type for use with
// apply.
func RoutePolicySpec() *RoutePolicySpecApplyConfiguration {
	return &RoutePolicySpecApplyConfiguration{}
}

// WithTargetRef sets the TargetRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TargetRef field is set to the value of the last call.
func (b *RoutePolicySpecApplyConfiguration) WithTargetRef(value *LocalPolicyTargetReferenceApplyConfiguration) *RoutePolicySpecApplyConfiguration {
	b.TargetRef = value
	return b
}

// WithTimeout sets the Timeout field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Timeout field is set to the value of the last call.
func (b *RoutePolicySpecApplyConfiguration) WithTimeout(value int) *RoutePolicySpecApplyConfiguration {
	b.Timeout = &value
	return b
}

// WithAI sets the AI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AI field is set to the value of the last call.
func (b *RoutePolicySpecApplyConfiguration) WithAI(value *AIApplyConfiguration) *RoutePolicySpecApplyConfiguration {
	b.AI = value
	return b
}
