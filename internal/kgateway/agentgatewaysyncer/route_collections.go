package agentgatewaysyncer

import (
	"iter"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/protomarshal"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
)

// TODO: support other route collections (TCP, TLS, etc.)
func ADPRouteCollection(
	httpRoutes krt.Collection[*gateway.HTTPRoute],
	inputs RouteContextInputs,
	krtopts krtutil.KrtOptions,
) krt.Collection[ADPResource] {
	routes := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gateway.HTTPRoute) []ADPResource {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gateway.HTTPRoute) iter.Seq2[ADPRoute, *ConfigError] {
			return func(yield func(ADPRoute, *ConfigError) bool) {
				for n, r := range route.Rules {
					// split the rule to make sure each rule has up to one match
					matches := slices.Reference(r.Matches)
					if len(matches) == 0 {
						matches = append(matches, nil)
					}
					for idx, m := range matches {
						if m != nil {
							r.Matches = []gateway.HTTPRouteMatch{*m}
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
			// for gateway routes, build one VS per gateway+host
			routes := gwResult.routes
			if len(routes) == 0 {
				continue
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

// buildGatewayRoutes contains common logic to build a set of routes with gateway semantics
func buildGatewayRoutes[T any](parentRefs []routeParentReference, convertRules func() T) T {
	return convertRules()
}

// RouteResult holds the result of a route collection
type RouteResult[I, IStatus any] struct {
	// VirtualServices are the primary output that configures the internal routing logic
	VirtualServices krt.Collection[*Config]
	// RouteAttachments holds information about parent attachment to routes, used for computed the `attachedRoutes` count.
	RouteAttachments krt.Collection[*RouteAttachment]
}

type GatewayAndListener struct {
	// To is assumed to be a Gateway
	To           types.NamespacedName
	ListenerName string
}

func (g GatewayAndListener) String() string {
	return g.To.String() + "/" + g.ListenerName
}

type RouteAttachment struct {
	From TypedResource
	// To is assumed to be a Gateway
	To           types.NamespacedName
	ListenerName string
}

func (r *RouteAttachment) ResourceName() string {
	return r.From.Kind.String() + "/" + r.From.Name.String() + "/" + r.To.String() + "/" + r.ListenerName
}

func (r *RouteAttachment) Equals(other RouteAttachment) bool {
	return r.From == other.From && r.To == other.To && r.ListenerName == other.ListenerName
}
