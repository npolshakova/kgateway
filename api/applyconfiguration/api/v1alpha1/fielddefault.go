// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// FieldDefaultApplyConfiguration represents a declarative configuration of the FieldDefault type for use
// with apply.
type FieldDefaultApplyConfiguration struct {
	Field    *string               `json:"field,omitempty"`
	Value    *runtime.RawExtension `json:"value,omitempty"`
	Override *bool                 `json:"override,omitempty"`
}

// FieldDefaultApplyConfiguration constructs a declarative configuration of the FieldDefault type for use with
// apply.
func FieldDefault() *FieldDefaultApplyConfiguration {
	return &FieldDefaultApplyConfiguration{}
}

// WithField sets the Field field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Field field is set to the value of the last call.
func (b *FieldDefaultApplyConfiguration) WithField(value string) *FieldDefaultApplyConfiguration {
	b.Field = &value
	return b
}

// WithValue sets the Value field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Value field is set to the value of the last call.
func (b *FieldDefaultApplyConfiguration) WithValue(value runtime.RawExtension) *FieldDefaultApplyConfiguration {
	b.Value = &value
	return b
}

// WithOverride sets the Override field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Override field is set to the value of the last call.
func (b *FieldDefaultApplyConfiguration) WithOverride(value bool) *FieldDefaultApplyConfiguration {
	b.Override = &value
	return b
}
