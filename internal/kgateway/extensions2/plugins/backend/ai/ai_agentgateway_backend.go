package ai

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/agentgateway/agentgateway/go/api"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

func ProcessAIBackendForAgentGateway(ctx krt.HandlerContext, be *v1alpha1.Backend, secrets krt.Collection[*corev1.Secret]) ([]*api.Backend, []*api.Policy, error) {
	if be.Spec.AI == nil {
		return nil, nil, fmt.Errorf("ai backend spec must not be nil for AI backend type")
	}

	// Extract the provider configuration
	var authPolicy *api.Policy
	var aiBackend *api.Backend

	if be.Spec.AI.LLM != nil {
		aiBackend, authPolicy = buildAIBackendFromLLM(ctx, be.Namespace, be.Name, be.Spec.AI.LLM, secrets)
	} else if be.Spec.AI.MultiPool != nil && len(be.Spec.AI.MultiPool.Priorities) > 0 &&
		len(be.Spec.AI.MultiPool.Priorities[0].Pool) > 0 {
		// For MultiPool, use the first provider from the first priority pool
		aiBackend, authPolicy = buildAIBackendFromLLM(ctx, be.Namespace, be.Name, &be.Spec.AI.MultiPool.Priorities[0].Pool[0], secrets)
	} else {
		return nil, nil, fmt.Errorf("AI backend has no valid LLM or MultiPool configuration")
	}

	return []*api.Backend{aiBackend}, []*api.Policy{authPolicy}, nil
}

// buildAIBackendFromLLM converts a kgateway LLMProvider to an agentgateway AIBackend
func buildAIBackendFromLLM(
	ctx krt.HandlerContext,
	namespace, name string,
	llm *v1alpha1.LLMProvider,
	secrets krt.Collection[*corev1.Secret]) (*api.Backend, *api.Policy) {
	beName := namespace + "/" + name
	// Create AIBackend structure with provider-specific configuration
	aiBackend := &api.AIBackend{}

	// Extract and set provider configuration based on the LLM provider type
	provider := llm.Provider

	var auth *api.BackendAuthPolicy
	if provider.OpenAI != nil {
		var model *wrappers.StringValue
		if provider.OpenAI.Model != nil {
			model = &wrappers.StringValue{
				Value: *provider.OpenAI.Model,
			}
		}
		aiBackend.Provider = &api.AIBackend_Openai{
			Openai: &api.AIBackend_OpenAI{
				Model: model,
			},
		}
		auth = buildAuthPolicy(ctx, &provider.OpenAI.AuthToken, secrets, namespace)
	} else if provider.AzureOpenAI != nil {
		aiBackend.Provider = &api.AIBackend_Openai{
			Openai: &api.AIBackend_OpenAI{},
		}
		auth = buildAuthPolicy(ctx, &provider.AzureOpenAI.AuthToken, secrets, namespace)
	} else if provider.Anthropic != nil {
		var model *wrappers.StringValue
		if provider.Anthropic.Model != nil {
			model = &wrappers.StringValue{
				Value: *provider.Anthropic.Model,
			}
		}
		aiBackend.Provider = &api.AIBackend_Anthropic_{
			Anthropic: &api.AIBackend_Anthropic{
				Model: model,
			},
		}
		auth = buildAuthPolicy(ctx, &provider.Anthropic.AuthToken, secrets, namespace)
	} else if provider.Gemini != nil {
		model := &wrappers.StringValue{
			Value: provider.Gemini.Model,
		}
		aiBackend.Provider = &api.AIBackend_Gemini_{
			Gemini: &api.AIBackend_Gemini{
				Model: model,
			},
		}
		auth = buildAuthPolicy(ctx, &provider.Gemini.AuthToken, secrets, namespace)
	} else if provider.VertexAI != nil {
		model := &wrappers.StringValue{
			Value: provider.VertexAI.Model,
		}
		aiBackend.Provider = &api.AIBackend_Vertex_{
			Vertex: &api.AIBackend_Vertex{
				Model: model,
			},
		}
		auth = buildAuthPolicy(ctx, &provider.VertexAI.AuthToken, secrets, namespace)
	} else if provider.Bedrock != nil {
		model := &wrappers.StringValue{
			Value: provider.Bedrock.Model,
		}
		region := wellknown.DefaultAWSRegion
		if provider.Bedrock.Region != nil {
			region = *provider.Bedrock.Region
		}
		var guardrailIdentifier, guardrailVersion *wrappers.StringValue
		if provider.Bedrock.Guardrail != nil {
			guardrailIdentifier = &wrappers.StringValue{
				Value: provider.Bedrock.Guardrail.GuardrailIdentifier,
			}
			guardrailVersion = &wrappers.StringValue{
				Value: provider.Bedrock.Guardrail.GuardrailVersion,
			}
		}

		aiBackend.Provider = &api.AIBackend_Bedrock_{
			Bedrock: &api.AIBackend_Bedrock{
				Model:               model,
				Region:              region,
				GuardrailIdentifier: guardrailIdentifier,
				GuardrailVersion:    guardrailVersion,
			},
		}
		// TODO: handle errors on report
		auth, _ = buildBedrockAuthPolicy(ctx, region, provider.Bedrock.Auth, secrets, namespace)
	}

	// Map common override configurations
	if llm.HostOverride != nil {
		aiBackend.Override = &api.AIBackend_Override{
			Host: llm.HostOverride.Host,
			Port: int32(llm.HostOverride.Port),
		}
	}

	return &api.Backend{
			Name: beName,
			Kind: &api.Backend_Ai{
				Ai: aiBackend,
			},
		}, &api.Policy{
			Name: fmt.Sprintf("auth-%s", beName),
			Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{
				Backend: beName,
			}},
			Spec: &api.PolicySpec{Kind: &api.PolicySpec_Auth{
				Auth: auth,
			}},
		}
}

// getSecretValue extracts a value from a Kubernetes secret, handling both Data and StringData fields.
// It prioritizes StringData over Data if both are present (consistent with Kubernetes behavior).
func getSecretValue(secret *corev1.Secret, key string) (string, bool) {
	// First check StringData (if it exists, it takes precedence)
	if secret.StringData != nil {
		if value, exists := secret.StringData[key]; exists {
			return value, true
		}
	}

	// Then check Data (base64-encoded)
	if secret.Data != nil {
		if value, exists := secret.Data[key]; exists && utf8.Valid(value) {
			return string(value), true
		}
	}

	return "", false
}

func buildBedrockAuthPolicy(ctx krt.HandlerContext, region string, auth *v1alpha1.AwsAuth, secrets krt.Collection[*corev1.Secret], namespace string) (*api.BackendAuthPolicy, error) {
	var errs []error
	if auth == nil {
		return nil, nil
	}

	switch auth.Type {
	case v1alpha1.AwsAuthTypeSecret:
		if auth.SecretRef == nil {
			return nil, nil
		}

		secretRef := auth.SecretRef
		secretKey := namespace + "/" + secretRef.Name
		secret := krt.FetchOne(ctx, secrets, krt.FilterKey(secretKey))
		if secret == nil {
			// Return nil auth policy if secret not found - this will be handled upstream
			return nil, nil
		}

		var accessKeyId, secretAccessKey string
		var sessionToken *string

		// Extract access key
		if value, exists := getSecretValue(*secret, wellknown.AccessKey); !exists {
			errs = append(errs, errors.New("access_key is missing or not a valid string"))
		} else {
			accessKeyId = value
		}

		// Extract secret key
		if value, exists := getSecretValue(*secret, wellknown.SecretKey); !exists {
			errs = append(errs, errors.New("secret_key is missing or not a valid string"))
		} else {
			secretAccessKey = value
		}

		// Extract session token (optional)
		if value, exists := getSecretValue(*secret, wellknown.SessionToken); exists {
			sessionToken = ptr.To(value)
		}

		return &api.BackendAuthPolicy{
			Kind: &api.BackendAuthPolicy_Aws{
				Aws: &api.Aws{
					Kind: &api.Aws_ExplicitConfig{
						ExplicitConfig: &api.AwsExplicitConfig{
							AccessKeyId:     accessKeyId,
							SecretAccessKey: secretAccessKey,
							SessionToken:    sessionToken,
							Region:          region,
						},
					},
				},
			},
		}, errors.Join(errs...)
	default:
		errs = append(errs, errors.New("unknown AWS auth type"))
		return nil, errors.Join(errs...)
	}
}

// buildAuthPolicy creates auth policy for the given auth token configuration
func buildAuthPolicy(ctx krt.HandlerContext, authToken *v1alpha1.SingleAuthToken, secrets krt.Collection[*corev1.Secret], namespace string) *api.BackendAuthPolicy {
	if authToken == nil {
		return nil
	}

	switch authToken.Kind {
	case v1alpha1.SecretRef:
		if authToken.SecretRef == nil {
			return nil
		}

		// Build the secret key in namespace/name format
		secretKey := namespace + "/" + authToken.SecretRef.Name
		secret := krt.FetchOne(ctx, secrets, krt.FilterKey(secretKey))
		if secret == nil {
			// Return nil auth policy if secret not found - this will be handled upstream
			return nil
		}

		// Extract the authorization key from the secret data
		authKey := ""
		if authValue, exists := getSecretValue(*secret, "Authorization"); exists {
			// Strip the "Bearer " prefix if present, as it will be added by the provider
			authValue = strings.TrimSpace(authValue)
			authKey = strings.TrimSpace(strings.TrimPrefix(authValue, "Bearer "))
		}

		if authKey == "" {
			return nil
		}

		return &api.BackendAuthPolicy{
			Kind: &api.BackendAuthPolicy_Key{
				Key: &api.Key{Secret: authKey},
			},
		}
	case v1alpha1.Passthrough:
		return &api.BackendAuthPolicy{
			Kind: &api.BackendAuthPolicy_Passthrough{},
		}
	default:
		return nil
	}
}
