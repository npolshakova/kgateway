package agentgatewaysyncer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// Bind represents an agentgateway bind resource
type Bind struct {
	*api.Bind
}

func (b Bind) ResourceName() string {
	return b.Key
}

func (b Bind) Equals(other Bind) bool {
	return b.Key == other.Key && b.Port == other.Port
}

var _ krt.Equaler[Bind] = new(Bind)

// Listener represents an agentgateway listener resource
type Listener struct {
	*api.Listener
}

func (l Listener) ResourceName() string {
	return l.Key
}

func (l Listener) Equals(other Listener) bool {
	return l.Key == other.Key && l.Name == other.Name && l.BindKey == other.BindKey &&
		l.GatewayName == other.GatewayName && l.Protocol == other.Protocol
}

var _ krt.Equaler[Listener] = new(Listener)

// Route represents an agentgateway route resource
type Route struct {
	*api.Route
	// ParentRefs contains the parent references this route is bound to
	ParentRefs []routeParentReference
}

func (r *Route) ResourceName() string {
	return r.Key
}

func (r *Route) Equals(other *Route) bool {
	return r.Key == other.Key && r.ListenerKey == other.ListenerKey &&
		r.RuleName == other.RuleName && r.RouteName == other.RouteName
}

var _ krt.Equaler[*Route] = new(Route)

// RouteCollection holds routes with an index for filtering by parent reference
type RouteCollection struct {
	routes     krt.Collection[*Route]
	routeIndex krt.Index[parentKey, *Route]
}

func (s *AgentGwSyncer) translateGateway(
	krtopts krtutil.KrtOptions,
	domainSuffix string,
) krt.Collection[agentGwXdsResources] {

	gatewaysCol := krt.NewCollection(s.commonCols.GatewayIndex.Gateways, func(kctx krt.HandlerContext, gw ir.Gateway) *ir.Gateway {
		if gw.Obj.Spec.GatewayClassName != wellknown.AgentGatewayClassName {
			return nil
		}
		return &gw
	}, krtopts.ToOptions("agentgateway")...)

	routeParents := BuildRouteParents(gatewaysCol)

	refGrantClient := krt.WrapClient(kclient.New[*gateway.ReferenceGrant](s.istioClient), krtopts.ToOptions("informer/ReferenceGrants")...)
	refGrants := BuildReferenceGrants(ReferenceGrantsCollection(refGrantClient, krtopts))

	routeInputs := RouteContextInputs{
		Grants:       refGrants,
		RouteParents: routeParents,
		DomainSuffix: domainSuffix,
		Services:     s.commonCols.Services,
		Namespaces:   s.commonCols.Namespaces,
		//ServiceEntries: inputs.ServiceEntries,
		//InferencePools: s.commonCols.InferencePools,
	}

	agentGatewayRoutes := ADPRouteCollection(
		s.httprouteCol,
		routeInputs,
		krtopts,
	)

	// k8s service -> service info ir
	// these are workloadapi-style services combined from kube services (todo: support services entries, backend types, etc.)
	agentGwServices := buildServicesCollection(s.commonCols.Services, domainSuffix)

	epSliceClient := kclient.NewFiltered[*discoveryv1.EndpointSlice](
		s.commonCols.Client,
		kclient.Filter{ObjectFilter: s.commonCols.Client.ObjectFilter()},
	)
	endpointSlices := krt.WrapClient(epSliceClient, s.commonCols.KrtOpts.ToOptions("EndpointSlices")...)

	pods := s.commonCols.PodWrapper

	// Create workload resources for agentgateway workload API
	agentGwWorkloads := buildWorkloadsCollection(
		pods,
		agentGwServices,
		endpointSlices,
		domainSuffix, s.clusterId,
	)

	s.waitForSync = append(s.waitForSync, []cache.InformerSynced{
		endpointSlices.HasSynced,
		pods.HasSynced,
		gatewaysCol.HasSynced,
		agentGwServices.HasSynced,
		agentGwWorkloads.HasSynced,
		s.httprouteCol.HasSynced,
	}...)

	gatewayResourcesXds := krt.NewCollection(gatewaysCol, func(kctx krt.HandlerContext, gw ir.Gateway) *agentGwXdsResources {
		gwNamespacedName := types.NamespacedName{Namespace: gw.Namespace, Name: gw.Name}

		// translate agw resources (Bind, Listener, Routes)
		resources := translateAGWResources(kctx, gw, agentGatewayRoutes)
		addrResources := translateAGWAddress(kctx, agentGwServices, agentGwWorkloads)

		// Create the resource wrappers
		var resourceVersion uint64
		for _, res := range resources {
			resourceVersion ^= res.(*envoyResourceWithCustomName).version
		}

		var addrVersion uint64
		for _, addrRes := range addrResources {
			addrVersion ^= addrRes.(*envoyResourceWithCustomName).version
		}

		result := &agentGwXdsResources{
			NamespacedName: gwNamespacedName,
			ResourceConfig: envoycache.NewResources(fmt.Sprintf("%d", resourceVersion), resources),
			AddressConfig:  envoycache.NewResources(fmt.Sprintf("%d", addrVersion), resources),
		}
		logger.Debug("created XDS resources for gateway with ID", "gwname", gw.Name, "resourceid", result.ResourceName())
		return result
	})

	return gatewayResourcesXds
}

func translateAGWResources(
	kctx krt.HandlerContext,
	gw ir.Gateway,
	agentGatewayRoutes krt.Collection[ADPResource],
) []envoytypes.Resource {
	var resources []envoytypes.Resource

	bindsForGw := translateBinds(gw)
	for _, bind := range bindsForGw {
		bindResource := &api.Resource{
			Kind: &api.Resource_Bind{
				Bind: bind.Bind,
			},
		}
		resources = append(resources, &envoyResourceWithCustomName{
			Message: bindResource,
			Name:    bind.Bind.GetKey(),
			version: utils.HashProto(bindResource),
		})
	}

	listenersForGw := translateListeners(gw)
	for _, listener := range listenersForGw {
		listenerResource := &api.Resource{
			Kind: &api.Resource_Listener{
				Listener: listener.Listener,
			},
		}
		resources = append(resources, &envoyResourceWithCustomName{
			Message: listenerResource,
			Name:    listener.Key,
			version: utils.HashProto(listenerResource),
		})
	}

	// routes for the gw - filter by gateway
	gwNamespacedName := types.NamespacedName{Namespace: gw.Namespace, Name: gw.Name}
	gwRoutes := krt.Fetch(kctx, agentGatewayRoutes, krt.FilterGeneric(func(obj any) bool {
		adpResource := obj.(ADPResource)
		return adpResource.Gateway == gwNamespacedName
	}))
	for _, adpResource := range gwRoutes {
		resources = append(resources, &envoyResourceWithCustomName{
			Message: adpResource.Resource,
			Name:    adpResource.ResourceName(),
			version: utils.HashProto(adpResource.Resource),
		})
	}
	return resources
}

func translateAGWAddress(
	kctx krt.HandlerContext,
	agentGwServices krt.Collection[ServiceInfo],
	agentGwWorkloads krt.Collection[WorkloadInfo],
) []envoytypes.Resource {
	var addrResources []envoytypes.Resource

	// service ir -> envoy resource for gw
	services := krt.Fetch(kctx, agentGwServices)
	for _, svc := range services {
		// Create Service resource
		serviceResource := svc.Service

		// Create Address resource wrapping the Service
		addressResource := &api.Address{
			Type: &api.Address_Service{
				Service: serviceResource,
			},
		}

		addrResources = append(addrResources, &envoyResourceWithCustomName{
			Message: addressResource,
			Name:    svc.ResourceName(),
			version: utils.HashProto(addressResource),
		})
	}

	// workload addresses for gw
	workloads := krt.Fetch(kctx, agentGwWorkloads)
	for _, workload := range workloads {
		// Create Address resource wrapping the Service
		addressResource := &api.Address{
			Type: &api.Address_Workload{
				Workload: workload.Workload,
			},
		}

		addrResources = append(addrResources, &envoyResourceWithCustomName{
			Message: addressResource,
			Name:    workload.ResourceName(),
			version: utils.HashProto(addressResource),
		})
	}
	return addrResources
}

func translateBinds(gw ir.Gateway) []*Bind {
	var bindsForGw []*Bind
	for _, listener := range gw.Listeners {
		port := uint32(listener.Port)
		bindKey := fmt.Sprintf("%d/%s/%s", port, gw.Namespace, gw.Name)
		bindForListener := &Bind{
			Bind: &api.Bind{
				Key:  bindKey,
				Port: port,
			},
		}
		bindsForGw = append(bindsForGw, bindForListener)
	}
	return bindsForGw
}

func translateListeners(gw ir.Gateway) []*Listener {
	var listenersForGw []*Listener
	for _, listener := range gw.Listeners {
		var hostname string
		if listener.Hostname != nil {
			hostname = string(*listener.Hostname)
		}

		l := &api.Listener{
			Key:         gw.ResourceName(),
			Name:        string(listener.Name),
			BindKey:     fmt.Sprint(listener.Port) + "/" + gw.Namespace + "/" + gw.Name,
			GatewayName: gw.Namespace + "/" + gw.Name,
			Hostname:    hostname,
		}

		switch listener.Protocol {
		case gwv1.HTTPProtocolType:
			l.Protocol = api.Protocol_HTTP
		case gwv1.HTTPSProtocolType:
			l.Protocol = api.Protocol_HTTPS
			if listener.TLS == nil {
				return nil
			}
			// TODO: handle tls cert ref resolution
		case gwv1.TLSProtocolType:
			l.Protocol = api.Protocol_TLS
			if listener.TLS == nil {
				return nil
			}
			// TODO: handle tls cert ref resolution
		case gwv1.TCPProtocolType:
			l.Protocol = api.Protocol_TCP
		default:
			return nil
		}
		listenersForGw = append(listenersForGw, &Listener{
			Listener: l,
		})
	}
	return listenersForGw
}

// RouteParents holds information about things routes can reference as parents.
type RouteParents struct {
	gateways     krt.Collection[ir.Gateway]
	gatewayIndex krt.Index[parentKey, ir.Gateway]
}

func (p RouteParents) fetch(ctx krt.HandlerContext, pk parentKey) []*parentInfo {
	// TODO: fetch is empty? add default type for kind?
	gateways := krt.Fetch(ctx, p.gateways, krt.FilterIndex(p.gatewayIndex, pk))
	var result []*parentInfo

	for _, gw := range gateways {
		// For each listener in the gateway, create a parentInfo
		for _, listener := range gw.Listeners {
			// Extract hostnames from the listener
			var hostnames []string
			if listener.Hostname != nil {
				hostnames = append(hostnames, string(*listener.Hostname))
			}

			// Extract allowed kinds from the listener's AllowedRoutes
			var allowedKinds []gwv1.RouteGroupKind
			if listener.AllowedRoutes != nil {
				allowedKinds = listener.AllowedRoutes.Kinds
			}

			// Create parentInfo for this listener
			pi := &parentInfo{
				InternalName: fmt.Sprintf("%s/%s", gw.Namespace, gw.Name),
				AllowedKinds: allowedKinds,
				Hostnames:    hostnames,
				SectionName:  listener.Name,
				Port:         listener.Port,
				Protocol:     listener.Protocol,
			}
			result = append(result, pi)
		}
	}

	return result
}

func BuildRouteParents(
	gateways krt.Collection[ir.Gateway],
) RouteParents {
	idx := krt.NewIndex(gateways, func(o ir.Gateway) []parentKey {
		return []parentKey{
			{
				Kind:      o.GetGroupKind(),
				Name:      o.Name,
				Namespace: o.Namespace,
			},
		}
	})
	return RouteParents{
		gateways:     gateways,
		gatewayIndex: idx,
	}
}

// getTargetName sanitizes the given resource name to ensure it matches the AgentGateway required pattern:
// ^[a-zA-Z0-9-]+$ by replacing slashes and removing invalid characters.
func getTargetName(resourceName string) string {
	var (
		invalidCharsRegex      = regexp.MustCompile(`[^a-zA-Z0-9-]+`)
		consecutiveDashesRegex = regexp.MustCompile(`-+`)
	)

	// Replace all invalid characters with dashes
	sanitized := invalidCharsRegex.ReplaceAllString(resourceName, "-")

	// Remove leading/trailing dashes and collapse consecutive dashes
	sanitized = strings.Trim(sanitized, "-")
	sanitized = consecutiveDashesRegex.ReplaceAllString(sanitized, "-")

	return sanitized
}
