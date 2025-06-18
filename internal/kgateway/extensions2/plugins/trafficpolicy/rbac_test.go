package trafficpolicy

import (
	"testing"

	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoy_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

func TestTranslateRbac(t *testing.T) {
	tests := []struct {
		name     string
		ns       string
		tpName   string
		rbac     *v1alpha1.Rbac
		expected *envoyauthz.RBACPerRoute
		wantErr  bool
	}{
		{
			name:   "allow action with single rule",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Rules: []v1alpha1.RbacRule{
					{
						Principal: v1alpha1.Principal{
							JWTPrincipals: []v1alpha1.JWTPrincipal{
								{
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
						Access: v1alpha1.AccessRule{
							PathPrefixes: []string{"/api"},
							Methods:      []string{"GET", "POST"},
						},
					},
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
										Rule: &envoycfgauthz.Permission_AndRules{
											AndRules: &envoycfgauthz.Permission_Set{
												Rules: []*envoycfgauthz.Permission{
													{
														Rule: &envoycfgauthz.Permission_Header{
															Header: &envoyroute.HeaderMatcher{
																Name: ":path",
																HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
																	StringMatch: &envoy_type_matcher_v3.StringMatcher{
																		MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
																			Prefix: "/api",
																		},
																	},
																},
															},
														},
													},
													{
														Rule: &envoycfgauthz.Permission_OrRules{
															OrRules: &envoycfgauthz.Permission_Set{
																Rules: []*envoycfgauthz.Permission{
																	{
																		Rule: &envoycfgauthz.Permission_Header{
																			Header: &envoyroute.HeaderMatcher{
																				Name: ":method",
																				HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
																					StringMatch: &envoy_type_matcher_v3.StringMatcher{
																						MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
																							Exact: "GET",
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Rule: &envoycfgauthz.Permission_Header{
																			Header: &envoyroute.HeaderMatcher{
																				Name: ":method",
																				HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
																					StringMatch: &envoy_type_matcher_v3.StringMatcher{
																						MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
																							Exact: "POST",
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "deny action with empty rules",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionDeny,
				Rules:  []v1alpha1.RbacRule{},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Rules: &envoycfgauthz.RBAC{
						Action:   envoycfgauthz.RBAC_DENY,
						Policies: map[string]*envoycfgauthz.Policy{},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "multiple rules with different JWT claims",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Rules: []v1alpha1.RbacRule{
					{
						Principal: v1alpha1.Principal{
							JWTPrincipals: []v1alpha1.JWTPrincipal{
								{
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
						Access: v1alpha1.AccessRule{
							PathPrefixes: []string{"/admin"},
						},
					},
					{
						Principal: v1alpha1.Principal{
							JWTPrincipals: []v1alpha1.JWTPrincipal{
								{
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
						Access: v1alpha1.AccessRule{
							PathPrefixes: []string{"/user"},
						},
					},
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
										Rule: &envoycfgauthz.Permission_Header{
											Header: &envoyroute.HeaderMatcher{
												Name: ":path",
												HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
													StringMatch: &envoy_type_matcher_v3.StringMatcher{
														MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
															Prefix: "/admin",
														},
													},
												},
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
																	Exact: "user",
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
										Rule: &envoycfgauthz.Permission_Header{
											Header: &envoyroute.HeaderMatcher{
												Name: ":path",
												HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
													StringMatch: &envoy_type_matcher_v3.StringMatcher{
														MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
															Prefix: "/user",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "nested JWT claims",
			ns:     "test-ns",
			tpName: "test-policy",
			rbac: &v1alpha1.Rbac{
				Action: v1alpha1.AuthorizationPolicyActionAllow,
				Rules: []v1alpha1.RbacRule{
					{
						Principal: v1alpha1.Principal{
							JWTPrincipals: []v1alpha1.JWTPrincipal{
								{
									Claims: []v1alpha1.JWTClaimMatch{
										{
											Name:    "metadata.role",
											Value:   "admin",
											Matcher: v1alpha1.ClaimMatcherExactString,
										},
									},
								},
							},
						},
						Access: v1alpha1.AccessRule{
							PathPrefixes: []string{"/admin"},
						},
					},
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
																Key: "metadata",
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
										Rule: &envoycfgauthz.Permission_Header{
											Header: &envoyroute.HeaderMatcher{
												Name: ":path",
												HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
													StringMatch: &envoy_type_matcher_v3.StringMatcher{
														MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
															Prefix: "/admin",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translateRbac(tt.ns, tt.tpName, tt.rbac)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			// Convert both to JSON for better diff output
			expectedJSON, err := protojson.Marshal(tt.expected)
			require.NoError(t, err)
			gotJSON, err := protojson.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedJSON), string(gotJSON))
		})
	}
}
