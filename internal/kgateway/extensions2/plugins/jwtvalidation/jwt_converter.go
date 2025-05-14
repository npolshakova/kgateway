package jwtvalidation

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
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/hashicorp/go-multierror"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

func convertJwtValidationConfig(krtctx krt.HandlerContext, policy *v1alpha1.JWTValidationPolicy, secrets *krtcollections.SecretIndex) (
	*jwtauthnv3.JwtAuthentication,
	error,
) {

	var errs []error

	// Add all unique providers to the filter config providers list (note: route-level providers should be unique because of the route name/namespace)
	uniqProviders := map[string]*jwtauthnv3.JwtProvider{}
	for name, provider := range policy.Spec.Providers {
		p, err := translateProvider(krtctx, provider, policy.GetNamespace(), secrets)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// TODO: do we need a prefix?
		uniqProviders[name] = p
	}

	requirements := map[string]*jwtauthnv3.JwtRequirement{}
	policyTargetName := buildPolicyTargetName(policy.GetName(), policy.GetNamespace(), policy.Spec.TargetRefs)
	req := buildJwtRequirementFromProviders(policyTargetName, uniqProviders)
	// TODO: do I need a prefix for route level providers
	requirements[ProviderName(policy.GetName(), policy.GetNamespace())] = req

	// Create the filter config for the given stage
	return &jwtauthnv3.JwtAuthentication{
		RequirementMap: requirements,
		Providers:      uniqProviders,
		FilterStateRules: &jwtauthnv3.FilterStateRule{
			// TODO: filter stage name
			//Name:     getFilterStateNameForStage(stage),
			Requires: make(map[string]*jwtauthnv3.JwtRequirement),
		},
	}, nil
}

func buildPolicyTargetName(name string, namespace string, refs []v1alpha1.LocalPolicyTargetReference) string {
	return fmt.Sprintf("%s_%s_%s", name, namespace, refs)
}

func translateProvider(krtctx krt.HandlerContext, j *v1alpha1.JWTProvider, policyNs string, secrets *krtcollections.SecretIndex) (*jwtauthnv3.JwtProvider, error) {
	outProvider := &jwtauthnv3.JwtProvider{
		Issuer:    j.Issuer,
		Audiences: j.Audiences,
		// TODO: add support for keep token? (agentgateway does not support it)
	}
	translateTokenSource(j, outProvider)

	err := translateJwks(krtctx, j.JWKS, secrets, policyNs, outProvider)
	return outProvider, err
}

func translateTokenSource(j *v1alpha1.JWTProvider, provider *jwtauthnv3.JwtProvider) {
	// TODO: add support for token source?
}

// ProviderName returns a unique name for a provider in the context of a route
func ProviderName(resourceName, providerName string) string {
	return fmt.Sprintf("%s_%s", resourceName, providerName)
}

//func buildJwtRequirementWithAllowMissingOrFailed(ctx context.Context, jwtReq *jwtauthnv3.JwtRequirement, validationConfig v1alpha1.JWTValidationPolicy, legacyAllowMissingOrFailed bool) *jwtauthnv3.JwtRequirement {
//	missingOrFailedReq := &jwtauthnv3.JwtRequirement{
//		RequiresType: &jwtauthnv3.JwtRequirement_AllowMissingOrFailed{
//			AllowMissingOrFailed: &empty.Empty{},
//		},
//	}
//	// TODO: support configuration
//
//	// check if legacy allow_missing_or_failed is set
//	if legacyAllowMissingOrFailed {
//		jwtReq = &jwtauthnv3.JwtRequirement{
//			RequiresType: &jwtauthnv3.JwtRequirement_RequiresAny{
//				// Requires Any will OR the two requirements
//				RequiresAny: &jwtauthnv3.JwtRequirementOrList{
//					Requirements: []*jwtauthnv3.JwtRequirement{
//						jwtReq,
//						missingOrFailedReq,
//					},
//				},
//			},
//		}
//	}
//
//	return jwtReq
//}

func translateJwks(krtctx krt.HandlerContext, jwkConfig v1alpha1.JWKS, secrets *krtcollections.SecretIndex, policyNs string, out *jwtauthnv3.JwtProvider) error {
	// TODO: make configurable?
	const RemoteJwksTimeoutSecs = 5

	var err error
	switch jwkConfig.Kind {
	case v1alpha1.Remote:
		out.JwksSourceSpecifier = &jwtauthnv3.JwtProvider_RemoteJwks{
			RemoteJwks: &jwtauthnv3.RemoteJwks{
				// TODO: Support CacheDuration and AsyncFetch
				HttpUri: &envoycore.HttpUri{
					Timeout:          &duration.Duration{Seconds: RemoteJwksTimeoutSecs},
					Uri:              jwkConfig.RemoteJwks.URL,
					HttpUpstreamType: &envoycore.HttpUri_Cluster{
						// TODO:
						//Cluster: TargetRefToClusterName(jwkConfig.RemoteJwks.TargetRef),
					},
				},
			},
		}
	case v1alpha1.Local:
		var jwkSource *jwtauthnv3.JwtProvider_LocalJwks
		switch jwkConfig.LocalJwks.Kind {
		case v1alpha1.JWKInline:
			if jwkConfig.LocalJwks.InlineKey != nil {
				return errors.New("inline key is required for inline jwks")
			}
			jwkSource, err = translateJwksInline(*jwkConfig.LocalJwks.InlineKey)
		case v1alpha1.JWKFile:
			if jwkConfig.LocalJwks.File != nil {
				return errors.New("file name is required for inline jwks")
			}
			jwkSource, err = translateJwksFile(*jwkConfig.LocalJwks.File)
		case v1alpha1.JWKSecret:
			// TODO: check unique?
			if jwkConfig.LocalJwks.SecretRef == nil {
				return errors.New("secret ref is required for local jwks")
			}
			secret, err := GetSecretIr(secrets, krtctx, jwkConfig.LocalJwks.SecretRef.Name, policyNs)
			if err != nil {
				return errors.New("failed to find secret " + jwkConfig.LocalJwks.SecretRef.Name)
			}

			jwkSource, err = translateJwksSecret(jwkConfig.LocalJwks.SecretRef, secret)
		}
		out.JwksSourceSpecifier = jwkSource
	default:
		return errors.New("unknown jwks source")
	}
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
		return nil, fmt.Errorf("failed to parse inline jwks", err)
	}

	keysetJson, err := json.Marshal(keyset)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize inline jwks", err)
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

func buildJwtRequirementFromProviders(routeName string, providersMap map[string]*jwtauthnv3.JwtProvider) *jwtauthnv3.JwtRequirement {
	var reqs []*jwtauthnv3.JwtRequirement
	for provider := range providersMap {
		providerName := provider
		if routeName != "" {
			// if we have a route name, we need to make sure the provider name is unique
			providerName = ProviderName(routeName, provider)
		}
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
