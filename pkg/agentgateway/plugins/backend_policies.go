package plugins

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/jwks"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/translator/sslutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	aiPolicySuffix                = ":ai"
	backendauthPolicySuffix       = ":backend-auth"
	tlsPolicySuffix               = ":tls"
	mcpAuthorizationPolicySuffix  = ":mcp-authorization"
	mcpAuthenticationPolicySuffix = ":mcp-authentication"
)

func translateBackendPolicyToAgw(
	ctx PolicyCtx,
	policy *v1alpha1.AgentgatewayPolicy,
	policyTarget *api.PolicyTarget,
) ([]AgwPolicy, error) {
	backend := policy.Spec.Backend
	if backend == nil {
		return nil, nil
	}
	agwPolicies := make([]AgwPolicy, 0)
	var errs []error

	policyName := getBackendPolicyName(policy.Namespace, policy.Name)

	if s := backend.HTTP; s != nil {
		pol, err := translateBackendHTTP(ctx, policy, policyName, policyTarget)
		if err != nil {
			logger.Error("error processing backend HTTP", "err", err)
			errs = append(errs, err)
		}
		agwPolicies = append(agwPolicies, pol...)
	}

	if s := backend.TLS; s != nil {
		pol, err := translateBackendTLS(ctx, policy, policyTarget)
		if err != nil {
			logger.Error("error processing backend TLS", "err", err)
			errs = append(errs, err)
		}
		agwPolicies = append(agwPolicies, pol...)
	}

	if s := backend.TCP; s != nil {
		pol, err := translateBackendTCP(ctx, policy, policyName, policyTarget)
		if err != nil {
			logger.Error("error processing backend TCP", "err", err)
			errs = append(errs, err)
		}
		agwPolicies = append(agwPolicies, pol...)
	}

	if s := backend.MCP; s != nil {
		pol := translateBackendMCP(ctx, policy, policyTarget)
		agwPolicies = append(agwPolicies, pol...)
	}

	if s := backend.AI; s != nil {
		pol, err := translateBackendAI(ctx, policy, policyName, policyTarget)
		if err != nil {
			logger.Error("error processing backend Tracing", "err", err)
			errs = append(errs, err)
		}
		agwPolicies = append(agwPolicies, pol...)
	}

	if s := backend.Auth; s != nil {
		pol, err := translateBackendAuth(ctx, policy, policyName, policyTarget)
		if err != nil {
			logger.Error("error processing backend Tracing", "err", err)
			errs = append(errs, err)
		}
		agwPolicies = append(agwPolicies, pol...)
	}

	return agwPolicies, errors.Join(errs...)
}

func translateBackendTCP(ctx PolicyCtx, policy *v1alpha1.AgentgatewayPolicy, name string, target *api.PolicyTarget) ([]AgwPolicy, error) {
	// TODO
	return nil, nil
}
func translateBackendTLS(ctx PolicyCtx, policy *v1alpha1.AgentgatewayPolicy, target *api.PolicyTarget) ([]AgwPolicy, error) {
	var errs []error

	// Build CA bundle from referenced ConfigMaps, if provided
	var caCert *wrapperspb.BytesValue
	if tls := policy.Spec.Backend.TLS; tls != nil && len(tls.CACertificateRefs) > 0 {
		var sb strings.Builder
		for _, ref := range tls.CACertificateRefs {
			nn := types.NamespacedName{Namespace: policy.Namespace, Name: ref.Name}
			cfgmap := krt.FetchOne(ctx.Krt, ctx.Collections.ConfigMaps, krt.FilterObjectName(nn))
			if cfgmap == nil {
				errs = append(errs, fmt.Errorf("ConfigMap %s not found", nn))
				continue
			}
			pem, err := sslutils.GetCACertFromConfigMap(ptr.Flatten(cfgmap))
			if err != nil {
				errs = append(errs, fmt.Errorf("error extracting CA cert from ConfigMap %s: %w", nn, err))
				continue
			}
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(pem)
		}
		if sb.Len() > 0 {
			caCert = wrapperspb.Bytes([]byte(sb.String()))
		}
	}

	// Map verify SANs to Hostname if provided (use first entry only)
	var hostname *wrapperspb.StringValue
	if tls := policy.Spec.Backend.TLS; tls != nil && len(tls.VerifySubjectAltNames) > 0 {
		hostname = wrapperspb.String(tls.VerifySubjectAltNames[0])
	}

	// Map insecure modes
	var insecure *wrapperspb.BoolValue
	if tls := policy.Spec.Backend.TLS; tls != nil && tls.InsecureSkipVerify != nil {
		switch *tls.InsecureSkipVerify {
		case v1alpha1.InsecureTLSModeAll:
			insecure = wrapperspb.Bool(true)
		case v1alpha1.InsecureTLSModeHostname:
			// Not directly supported in agentgateway API; fall back to default verification
		}
	}

	tlsPolicy := &api.Policy{
		Name:   policy.Namespace + "/" + policy.Name + tlsPolicySuffix + attachmentName(target),
		Target: target,
		Kind: &api.Policy_Backend{
			Backend: &api.BackendPolicySpec{
				Kind: &api.BackendPolicySpec_BackendTls{
					BackendTls: &api.BackendPolicySpec_BackendTLS{
						Root:     caCert,
						Cert:     nil,
						Key:      nil,
						Insecure: insecure,
						Hostname: hostname,
					},
				},
			}},
	}

	logger.Debug("generated TLS policy",
		"policy", policy.Name,
		"agentgateway_policy", tlsPolicy.Name)

	return []AgwPolicy{{Policy: tlsPolicy}}, errors.Join(errs...)
}
func translateBackendHTTP(ctx PolicyCtx, policy *v1alpha1.AgentgatewayPolicy, name string, target *api.PolicyTarget) ([]AgwPolicy, error) {
	// TODO
	return nil, nil
}

func translateBackendMCP(ctx PolicyCtx, policy *v1alpha1.AgentgatewayPolicy, target *api.PolicyTarget) []AgwPolicy {
	backend := policy.Spec.Backend
	if backend == nil || backend.MCP == nil {
		return nil
	}

	var mcpPolicies []AgwPolicy

	if backend.MCP.Authorization != nil {
		mcpPolicies = append(mcpPolicies, translateMCPAuthzPolicy(backend, policy.Name, policy.Namespace, target)...)
	}

	if backend.MCP.Authentication != nil {
		mcpPolicies = append(mcpPolicies, translateMCPAuthnPolicy(ctx, backend, policy.Name, policy.Namespace, target)...)
	}

	return mcpPolicies
}

func translateMCPAuthzPolicy(backend *v1alpha1.AgentgatewayPolicyBackend, policyName, policyNs string, target *api.PolicyTarget) []AgwPolicy {
	authzPolicy := backend.MCP.Authorization
	if authzPolicy == nil {
		return nil
	}
	var allowPolicies, denyPolicies []string
	if authzPolicy.Action == v1alpha1.AuthorizationPolicyActionDeny {
		denyPolicies = append(denyPolicies, cast(authzPolicy.Policy.MatchExpressions)...)
	} else {
		allowPolicies = append(allowPolicies, cast(authzPolicy.Policy.MatchExpressions)...)
	}
	mcpAuthorization := &api.BackendPolicySpec_McpAuthorization{
		Allow: allowPolicies,
		Deny:  denyPolicies,
	}

	mcpAuthorizationPolicy := &api.Policy{
		Name:   policyNs + "/" + policyName + mcpAuthorizationPolicySuffix + attachmentName(target),
		Target: target,
		Kind: &api.Policy_Backend{
			Backend: &api.BackendPolicySpec{
				Kind: &api.BackendPolicySpec_McpAuthorization_{
					McpAuthorization: mcpAuthorization,
				},
			}},
	}

	logger.Debug("generated MCP authorization policy",
		"policy", policyName,
		"agentgateway_policy", mcpAuthorizationPolicy.Name)

	return []AgwPolicy{{Policy: mcpAuthorizationPolicy}}
}

func translateMCPAuthnPolicy(ctx PolicyCtx, backend *v1alpha1.AgentgatewayPolicyBackend, policyName, policyNs string, target *api.PolicyTarget) []AgwPolicy {
	authnPolicy := backend.MCP.Authentication
	if authnPolicy == nil {
		return nil
	}

	var errs []error
	var idp api.BackendPolicySpec_McpAuthentication_McpIDP
	if authnPolicy.McpIDP == v1alpha1.Auth0 {
		idp = api.BackendPolicySpec_McpAuthentication_AUTH0
	} else if authnPolicy.McpIDP == v1alpha1.Keycloak {
		idp = api.BackendPolicySpec_McpAuthentication_KEYCLOAK
	}

	// TODO: share logic with jwt translation
	if _, err := url.Parse(authnPolicy.JWKS.JwksUri); err != nil {
		logger.Error("invalid jwks url in JWTAuthentication policy", "jwks_uri", authnPolicy.JWKS.JwksUri)
		errs = append(errs, fmt.Errorf("invalid jwks url in JWTAuthentication policy %w", err))
		return nil
	}
	jwksStoreName := jwks.JwksConfigMapNamespacedName(authnPolicy.JWKS.JwksUri)
	if jwksStoreName == nil {
		logger.Error("jwks store name not found", "jwks_uri", authnPolicy.JWKS.JwksUri)
		errs = append(errs, fmt.Errorf("jwks store hasn't been initialized"))
		return nil
	}
	jwksCM := ptr.Flatten(krt.FetchOne(ctx.Krt, ctx.Collections.ConfigMaps, krt.FilterObjectName(*jwksStoreName)))
	if jwksCM == nil {
		logger.Error("jwks ConfigMap not found", "name", jwksStoreName.Name, "namespace", jwksStoreName.Namespace)
		errs = append(errs, fmt.Errorf("jwks ConfigMap isn't available"))
		return nil
	}
	jwksForUri, err := jwks.JwksFromConfigMap(jwksCM)
	if err != nil {
		logger.Error("error deserializing jwks ConfigMap", "name", jwksStoreName.Name, "namespace", jwksStoreName.Namespace, "error", err)
		errs = append(errs, fmt.Errorf("error deserializing jwks ConfigMap %w", err))
		return nil
	}
	translatedInlineJwks, ok := jwksForUri[authnPolicy.JWKS.JwksUri]
	if !ok {
		logger.Error("jwks is not available in the jwks ConfigMap", "uri", authnPolicy.JWKS.JwksUri)
		errs = append(errs, fmt.Errorf("jwks %s is not available in the jwks ConfigMap", authnPolicy.JWKS.JwksUri))
		return nil
	}

	var mode api.BackendPolicySpec_McpAuthentication_Mode
	if authnPolicy.Mode == v1alpha1.JWTAuthenticationModeOptional {
		mode = api.BackendPolicySpec_McpAuthentication_OPTIONAL
	} else if authnPolicy.Mode == v1alpha1.JWTAuthenticationModePermissive {
		mode = api.BackendPolicySpec_McpAuthentication_PERMISSIVE
	} else if authnPolicy.Mode == v1alpha1.JWTAuthenticationModeStrict {
		mode = api.BackendPolicySpec_McpAuthentication_STRICT
	}

	mcpAuthn := &api.BackendPolicySpec_McpAuthentication{
		Issuer:    authnPolicy.Issuer,
		Audiences: authnPolicy.Audiences,
		Provider:  idp,
		ResourceMetadata: &api.BackendPolicySpec_McpAuthentication_ResourceMetadata{
			Extra: authnPolicy.ResourceMetadata,
		},
		JwksInline: translatedInlineJwks,
		Mode:       mode,
	}

	mcpAuthnPolicy := &api.Policy{
		Name:   policyNs + "/" + policyName + mcpAuthenticationPolicySuffix + attachmentName(target),
		Target: target,
		Kind: &api.Policy_Backend{
			Backend: &api.BackendPolicySpec{
				Kind: &api.BackendPolicySpec_McpAuthentication_{
					McpAuthentication: mcpAuthn,
				},
			}},
	}

	logger.Debug("generated MCP authentication policy",
		"policy", policyName,
		"agentgateway_policy", mcpAuthnPolicy.Name)

	return []AgwPolicy{{Policy: mcpAuthnPolicy}}
}

// translateBackendAI processes AI configuration and creates corresponding Agw policies
func translateBackendAI(ctx PolicyCtx, agwPolicy *v1alpha1.AgentgatewayPolicy, name string, policyTarget *api.PolicyTarget) ([]AgwPolicy, error) {
	var errs []error
	aiSpec := agwPolicy.Spec.Backend.AI

	translatedAIPolicy := &api.BackendPolicySpec_Ai{}
	if aiSpec.PromptEnrichment != nil {
		translatedAIPolicy.Prompts = processPromptEnrichment(aiSpec.PromptEnrichment)
	}

	for _, def := range aiSpec.Defaults {
		val, err := toJSONValue(def.Value)
		if err != nil {
			logger.Error("error parsing field value", "field", def.Field, "error", err)
			errs = append(errs, err)
			continue
		}
		if def.Override {
			if translatedAIPolicy.Overrides == nil {
				translatedAIPolicy.Overrides = make(map[string]string)
			}
			translatedAIPolicy.Overrides[def.Field] = val
		} else {
			if translatedAIPolicy.Defaults == nil {
				translatedAIPolicy.Defaults = make(map[string]string)
			}
			translatedAIPolicy.Defaults[def.Field] = val
		}
	}

	if aiSpec.PromptGuard != nil {
		if translatedAIPolicy.PromptGuard == nil {
			translatedAIPolicy.PromptGuard = &api.BackendPolicySpec_Ai_PromptGuard{}
		}
		//if aiSpec.PromptGuard.Request != nil {
		//	translatedAIPolicy.PromptGuard.Request = processRequestGuard(ctx.Krt, ctx.Collections.Secrets, agwPolicy.Namespace, aiSpec.PromptGuard.Request)
		//}
		//
		//if aiSpec.PromptGuard.Response != nil {
		//	translatedAIPolicy.PromptGuard.Response = processResponseGuard(aiSpec.PromptGuard.Response)
		//}
	}

	if aiSpec.ModelAliases != nil {
		translatedAIPolicy.ModelAliases = aiSpec.ModelAliases
	}

	if aiSpec.PromptCaching != nil {
		translatedAIPolicy.PromptCaching = &api.BackendPolicySpec_Ai_PromptCaching{
			CacheSystem:   aiSpec.PromptCaching.CacheSystem,
			CacheMessages: aiSpec.PromptCaching.CacheMessages,
			CacheTools:    aiSpec.PromptCaching.CacheTools,
		}
		translatedAIPolicy.PromptCaching.MinTokens = ptr.Of(uint32(aiSpec.PromptCaching.MinTokens)) //nolint:gosec // G115: MinTokens is validated by kubebuilder to be >= 0
	}

	aiPolicy := &api.Policy{
		Name:   name + aiPolicySuffix + attachmentName(policyTarget),
		Target: policyTarget,
		Kind: &api.Policy_Backend{
			Backend: &api.BackendPolicySpec{
				Kind: &api.BackendPolicySpec_Ai_{
					Ai: translatedAIPolicy,
				},
			},
		},
	}

	logger.Debug("generated AI policy",
		"policy", agwPolicy.Name,
		"agentgateway_policy", aiPolicy.Name)

	return []AgwPolicy{{Policy: aiPolicy}}, errors.Join(errs...)
}

func translateBackendAuth(ctx PolicyCtx, policy *v1alpha1.AgentgatewayPolicy, name string, target *api.PolicyTarget) ([]AgwPolicy, error) {
	var errs []error
	auth := policy.Spec.Backend.Auth

	if auth == nil {
		return nil, nil
	}

	var translatedAuth *api.BackendAuthPolicy

	if auth.InlineKey != nil && *auth.InlineKey != "" {
		translatedAuth = &api.BackendAuthPolicy{
			Kind: &api.BackendAuthPolicy_Key{
				Key: &api.Key{Secret: *auth.InlineKey},
			},
		}
	} else if auth.SecretRef != nil {
		// Resolve secret and extract Authorization value
		secret, err := kubeutils.GetSecret(ctx.Collections.Secrets, ctx.Krt, auth.SecretRef.Name, policy.Namespace)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get secret %s/%s: %w", policy.Namespace, auth.SecretRef.Name, err))
		} else {
			if authKey, ok := kubeutils.GetSecretAuth(secret); ok {
				translatedAuth = &api.BackendAuthPolicy{
					Kind: &api.BackendAuthPolicy_Key{
						Key: &api.Key{Secret: authKey},
					},
				}
			} else {
				errs = append(errs, fmt.Errorf("secret %s/%s missing Authorization value", policy.Namespace, auth.SecretRef.Name))
			}
		}
	} else {
		errs = append(errs, fmt.Errorf("backend auth requires either inline key or secretRef"))
	}

	if translatedAuth == nil {
		return nil, errors.Join(errs...)
	}

	authPolicy := &api.Policy{
		Name:   name + backendauthPolicySuffix + attachmentName(target),
		Target: target,
		Kind: &api.Policy_Backend{
			Backend: &api.BackendPolicySpec{
				Kind: &api.BackendPolicySpec_Auth{
					Auth: translatedAuth,
				},
			},
		},
	}
	logger.Debug("generated backend auth policy",
		"policy", policy.Name,
		"agentgateway_policy", authPolicy.Name)

	return []AgwPolicy{{Policy: authPolicy}}, errors.Join(errs...)
}
