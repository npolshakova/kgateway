package trafficpolicy

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	"github.com/go-jose/go-jose/v3"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

const (
	JwtFilterName = "envoy.filters.http.jwt_authn"
)

type JwtIr struct {
	jwtConfig        *jwtauthnv3.JwtAuthentication
	requirementsName string
}

func (j *JwtIr) Equals(other *JwtIr) bool {
	if j == nil && other == nil {
		return true
	}
	if j == nil || other == nil {
		return false
	}

	return proto.Equal(j.jwtConfig, other.jwtConfig)
}

// handleJwt configures the filter JwtAuthentication and per-route JWT configuration for a specific route
func (p *trafficPolicyPluginGwPass) handleJwt(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, jwtIr *JwtIr) {
	if jwtIr == nil || jwtIr.jwtConfig == nil {
		return
	}

	// Add a filter to the chain. When having a jwt policy for a route we need to also have a
	// global jwt http filter in the chain otherwise it will be ignored.
	if p.jwtInChain == nil {
		p.jwtInChain = make(map[string]*jwtauthnv3.JwtAuthentication)
	}
	if _, ok := p.jwtInChain[fcn]; !ok {
		p.jwtInChain[fcn] = jwtIr.jwtConfig
	}

	perRouteConfig, err := translatePerRouteConfig(jwtIr.requirementsName)
	if err != nil {
		// Log error but don't fail the route configuration
		return
	}

	// Add the per-route JWT configuration to the typed filter config
	pCtxTypedFilterConfig.AddTypedConfig(JwtFilterName, perRouteConfig)
}

func translatePerRouteConfig(requirementsName string) (*jwtauthnv3.PerRouteConfig, error) {
	perRouteConfig := &jwtauthnv3.PerRouteConfig{
		RequirementSpecifier: &jwtauthnv3.PerRouteConfig_RequirementName{
			RequirementName: requirementsName,
		},
	}
	return perRouteConfig, nil
}

// constructJwt translates the jwt spec into an envoy jwt policy and stores it in the traffic policy IR
func constructJwt(krtctx krt.HandlerContext, policy *v1alpha1.TrafficPolicy, out *trafficPolicySpecIr, secrets *krtcollections.SecretIndex) error {
	spec := policy.Spec
	if spec.JWT == nil {
		return nil
	}
	jwtIr, err := convertJwtValidationConfig(krtctx, policy.Name, policy.Namespace, spec.JWT, secrets)
	if err != nil {
		return err
	}
	out.jwt = jwtIr
	return nil
}

// Validate performs validation on the jwt component.
func (j *JwtIr) Validate() error {
	return j.validate()
}

func (j *JwtIr) validate() error {
	if j == nil {
		return nil
	}

	var errs []error
	if j.requirementsName == "" {
		errs = append(errs, errors.New("requirementsName is empty"))
	}

	err := j.jwtConfig.Validate()
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)

}

func convertJwtValidationConfig(krtctx krt.HandlerContext, policyName, policyNs string, policy *v1alpha1.JWTValidation, secrets *krtcollections.SecretIndex) (
	*JwtIr,
	error,
) {
	uniqProviders := make(map[string]*jwtauthnv3.JwtProvider)
	policyNameNamespace := fmt.Sprintf("%s_%s", policyName, policyNs)

	for providerName, provider := range policy.Providers {
		providerNameForPolicy := ProviderName(policyNameNamespace, providerName)
		jwtProvider, err := translateProvider(krtctx, provider, providerNameForPolicy, policyNs, secrets)
		if err != nil {
			return nil, err
		}
		uniqProviders[providerNameForPolicy] = jwtProvider
	}

	requirementsName := fmt.Sprintf("%s_requirements", policyNameNamespace)
	requirements := make(map[string]*jwtauthnv3.JwtRequirement)
	requirements[requirementsName] = buildJwtRequirementFromProviders(uniqProviders)

	// Note: only one stage is supported for now so we're not setting the filter state rules here
	jwtConfig := &jwtauthnv3.JwtAuthentication{
		RequirementMap: requirements,
		Providers:      uniqProviders,
	}

	return &JwtIr{
		jwtConfig:        jwtConfig,
		requirementsName: requirementsName,
	}, nil
}

// ProviderName returns a unique name for a provider in the context of a route
func ProviderName(resourceName, providerName string) string {
	return fmt.Sprintf("%s_%s", resourceName, providerName)
}

func translateProvider(krtctx krt.HandlerContext, provider v1alpha1.JWTProvider, providerNameForPolicy, policyNs string, secrets *krtcollections.SecretIndex) (*jwtauthnv3.JwtProvider, error) {
	var claimToHeaders []*jwtauthnv3.JwtClaimToHeader
	for _, claim := range provider.ClaimsToHeaders {
		claimToHeaders = append(claimToHeaders, &jwtauthnv3.JwtClaimToHeader{
			ClaimName:  claim.Name,
			HeaderName: claim.Header,
		})
	}
	var shouldForward bool
	if provider.KeepToken != nil && *provider.KeepToken == v1alpha1.TokenForward {
		shouldForward = true
	}
	jwtProvider := &jwtauthnv3.JwtProvider{
		Issuer:            provider.Issuer,
		Audiences:         provider.Audiences,
		PayloadInMetadata: providerNameForPolicy,
		ClaimToHeaders:    claimToHeaders,
		Forward:           shouldForward,
		// TODO(npolshak): Do we want to set NormalizePayload  to support https://datatracker.ietf.org/doc/html/rfc8693#name-scope-scopes-claim
	}
	translateTokenSource(provider, jwtProvider)
	err := translateJwks(krtctx, provider.JWKS, secrets, policyNs, jwtProvider)

	if err != nil {
		return nil, err
	}
	return jwtProvider, nil
}

func translateTokenSource(provider v1alpha1.JWTProvider, out *jwtauthnv3.JwtProvider) {
	if provider.TokenSource == nil {
		return
	}
	if provider.TokenSource.HeaderSource != nil {
		if headers := provider.TokenSource.HeaderSource; len(headers) != 0 {
			for _, header := range headers {
				var headerStr, prefixStr string
				if header.Header != nil {
					headerStr = *header.Header
				}
				if header.Prefix != nil {
					prefixStr = *header.Prefix
				}
				out.FromHeaders = append(out.GetFromHeaders(), &jwtauthnv3.JwtHeader{
					Name:        headerStr,
					ValuePrefix: prefixStr,
				})
			}
		}
	}
	if provider.TokenSource.QueryParams != nil {
		out.FromParams = provider.TokenSource.QueryParams
	}
}

func translateJwks(krtctx krt.HandlerContext, jwkConfig v1alpha1.JWKS, secrets *krtcollections.SecretIndex, policyNs string, out *jwtauthnv3.JwtProvider) error {
	var err error
	var secret *ir.Secret
	var jwkSource *jwtauthnv3.JwtProvider_LocalJwks
	if jwkConfig.LocalJWKS.File != nil {
		jwkSource, err = translateJwksFile(*jwkConfig.LocalJWKS.File)
	} else if jwkConfig.LocalJWKS.InlineKey != nil {
		jwkSource, err = translateJwksInline(*jwkConfig.LocalJWKS.InlineKey)
	} else if jwkConfig.LocalJWKS.SecretRef != nil {
		secret, err = GetSecretIr(secrets, krtctx, jwkConfig.LocalJWKS.SecretRef.Name, policyNs)
		if err != nil {
			return errors.New("failed to find secret " + jwkConfig.LocalJWKS.SecretRef.Name)
		}
		jwkSource, err = translateJwksSecret(jwkConfig.LocalJWKS.SecretRef, secret)
	}
	out.JwksSourceSpecifier = jwkSource
	return err
}

func translateJwksSecret(ref *corev1.LocalObjectReference, secret *ir.Secret) (*jwtauthnv3.JwtProvider_LocalJwks, error) {
	k8sSecret := secret.Obj.(*corev1.Secret)
	secretKey := k8sSecret.Data[ref.Name]
	if secretKey == nil {
		return nil, errors.New("secret key not found")
	}
	return translateJwksInline(string(secretKey))
}

func translateJwksFile(filename string) (*jwtauthnv3.JwtProvider_LocalJwks, error) {
	return &jwtauthnv3.JwtProvider_LocalJwks{
		LocalJwks: &envoycore.DataSource{
			Specifier: &envoycore.DataSource_Filename{
				Filename: filename,
			},
		},
	}, nil
}

func translateJwksInline(inlineKey string) (*jwtauthnv3.JwtProvider_LocalJwks, error) {
	keyset, err := TranslateKey(inlineKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inline jwks: %v", err)
	}

	keysetJson, err := json.Marshal(keyset)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize inline jwks: %v", err)
	}

	return &jwtauthnv3.JwtProvider_LocalJwks{
		LocalJwks: &envoycore.DataSource{
			Specifier: &envoycore.DataSource_InlineString{
				InlineString: string(keysetJson),
			},
		},
	}, nil
}

func TranslateKey(key string) (*jose.JSONWebKeySet, error) {
	// key can be an individual key, a key set or a pem block public key:
	// is it a pem block?
	var multierr error
	ks, err := parsePem(key)
	if err == nil {
		return ks, nil
	}
	multierr = multierror.Append(multierr, fmt.Errorf("PEM %v", err))

	ks, err = parseKeySet(key)
	if err == nil {
		if len(ks.Keys) != 0 {
			return ks, nil
		}
		err = errors.New("no keys in set")
	}
	multierr = multierror.Append(multierr, fmt.Errorf("JWKS %v", err))

	ks, err = parseKey(key)
	if err == nil {
		return ks, nil
	}
	multierr = multierror.Append(multierr, fmt.Errorf("JWK %v", err))

	return nil, fmt.Errorf("cannot parse local jwks: %v", multierr)
}

func parseKeySet(key string) (*jose.JSONWebKeySet, error) {
	var keyset jose.JSONWebKeySet
	err := json.Unmarshal([]byte(key), &keyset)
	return &keyset, err
}

func parseKey(key string) (*jose.JSONWebKeySet, error) {
	var jwk jose.JSONWebKey
	err := json.Unmarshal([]byte(key), &jwk)
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}, err
}

func parsePem(key string) (*jose.JSONWebKeySet, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	var err error
	var publicKey interface{}
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		publicKey, err = x509.ParsePKIXPublicKey(block.Bytes) // Parses both RS256 and PS256
		if err != nil {
			return nil, err
		}
	}

	alg := ""
	switch publicKey.(type) {
	// RS256 implied for hash
	case *rsa.PublicKey:
		alg = "RS256"

	case *ecdsa.PublicKey:
		alg = "ES256"

	case ed25519.PublicKey:
		alg = "EdDSA"

	default:
		// HS256 is not supported as this is only used by HMAC, which doesn't use public keys
		return nil, errors.New("unsupported public key. only RSA, ECDSA, and Ed25519 public keys are supported in PEM format")
	}

	jwk := jose.JSONWebKey{
		Key:       publicKey,
		Algorithm: alg,
		Use:       "sig",
	}
	keySet := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}
	return keySet, nil
}

func buildJwtRequirementFromProviders(providersMap map[string]*jwtauthnv3.JwtProvider) *jwtauthnv3.JwtRequirement {
	var reqs []*jwtauthnv3.JwtRequirement
	for providerName := range providersMap {
		reqs = append(reqs, &jwtauthnv3.JwtRequirement{
			RequiresType: &jwtauthnv3.JwtRequirement_ProviderName{
				ProviderName: providerName,
			},
		})
	}

	// sort for idempotency
	sort.Slice(reqs, func(i, j int) bool { return reqs[i].GetProviderName() < reqs[j].GetProviderName() })

	// if there is only one requirement, return it directly
	if len(reqs) == 1 {
		return reqs[0]
	}
	// if there are multiple requirements, return a RequiresAny requirement. Requires Any will OR the requirements
	return &jwtauthnv3.JwtRequirement{
		RequiresType: &jwtauthnv3.JwtRequirement_RequiresAny{
			RequiresAny: &jwtauthnv3.JwtRequirementOrList{
				Requirements: reqs,
			},
		},
	}
}

func GetSecretIr(secrets *krtcollections.SecretIndex, krtctx krt.HandlerContext, secretName, ns string) (*ir.Secret, error) {
	secretRef := gwv1.SecretObjectReference{
		Name: gwv1.ObjectName(secretName),
	}
	from := krtcollections.From{
		GroupKind: wellknown.BackendGVK.GroupKind(),
		Namespace: ns,
	}
	secret, err := secrets.GetSecret(krtctx, from, secretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find secret %s: %v", secretName, err)
	}
	return secret, nil
}
