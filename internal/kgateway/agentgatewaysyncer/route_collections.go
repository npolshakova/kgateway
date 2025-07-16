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
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
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

	// Create an index on routes by gateway for efficient lookup
	routesByGateway := krt.NewIndex(routes, func(adpResource ADPResource) []types.NamespacedName {
		return []types.NamespacedName{adpResource.Gateway}
	})

	// Set attached routes for each gateway after routes are built
	// Create a collection to trigger route counting when both gateways and routes are ready
	_ = krt.NewManyCollection(gateways, func(krtctx krt.HandlerContext, gateway Gateway) []struct{} {
		logger.Debug("setting attached routes for gateway", "gateway", gateway.ResourceName())

		// Find the corresponding gwv1.Gateway object
		gatewayObj := findGatewayObject(krtctx, gatewayObjs, gateway)
		if gatewayObj == nil {
			logger.Debug("could not find corresponding gwv1.Gateway object", "gateway", gateway.ResourceName())
			return []struct{}{}
		}
		setAttachedRoutes(&gateway, krtctx, routes, routesByGateway, gatewayObj, rep)

		// Return empty slice since this is just for side effects
		return []struct{}{}
	}, krtopts.ToOptions("AttachedRouteSetter")...)

	return routes
}

// Helper function to find the corresponding gwv1.Gateway object for a Gateway
func findGatewayObject(krtctx krt.HandlerContext, gatewayObjs krt.Collection[*gwv1.Gateway], gateway Gateway) *gwv1.Gateway {
	// Find the gwv1.Gateway that matches this Gateway's parent namespace and name
	allGateways := krt.Fetch(krtctx, gatewayObjs)
	for _, gw := range allGateways {
		if gw.Namespace == gateway.parent.Namespace && gw.Name == gateway.parent.Name {
			return gw
		}
	}
	return nil
}

// Helper function to extract listener name from ListenerKey
// ListenerKey format: gwName-AgentgatewayName-lName
func extractListenerNameFromKey(listenerKey string) string {
	// The listener name is the part after the last occurrence of AgentgatewayName-
	prefix := AgentgatewayName + "-"
	idx := strings.LastIndex(listenerKey, prefix)
	if idx == -1 {
		// Fallback: if the expected pattern is not found, return the original key
		return listenerKey
	}
	return listenerKey[idx+len(prefix):]
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

func setAttachedRoutes(gateway *Gateway, krtctx krt.HandlerContext, routes krt.Collection[ADPResource], routesByGateway krt.Index[types.NamespacedName, ADPResource], gatewayObj *gwv1.Gateway, reporter reports.Reporter) {
	// In agentgatewaysyncer, each Gateway represents a single listener
	// We need to count routes attached to this specific gateway/listener and set the count

	// Get the gateway reporter
	gwReporter := reporter.Gateway(gatewayObj)

	// Count routes that are attached to this specific listener
	// Routes have a ListenerKey that matches the listener name
	listenerName := string(gateway.parentInfo.SectionName)

	// Use the index to find all routes attached to this specific gateway
	gatewayNamespacedName := types.NamespacedName{
		Namespace: gateway.parent.Namespace,
		Name:      gateway.parent.Name,
	}
	routesForGateway := krt.Fetch(krtctx, routes, krt.FilterIndex(routesByGateway, gatewayNamespacedName))

	routeCount := 0
	for _, adpResource := range routesForGateway {
		if adpResource.Resource != nil {
			if routeRes := adpResource.Resource.GetRoute(); routeRes != nil {
				// Check if this route is attached to our listener
				logger.Debug("checking route", "route", routeRes.ListenerKey, "listener", listenerName)
				// Extract listener name from ListenerKey (format: gwName-AgentgatewayName-lName)
				// TODO: fix this
				extractedListenerName := extractListenerNameFromKey(routeRes.ListenerKey)
				if extractedListenerName == listenerName {
					// Also verify the gateway matches
					logger.Debug("checking gw for route ns", "adp", adpResource.Gateway.Namespace, "parent", gateway.parent.Namespace)
					logger.Debug("checking gw for route name", "adp", adpResource.Gateway.Name, "parent", gateway.parent.Name)
					if adpResource.Gateway.Namespace == gateway.parent.Namespace &&
						adpResource.Gateway.Name == gateway.parent.Name {
						routeCount++
					}
				}
			}
		}
	}

	// Find the corresponding listener in the Gateway object
	var targetListener *gwv1.Listener
	for _, listener := range gatewayObj.Spec.Listeners {
		logger.Debug("checking section name", "gw", gateway.parentInfo.SectionName, "listener", listener.Name)
		if listener.Name == gateway.parentInfo.SectionName {
			targetListener = &listener
			break
		}
	}

	if targetListener != nil {
		// Set the attached routes count for this listener
		logger.Debug("setting attached routes", "listener", targetListener.Name, "count", routeCount)
		gwReporter.Listener(targetListener).SetAttachedRoutes(uint(routeCount))
	}
}
