package v1alpha1

import "github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"

// The type of route to the LLM provider API.
type RouteType int

const (
	// The LLM generates the full response before responding to a client.
	CHAT RouteType = iota
	// Stream responses to a client, which allows the LLM to stream out tokens as they are generated.
	CHAT_STREAMING
)

// An entry for a message to prepend or append to each prompt.
type Message struct {
	// Role of the message. The available roles depend on the backend
	// LLM provider model, such as `SYSTEM` or `USER` in the OpenAI API.
	role string `json:"role,omitempty"`
	// String content of the message.
	content string `json:"content,omitempty"`
}

// Enrich requests sent to the LLM provider by appending and prepending system prompts.
// This can be configured only for LLM providers that use the CHAT API type.
//
// Prompt enrichment allows you to add additional context to the prompt before sending it to the model.
// Unlike RAG or other dynamic context methods, prompt enrichment is static and is applied to every request.
//
// **Note**: Some providers, including Anthropic, do not support SYSTEM role messages, and instead have a dedicated
// system field in the input JSON. In this case, use the [`defaults` setting](#fielddefault) to set the system field.
//
// The following example prepends a system prompt of `Answer all questions in French.`
// and appends `Describe the painting as if you were a famous art critic from the 17th century.`
// to each request that is sent to the `openai` HTTPRoute.
// ```yaml
// apiVersion: gateway.solo.io/v1
// kind: RouteOption
// metadata:
//
//	name: openai-opt
//	namespace: gloo-system
//
// spec:
//
//	targetRefs:
//	- group: gateway.networking.k8s.io
//	  kind: HTTPRoute
//	  name: openai
//	options:
//	  ai:
//	    promptEnrichment:
//	      prepend:
//	      - role: SYSTEM
//	        content: "Answer all questions in French."
//	      append:
//	      - role: USER
//	        content: "Describe the painting as if you were a famous art critic from the 17th century."
//
// ```
type AIPromptEnrichment struct {
	// A list of messages to be prepended to the prompt sent by the client.
	prepend []Message `json:"prepend,omitempty"`
	// A list of messages to be appended to the prompt sent by the client.
	append []Message `json:"append,omitempty"`
}

// Built-in regex patterns for specific types of strings in prompts.
// For example, if you specify `CREDIT_CARD`, any credit card numbers
// in the request or response are matched.
type BuiltIn int

const (
	// Default regex matching for Social Security numbers.
	SSN BuiltIn = iota
	// Default regex matching for credit card numbers.
	CREDIT_CARD
	// Default regex matching for phone numbers.
	PHONE_NUMBER
	// Default regex matching for email addresses.
	EMAIL
)

// Regular expression (regex) matching for prompt guards and data masking.
type RegexMatch struct {
	// The regex pattern to match against the request or response.
	Pattern string `json:"pattern,omitempty"`
	// An optional name for this match, which can be used for debugging purposes.
	Name string `json:"name,omitempty"`
}

// The action to take if a regex pattern is matched in a request or response.
// This setting applies only to request matches. Response matches are always masked by default.
type Action int

const (
	// Mask the matched data in the request.
	MASK Action = iota
	// Reject the request if the regex matches content in the request.
	REJECT
)

// Regular expression (regex) matching for prompt guards and data masking.
type Regex struct {
	// A list of regex patterns to match against the request or response.
	// Matches and built-ins are additive.
	Matches []RegexMatch `json:"regexMatch,omitempty"`
	// A list of built-in regex patterns to match against the request or response.
	// Matches and built-ins are additive.
	Builtins []BuiltIn `json:"builtins,omitempty"`
	// The action to take if a regex pattern is matched in a request or response.
	// This setting applies only to request matches. Response matches are always masked by default.
	Action Action `json:"action,omitempty"`
}

// The header string match type.
type MatchType int64

const (
	// The string must match exactly the specified string.
	EXACT MatchType = iota
	// The string must have the specified prefix.
	PREFIX
	// The string must have the specified suffix.
	SUFFIX
	// The header string must contain the specified string.
	CONTAINS
	// The string must match the specified [RE2-style regular expression](https://github.com/google/re2/wiki/) pattern.
	REGEX
)

// Describes how to match a given string in HTTP headers. Match is case-sensitive.
type HeaderMatch struct {
	// The header key string to match against.
	Key string `json:"key,omitempty"`
	// The type of match to use.
	MatchType MatchType `json:"matchType,omitempty"`
}

// Configure a webhook to forward requests or responses to for prompt guarding.
type Webhook struct {
	// Host to send the traffic to.
	Host string `json:"host,omitempty"`

	// Port to send the traffic to
	port uint32 `json:"port,omitempty"`

	// Headers to forward with the request to the webhook.
	forwardHeaders []HeaderMatch `json:"forwardHeaders,omitempty"`
}

// A custom response to return to the client if request content
// is matched against a regex pattern and the action is `REJECT`.
type CustomResponse struct {
	// A custom response message to return to the client. If not specified, defaults to
	// "The request was rejected due to inappropriate content".
	message string `json:"message,omitempty"`

	// The status code to return to the client.
	StatusCode uint32 `json:"statusCode,omitempty"`
}

// Configure an OpenAI moderation endpoint.
type OpenAIModeration struct {
	// The name of the OpenAI moderation model to use. Defaults to
	// [`omni-moderation-latest`](https://platform.openai.com/docs/guides/moderation).
	Model string `json:"model,omitempty"`

	// The authorization token that the AI gateway uses
	// to access the OpenAI moderation model.
	// TODO: share with upstream in separate ai file
	AuthToken v1alpha1.SingleAuthToken `json:"authToken,omitempty"`
}

// Pass prompt data through an external moderation model endpoint,
// which compares the request prompt input to predefined content rules.
// Any requests that are routed through Gloo AI Gateway pass through the
// moderation model that you specify. If the content is identified as harmful
// according to the model's content rules, the request is automatically rejected.
//
// You can configure an moderation endpoint either as a standalone prompt guard setting
// or in addition to other request and response guard settings.
type Moderation struct {
	// Pass prompt data through an external moderation model endpoint,
	// which compares the request prompt input to predefined content rules.
	// Configure an OpenAI moderation endpoint.
	OpenAIModeration OpenAIModeration `json:"openAIModeration,omitempty"`
}

// Prompt guards to apply to requests sent by the client.
type Request struct {

	// A custom response message to return to the client. If not specified, defaults to
	// "The request was rejected due to inappropriate content".
	CustomResponse CustomResponse `json:"customResponse,omitempty"`

	// Regular expression (regex) matching for prompt guards and data masking.
	Regex Regex `json:"regex,omitempty"`

	// Configure a webhook to forward requests to for prompt guarding.
	Webhook Webhook `json:"webhook,omitempty"`

	// Pass prompt data through an external moderation model endpoint,
	// which compares the request prompt input to predefined content rules.
	Moderation Moderation `json:"moderation,omitempty"`
}

// Prompt guards to apply to responses returned by the LLM provider.
type Response struct {
	// Regular expression (regex) matching for prompt guards and data masking.
	Regex Regex `json:"regex,omitempty"`

	// Configure a webhook to forward responses to for prompt guarding.
	Webhook Webhook `json:"webhook,omitempty"`
}

// Set up prompt guards to block unwanted requests to the LLM provider and mask sensitive data.
// Prompt guards can be used to reject requests based on the content of the prompt, as well as
// mask responses based on the content of the response.
//
// This example rejects any request prompts that contain
// the string "credit card", and masks any credit card numbers in the response.
// ```yaml
// promptGuard:
//
//	request:
//	  customResponse:
//	    message: "Rejected due to inappropriate content"
//	  regex:
//	    action: REJECT
//	    matches:
//	    - pattern: "credit card"
//	      name: "CC"
//	response:
//	  regex:
//	    builtins:
//	    - CREDIT_CARD
//	    action: MASK
//
// ```
type AIPromptGuard struct {
	// Prompt guards to apply to requests sent by the client.
	Request Request `json:"request,omitempty"`
	// Prompt guards to apply to responses returned by the LLM provider.
	Response Response `json:"response,omitempty"`
}

// Provide defaults to merge with user input fields.
// Defaults do _not_ override the user input fields, unless you explicitly set `override` to `true`.
//
// Example overriding the system field for Anthropic:
// ```yaml
// # Anthropic doesn't support a system chat type
// defaults:
//   - field: "system"
//     value: "answer all questions in french"
//
// ```
//
// Example setting the temperature and overriding `max_tokens`:
// ```yaml
// defaults:
//   - field: "temperature"
//     value: 0.5
//   - field: "max_tokens"
//     value: 100
//
// ```
type FieldDefault struct {
	// The name of the field.
	Field string `json:"field,omitempty"`
	// The field default value, which can be any JSON Data Type.
	Value interface{} `json:"value,omitempty"`
	// Whether to override the field's value if it already exists.
	// Defaults to false.
	Override bool `json:"override,omitempty"`
}

// When you deploy the Gloo AI Gateway, you can use the `spec.options.ai` section
// of the RouteOptions resource to configure the behavior of the LLM provider
// on the level of individual routes. These route settings, such as prompt enrichment,
// retrieval augmented generation (RAG), and semantic caching, are applicable only
// for routes that send requests to an LLM provider backend.
//
// For more information about the RouteOptions resource, see the
// [API reference]({{% versioned_link_path fromRoot="/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/route_options.proto.sk/" %}}).
type AI struct {

	// Enrich requests sent to the LLM provider by appending and prepending system prompts.
	// This can be configured only for LLM providers that use the `CHAT` API route type.
	PromptEnrichment AIPromptEnrichment `json:"prompt_enrichment,omitempty"`

	// Set up prompt guards to block unwanted requests to the LLM provider and mask sensitive data.
	// Prompt guards can be used to reject requests based on the content of the prompt, as well as
	// mask responses based on the content of the response.
	PromptGuard AIPromptGuard `json:"prompt_guard,omitempty"`

	// Provide defaults to merge with user input fields.
	// Defaults do _not_ override the user input fields, unless you explicitly set `override` to `true`.
	defaults []FieldDefault `json:"defaults,omitempty"`

	// The type of route to the LLM provider API. Currently, `CHAT` and `CHAT_STREAMING` are supported.
	RouteType RouteType `json:"route_type,omitempty"`
}
