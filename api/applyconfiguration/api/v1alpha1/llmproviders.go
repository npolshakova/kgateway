// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// LLMProvidersApplyConfiguration represents a declarative configuration of the LLMProviders type for use
// with apply.
type LLMProvidersApplyConfiguration struct {
	OpenAI      *OpenAIConfigApplyConfiguration      `json:"openai,omitempty"`
	AzureOpenAI *AzureOpenAIConfigApplyConfiguration `json:"azureopenai,omitempty"`
	Anthropic   *AnthropicConfigApplyConfiguration   `json:"anthropic,omitempty"`
	Gemini      *GeminiConfigApplyConfiguration      `json:"gemini,omitempty"`
	VertexAI    *VertexAIConfigApplyConfiguration    `json:"vertexai,omitempty"`
	Mistral     *MistralConfigApplyConfiguration     `json:"mistral,omitempty"`
}

// LLMProvidersApplyConfiguration constructs a declarative configuration of the LLMProviders type for use with
// apply.
func LLMProviders() *LLMProvidersApplyConfiguration {
	return &LLMProvidersApplyConfiguration{}
}

// WithOpenAI sets the OpenAI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OpenAI field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithOpenAI(value *OpenAIConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.OpenAI = value
	return b
}

// WithAzureOpenAI sets the AzureOpenAI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AzureOpenAI field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithAzureOpenAI(value *AzureOpenAIConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.AzureOpenAI = value
	return b
}

// WithAnthropic sets the Anthropic field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Anthropic field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithAnthropic(value *AnthropicConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.Anthropic = value
	return b
}

// WithGemini sets the Gemini field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Gemini field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithGemini(value *GeminiConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.Gemini = value
	return b
}

// WithVertexAI sets the VertexAI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VertexAI field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithVertexAI(value *VertexAIConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.VertexAI = value
	return b
}

// WithMistral sets the Mistral field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Mistral field is set to the value of the last call.
func (b *LLMProvidersApplyConfiguration) WithMistral(value *MistralConfigApplyConfiguration) *LLMProvidersApplyConfiguration {
	b.Mistral = value
	return b
}
