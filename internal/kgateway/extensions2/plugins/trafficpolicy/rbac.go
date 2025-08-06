package trafficpolicy

import (
	"fmt"
	"strconv"
	"strings"

	cncfcorev3 "github.com/cncf/xds/go/xds/core/v3"
	cncfmatcherv3 "github.com/cncf/xds/go/xds/type/matcher/v3"
	cncftypev3 "github.com/cncf/xds/go/xds/type/v3"
	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoy_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/util/errors"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
)

const (
	nestedClaimsDelimiter = "."
)

// rbacIr is the internal representation of an RBAC policy.
type rbacIr struct {
	rbacConfig *envoyauthz.RBACPerRoute
}

func (r *rbacIr) Equals(other *rbacIr) bool {
	if r == nil && other == nil {
		return true
	}
	if r == nil || other == nil {
		return false
	}
	return proto.Equal(r.rbacConfig, other.rbacConfig)
}

// handleRbac configures the RBAC filter and per-route RBAC configuration for a specific route
func (p *trafficPolicyPluginGwPass) handleRbac(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, rbacIr *rbacIr) {
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

// constructRbac translates the RBAC spec into an envoy RBAC policy and stores it in the traffic policy IR
func constructRbac(krtctx krt.HandlerContext, policy *v1alpha1.TrafficPolicy, out *trafficPolicySpecIr, fetchGatewayExtension FetchGatewayExtensionFunc) error {
	spec := policy.Spec
	if spec.RBAC == nil {
		return nil
	}

	rbacConfig, err := translateRbac(krtctx, policy.Namespace, policy.Name, spec.RBAC, spec.JWT, fetchGatewayExtension)
	if err != nil {
		return err
	}

	out.rbac = &rbacIr{
		rbacConfig: rbacConfig,
	}
	return nil
}

func translateRbac(krtctx krt.HandlerContext, tpNs, tpName string, rbac *v1alpha1.Rbac, jwtAuthn *v1alpha1.JWTValidation, fetchGatewayExtension FetchGatewayExtensionFunc) (*envoyauthz.RBACPerRoute, error) {
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
	for ruleIdx, rule := range rbac.Policies {
		var err error
		pName := policyName(tpNs, tpName, ruleIdx)
		policies[pName], err = translatePolicy(krtctx, tpNs, rule, jwtAuthn, fetchGatewayExtension)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func translatePolicy(krtctx krt.HandlerContext, namespace string, rule v1alpha1.RbacPolicy, jwtAuthn *v1alpha1.JWTValidation, fetchGatewayExtension FetchGatewayExtensionFunc) (*envoycfgauthz.Policy, error) {
	outPolicy := &envoycfgauthz.Policy{}
	var errs []error

	if rule.Principals != nil {
		for _, principal := range rule.Principals {
			// TODO: support other principal types (e.g. OIDC, etc.)
			if principal.JWTPrincipals == nil {
				// any principal
				outPolicy.Principals = append(outPolicy.GetPrincipals(), &envoycfgauthz.Principal{
					Identifier: &envoycfgauthz.Principal_Any{
						Any: true,
					},
				})
			} else {
				for jwtPrincipalName, jwtPrincipal := range principal.JWTPrincipals {
					if jwtPrincipal == nil {
						errs = append(errs, fmt.Errorf("jwt principal %s not found", jwtPrincipalName))
						continue
					}
					provider, err := fetchGatewayExtension(krtctx, &jwtAuthn.ExtensionRef, namespace)
					if err != nil {
						return nil, fmt.Errorf("jwt: %w", err)
					}
					// check if providerName on rbac is in the list of providers
					if _, ok := provider.JwtProviders[jwtPrincipalName]; !ok {
						errs = append(errs, fmt.Errorf("jwt provider %s not found", jwtPrincipalName))
						continue
					}
					outPrincipal, err := translateJwtPrincipal(provider.ResourceName(), jwtPrincipal)
					if err != nil {
						return nil, err
					}
					if outPrincipal != nil {
						outPolicy.Principals = append(outPolicy.GetPrincipals(), outPrincipal)
					}
				}
			}
		}
	}

	var allPermissions []*envoycfgauthz.Permission
	if rule.Conditions != nil {
		if rule.Conditions.CelMatchExpression != nil {
			celMatchInput, err := utils.MessageToAny(&cncfmatcherv3.HttpAttributesCelMatchInput{})
			if err != nil {
				errs = append(errs, err)
			}
			celMatchInputConfig := &cncfcorev3.TypedExtensionConfig{
				Name:        "envoy.matching.inputs.cel_data_input",
				TypedConfig: celMatchInput,
			}
			var matchers []*cncfmatcherv3.Matcher_MatcherList_FieldMatcher

			// TODO: handle single vs. list case separately
			for _, cel := range rule.Conditions.CelMatchExpression {
				typedCelMatchAny, marshalErr := utils.MessageToAny(&cncfmatcherv3.CelMatcher{
					ExprMatch: &cncftypev3.CelExpression{
						CelExprString: cel,
					},
				})
				if marshalErr != nil {
					errs = append(errs, err)
					continue
				}
				typedCelMatchConfig := &cncfcorev3.TypedExtensionConfig{
					// TODO: use unique user-defined name?
					Name:        "envoy.matching.matchers.cel_matcher",
					TypedConfig: typedCelMatchAny,
				}
				predicate := &cncfmatcherv3.Matcher_MatcherList_Predicate{
					MatchType: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate_{
						SinglePredicate: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate{
							Input: celMatchInputConfig,
							Matcher: &cncfmatcherv3.Matcher_MatcherList_Predicate_SinglePredicate_CustomMatch{
								CustomMatch: typedCelMatchConfig,
							},
						},
					},
				}
				matcher := &cncfmatcherv3.Matcher_MatcherList_FieldMatcher{
					Predicate: predicate,
					// TODO: on match action?
				}

				matchers = append(matchers, matcher)
			}

			celMatcher := &cncfmatcherv3.Matcher{
				MatcherType: &cncfmatcherv3.Matcher_MatcherList_{
					MatcherList: &cncfmatcherv3.Matcher_MatcherList{
						Matchers: matchers,
					},
				},
			}
			typedCelMatchAny, marshalErr := utils.MessageToAny(celMatcher)
			if marshalErr != nil {
				// failed to marshal cel matcher, return error (nothing to do)
				return nil, marshalErr
			}
			typedCelMatcheConfig := &envoycorev3.TypedExtensionConfig{
				Name:        "cel-matcher",
				TypedConfig: typedCelMatchAny,
			}
			allPermissions = append(allPermissions, &envoycfgauthz.Permission{
				Rule: &envoycfgauthz.Permission_Matcher{
					Matcher: typedCelMatcheConfig,
				},
			})
		}
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

	return outPolicy, errors.NewAggregate(errs)
}

func translateJwtPrincipal(jwtProviderName string, jwtPrincipal *v1alpha1.JWTPrincipal) (*envoycfgauthz.Principal, error) {
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
						Filter: jwtFilterName(jwtProviderName),
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
func getPath(claim v1alpha1.JWTClaimMatch, jwtPrincipal *v1alpha1.JWTPrincipal) []*envoy_matcher_v3.MetadataMatcher_PathSegment {
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
		var path []*envoy_matcher_v3.MetadataMatcher_PathSegment
		path = make([]*envoy_matcher_v3.MetadataMatcher_PathSegment, len(substrings)+1)
		path[0] = &envoy_matcher_v3.MetadataMatcher_PathSegment{
			Segment: &envoy_matcher_v3.MetadataMatcher_PathSegment_Key{
				Key: PayloadInMetadata, // this needs to match jwt payload in metadata
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
					Key: PayloadInMetadata, // this needs to match jwt payload in metadata
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
