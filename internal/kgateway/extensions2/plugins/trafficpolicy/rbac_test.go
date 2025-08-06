package trafficpolicy

import (
	"fmt"
	"testing"

	cncfmatcherv3 "github.com/cncf/xds/go/xds/type/matcher/v3"
	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoy_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	proto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// extractCELExpressions extracts CEL expressions from a Permission_Matcher
func extractCELExpressions(t *testing.T, perm *envoycfgauthz.Permission) []string {
	matcherPerm, ok := perm.Rule.(*envoycfgauthz.Permission_Matcher)
	require.True(t, ok, "Expected Permission_Matcher")

	// Unmarshal the TypedConfig to get the Matcher
	var matcher cncfmatcherv3.Matcher
	err := anypb.UnmarshalTo(matcherPerm.Matcher.TypedConfig, &matcher, proto.UnmarshalOptions{})
	require.NoError(t, err)

	// Extract CEL expressions from the matcher list
	matcherList := matcher.GetMatcherList()
	require.NotNil(t, matcherList)

	var celExpressions []string
	for _, fieldMatcher := range matcherList.Matchers {
		predicate := fieldMatcher.Predicate
		singlePred := predicate.GetSinglePredicate()
		require.NotNil(t, singlePred)

		customMatch := singlePred.GetCustomMatch()
		require.NotNil(t, customMatch)

		var celMatcher cncfmatcherv3.CelMatcher
		err := anypb.UnmarshalTo(customMatch.TypedConfig, &celMatcher, proto.UnmarshalOptions{})
		require.NoError(t, err)

		exprMatch := celMatcher.GetExprMatch()
		require.NotNil(t, exprMatch)

		celExpressions = append(celExpressions, exprMatch.CelExprString)
	}

	return celExpressions
}

// mockGatewayExtensionStore is a map of extension names to their mock implementations
var mockGatewayExtensionStore = map[string]*TrafficPolicyGatewayExtensionIR{
	"test-provider": {
		Name:    "test-provider",
		ExtType: v1alpha1.GatewayExtensionTypeJWTProvider,
		JwtProviders: map[string]v1alpha1.JWTProvider{
			"test-provider": {
				Issuer: "test-issuer",
				JWKS: v1alpha1.JWKS{
					LocalJWKS: &v1alpha1.LocalJWKS{
						InlineKey: func() *string { s := "test-key"; return &s }(),
					},
				},
			},
		},
	},
}

// Mock function for fetchGatewayExtension that returns providers based on extension name
func fetchGatewayExtension(krtctx krt.HandlerContext, extensionRef *corev1.LocalObjectReference, ns string) (*TrafficPolicyGatewayExtensionIR, error) {
	if extensionRef == nil {
		return nil, fmt.Errorf("extension reference is nil")
	}

	// Look up the mock extension in our store
	if ext, exists := mockGatewayExtensionStore[extensionRef.Name]; exists {
		return ext, nil
	}

	// If not found, return a generic mock that should work for most test cases
	return &TrafficPolicyGatewayExtensionIR{
		Name:    extensionRef.Name,
		ExtType: v1alpha1.GatewayExtensionTypeJWTProvider,
		JwtProviders: map[string]v1alpha1.JWTProvider{
			extensionRef.Name: {
				Issuer: "mock-issuer",
				JWKS: v1alpha1.JWKS{
					LocalJWKS: &v1alpha1.LocalJWKS{
						InlineKey: func() *string { s := "mock-key"; return &s }(),
					},
				},
			},
		},
	}, nil
}

func TestTranslateRbac(t *testing.T) {
	tests := []struct {
		name             string
		ns               string
		tpName           string
		rbac             *v1alpha1.Rbac
		jwt              *v1alpha1.JWTValidation
		expected         *envoyauthz.RBACPerRoute
		expectedCELRules map[string][]string // policy name -> expected CEL expressions
		wantErr          bool
	}{
		{
			name:   "allow action with single rule",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Policies: []v1alpha1.RbacPolicy{
					{
						Name: "policy-0",
						Principals: []v1alpha1.Principal{
							{
								JWTPrincipals: map[string]*v1alpha1.JWTPrincipal{
									"test-provider": {
										Claims: []v1alpha1.JWTClaimMatch{
											{
												Name:    "sub",
												Value:   "test-user",
												Matcher: v1alpha1.ClaimMatcherExactString,
											},
										},
									},
								},
							},
						},
						Conditions: &v1alpha1.CELConditions{
							CelMatchExpression: []string{"request.auth.claims.groups == 'group1'", "request.auth.claims.groups == 'group2'"},
						},
					},
				},
			},
			jwt: &v1alpha1.JWTValidation{
				ExtensionRef: corev1.LocalObjectReference{
					Name: "test-provider",
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Rules: &envoycfgauthz.RBAC{
						Action: envoycfgauthz.RBAC_ALLOW,
						Policies: map[string]*envoycfgauthz.Policy{
							"ns[test-ns]-policy[test-policy]-rule[0]": {
								Principals: []*envoycfgauthz.Principal{
									{
										Identifier: &envoycfgauthz.Principal_SourcedMetadata{
											SourcedMetadata: &envoycfgauthz.SourcedMetadata{
												MetadataMatcher: &envoy_matcher_v3.MetadataMatcher{
													Filter: "envoy.filters.http.jwt_authn",
													Path: []*envoy_matcher_v3.MetadataMatcher_PathSegment{
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: PayloadInMetadata,
															},
														},
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: "sub",
															},
														},
													},
													Value: &envoy_matcher_v3.ValueMatcher{
														MatchPattern: &envoy_matcher_v3.ValueMatcher_StringMatch{
															StringMatch: &envoy_matcher_v3.StringMatcher{
																MatchPattern: &envoy_matcher_v3.StringMatcher_Exact{
																	Exact: "test-user",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Permissions: []*envoycfgauthz.Permission{
									{
										Rule: &envoycfgauthz.Permission_Matcher{
											Matcher: &envoycorev3.TypedExtensionConfig{
												Name: "cel-matcher",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedCELRules: map[string][]string{
				"ns[test-ns]-policy[test-policy]-rule[0]": {"request.auth.claims.groups == 'group1'", "request.auth.claims.groups == 'group2'"},
			},
			wantErr: false,
		},
		{
			name:   "deny action with empty rules",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action:   v1alpha1.AuthorizationPolicyActionDeny,
				Policies: []v1alpha1.RbacPolicy{},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Rules: &envoycfgauthz.RBAC{
						Action:   envoycfgauthz.RBAC_DENY,
						Policies: map[string]*envoycfgauthz.Policy{},
					},
				},
			},
			expectedCELRules: map[string][]string{},
			wantErr:          false,
		},
		{
			name:   "multiple rules with different JWT claims",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Policies: []v1alpha1.RbacPolicy{
					{
						Name: "policy-0",
						Principals: []v1alpha1.Principal{
							{
								JWTPrincipals: map[string]*v1alpha1.JWTPrincipal{
									"test-provider": {
										Claims: []v1alpha1.JWTClaimMatch{
											{
												Name:    "role",
												Value:   "admin",
												Matcher: v1alpha1.ClaimMatcherExactString,
											},
										},
									},
								},
							},
						},
						Conditions: &v1alpha1.CELConditions{
							CelMatchExpression: []string{"request.auth.claims.groups == 'group1'"},
						},
					},
					{
						Name: "policy-1",
						Principals: []v1alpha1.Principal{
							{
								JWTPrincipals: map[string]*v1alpha1.JWTPrincipal{
									"test-provider": {
										Claims: []v1alpha1.JWTClaimMatch{
											{
												Name:    "role",
												Value:   "user",
												Matcher: v1alpha1.ClaimMatcherExactString,
											},
										},
									},
								},
							},
						},
						Conditions: &v1alpha1.CELConditions{
							CelMatchExpression: []string{"request.auth.claims.groups == 'group2'"},
						},
					},
				},
			},
			jwt: &v1alpha1.JWTValidation{
				ExtensionRef: corev1.LocalObjectReference{
					Name: "test-provider",
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Rules: &envoycfgauthz.RBAC{
						Action: envoycfgauthz.RBAC_ALLOW,
						Policies: map[string]*envoycfgauthz.Policy{
							"ns[test-ns]-policy[test-policy]-rule[0]": {
								Principals: []*envoycfgauthz.Principal{
									{
										Identifier: &envoycfgauthz.Principal_SourcedMetadata{
											SourcedMetadata: &envoycfgauthz.SourcedMetadata{
												MetadataMatcher: &envoy_matcher_v3.MetadataMatcher{
													Filter: "envoy.filters.http.jwt_authn",
													Path: []*envoy_matcher_v3.MetadataMatcher_PathSegment{
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: PayloadInMetadata,
															},
														},
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: "role",
															},
														},
													},
													Value: &envoy_matcher_v3.ValueMatcher{
														MatchPattern: &envoy_matcher_v3.ValueMatcher_StringMatch{
															StringMatch: &envoy_matcher_v3.StringMatcher{
																MatchPattern: &envoy_matcher_v3.StringMatcher_Exact{
																	Exact: "admin",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Permissions: []*envoycfgauthz.Permission{
									{
										Rule: &envoycfgauthz.Permission_Matcher{
											Matcher: &envoycorev3.TypedExtensionConfig{
												Name: "cel-matcher",
											},
										},
									},
								},
							},
							"ns[test-ns]-policy[test-policy]-rule[1]": {
								Principals: []*envoycfgauthz.Principal{
									{
										Identifier: &envoycfgauthz.Principal_SourcedMetadata{
											SourcedMetadata: &envoycfgauthz.SourcedMetadata{
												MetadataMatcher: &envoy_matcher_v3.MetadataMatcher{
													Filter: "envoy.filters.http.jwt_authn",
													Path: []*envoy_matcher_v3.MetadataMatcher_PathSegment{
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: PayloadInMetadata,
															},
														},
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: "role",
															},
														},
													},
													Value: &envoy_matcher_v3.ValueMatcher{
														MatchPattern: &envoy_matcher_v3.ValueMatcher_StringMatch{
															StringMatch: &envoy_matcher_v3.StringMatcher{
																MatchPattern: &envoy_matcher_v3.StringMatcher_Exact{
																	Exact: "admin",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Permissions: []*envoycfgauthz.Permission{
									{
										Rule: &envoycfgauthz.Permission_Matcher{
											Matcher: &envoycorev3.TypedExtensionConfig{
												Name: "cel-matcher",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedCELRules: map[string][]string{
				"ns[test-ns]-policy[test-policy]-rule[0]": {"request.auth.claims.groups == 'group1'"},
				"ns[test-ns]-policy[test-policy]-rule[1]": {"request.auth.claims.groups == 'group2'"},
			},
			wantErr: false,
		},
		{
			name:   "nested JWT claims",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Policies: []v1alpha1.RbacPolicy{
					{
						Name: "policy-0",
						Principals: []v1alpha1.Principal{
							{
								JWTPrincipals: map[string]*v1alpha1.JWTPrincipal{
									"test-provider": {
										Claims: []v1alpha1.JWTClaimMatch{
											{
												Name:    "email",
												Value:   "dev2@kgateway.io", // route requires dev2
												Matcher: v1alpha1.ClaimMatcherContains,
											},
										},
									},
								},
							},
						},
						Conditions: &v1alpha1.CELConditions{
							CelMatchExpression: []string{"request.auth.claims.groups == 'group1'"},
						},
					},
				},
			},
			jwt: &v1alpha1.JWTValidation{
				ExtensionRef: corev1.LocalObjectReference{
					Name: "test-provider",
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Rules: &envoycfgauthz.RBAC{
						Action: envoycfgauthz.RBAC_ALLOW,
						Policies: map[string]*envoycfgauthz.Policy{
							"ns[test-ns]-policy[test-policy]-rule[0]": {
								Principals: []*envoycfgauthz.Principal{
									{
										Identifier: &envoycfgauthz.Principal_SourcedMetadata{
											SourcedMetadata: &envoycfgauthz.SourcedMetadata{
												MetadataMatcher: &envoy_matcher_v3.MetadataMatcher{
													Filter: "envoy.filters.http.jwt_authn",
													Path: []*envoy_matcher_v3.MetadataMatcher_PathSegment{
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: PayloadInMetadata,
															},
														},
														{
															Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
																Key: "email",
															},
														},
													},
													Value: &envoy_matcher_v3.ValueMatcher{
														MatchPattern: &envoy_matcher_v3.ValueMatcher_StringMatch{
															StringMatch: &envoy_matcher_v3.StringMatcher{
																MatchPattern: &envoy_matcher_v3.StringMatcher_Contains{
																	Contains: "dev2@kgateway.io",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Permissions: []*envoycfgauthz.Permission{
									{
										Rule: &envoycfgauthz.Permission_Matcher{
											Matcher: &envoycorev3.TypedExtensionConfig{
												Name: "cel-matcher",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedCELRules: map[string][]string{
				"ns[test-ns]-policy[test-policy]-rule[0]": {"request.auth.claims.groups == 'group1'"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			krtctx := krt.TestingDummyContext{}
			got, err := translateRbac(krtctx, tt.ns, tt.tpName, tt.rbac, tt.jwt, fetchGatewayExtension)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Basic structure validation
			assert.Equal(t, tt.expected.Rbac.Rules.Action, got.Rbac.Rules.Action)
			assert.Equal(t, len(tt.expected.Rbac.Rules.Policies), len(got.Rbac.Rules.Policies))

			// For each policy, validate the structure but not the exact TypedConfig content
			for policyName, expectedPolicy := range tt.expected.Rbac.Rules.Policies {
				gotPolicy, exists := got.Rbac.Rules.Policies[policyName]
				require.True(t, exists, "Policy %s should exist", policyName)

				// Validate principals
				assert.Equal(t, len(expectedPolicy.Principals), len(gotPolicy.Principals))

				// Validate permissions structure (but not exact content for Matcher type)
				assert.Equal(t, len(expectedPolicy.Permissions), len(gotPolicy.Permissions))
				for i, expectedPerm := range expectedPolicy.Permissions {
					gotPerm := gotPolicy.Permissions[i]
					// Check that we have the right permission rule type
					switch expectedPerm.Rule.(type) {
					case *envoycfgauthz.Permission_AndRules:
						_, ok := gotPerm.Rule.(*envoycfgauthz.Permission_AndRules)
						assert.True(t, ok, "Expected AndRules permission")
					case *envoycfgauthz.Permission_Matcher:
						_, ok := gotPerm.Rule.(*envoycfgauthz.Permission_Matcher)
						assert.True(t, ok, "Expected Matcher permission")
					case *envoycfgauthz.Permission_Any:
						_, ok := gotPerm.Rule.(*envoycfgauthz.Permission_Any)
						assert.True(t, ok, "Expected Any permission")
					}
				}

				// Validate CEL expressions if this policy should have them
				if expectedCELs, ok := tt.expectedCELRules[policyName]; ok {
					require.Equal(t, 1, len(gotPolicy.Permissions), "Expected exactly one permission for CEL validation")
					actualCELs := extractCELExpressions(t, gotPolicy.Permissions[0])
					assert.ElementsMatch(t, expectedCELs, actualCELs, "CEL expressions should match")
				}
			}
		})
	}
}
