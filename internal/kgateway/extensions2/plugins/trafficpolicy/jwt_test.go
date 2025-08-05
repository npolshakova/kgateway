package trafficpolicy

import (
	"testing"

	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
)

func TestTranslateKey(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		expectedError bool
		expectedKeys  int
	}{
		{
			name: "valid JWKS",
			key: `{
				"keys": [
					{
						"kty": "RSA",
						"kid": "test-key",
						"use": "sig",
						"alg": "RS256",
						"n": "test-n",
						"e": "AQAB"
					}
				]
			}`,
			expectedError: false,
			expectedKeys:  1,
		},
		{
			name: "valid single JWK",
			key: `{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n": "test-n",
				"e": "AQAB"
			}`,
			expectedError: false,
			expectedKeys:  1,
		},
		{
			name:          "invalid JSON",
			key:           "{invalid json}",
			expectedError: true,
			expectedKeys:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyset, err := TranslateKey(tt.key)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedKeys, len(keyset.Keys))
		})
	}
}

func TestBuildJwtRequirementFromProviders(t *testing.T) {
	tests := []struct {
		name          string
		routeName     string
		providers     map[string]*jwtauthnv3.JwtProvider
		expectedType  string
		expectedCount int
	}{
		{
			name:      "single provider",
			routeName: "test-route",
			providers: map[string]*jwtauthnv3.JwtProvider{
				"provider1": {Issuer: "test-issuer"},
			},
			expectedType:  "provider_name",
			expectedCount: 1,
		},
		{
			name:      "multiple providers",
			routeName: "test-route",
			providers: map[string]*jwtauthnv3.JwtProvider{
				"provider1": {Issuer: "test-issuer-1"},
				"provider2": {Issuer: "test-issuer-2"},
			},
			expectedType:  "requires_any",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := buildJwtRequirementFromProviders(tt.providers)
			if tt.expectedType == "provider_name" {
				assert.NotNil(t, req.GetProviderName())
				assert.Equal(t, "provider1", req.GetProviderName())
			} else {
				assert.NotNil(t, req.GetRequiresAny())
				assert.Equal(t, tt.expectedCount, len(req.GetRequiresAny().Requirements))
			}
		})
	}
}

func TestTranslateJwksSecret(t *testing.T) {
	tests := []struct {
		name          string
		secret        *corev1.Secret
		ref           *corev1.LocalObjectReference
		expectedError bool
	}{
		{
			name: "valid secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-secret",
				},
				Data: map[string][]byte{
					"test-key": []byte(`{"keys":[{"kty":"RSA","kid":"test-key","use":"sig","alg":"RS256","n":"test-n","e":"AQAB"}]}`),
				},
			},
			ref: &corev1.LocalObjectReference{
				Name: "test-key",
			},
			expectedError: false,
		},
		{
			name: "missing key in secret",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-secret",
				},
				Data: map[string][]byte{},
			},
			ref: &corev1.LocalObjectReference{
				Name: "test-key",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretIr := &ir.Secret{
				Obj: tt.secret,
			}
			jwks, err := translateJwksSecret(tt.ref, secretIr)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, jwks)
			assert.NotNil(t, jwks.LocalJwks)
		})
	}
}

func TestConvertJwtValidationConfig(t *testing.T) {
	tests := []struct {
		name           string
		providers      map[string]v1alpha1.JWTProvider
		expectedError  bool
		expectedConfig *jwtauthnv3.JwtAuthentication
	}{
		{
			name: "basic provider with inline JWKS",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer: "test-issuer",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key","use":"sig","alg":"RS256","n":"test-n","e":"AQAB"}]}`),
						},
					},
					ClaimsToHeaders: []v1alpha1.JWTClaimToHeader{
						{
							Name:   "sub",
							Header: "X-Subject",
						},
					},
					KeepToken: ptr.To(v1alpha1.TokenForward),
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_test-provider": {
						Issuer:            "test-issuer",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
						ClaimToHeaders: []*jwtauthnv3.JwtClaimToHeader{
							{
								ClaimName:  "sub",
								HeaderName: "X-Subject",
							},
						},
						Forward: true,
					},
				},
			},
		},
		{
			name: "provider with file JWKS",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer: "test-issuer",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							File: ptr.To("/path/to/jwks.json"),
						},
					},
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_test-provider": {
						Issuer:            "test-issuer",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
					},
				},
			},
		},
		{
			name: "missing inline key for inline JWKS",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer: "test-issuer",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To("abc"),
						},
					},
				},
			},
			expectedError:  true,
			expectedConfig: nil,
		},
		{
			name: "multiple providers",
			providers: map[string]v1alpha1.JWTProvider{
				"provider1": {
					Issuer: "test-issuer-1",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key-1","use":"sig","alg":"RS256","n":"test-n-1","e":"AQAB"}]}`),
						},
					},
				},
				"provider2": {
					Issuer: "test-issuer-2",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key-2","use":"sig","alg":"RS256","n":"test-n-2","e":"AQAB"}]}`),
						},
					},
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_provider1": {
						Issuer:            "test-issuer-1",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
					},
					"test-policy_test-ns_provider2": {
						Issuer:            "test-issuer-2",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
					},
				},
			},
		},
		{
			name: "provider with audiences",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer:    "test-issuer",
					Audiences: []string{"aud1", "aud2"},
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key","use":"sig","alg":"RS256","n":"test-n","e":"AQAB"}]}`),
						},
					},
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_test-provider": {
						Issuer:            "test-issuer",
						Audiences:         []string{"aud1", "aud2"},
						PayloadInMetadata: PayloadInMetadata,
					},
				},
			},
		},
		{
			name: "provider with token source",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer: "test-issuer",
					TokenSource: &v1alpha1.JWTTokenSource{
						HeaderSource: []v1alpha1.HeaderSource{
							{
								Header: ptr.To("Authorization"),
							},
						},
					},
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key","use":"sig","alg":"RS256","n":"test-n","e":"AQAB"}]}`),
						},
					},
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_test-provider": {
						Issuer:            "test-issuer",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
					},
				},
			},
		},
		{
			name: "provider with remove token",
			providers: map[string]v1alpha1.JWTProvider{
				"test-provider": {
					Issuer: "test-issuer",
					JWKS: v1alpha1.JWKS{
						LocalJWKS: &v1alpha1.LocalJWKS{
							InlineKey: ptr.To(`{"keys":[{"kty":"RSA","kid":"test-key","use":"sig","alg":"RS256","n":"test-n","e":"AQAB"}]}`),
						},
					},
					KeepToken: ptr.To(v1alpha1.TokenRemove),
				},
			},
			expectedError: false,
			expectedConfig: &jwtauthnv3.JwtAuthentication{
				Providers: map[string]*jwtauthnv3.JwtProvider{
					"test-policy_test-ns_test-provider": {
						Issuer:            "test-issuer",
						Audiences:         nil,
						PayloadInMetadata: PayloadInMetadata,
						Forward:           false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := resolveJwtProviders(nil, "test-policy", "test-ns", tt.providers)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, config)
			assert.Equal(t, len(tt.expectedConfig.Providers), len(config.Providers))
			for providerName, expectedProvider := range tt.expectedConfig.Providers {
				actualProvider, ok := config.Providers[providerName]
				require.True(t, ok, "provider %s not found", providerName)
				assert.Equal(t, expectedProvider.Issuer, actualProvider.Issuer)
				assert.Equal(t, expectedProvider.Audiences, actualProvider.Audiences)
				assert.Equal(t, expectedProvider.PayloadInMetadata, actualProvider.PayloadInMetadata)
				assert.Equal(t, expectedProvider.Forward, actualProvider.Forward)

				// Check claim to headers
				assert.Equal(t, len(expectedProvider.ClaimToHeaders), len(actualProvider.ClaimToHeaders))
				for i, expectedClaim := range expectedProvider.ClaimToHeaders {
					actualClaim := actualProvider.ClaimToHeaders[i]
					assert.Equal(t, expectedClaim.ClaimName, actualClaim.ClaimName)
					assert.Equal(t, expectedClaim.HeaderName, actualClaim.HeaderName)
				}

				// Check token source
				if expectedProvider.FromHeaders != nil {
					assert.Equal(t, len(expectedProvider.FromHeaders), len(actualProvider.FromHeaders))
					for i, expectedHeader := range expectedProvider.FromHeaders {
						actualHeader := actualProvider.FromHeaders[i]
						assert.Equal(t, expectedHeader.Name, actualHeader.Name)
						assert.Equal(t, expectedHeader.ValuePrefix, actualHeader.ValuePrefix)
					}
				}
				assert.Equal(t, expectedProvider.FromParams, actualProvider.FromParams)

				// Check JWKS source
				if expectedProvider.JwksSourceSpecifier != nil {
					assert.NotNil(t, actualProvider.JwksSourceSpecifier)
					expectedJwks := expectedProvider.JwksSourceSpecifier.(*jwtauthnv3.JwtProvider_LocalJwks)
					actualJwks := actualProvider.JwksSourceSpecifier.(*jwtauthnv3.JwtProvider_LocalJwks)
					assert.NotNil(t, expectedJwks)
					assert.NotNil(t, actualJwks)
					assert.NotNil(t, actualJwks.LocalJwks)
				}
			}
		})
	}
}
