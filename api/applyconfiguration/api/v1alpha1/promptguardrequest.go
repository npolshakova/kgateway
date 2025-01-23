// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// PromptguardRequestApplyConfiguration represents a declarative configuration of the PromptguardRequest type for use
// with apply.
type PromptguardRequestApplyConfiguration struct {
	CustomResponse *CustomResponseApplyConfiguration `json:"customResponse,omitempty"`
	Regex          *RegexApplyConfiguration          `json:"regex,omitempty"`
	Webhook        *WebhookApplyConfiguration        `json:"webhook,omitempty"`
	Moderation     *ModerationApplyConfiguration     `json:"moderation,omitempty"`
}

// PromptguardRequestApplyConfiguration constructs a declarative configuration of the PromptguardRequest type for use with
// apply.
func PromptguardRequest() *PromptguardRequestApplyConfiguration {
	return &PromptguardRequestApplyConfiguration{}
}

// WithCustomResponse sets the CustomResponse field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CustomResponse field is set to the value of the last call.
func (b *PromptguardRequestApplyConfiguration) WithCustomResponse(value *CustomResponseApplyConfiguration) *PromptguardRequestApplyConfiguration {
	b.CustomResponse = value
	return b
}

// WithRegex sets the Regex field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Regex field is set to the value of the last call.
func (b *PromptguardRequestApplyConfiguration) WithRegex(value *RegexApplyConfiguration) *PromptguardRequestApplyConfiguration {
	b.Regex = value
	return b
}

// WithWebhook sets the Webhook field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Webhook field is set to the value of the last call.
func (b *PromptguardRequestApplyConfiguration) WithWebhook(value *WebhookApplyConfiguration) *PromptguardRequestApplyConfiguration {
	b.Webhook = value
	return b
}

// WithModeration sets the Moderation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Moderation field is set to the value of the last call.
func (b *PromptguardRequestApplyConfiguration) WithModeration(value *ModerationApplyConfiguration) *PromptguardRequestApplyConfiguration {
	b.Moderation = value
	return b
}
