// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// GeminiConfigApplyConfiguration represents a declarative configuration of the GeminiConfig type for use
// with apply.
type GeminiConfigApplyConfiguration struct {
	AuthToken  *SingleAuthTokenApplyConfiguration `json:"authToken,omitempty"`
	Model      *string                            `json:"model,omitempty"`
	ApiVersion *string                            `json:"apiVersion,omitempty"`
}

// GeminiConfigApplyConfiguration constructs a declarative configuration of the GeminiConfig type for use with
// apply.
func GeminiConfig() *GeminiConfigApplyConfiguration {
	return &GeminiConfigApplyConfiguration{}
}

// WithAuthToken sets the AuthToken field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AuthToken field is set to the value of the last call.
func (b *GeminiConfigApplyConfiguration) WithAuthToken(value *SingleAuthTokenApplyConfiguration) *GeminiConfigApplyConfiguration {
	b.AuthToken = value
	return b
}

// WithModel sets the Model field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Model field is set to the value of the last call.
func (b *GeminiConfigApplyConfiguration) WithModel(value string) *GeminiConfigApplyConfiguration {
	b.Model = &value
	return b
}

// WithApiVersion sets the ApiVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ApiVersion field is set to the value of the last call.
func (b *GeminiConfigApplyConfiguration) WithApiVersion(value string) *GeminiConfigApplyConfiguration {
	b.ApiVersion = &value
	return b
}
