// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// McpTargetApplyConfiguration represents a declarative configuration of the McpTarget type for use
// with apply.
type McpTargetApplyConfiguration struct {
	Name     *string                  `json:"name,omitempty"`
	Host     *string                  `json:"host,omitempty"`
	Port     *int32                   `json:"port,omitempty"`
	Protocol *apiv1alpha1.MCPProtocol `json:"protocol,omitempty"`
}

// McpTargetApplyConfiguration constructs a declarative configuration of the McpTarget type for use with
// apply.
func McpTarget() *McpTargetApplyConfiguration {
	return &McpTargetApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *McpTargetApplyConfiguration) WithName(value string) *McpTargetApplyConfiguration {
	b.Name = &value
	return b
}

// WithHost sets the Host field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Host field is set to the value of the last call.
func (b *McpTargetApplyConfiguration) WithHost(value string) *McpTargetApplyConfiguration {
	b.Host = &value
	return b
}

// WithPort sets the Port field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Port field is set to the value of the last call.
func (b *McpTargetApplyConfiguration) WithPort(value int32) *McpTargetApplyConfiguration {
	b.Port = &value
	return b
}

// WithProtocol sets the Protocol field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Protocol field is set to the value of the last call.
func (b *McpTargetApplyConfiguration) WithProtocol(value apiv1alpha1.MCPProtocol) *McpTargetApplyConfiguration {
	b.Protocol = &value
	return b
}
