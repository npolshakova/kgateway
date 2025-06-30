package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

var routeLogger = logging.New("agentgateway/route-collection")

// AgentRouteResource represents a translated route resource for agentgateway
type AgentRouteResource struct {
	types.NamespacedName
	Route *api.Route
}

func (r AgentRouteResource) ResourceName() string {
	return fmt.Sprintf("%s/%s", r.Namespace, r.Name)
}

func (r AgentRouteResource) Equals(other AgentRouteResource) bool {
	if r.NamespacedName != other.NamespacedName {
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
	Krt           krt.HandlerContext
	RouteParents  RouteParents
	DomainSuffix  string
	Services      krt.Collection[*corev1.Service]
	ServicesIndex krt.Index[string, *corev1.Service]
	Namespaces    krt.Collection[krtcollections.NamespaceMetadata]
}

// RouteContextInputs defines the inputs needed for route processing
type RouteContextInputs struct {
	RouteParents  RouteParents
	DomainSuffix  string
	Services      krt.Collection[*corev1.Service]
	ServicesIndex krt.Index[string, *corev1.Service]
	Namespaces    krt.Collection[krtcollections.NamespaceMetadata]
}

func (i RouteContextInputs) WithCtx(krtctx krt.HandlerContext) RouteContext {
	return RouteContext{
		Krt:           krtctx,
		RouteParents:  i.RouteParents,
		Services:      i.Services,
		ServicesIndex: i.ServicesIndex,
		Namespaces:    i.Namespaces,
		DomainSuffix:  i.DomainSuffix,
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
