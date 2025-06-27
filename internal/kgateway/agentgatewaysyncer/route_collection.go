package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"
	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

var routeLogger = logging.New("agentgateway/route-collection")

// AgentRouteResource represents a translated route resource for agentgateway
type AgentRouteResource struct {
	types.NamespacedName
	Route *api.Route
	Valid bool
}

func (r AgentRouteResource) ResourceName() string {
	return fmt.Sprintf("%s/%s", r.Namespace, r.Name)
}

func (r AgentRouteResource) Equals(other AgentRouteResource) bool {
	if r.NamespacedName != other.NamespacedName || r.Valid != other.Valid {
		return false
	}
	if r.Route == nil && other.Route != nil || r.Route != nil && other.Route == nil {
		return false
	}
	if r.Route != nil && !proto.Equal(r.Route, other.Route) {
		return false
	}
	return true
}

// RouteContext defines the context for route processing
type RouteContext struct {
	Krt                  krt.HandlerContext
	AgentGatewayResource krt.Collection[AgentGatewayResource]
	Services             krt.Collection[ServiceInfo]
	Workloads            krt.Collection[WorkloadInfo]
	DomainSuffix         string
}

// RouteContextInputs defines the inputs needed for route processing
type RouteContextInputs struct {
	AgentGatewayResource krt.Collection[AgentGatewayResource]
	Services             krt.Collection[ServiceInfo]
	Workloads            krt.Collection[WorkloadInfo]
	DomainSuffix         string
}

func (i RouteContextInputs) WithCtx(krtctx krt.HandlerContext) RouteContext {
	return RouteContext{
		Krt:                  krtctx,
		AgentGatewayResource: i.AgentGatewayResource,
		Services:             i.Services,
		Workloads:            i.Workloads,
		DomainSuffix:         i.DomainSuffix,
	}
}

// RouteResult holds the result of a route collection
type RouteResult[T any] struct {
	Routes krt.Collection[AgentRouteResource]
	Input  krt.Collection[T]
}

// AgentHTTPRouteCollection creates a collection that translates HTTPRoute resources to agentgateway routes
func AgentHTTPRouteCollection(
	httpRoutes krt.Collection[*gwv1.HTTPRoute],
	inputs RouteContextInputs,
	krtopts krtutil.KrtOptions,
) RouteResult[*gwv1.HTTPRoute] {
	routes := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gwv1.HTTPRoute) []AgentRouteResource {
		ctx := inputs.WithCtx(krtctx)
		route := obj.Spec

		var result []AgentRouteResource

		// Get all gateways that this route can attach to
		gateways := krt.Fetch(ctx.Krt, ctx.AgentGatewayResource)
		for _, gw := range gateways {
			if !gw.Valid {
				continue
			}

			// Check if this route can attach to this gateway
			if !canRouteAttachToGateway(obj, gw) {
				continue
			}

			// Process each rule in the route
			for ruleIndex, rule := range route.Rules {
				// Split rules by matches to ensure each rule has at most one match
				matches := slices.Reference(rule.Matches)
				if len(matches) == 0 {
					matches = append(matches, nil)
				}

				for matchIndex, match := range matches {
					// Create a route for each match
					agentRoute := convertHTTPRouteToAgentRoute(ctx, rule, obj, ruleIndex, matchIndex, match, gw)
					if agentRoute != nil {
						result = append(result, *agentRoute)
					}
				}
			}
		}

		return result
	}, krtopts.ToOptions("agentgateway-http-route")...)

	return RouteResult[*gwv1.HTTPRoute]{
		Routes: routes,
		Input:  httpRoutes,
	}
}

// canRouteAttachToGateway checks if a route can attach to a specific gateway
func canRouteAttachToGateway(route *gwv1.HTTPRoute, gateway AgentGatewayResource) bool {
	// Check if the route references this gateway
	for _, parentRef := range route.Spec.ParentRefs {
		if parentRef.Name == gwv1.ObjectName(gateway.Name) {
			// Check namespace
			if parentRef.Namespace != nil && *parentRef.Namespace != gwv1.Namespace(gateway.Namespace) {
				continue
			}
			// Check gateway class
			if parentRef.SectionName != nil {
				// TODO: implement section name matching
				continue
			}
			return true
		}
	}
	return false
}

// convertHTTPRouteToAgentRoute converts an HTTPRoute rule to an agentgateway route
func convertHTTPRouteToAgentRoute(
	ctx RouteContext,
	rule gwv1.HTTPRouteRule,
	route *gwv1.HTTPRoute,
	ruleIndex int,
	matchIndex int,
	match *gwv1.HTTPRouteMatch,
	gateway AgentGatewayResource,
) *AgentRouteResource {

	// Create route key
	routeKey := fmt.Sprintf("route-%s-%d-%d-%s", route.Name, ruleIndex, matchIndex, gateway.Name)

	// Create path matches
	var matches []*api.RouteMatch
	if match != nil && match.Path != nil {
		pathMatch := convertPathMatch(match.Path)
		if pathMatch != nil {
			matches = append(matches, pathMatch)
		}
	} else {
		// Default match for all paths
		matches = append(matches, &api.RouteMatch{
			Path: &api.PathMatch{
				Kind: &api.PathMatch_PathPrefix{
					PathPrefix: "/",
				},
			},
		})
	}

	// Create header matches
	if match != nil && len(match.Headers) > 0 {
		for _, header := range match.Headers {
			headerMatch := convertHeaderMatch(header)
			if headerMatch != nil {
				matches = append(matches, headerMatch)
			}
		}
	}

	// Create backends
	backends := make([]*api.RouteBackend, 0)
	for _, backend := range rule.BackendRefs {
		backendRoutes := convertBackendRef(ctx, backend, route.Namespace)
		backends = append(backends, backendRoutes...)
	}

	// Create filters
	var filters []*api.RouteFilter
	if len(rule.Filters) > 0 {
		for _, filter := range rule.Filters {
			routeFilter := convertHTTPFilter(filter)
			if routeFilter != nil {
				filters = append(filters, routeFilter)
			}
		}
	}

	// Create the route
	agentRoute := &api.Route{
		Key:         routeKey,
		ListenerKey: gateway.Listener.Key,
		RuleName:    fmt.Sprintf("rule-%d", ruleIndex),
		RouteName:   fmt.Sprintf("route-%d", matchIndex),
		Matches:     matches,
		Filters:     filters,
		Backends:    backends,
	}

	return &AgentRouteResource{
		NamespacedName: types.NamespacedName{Namespace: route.Namespace, Name: route.Name},
		Route:          agentRoute,
		Valid:          true,
	}
}

// convertPathMatch converts a Gateway API path match to agentgateway path match
func convertPathMatch(path *gwv1.HTTPPathMatch) *api.RouteMatch {
	if path == nil {
		return nil
	}

	var pathMatch *api.PathMatch

	switch {
	case *path.Type == gwv1.PathMatchExact && path.Value != nil:
		pathMatch = &api.PathMatch{
			Kind: &api.PathMatch_Exact{
				Exact: *path.Value,
			},
		}
	case *path.Type == gwv1.PathMatchPathPrefix && path.Value != nil:
		pathMatch = &api.PathMatch{
			Kind: &api.PathMatch_PathPrefix{
				PathPrefix: *path.Value,
			},
		}
	case *path.Type == gwv1.PathMatchRegularExpression && path.Value != nil:
		pathMatch = &api.PathMatch{
			Kind: &api.PathMatch_Regex{
				Regex: *path.Value,
			},
		}
	}

	if pathMatch == nil {
		return nil
	}

	return &api.RouteMatch{
		Path: pathMatch,
	}
}

// convertHeaderMatch converts a Gateway API header match to agentgateway header match
func convertHeaderMatch(header gwv1.HTTPHeaderMatch) *api.RouteMatch {
	// TODO: implement header matching when agentgateway API supports it
	routeLogger.Debug("header matching not yet implemented")
	return nil
}

// convertBackendRef converts a Gateway API backend reference to a list of agentgateway backends
func convertBackendRef(ctx RouteContext, backend gwv1.HTTPBackendRef, routeNamespace string) []*api.RouteBackend {
	var result []*api.RouteBackend

	// Get the service name and namespace
	serviceName := string(backend.Name)
	serviceNamespace := routeNamespace
	if backend.Namespace != nil {
		serviceNamespace = string(*backend.Namespace)
	}

	// Find the service
	services := krt.Fetch(ctx.Krt, ctx.Services)
	var targetService *ServiceInfo
	for _, svc := range services {
		if svc.Service.Name == serviceName && svc.Service.Namespace == serviceNamespace {
			targetService = &svc
			break
		}
	}

	if targetService == nil {
		routeLogger.Debug("service not found for backend", "service", fmt.Sprintf("%s/%s", serviceNamespace, serviceName))
		return nil
	}

	// Determine weight
	weight := int32(1)
	if backend.Weight != nil {
		weight = *backend.Weight
	}

	svcNamespacedName := targetService.ResourceName()

	// Determine port(s)
	var port int32
	if backend.Port != nil {
		port = int32(*backend.Port)
		result = append(result, &api.RouteBackend{
			Kind: &api.RouteBackend_Service{
				Service: svcNamespacedName,
			},
			Weight: weight,
			Port:   port,
		})
	} else {
		for _, svcPort := range targetService.Service.Ports {
			if svcPort != nil {
				result = append(result, &api.RouteBackend{
					Kind: &api.RouteBackend_Service{
						Service: svcNamespacedName,
					},
					Weight: weight,
					Port:   int32(svcPort.ServicePort),
				})
			}
		}
	}

	return result
}

// convertHTTPFilter converts a Gateway API HTTP filter to agentgateway filter
func convertHTTPFilter(filter gwv1.HTTPRouteFilter) *api.RouteFilter {
	switch filter.Type {
	case gwv1.HTTPRouteFilterRequestRedirect:
		if filter.RequestRedirect == nil {
			return nil
		}
		return convertRedirectFilter(filter.RequestRedirect)
	case gwv1.HTTPRouteFilterURLRewrite:
		if filter.URLRewrite == nil {
			return nil
		}
		return convertURLRewriteFilter(filter.URLRewrite)
	case gwv1.HTTPRouteFilterRequestHeaderModifier:
		if filter.RequestHeaderModifier == nil {
			return nil
		}
		return convertHeaderModifierFilter(filter.RequestHeaderModifier)
	default:
		routeLogger.Debug("unsupported filter type", "type", filter.Type)
		return nil
	}
}

// convertRedirectFilter converts a redirect filter
func convertRedirectFilter(redirect *gwv1.HTTPRequestRedirectFilter) *api.RouteFilter {
	// TODO: implement redirect filter conversion
	routeLogger.Debug("redirect filter not yet implemented")
	return nil
}

// convertURLRewriteFilter converts a URL rewrite filter
func convertURLRewriteFilter(rewrite *gwv1.HTTPURLRewriteFilter) *api.RouteFilter {
	if rewrite == nil {
		return nil
	}

	if rewrite.Path == nil {
		return nil
	}

	var pathRewrite *api.UrlRewrite
	switch {
	case rewrite.Path.ReplaceFullPath != nil:
		pathRewrite = &api.UrlRewrite{
			Path: &api.UrlRewrite_Full{
				Full: *rewrite.Path.ReplaceFullPath,
			},
		}
	case rewrite.Path.ReplacePrefixMatch != nil:
		pathRewrite = &api.UrlRewrite{
			Path: &api.UrlRewrite_Prefix{
				Prefix: *rewrite.Path.ReplacePrefixMatch,
			},
		}
	}

	if pathRewrite == nil {
		return nil
	}

	return &api.RouteFilter{
		Kind: &api.RouteFilter_UrlRewrite{
			UrlRewrite: pathRewrite,
		},
	}
}

// convertHeaderModifierFilter converts a header modifier filter
func convertHeaderModifierFilter(modifier *gwv1.HTTPHeaderFilter) *api.RouteFilter {
	// TODO: implement header modifier filter conversion
	routeLogger.Debug("header modifier filter not yet implemented")
	return nil
}

//// AgentGRPCRouteCollection creates a collection that translates GRPCRoute resources to agentgateway routes
//func AgentGRPCRouteCollection(
//	grpcRoutes krt.Collection[*gwv1.GRPCRoute],
//	inputs RouteContextInputs,
//	krtopts krtutil.KrtOptions,
//) RouteResult[*gwv1.GRPCRoute] {
//	// TODO: implement GRPC route collection
//	routeLogger.Debug("GRPC route collection not yet implemented")
//
//	emptyRoutes := krt.NewCollection(grpcRoutes, func(krtctx krt.HandlerContext, obj *gwv1.GRPCRoute) *AgentRouteResource {
//		return nil
//	}, krtopts.ToOptions("agentgateway-grpc-route")...)
//
//	return RouteResult[*gwv1.GRPCRoute]{
//		Routes: emptyRoutes,
//		Input:  grpcRoutes,
//	}
//}
//
//// AgentTCPRouteCollection creates a collection that translates TCPRoute resources to agentgateway routes
//func AgentTCPRouteCollection(
//	tcpRoutes krt.Collection[any], // Using any for now since TCPRoute type is not available
//	inputs RouteContextInputs,
//	krtopts krtutil.KrtOptions,
//) RouteResult[any] {
//	// TODO: implement TCP route collection
//	routeLogger.Debug("TCP route collection not yet implemented")
//
//	emptyRoutes := krt.NewCollection(tcpRoutes, func(krtctx krt.HandlerContext, obj any) *AgentRouteResource {
//		return nil
//	}, krtopts.ToOptions("agentgateway-tcp-route")...)
//
//	return RouteResult[any]{
//		Routes: emptyRoutes,
//		Input:  tcpRoutes,
//	}
//}
//
//// AgentTLSRouteCollection creates a collection that translates TLSRoute resources to agentgateway routes
//func AgentTLSRouteCollection(
//	tlsRoutes krt.Collection[any], // Using any for now since TLSRoute type is not available
//	inputs RouteContextInputs,
//	krtopts krtutil.KrtOptions,
//) RouteResult[any] {
//	// TODO: implement TLS route collection
//	routeLogger.Debug("TLS route collection not yet implemented")
//
//	emptyRoutes := krt.NewCollection(tlsRoutes, func(krtctx krt.HandlerContext, obj any) *AgentRouteResource {
//		return nil
//	}, krtopts.ToOptions("agentgateway-tls-route")...)
//
//	return RouteResult[any]{
//		Routes: emptyRoutes,
//		Input:  tlsRoutes,
//	}
//}
