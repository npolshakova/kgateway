package agentgatewaysyncer

import (
	"iter"
	"strings"

	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/protomarshal"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
)

// ADPRouteCollection creates the collection of translated routes
func ADPRouteCollection(
	httpRouteCol krt.Collection[*gwv1.HTTPRoute],
	grpcRouteCol krt.Collection[*gwv1.GRPCRoute],
	tcpRouteCol krt.Collection[*gwv1alpha2.TCPRoute],
	tlsRouteCol krt.Collection[*gwv1alpha2.TLSRoute],
	gateways krt.Collection[Gateway],
	gatewayObjs krt.Collection[*gwv1.Gateway],
	inputs RouteContextInputs,
	krtopts krtutil.KrtOptions,
	rm reports.ReportMap,
	rep reporter.Reporter,
) krt.Collection[ADPResource] {
	httpRouteAttachments := gatewayRouteAttachmentCountCollection(inputs, httpRouteCol, wellknown.HTTPRouteGVK, krtopts)
	httpRoutes := krt.NewManyCollection(httpRouteCol, func(krtctx krt.HandlerContext, obj *gwv1.HTTPRoute) []ADPResource {
		logger.Debug("translating HTTPRoute", "route_name", obj.GetName(), "resource_version", obj.GetResourceVersion())

		ctx := inputs.WithCtx(krtctx)
		routeReporter := rep.Route(obj)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gwv1.HTTPRoute) iter.Seq2[ADPRoute, *reporter.RouteCondition] {
			return func(yield func(ADPRoute, *reporter.RouteCondition) bool) {
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
			// Always create a route reporter entry for the parent ref
			parentRefReporter := routeReporter.ParentRef(&parent.OriginalReference)

			// for gwv1beta1 routes, build one VS per gwv1beta1+host
			routes := gwResult.routes
			if len(routes) == 0 {
				logger.Debug("no routes for parent", "route_name", obj.GetName(), "parent", parent.ParentKey)
				continue
			}
			if gwResult.error != nil {
				parentRefReporter.SetCondition(*gwResult.error)
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
				return toResourceWithReports(gw, ADPRoute{Route: inner}, rm)
			})...)
		}
		return res
	}, krtopts.ToOptions("ADPHTTPRoutes")...)

	grpcRouteAttachments := gatewayRouteAttachmentCountCollection(inputs, grpcRouteCol, wellknown.GRPCRouteGVK, krtopts)
	grpcRoutes := krt.NewManyCollection(grpcRouteCol, func(krtctx krt.HandlerContext, obj *gwv1.GRPCRoute) []ADPResource {
		logger.Debug("translating GRPCRoute", "route_name", obj.GetName(), "resource_version", obj.GetResourceVersion())

		ctx := inputs.WithCtx(krtctx)
		routeReporter := rep.Route(obj)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gwv1.GRPCRoute) iter.Seq2[ADPRoute, *reporter.RouteCondition] {
			return func(yield func(ADPRoute, *reporter.RouteCondition) bool) {
				for n, r := range route.Rules {
					// Convert the entire rule with all matches at once
					res, err := convertGRPCRouteToADP(ctx, r, obj, n)
					if !yield(ADPRoute{Route: res}, err) {
						return
					}
				}
			}
		})

		var res []ADPResource
		for _, parent := range filteredReferences(parentRefs) {
			// Always create a route reporter entry for the parent ref
			parentRefReporter := routeReporter.ParentRef(&parent.OriginalReference)

			// for gwv1beta1 routes, build one VS per gwv1beta1+host
			routes := gwResult.routes
			if len(routes) == 0 {
				logger.Debug("no routes for parent", "route_name", obj.GetName(), "parent", parent.ParentKey)
				continue
			}
			if gwResult.error != nil {
				parentRefReporter.SetCondition(*gwResult.error)
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
				return toResourceWithReports(gw, ADPRoute{Route: inner}, rm)
			})...)
		}
		return res
	}, krtopts.ToOptions("ADPGRPCRoutes")...)

	tcpRouteAttachments := gatewayRouteAttachmentCountCollection(inputs, tcpRouteCol, wellknown.TCPRouteGVK, krtopts)
	tcpRoutes := krt.NewManyCollection(tcpRouteCol, func(krtctx krt.HandlerContext, obj *gwv1alpha2.TCPRoute) []ADPResource {
		logger.Debug("translating TCPRoute", "route_name", obj.GetName(), "resource_version", obj.GetResourceVersion())

		ctx := inputs.WithCtx(krtctx)
		routeReporter := rep.Route(obj)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gwv1alpha2.TCPRoute) iter.Seq2[ADPRoute, *reporter.RouteCondition] {
			return func(yield func(ADPRoute, *reporter.RouteCondition) bool) {
				for n, r := range route.Rules {
					// Convert the entire rule with all matches at once
					res, err := convertTCPRouteToADP(ctx, r, obj, n)
					if !yield(ADPRoute{Route: res}, err) {
						return
					}
				}
			}
		})

		var res []ADPResource
		for _, parent := range filteredReferences(parentRefs) {
			// Always create a route reporter entry for the parent ref
			parentRefReporter := routeReporter.ParentRef(&parent.OriginalReference)

			// for gwv1beta1 routes, build one VS per gwv1beta1+host
			routes := gwResult.routes
			if len(routes) == 0 {
				logger.Debug("no routes for parent", "route_name", obj.GetName(), "parent", parent.ParentKey)
				continue
			}
			if gwResult.error != nil {
				parentRefReporter.SetCondition(*gwResult.error)
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
				return toResourceWithReports(gw, ADPRoute{Route: inner}, rm)
			})...)
		}
		return res
	}, krtopts.ToOptions("ADPTCPRoutes")...)

	tlsRouteAttachments := gatewayRouteAttachmentCountCollection(inputs, tlsRouteCol, wellknown.TLSRouteGVK, krtopts)
	tlsRoutes := krt.NewManyCollection(tlsRouteCol, func(krtctx krt.HandlerContext, obj *gwv1alpha2.TLSRoute) []ADPResource {
		logger.Debug("translating TLSRoute", "route_name", obj.GetName(), "resource_version", obj.GetResourceVersion())

		ctx := inputs.WithCtx(krtctx)
		routeReporter := rep.Route(obj)
		route := obj.Spec
		parentRefs, gwResult := computeRoute(ctx, obj, func(obj *gwv1alpha2.TLSRoute) iter.Seq2[ADPRoute, *reporter.RouteCondition] {
			return func(yield func(ADPRoute, *reporter.RouteCondition) bool) {
				for n, r := range route.Rules {
					// Convert the entire rule with all matches at once
					res, err := convertTLSRouteToADP(ctx, r, obj, n)
					if !yield(ADPRoute{Route: res}, err) {
						return
					}
				}
			}
		})

		var res []ADPResource
		for _, parent := range filteredReferences(parentRefs) {
			// Always create a route reporter entry for the parent ref
			parentRefReporter := routeReporter.ParentRef(&parent.OriginalReference)

			// for gwv1beta1 routes, build one VS per gwv1beta1+host
			routes := gwResult.routes
			if len(routes) == 0 {
				logger.Debug("no routes for parent", "route_name", obj.GetName(), "parent", parent.ParentKey)
				continue
			}
			if gwResult.error != nil {
				parentRefReporter.SetCondition(*gwResult.error)
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
				return toResourceWithReports(gw, ADPRoute{Route: inner}, rm)
			})...)
		}
		return res
	}, krtopts.ToOptions("ADPTLSRoutes")...)

	routes := krt.JoinCollection([]krt.Collection[ADPResource]{httpRoutes, grpcRoutes, tcpRoutes, tlsRoutes}, krtopts.ToOptions("ADPRoutes")...)

	routeAttachments := krt.JoinCollection([]krt.Collection[*RouteAttachment]{
		httpRouteAttachments,
		grpcRouteAttachments,
		tcpRouteAttachments,
		tlsRouteAttachments,
	}, krtopts.ToOptions("RouteAttachments")...)
	routeAttachmentsIndex := krt.NewIndex(routeAttachments, func(o *RouteAttachment) []types.NamespacedName {
		return []types.NamespacedName{o.To}
	})
	FinalGatewayStatusCollectionAttachedRoutes(gatewayObjs, routeAttachments, routeAttachmentsIndex, krtopts, rep)

	return routes
}

type conversionResult[O any] struct {
	error  *reporter.RouteCondition
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
) iter.Seq2[O, *reporter.RouteCondition],
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
	Backends       *krtcollections.BackendIndex
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

// gatewayRouteAttachmentCountCollection holds the generic logic to determine the parents a route is attached to, used for
// computing the aggregated `attachedRoutes` status in Gateway.
func gatewayRouteAttachmentCountCollection[T controllers.Object](
	inputs RouteContextInputs,
	col krt.Collection[T],
	kind schema.GroupVersionKind,
	krtopts krtutil.KrtOptions,
) krt.Collection[*RouteAttachment] {
	return krt.NewManyCollection(col, func(krtctx krt.HandlerContext, obj T) []*RouteAttachment {
		ctx := inputs.WithCtx(krtctx)
		from := TypedResource{
			Kind: kind,
			Name: config.NamespacedName(obj),
		}

		parentRefs := extractParentReferenceInfo(ctx, inputs.RouteParents, obj)
		return slices.MapFilter(filteredReferences(parentRefs), func(e routeParentReference) **RouteAttachment {
			if e.ParentKey.Kind != wellknown.GatewayGVK {
				return nil
			}
			return ptr.Of(&RouteAttachment{
				From: from,
				To: types.NamespacedName{
					Name:      e.ParentKey.Name,
					Namespace: e.ParentKey.Namespace,
				},
				ListenerName: string(e.ParentSection),
			})
		})
	}, krtopts.ToOptions(kind.Kind+"/count")...)
}

// gatewayStatusUpdate is a simple wrapper type for Gateway status updates that implements ResourceNamer
type gatewayStatusUpdate struct {
	gateway *gwv1.Gateway
}

// ResourceName implements krt.ResourceNamer interface
func (g gatewayStatusUpdate) ResourceName() string {
	return g.gateway.Namespace + "/" + g.gateway.Name
}

// Equals implements krt.Equaler interface
func (g gatewayStatusUpdate) Equals(other gatewayStatusUpdate) bool {
	return g.gateway.Namespace == other.gateway.Namespace &&
		g.gateway.Name == other.gateway.Name &&
		g.gateway.ResourceVersion == other.gateway.ResourceVersion
}

// FinalGatewayStatusCollectionAttachedRoutes finalizes a Gateway status. There is a circular logic between Gateways and Routes to determine
// the attachedRoute count, so we first build a partial Gateway status, then once routes are computed we finalize it with
// the attachedRoute count.
func FinalGatewayStatusCollectionAttachedRoutes(
	gateways krt.Collection[*gwv1.Gateway],
	routeAttachments krt.Collection[*RouteAttachment],
	routeAttachmentsIndex krt.Index[types.NamespacedName, *RouteAttachment],
	krtopts krtutil.KrtOptions,
	rep reports.Reporter,
) {
	_ = krt.NewCollection(
		gateways,
		func(ctx krt.HandlerContext, obj *gwv1.Gateway) *gatewayStatusUpdate {
			routeAttachmentsForGw := krt.Fetch(ctx, routeAttachments, krt.FilterIndex(routeAttachmentsIndex, types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}))
			counts := map[string]uint{}
			for _, r := range routeAttachmentsForGw {
				counts[r.ListenerName]++
			}
			for _, listener := range obj.Spec.Listeners {
				rep.Gateway(obj).Listener(&listener).SetAttachedRoutes(counts[string(listener.Name)])
			}
			// Return a wrapper instead of the raw Gateway object
			return &gatewayStatusUpdate{gateway: obj}
		}, krtopts.ToOptions("GatewayFinalStatus")...)
}
