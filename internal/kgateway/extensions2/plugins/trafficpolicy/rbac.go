package trafficpolicy

import (
	"fmt"
	"strconv"
	"strings"

	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoy_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

const (
	nestedClaimsDelimiter = "."
)

// RbacIr is the internal representation of an RBAC policy.
type RbacIr struct {
	rbacConfig *envoyauthz.RBACPerRoute
}

func (r *RbacIr) Equals(other *RbacIr) bool {
	if r == nil && other == nil {
		return true
	}
	if r == nil || other == nil {
		return false
	}
	return proto.Equal(r.rbacConfig, other.rbacConfig)
}

// handleRbac configures the RBAC filter and per-route RBAC configuration for a specific route
func (p *trafficPolicyPluginGwPass) handleRbac(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, rbacIr *RbacIr) {
	if rbacIr == nil || rbacIr.rbacConfig == nil {
		return
	}

	// Add a filter to the chain. When having a jwt policy for a route we need to also have a
	// global rbac http filter in the chain otherwise it will be ignored.
	if p.rbacInChain == nil {
		p.rbacInChain = make(map[string]*envoyauthz.RBAC)
	}
	if _, ok := p.rbacInChain[fcn]; !ok {
		p.rbacInChain[fcn] = &envoyauthz.RBAC{}
	}

	// Add the per-route RBAC configuration to the typed filter config
	pCtxTypedFilterConfig.AddTypedConfig(rbacFilterNamePrefix, rbacIr.rbacConfig)
}

// rbacForPolicy translates the RBAC spec into an envoy RBAC policy and stores it in the traffic policy IR
func (b *TrafficPolicyBuilder) rbacForPolicy(krtctx krt.HandlerContext, policy v1alpha1.TrafficPolicy, out *trafficPolicySpecIr) error {
	spec := policy.Spec
	if spec.RBAC == nil {
		return nil
	}

	rbacConfig, err := translateRbac(policy.Namespace, policy.Name, spec.RBAC)
	if err != nil {
		return err
	}

	out.rbac = &RbacIr{
		rbacConfig: rbacConfig,
	}
	return nil
}

func translateRbac(tpNs, tpName string, rbac *v1alpha1.Rbac) (*envoyauthz.RBACPerRoute, error) {
	policies := make(map[string]*envoycfgauthz.Policy)

	action := envoycfgauthz.RBAC_ALLOW
	if rbac.Action == v1alpha1.AuthorizationPolicyActionDeny {
		action = envoycfgauthz.RBAC_DENY
	}

	res := &envoyauthz.RBACPerRoute{
		Rbac: &envoyauthz.RBAC{
			Rules: &envoycfgauthz.RBAC{
				Action:   action,
				Policies: policies,
			},
		},
	}
	for ruleIdx, rule := range rbac.Rules {
		var err error
		policyName := policyName(tpNs, tpName, ruleIdx)
		policies[policyName], err = translatePolicy(rule)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func translatePolicy(rule v1alpha1.RbacRule) (*envoycfgauthz.Policy, error) {
	outPolicy := &envoycfgauthz.Policy{}

	if rule.Principal.JWTPrincipals != nil {
		jwtPrincipals := rule.Principal.JWTPrincipals
		for _, jwtPrincipal := range jwtPrincipals {
			outPrincipal, err := translateJwtPrincipal(jwtPrincipal)
			if err != nil {
				return nil, err
			}
			if outPrincipal != nil {
				outPolicy.Principals = append(outPolicy.GetPrincipals(), outPrincipal)
			}
		}
	}

	var allPermissions []*envoycfgauthz.Permission
	permission := rule.Access
	if permission.PathPrefixes != nil {
		for _, path := range permission.PathPrefixes {
			allPermissions = append(allPermissions, &envoycfgauthz.Permission{
				Rule: &envoycfgauthz.Permission_Header{
					Header: &envoyroute.HeaderMatcher{
						Name: ":path",
						HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
							StringMatch: &envoy_type_matcher_v3.StringMatcher{
								MatchPattern: &envoy_type_matcher_v3.StringMatcher_Prefix{
									Prefix: path,
								},
							},
						},
					},
				},
			})
		}
	}

	if len(permission.Methods) != 0 {
		allPermissions = append(allPermissions, translateMethods(permission.Methods))
	}

	if len(allPermissions) == 0 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_Any{
				Any: true,
			},
		}}
	} else if len(allPermissions) == 1 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{allPermissions[0]}
	} else {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_AndRules{
				AndRules: &envoycfgauthz.Permission_Set{
					Rules: allPermissions,
				},
			},
		}}
	}

	return outPolicy, nil
}

func translateJwtPrincipal(jwtPrincipal v1alpha1.JWTPrincipal) (*envoycfgauthz.Principal, error) {
	var jwtPrincipals []*envoycfgauthz.Principal
	claims := jwtPrincipal.Claims
	// sort for idempotency
	for _, claim := range claims {
		valueMatcher, err := GetValueMatcher(claim.Value, claim.Matcher)
		if err != nil {
			return nil, err
		}
		claimPrincipal := &envoycfgauthz.Principal{
			Identifier: &envoycfgauthz.Principal_SourcedMetadata{
				SourcedMetadata: &envoycfgauthz.SourcedMetadata{
					MetadataMatcher: &envoy_matcher_v3.MetadataMatcher{
						Filter: "envoy.filters.http.jwt_authn",
						Path:   getPath(claim, jwtPrincipal),
						Value:  valueMatcher,
					},
				},
			},
		}
		jwtPrincipals = append(jwtPrincipals, claimPrincipal)
	}

	if len(jwtPrincipals) == 0 {
		logger.Info("RBAC JWT Principal with zero claims - ignoring")
		return nil, nil
	} else if len(jwtPrincipals) == 1 {
		return jwtPrincipals[0], nil
	}
	return &envoycfgauthz.Principal{
		Identifier: &envoycfgauthz.Principal_AndIds{
			AndIds: &envoycfgauthz.Principal_Set{
				Ids: jwtPrincipals,
			},
		},
	}, nil
}

// For rbac config the PayloadInMetadata must match the value set in the jwt filter.
func getPath(claim v1alpha1.JWTClaimMatch, jwtPrincipal v1alpha1.JWTPrincipal) []*envoy_matcher_v3.MetadataMatcher_PathSegment {
	// If the claim name contains the nestedClaimsDelimiter then it's a nested claim, and the path
	// should contain a segment for each layer of nesting, for example:
	// {
	//   "sub": "1234567890",
	//   "name": "John Doe",
	//   "iat": 1516239022,
	//   "metadata": {
	//     "role": [
	//       "user",
	//       "editor",
	//       "admin"
	//     ]
	//   }
	// }
	// The nested claim name "role" would get a [metadata] segment and a [role] segment.
	// The claim name "sub" would only have a single [sub] segment.
	if strings.Contains(claim.Name, nestedClaimsDelimiter) {
		substrings := strings.Split(claim.Name, nestedClaimsDelimiter)
		path := make([]*envoy_matcher_v3.MetadataMatcher_PathSegment, len(substrings)+1)
		path[0] = &envoy_matcher_v3.MetadataMatcher_PathSegment{
			Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
				Key: jwtPrincipal.Provider, // this needs to match jwt provider name
			},
		}
		for i, substring := range substrings {
			path[i+1] = &envoy_matcher_v3.MetadataMatcher_PathSegment{
				Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
					Key: substring,
				},
			}
		}
		return path
	} else {
		return []*envoy_matcher_v3.MetadataMatcher_PathSegment{
			{
				Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
					Key: jwtPrincipal.Provider, // this needs to match jwt provider name
				},
			},
			{
				Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
					Key: claim.Name,
				},
			},
		}
	}
}

func GetValueMatcher(value string, claimMatcher v1alpha1.ClaimMatcher) (*envoy_matcher_v3.ValueMatcher, error) {
	switch claimMatcher {
	case v1alpha1.ClaimMatcherExactString:
		return getExactStringValueMatcher(value), nil
	case v1alpha1.ClaimMatcherBoolean:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("value cannot be parsed to a bool to use ClaimMatcher.BOOLEAN: %v", value)
		}
		return &envoy_matcher_v3.ValueMatcher{
			MatchPattern: &envoy_matcher_v3.ValueMatcher_BoolMatch{
				BoolMatch: boolValue,
			},
		}, nil
	case v1alpha1.ClaimMatcherContains:
		return &envoy_matcher_v3.ValueMatcher{
			MatchPattern: &envoy_matcher_v3.ValueMatcher_ListMatch{
				ListMatch: &envoy_matcher_v3.ListMatcher{
					MatchPattern: &envoy_matcher_v3.ListMatcher_OneOf{
						OneOf: getExactStringValueMatcher(value),
					},
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("no implementation defined for ClaimMatcher: %v", claimMatcher)
	}
}

func translateMethods(methods []string) *envoycfgauthz.Permission {
	var allPermissions []*envoycfgauthz.Permission
	for _, method := range methods {
		allPermissions = append(allPermissions, &envoycfgauthz.Permission{
			Rule: &envoycfgauthz.Permission_Header{
				Header: &envoyroute.HeaderMatcher{
					Name: ":method",
					HeaderMatchSpecifier: &envoyroute.HeaderMatcher_StringMatch{
						StringMatch: &envoy_type_matcher_v3.StringMatcher{
							MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
								Exact: method,
							},
						},
					},
				},
			},
		})
	}

	if len(allPermissions) == 1 {
		return allPermissions[0]
	}

	return &envoycfgauthz.Permission{
		Rule: &envoycfgauthz.Permission_OrRules{
			OrRules: &envoycfgauthz.Permission_Set{
				Rules: allPermissions,
			},
		},
	}
}

func getExactStringValueMatcher(value string) *envoy_matcher_v3.ValueMatcher {
	return &envoy_matcher_v3.ValueMatcher{
		MatchPattern: &envoy_matcher_v3.ValueMatcher_StringMatch{
			StringMatch: &envoy_matcher_v3.StringMatcher{
				MatchPattern: &envoy_matcher_v3.StringMatcher_Exact{
					Exact: value,
				},
			},
		},
	}
}

func policyName(namespace, name string, rule int) string {
	return fmt.Sprintf("ns[%s]-policy[%s]-rule[%d]", namespace, name, rule)
}
