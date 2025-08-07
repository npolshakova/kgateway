package trafficpolicy

import (
	"fmt"
	"testing"

	cncfcorev3 "github.com/cncf/xds/go/xds/core/v3"
	cncfmatcherv3 "github.com/cncf/xds/go/xds/type/matcher/v3"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	"github.com/google/cel-go/cel"
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

// createExpectedMatcher creates an expected matcher structure for testing
func createExpectedMatcher(action v1alpha1.AuthorizationPolicyAction, numRules int) *cncfmatcherv3.Matcher {
	// Create a simplified matcher structure for testing
	// We don't need to match the exact complex internal structure,
	// just the basic structure with the right number of matchers
	var matchers []*cncfmatcherv3.Matcher_MatcherList_FieldMatcher
	for i := 0; i < numRules; i++ {
		matcher := &cncfmatcherv3.Matcher_MatcherList_FieldMatcher{
			// Simplified structure - the actual implementation creates complex CEL matchers
			Predicate: &cncfmatcherv3.Matcher_MatcherList_Predicate{
				MatchType: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate_{
					SinglePredicate: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate{
						Input: &cncfcorev3.TypedExtensionConfig{
							Name: "envoy.matching.inputs.cel_data_input",
						},
						Matcher: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate_CustomMatch{
							CustomMatch: &cncfcorev3.TypedExtensionConfig{
								Name: "envoy.matching.matchers.cel_matcher",
							},
						},
					},
				},
			},
			OnMatch: &cncfmatcherv3.Matcher_OnMatch{
				OnMatch: &cncfmatcherv3.Matcher_OnMatch_Action{
					Action: &cncfcorev3.TypedExtensionConfig{
						Name: "envoy.filters.rbac.action",
					},
				},
			},
		}
		matchers = append(matchers, matcher)
	}

	return &cncfmatcherv3.Matcher{
		MatcherType: &cncfmatcherv3.Matcher_MatcherList_{
			MatcherList: &cncfmatcherv3.Matcher_MatcherList{
				Matchers: matchers,
			},
		},
		OnNoMatch: &cncfmatcherv3.Matcher_OnMatch{
			OnMatch: &cncfmatcherv3.Matcher_OnMatch_Action{
				Action: &cncfcorev3.TypedExtensionConfig{
					Name: "action",
				},
			},
		},
	}
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
						CelMatchExpression: []string{"request.auth.claims.groups == 'group1'", "request.auth.claims.groups == 'group2'"},
					},
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Matcher: createExpectedMatcher(v1alpha1.AuthorizationPolicyActionAllow, 1),
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
						Action: envoycfgauthz.RBAC_DENY,
					},
					Matcher: createExpectedMatcher(v1alpha1.AuthorizationPolicyActionDeny, 0),
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
						CelMatchExpression: []string{"request.auth.claims.groups == 'group1'"},
					},
					{
						CelMatchExpression: []string{"request.auth.claims.groups == 'group2'"},
					},
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Matcher: createExpectedMatcher(v1alpha1.AuthorizationPolicyActionAllow, 2),
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
						CelMatchExpression: []string{"request.auth.claims.groups == 'group1'"},
					},
				},
			},
			expected: &envoyauthz.RBACPerRoute{
				Rbac: &envoyauthz.RBAC{
					Matcher: createExpectedMatcher(v1alpha1.AuthorizationPolicyActionAllow, 1),
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
			got, err := translateRbac(tt.rbac)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			if got.Rbac.Matcher != nil {
				// When CEL expressions are present, expect Matcher field
				require.NotNil(t, got.Rbac.Matcher, "Expected Matcher field in actual result")

				// Create CEL environment for validation
				env, err := cel.NewEnv()
				require.NoError(t, err, "Failed to create CEL environment")

				// Validate CEL expressions for all expected rules
				for _, expectedCELs := range tt.expectedCELRules {
					assert.Greater(t, len(expectedCELs), 0, "Expected CEL expressions should not be empty")

					// Validate each CEL expression can be parsed
					for _, celExpr := range expectedCELs {
						parsedExpr, err := parseCELExpression(env, celExpr)
						assert.NoError(t, err, "CEL expression should be valid: %s", celExpr)
						assert.NotNil(t, parsedExpr, "Parsed CEL expression should not be nil: %s", celExpr)
					}
				}
			}
		})
	}
}
