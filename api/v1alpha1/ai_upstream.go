package v1alpha1

import corev1 "k8s.io/api/core/v1"

type AIUpstream struct {
	// Send requests to a custom host and port, such as to proxy the request,
	// or to use a different backend that is API-compliant with the upstream version.
	CustomHost Host `json:"customHost,omitempty"`

	// The LLM configures the AI gateway to use a single LLM provider backend.
	LLM LLMProviders `json:"llm,inline"`
	// The MultiPool configures the backends for multiple hosts or models from the same provider in one Upstream resource.
	MultiPool *MultiPoolConfig `json:"multipool,omitempty"`
}

// LLMProviders configures the AI gateway to use a single LLM provider backend.
// +kubebuilder:validation:XValidation:message="There must one and only one LLMProviders type set",rule="1 == (self.openai != null?1:0) + (self.azureopenai != null?1:0) + (self.anthropic != null?1:0) + (self.gemini != null?1:0) + (self.vertexai != null?1:0) + (self.mistral != null?1:0)"
type LLMProviders struct {
	OpenAI      *OpenAIConfig      `json:"openai,omitempty"`
	AzureOpenAI *AzureOpenAIConfig `json:"azureopenai,omitempty"`
	Anthropic   *AnthropicConfig   `json:"anthropic,omitempty"`
	Gemini      *GeminiConfig      `json:"gemini,omitempty"`
	VertexAI    *VertexAIConfig    `json:"vertexai,omitempty"`
	Mistral     *MistralConfig     `json:"mistral,omitempty"`
}

type SingleAuthTokenKind string

const (
	// Inline provides the token directly in the configuration for the Upstream.
	Inline SingleAuthTokenKind = "INLINE"

	// SecretRef provides the token directly in the configuration for the Upstream.
	SecretRef SingleAuthTokenKind = "SECRETREF"

	// Passthrough the existing token. This token can either
	// come directly from the client, or be generated by an OIDC flow
	// early in the request lifecycle. This option is useful for
	// backends which have federated identity setup and can re-use
	// the token from the client.
	// Currently, this token must exist in the `Authorization` header.
	Passthrough SingleAuthTokenKind = "PASSTHROUGH"
)

// SingleAuthToken configures the authorization token that the AI gateway uses to access the LLM provider API.
// This token is automatically sent in a request header, depending on the LLM provider.
type SingleAuthToken struct {
	// Kind specifies which type of authorization token is being used.
	// Must be one of: "INLINE", "SECRETREF", "PASSTHROUGH".
	// +kubebuilder:validation:Enum=INLINE;SECRETREF;PASSTHROUGH
	Kind SingleAuthTokenKind `json:"kind"`

	// Provide the token directly in the configuration for the Upstream.
	// This option is the least secure. Only use this option for quick tests such as trying out AI Gateway.
	// +kubebuilder:validation:XValidation:message="There must one and only one SingleAuthToken type set",rule="1 == (self.inline != null?1:0) + (self.secretRef != null?1:0)"
	// +kubebuilder:validation:XValidation:message="Inline token must be set when kind is INLINE",rule="self.kind != 'INLINE' || self.inline != null"
	Inline string `json:"inline,omitempty"`

	// Store the API key in a Kubernetes secret in the same namespace as the Upstream.
	// Then, refer to the secret in the Upstream configuration. This option is more secure than an inline token,
	// because the API key is encoded and you can restrict access to secrets through RBAC rules.
	// You might use this option in proofs of concept, controlled development and staging environments,
	// or well-controlled prod environments that use secrets.
	// +kubebuilder:validation:XValidation:message="There must be one and only one SingleAuthToken type set",rule="1 == (self.inline != null ? 1 : 0) + (self.secretRef != null ? 1 : 0)"
	// +kubebuilder:validation:XValidation:message="SecretRef must be set when kind is SECRETREF",rule="self.kind != 'SECRETREF' || self.secretRef != null"
	SecretRef corev1.LocalObjectReference `json:"secretRef,omitempty"`
}

// OpenAIConfig settings for the [OpenAI](https://platform.openai.com/docs/api-reference/streaming) LLM provider.
type OpenAIConfig struct {
	// The authorization token that the AI gateway uses to access the OpenAI API.
	// This token is automatically sent in the `Authorization` header of the
	// request and prefixed with `Bearer`.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`
	// Optional: Send requests to a custom host and port, such as to proxy the request,
	// or to use a different backend that is API-compliant with the upstream version.
	CustomHost Host `json:"customHost,omitempty"`
	// Optional: Override the model name, such as `gpt-4o-mini`.
	// If unset, the model name is taken from the request.
	// This setting can be useful when setting up model failover within the same LLM provider.
	Model string `json:"model,omitempty"`
}

// AzureOpenAIConfig settings for the [Azure OpenAI](https://learn.microsoft.com/en-us/azure/ai-services/openai/) LLM provider.
type AzureOpenAIConfig struct {
	// The authorization token that the AI gateway uses to access the Azure OpenAI API.
	// This token is automatically sent in the `api-key` header of the request.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`

	// The endpoint for the Azure OpenAI API to use, such as `my-endpoint.openai.azure.com`.
	// If the scheme is included, it is stripped.
	Endpoint string `json:"endpoint,omitempty"`

	// The name of the Azure OpenAI model deployment to use.
	// For more information, see the [Azure OpenAI model docs](https://learn.microsoft.com/en-us/azure/ai-services/openai/concepts/models).
	DeploymentName string `json:"deploymentName,omitempty"`

	// The version of the Azure OpenAI API to use.
	// For more information, see the [Azure OpenAI API version reference](https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#api-specs).
	ApiVersion string `json:"apiVersion,omitempty"`
}

// GeminiConfig settings for the [Gemini](https://ai.google.dev/gemini-api/docs) LLM provider.
type GeminiConfig struct {
	// The authorization token that the AI gateway uses to access the Gemini API.
	// This token is automatically sent in the `key` query parameter of the request.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`

	// The Gemini model to use.
	// For more information, see the [Gemini models docs](https://ai.google.dev/gemini-api/docs/models/gemini).
	Model string `json:"model,omitempty"`

	// The version of the Gemini API to use.
	// For more information, see the [Gemini API version docs](https://ai.google.dev/gemini-api/docs/api-versions).
	ApiVersion string `json:"apiVersion,omitempty"`
}

// Publisher configures the type of publisher model to use for VertexAI. Currently, only Google is supported.
type Publisher int

const GOOGLE Publisher = iota

// VertexAIConfig settings for the [Vertex AI](https://cloud.google.com/vertex-ai/docs) LLM provider.
// To find the values for the project ID, project location, and publisher, you can check the fields of an API request, such as
// `https://{LOCATION}-aiplatform.googleapis.com/{VERSION}/projects/{PROJECT_ID}/locations/{LOCATION}/publishers/{PROVIDER}/<model-path>`.
type VertexAIConfig struct {
	// The authorization token that the AI gateway uses to access the Vertex AI API.
	// This token is automatically sent in the `key` header of the request.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`

	// The Vertex AI model to use.
	// For more information, see the [Vertex AI model docs](https://cloud.google.com/vertex-ai/generative-ai/docs/learn/models).
	Model string `json:"model,omitempty"`

	// The version of the Vertex AI API to use.
	// For more information, see the [Vertex AI API reference](https://cloud.google.com/vertex-ai/docs/reference#versions).
	ApiVersion string `json:"apiVersion,omitempty"`

	// The ID of the Google Cloud Project that you use for the Vertex AI.
	ProjectId string `json:"projectId,omitempty"`

	// The location of the Google Cloud Project that you use for the Vertex AI.
	Location string `json:"location,omitempty"`

	// Optional: The model path to route to. Defaults to the Gemini model path, `generateContent`.
	ModelPath string `json:"modelPath,omitempty"`

	// The type of publisher model to use. Currently, only Google is supported.
	Publisher Publisher `json:"publisher,omitempty"`
}

// MistralConfig configures the settings for the [Mistral AI](https://docs.mistral.ai/getting-started/quickstart/) LLM provider.
type MistralConfig struct {
	// The authorization token that the AI gateway uses to access the OpenAI API.
	// This token is automatically sent in the `Authorization` header of the
	// request and prefixed with `Bearer`.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`
	// Optional: Send requests to a custom host and port, such as to proxy the request,
	// or to use a different backend that is API-compliant with the upstream version.
	CustomHost Host `json:"customHost,omitempty"`
	// Optional: Override the model name.
	// If unset, the model name is taken from the request.
	// This setting can be useful when testing model failover scenarios.
	Model string `json:"model,omitempty"`
}

// AnthropicConfig settings for the [Anthropic](https://docs.anthropic.com/en/release-notes/api) LLM provider.
type AnthropicConfig struct {
	// The authorization token that the AI gateway uses to access the Anthropic API.
	// This token is automatically sent in the `x-api-key` header of the request.
	AuthToken SingleAuthToken `json:"authToken,omitempty"`
	// Optional: Send requests to a custom host and port, such as to proxy the request,
	// or to use a different backend that is API-compliant with the upstream version.
	CustomHost Host `json:"customHost,omitempty"`
	// Optional: A version header to pass to the Anthropic API.
	// For more information, see the [Anthropic API versioning docs](https://docs.anthropic.com/en/api/versioning).
	Version string `json:"apiVersion,omitempty"`
	// Optional: Override the model name.
	// If unset, the model name is taken from the request.
	// This setting can be useful when testing model failover scenarios.
	Model string `json:"model,omitempty"`
}

// Priority configures the priority of the backend endpoints.
type Priority struct {
	// A list of LLM provider backends within a single endpoint pool entry.
	Pool []LLMProviders `json:"pool,omitempty"`
}

// MultiPoolConfig configures the backends for multiple hosts or models from the same provider in one Upstream resource.
// This method can be useful for creating one logical endpoint that is backed
// by multiple hosts or models.
//
// In the `priorities` section, the order of `pool` entries defines the priority of the backend endpoints.
// The `pool` entries can either define a list of backends or a single backend.
// Note: Only two levels of nesting are permitted. Any nested entries after the second level are ignored.
//
// ```yaml
// multi:
//
//	priorities:
//	- pool:
//	  - azureOpenai:
//	      deploymentName: gpt-4o-mini
//	      apiVersion: 2024-02-15-preview
//	      endpoint: ai-gateway.openai.azure.com
//	      authToken:
//	        secretRef:
//	          name: azure-secret
//	          namespace: kgateway-system
//	- pool:
//	  - azureOpenai:
//	      deploymentName: gpt-4o-mini-2
//	      apiVersion: 2024-02-15-preview
//	      endpoint: ai-gateway-2.openai.azure.com
//	      authToken:
//	        secretRef:
//	          name: azure-secret-2
//	          namespace: kgateway-system
//
// ```
type MultiPoolConfig struct {
	// The priority list of backend pools. Each entry represents a set of LLM provider backends.
	// The order defines the priority of the backend endpoints.
	Priorities []Priority `json:"priorities,omitempty"`
}
