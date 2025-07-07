package agentgatewaysyncer

import (
	"iter"
	"strings"

	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/protomarshal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
)

// TODO: support other route collections (TCP, TLS, etc.)
func ADPRouteCollection(
	httpRoutes krt.Collection[*gwv1.HTTPRoute],
	inputs RouteContextInputs,
	krtopts krtutil.KrtOptions,
) krt.Collection[ADPResource] {
	routes := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gwv1.HTTPRoute) []ADPResource {
		rm := reports.NewReportMap()
		rep := reports.NewReporter(&rm)
		logger.Debug("translating HTTPRoute", "route_name", obj.GetName(), "resource_version", obj.GetResourceVersion())

		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gwv1.HTTPRoute) iter.Seq2[ADPRoute, *ConfigError] {
			return func(yield func(ADPRoute, *ConfigError) bool) {
				for n, r := range route.Rules {
					// split the rule to make sure each rule has up to one match
					matches := slices.Reference(r.Matches)
					if len(matches) == 0 {
						matches = append(matches, nil)
					}
					for idx, m := range matches {
						if m != nil {
							r.Matches = []gwv1.HTTPRouteMatch{*m}
						}
						res, err := convertHTTPRouteToADP(ctx, r, obj, n, idx)
						if !yield(ADPRoute{Route: res}, err) {
							return
						}
					}
				}
			}
		})

		var res []ADPResource
		for _, parent := range filteredReferences(parentRefs) {
			// for gwv1beta1 routes, build one VS per gwv1beta1+host
			routes := gwResult.routes
			if len(routes) == 0 {
				continue
			}
			if gwResult.error != nil {
				rep.Route(obj).ParentRef(&parent.OriginalReference).SetCondition(reporter.RouteCondition{
					Type:    gwv1beta1.RouteConditionResolvedRefs, // TODO: check type
					Status:  metav1.ConditionFalse,
					Reason:  gwv1beta1.RouteConditionReason(gwResult.error.Reason),
					Message: gwResult.error.Message,
				})
			}

			gw := types.NamespacedName{
				Namespace: parent.ParentKey.Namespace,
				Name:      parent.ParentKey.Name,
			}
			res = append(res, slices.Map(routes, func(e ADPRoute) ADPResource {
				inner := protomarshal.Clone(e.Route)
				_, name, _ := strings.Cut(parent.InternalName, "/")
				inner.ListenerKey = name
				inner.Key = inner.GetKey() + "." + string(parent.ParentSection)
				return toResource(gw, ADPRoute{Route: inner})
			})...)
		}
		return res
	}, krtopts.ToOptions("ADPRoutes")...)

	return routes
}

type conversionResult[O any] struct {
	error  *ConfigError
	routes []O
}

// IsNil works around comparing generic types
func IsNil[O comparable](o O) bool {
	var t O
	return o == t
}

// computeRoute holds the common route building logic shared amongst all types
func computeRoute[T controllers.Object, O comparable](ctx RouteContext, obj T, translator func(
	obj T,
) iter.Seq2[O, *ConfigError],
) ([]routeParentReference, conversionResult[O]) {
	parentRefs := extractParentReferenceInfo(ctx, ctx.RouteParents, obj)

	convertRules := func() conversionResult[O] {
		res := conversionResult[O]{}
		for vs, err := range translator(obj) {
			// This was a hard error
			if IsNil(vs) {
				res.error = err
				return conversionResult[O]{error: err}
			}
			// Got an error but also routes
			if err != nil {
				res.error = err
			}
			res.routes = append(res.routes, vs)
		}
		return res
	}
	gwResult := buildGatewayRoutes(parentRefs, convertRules)

	return parentRefs, gwResult
}

// RouteContext defines a common set of inputs to a route collection. This should be built once per route translation and
// not shared outside of that.
// The embedded RouteContextInputs is typically based into a collection, then translated to a RouteContext with RouteContextInputs.WithCtx().
type RouteContext struct {
	Krt krt.HandlerContext
	RouteContextInputs
}

type RouteContextInputs struct {
	Grants         ReferenceGrants
	RouteParents   RouteParents
	DomainSuffix   string
	Services       krt.Collection[*corev1.Service]
	InferencePools krt.Collection[*inf.InferencePool]
	Namespaces     krt.Collection[*corev1.Namespace]
	ServiceEntries krt.Collection[*networkingclient.ServiceEntry]
}

func (i RouteContextInputs) WithCtx(krtctx krt.HandlerContext) RouteContext {
	return RouteContext{
		Krt:                krtctx,
		RouteContextInputs: i,
	}
}

type RouteWithKey struct {
	*Config
	Key string
}

func (r RouteWithKey) ResourceName() string {
	return config.NamespacedName(r.Config).String()
}

func (r RouteWithKey) Equals(o RouteWithKey) bool {
	return r.Config.Equals(o.Config)
}

// buildGatewayRoutes contains common logic to build a set of routes with gwv1beta1 semantics
func buildGatewayRoutes[T any](parentRefs []routeParentReference, convertRules func() T) T {
	return convertRules()
}
