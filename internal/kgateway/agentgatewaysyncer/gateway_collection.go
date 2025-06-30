package agentgatewaysyncer

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"google.golang.org/protobuf/types/known/durationpb"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/protomarshal"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
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

	// Create services index for efficient key-based lookups
	servicesIndex := krt.NewIndex(s.commonCols.Services, func(svc *corev1.Service) []string {
		return []string{svc.Namespace + "/" + svc.Name}
	})

	routeInputs := RouteContextInputs{
		RouteParents:  routeParents,
		DomainSuffix:  domainSuffix,
		Services:      s.commonCols.Services,
		ServicesIndex: servicesIndex,
		Namespaces:    s.commonCols.Namespaces,
	}

	agentGatewayRoutes := agentGatewayRouteCollection(
		s.httprouteCol,
		routeInputs,
	)

	// k8s service -> service info ir
	// these are workloadapi-style services combined from kube services (todo: support services entries, backend types, etc.)
	agentGwServices := buildServicesCollection(s.commonCols.Services, domainSuffix)

	epSliceClient := kclient.NewFiltered[*discoveryv1.EndpointSlice](
		s.commonCols.Client,
		kclient.Filter{ObjectFilter: s.commonCols.Client.ObjectFilter()},
	)
	endpointSlices := krt.WrapClient(epSliceClient, s.commonCols.KrtOpts.ToOptions("EndpointSlices")...)

	podsClient := kclient.NewFiltered[*corev1.Pod](
		s.commonCols.Client,
		kclient.Filter{ObjectFilter: s.commonCols.Client.ObjectFilter()},
	)
	pods := krt.WrapClient(podsClient, s.commonCols.KrtOpts.ToOptions("Pods")...)

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
	agentGatewayRoutes RouteCollection,
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

	// routes for the gw - filter by parent key
	gwParentKey := parentKey{
		Kind:      gw.Obj.GetObjectKind().GroupVersionKind(),
		Name:      gw.Name,
		Namespace: gw.Namespace,
	}
	gwRoutes := krt.Fetch(kctx, agentGatewayRoutes.routes, krt.FilterIndex(agentGatewayRoutes.routeIndex, gwParentKey))
	for _, route := range gwRoutes {
		routeResource := &api.Resource{
			Kind: &api.Resource_Route{
				Route: route.Route,
			},
		}
		resources = append(resources, &envoyResourceWithCustomName{
			Message: routeResource,
			Name:    route.Route.Key,
			version: utils.HashProto(routeResource),
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

func agentGatewayRouteCollection(
	httpRoutes krt.Collection[*gwv1.HTTPRoute],
	inputs RouteContextInputs,
) RouteCollection {
	routes := krt.NewManyCollection(httpRoutes, func(krtctx krt.HandlerContext, obj *gwv1.HTTPRoute) []*Route {
		ctx := inputs.WithCtx(krtctx)
		var results []*Route

		parentRefs := extractParentReferenceInfo(ctx, ctx.RouteParents, obj)
		filteredParentRefs := filteredReferences(parentRefs)
		for n, r := range obj.Spec.Rules {
			// split the rule to make sure each rule has up to one match
			matches := slices.Reference(r.Matches)
			if len(matches) == 0 {
				matches = append(matches, nil)
			}
			for idx, m := range matches {
				if m != nil {
					r.Matches = []gwv1.HTTPRouteMatch{*m}
				}
				apiRoute, err := convertHTTPRouteToAGWRoute(ctx, r, obj, n, idx)
				if err != nil {
					// TODO: append err on result?
					logger.Error("error converting http route", "err", err.Message)
				}
				if len(filteredParentRefs) > 0 {
					for _, parent := range filteredParentRefs {
						_, listenerName, _ := strings.Cut(parent.InternalName, "/")
						inner := protomarshal.Clone(apiRoute)
						inner.ListenerKey = listenerName
						inner.Key = inner.Key + "." + string(parent.ParentSection)
						results = append(results, &Route{
							Route:      inner,
							ParentRefs: parentRefs,
						})
					}
				} else {
					results = append(results, &Route{
						Route:      apiRoute,
						ParentRefs: parentRefs,
					})
				}
			}
		}

		return results
	}, krt.WithName("Routes"))

	// Create an index on routes by parent reference
	routesWithIndex := krt.NewIndex(routes, func(route *Route) []parentKey {
		var keys []parentKey
		for _, parentRef := range route.ParentRefs {
			keys = append(keys, parentRef.ParentKey)
		}
		return keys
	})

	return RouteCollection{
		routes:     routes,
		routeIndex: routesWithIndex,
	}
}

func filteredReferences(parents []routeParentReference) []routeParentReference {
	ret := make([]routeParentReference, 0, len(parents))
	for _, p := range parents {
		if p.DeniedReason != nil {
			// We should filter this out
			continue
		}
		ret = append(ret, p)
	}
	// To ensure deterministic order, sort them
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].InternalName < ret[j].InternalName
	})
	return ret
}

type ParentErrorReason string

const (
	ParentErrorNotAccepted       = ParentErrorReason(gwv1.RouteReasonNoMatchingParent)
	ParentErrorNotAllowed        = ParentErrorReason(gwv1.RouteReasonNotAllowedByListeners)
	ParentErrorNoHostname        = ParentErrorReason(gwv1.RouteReasonNoMatchingListenerHostname)
	ParentErrorParentRefConflict = ParentErrorReason("ParentRefConflict")
	ParentNoError                = ParentErrorReason("")
)

// ParentError represents that a parent could not be referenced
type ParentError struct {
	Reason  ParentErrorReason
	Message string
}

// routeParentReference holds information about a route's parent reference
type routeParentReference struct {
	// InternalName refers to the internal name of the parent we can reference it by. For example  "my-ns/my-gateway"
	InternalName string
	// InternalKind is the Group/Kind of the parent
	InternalKind schema.GroupVersionKind
	// DeniedReason, if present, indicates why the reference was not valid
	DeniedReason *ParentError
	// OriginalReference contains the original reference
	OriginalReference gwv1.ParentReference
	// Hostname is the hostname match of the parent, if any
	Hostname        string
	BannedHostnames sets.Set[string]
	ParentKey       parentKey
	ParentSection   gwv1.SectionName
}

type ConfigErrorReason = string

const (
	// InvalidDestination indicates an issue with the destination
	InvalidDestination ConfigErrorReason = "InvalidDestination"
	InvalidAddress     ConfigErrorReason = ConfigErrorReason(gwv1.GatewayReasonUnsupportedAddress)
	// InvalidDestinationPermit indicates a destination was not permitted
	InvalidDestinationPermit ConfigErrorReason = ConfigErrorReason(gwv1.RouteReasonRefNotPermitted)
	// InvalidDestinationKind indicates an issue with the destination kind
	InvalidDestinationKind ConfigErrorReason = ConfigErrorReason(gwv1.RouteReasonInvalidKind)
	// InvalidDestinationNotFound indicates a destination does not exist
	InvalidDestinationNotFound ConfigErrorReason = ConfigErrorReason(gwv1.RouteReasonBackendNotFound)
	// InvalidFilter indicates an issue with the filters
	InvalidFilter ConfigErrorReason = "InvalidFilter"
	// InvalidTLS indicates an issue with TLS settings
	InvalidTLS ConfigErrorReason = ConfigErrorReason(gwv1.ListenerReasonInvalidCertificateRef)
	// InvalidListenerRefNotPermitted indicates a listener reference was not permitted
	InvalidListenerRefNotPermitted ConfigErrorReason = ConfigErrorReason(gwv1.ListenerReasonRefNotPermitted)
	// InvalidConfiguration indicates a generic error for all other invalid configurations
	InvalidConfiguration ConfigErrorReason = "InvalidConfiguration"
	DeprecateFieldUsage  ConfigErrorReason = "DeprecatedField"
)

// ConfigError represents an invalid configuration that will be reported back to the user.
type ConfigError struct {
	Reason  ConfigErrorReason
	Message string
}

// normalizeReference takes a generic Group/Kind (the API uses a few variations) and converts to a known GroupVersionKind.
// Defaults for the group/kind are also passed.
func normalizeReference[G ~string, K ~string](group *G, kind *K, def schema.GroupVersionKind) schema.GroupVersionKind {
	k := def.Kind
	if kind != nil {
		k = string(*kind)
	}
	g := def.Group
	if group != nil {
		g = string(*group)
	}
	s, f := collections.All.FindByGroupKind(config.GroupVersionKind{Group: g, Kind: k})
	if f {
		return schema.GroupVersionKind{
			Group:   s.GroupVersionKind().Group,
			Version: s.GroupVersionKind().Version,
			Kind:    s.GroupVersionKind().Kind,
		}
	}
	// TODO: check this
	// If not found in collections (e.g., Gateway API types), use the default GVK
	// This ensures Gateway API types get the correct version
	return def
}

func defaultString[T ~string](s *T, def string) string {
	if s == nil {
		return def
	}
	return string(*s)
}

type parentReference struct {
	parentKey

	SectionName gwv1.SectionName
	Port        gwv1.PortNumber
}

// parentPortInAllowSet checks if the given port exists in any of the allowed listener sets
func parentPortInAllowSet(port gwv1.PortNumber, allowedListenerSets []ir.ListenerSet) bool {
	for _, ls := range allowedListenerSets {
		for _, listener := range ls.Listeners {
			if listener.Port == port {
				return true
			}
		}
	}
	return false
}

// parentPortNotInDenySet checks if the given port does NOT exist in any of the denied listener sets
func parentPortNotInDenySet(port gwv1.PortNumber, deniedListenerSets []ir.ListenerSet) bool {
	for _, ls := range deniedListenerSets {
		for _, listener := range ls.Listeners {
			if listener.Port == port {
				return false
			}
		}
	}
	return true
}

// parentSectionNameInAllowSet checks if the given section name exists in any of the allowed listener sets
func parentSectionNameInAllowSet(sectionName gwv1.SectionName, allowedListenerSets []ir.ListenerSet) bool {
	for _, ls := range allowedListenerSets {
		for _, listener := range ls.Listeners {
			if string(listener.Name) == string(sectionName) {
				return true
			}
		}
	}
	return false
}

// parentSectionNameNotInDenySet checks if the given section name does NOT exist in any of the denied listener sets
func parentSectionNameNotInDenySet(sectionName gwv1.SectionName, deniedListenerSets []ir.ListenerSet) bool {
	for _, ls := range deniedListenerSets {
		for _, listener := range ls.Listeners {
			if string(listener.Name) == string(sectionName) {
				return false
			}
		}
	}
	return true
}

func referenceAllowed(
	ctx RouteContext,
	parentGw *ir.Gateway,
	routeKind schema.GroupVersionKind,
	parentRef parentReference,
	hostnames []gwv1.Hostname,
	localNamespace string,
) *ParentError {
	if parentRef.Kind == wellknown.ServiceGVK {

		key := parentRef.Namespace + "/" + parentRef.Name
		svc := ptr.Flatten(krt.FetchOne(ctx.Krt, ctx.Services, krt.FilterIndex(ctx.ServicesIndex, key)))

		// check that the referenced svc exists
		if svc == nil {
			return &ParentError{
				Reason:  ParentErrorNotAccepted,
				Message: fmt.Sprintf("parent service: %q not found", parentRef.Name),
			}
		}

	} else {
		// First, check section and port apply. This must come first
		if parentRef.Port != 0 && (!parentPortInAllowSet(parentRef.Port, parentGw.AllowedListenerSets) || !parentPortNotInDenySet(parentRef.Port, parentGw.DeniedListenerSets)) {
			return &ParentError{
				Reason:  ParentErrorNotAccepted,
				Message: fmt.Sprintf("port %v not found in allowed listener sets or is in denied listener sets", parentRef.Port),
			}
		}
		if len(parentRef.SectionName) > 0 && (!parentSectionNameInAllowSet(parentRef.SectionName, parentGw.AllowedListenerSets) || !parentSectionNameNotInDenySet(parentRef.SectionName, parentGw.DeniedListenerSets)) {
			return &ParentError{
				Reason:  ParentErrorNotAccepted,
				Message: fmt.Sprintf("sectionName %q not found in allowed listener sets or is in denied listener sets", parentRef.SectionName),
			}
		}

		// TODO: Next check the hostnames are a match. This is a bi-directional wildcard match. Only one route
		// hostname must match for it to be allowed (but the others will be filtered at runtime)
		// If either is empty its treated as a wildcard which always matches

	}
	// TODO: Also make sure this route kind is allowed
	return nil
}

func convertHTTPRouteToAGWRoute(ctx RouteContext, r gwv1.HTTPRouteRule,
	obj *gwv1.HTTPRoute, pos int, matchPos int,
) (*api.Route, *ConfigError) {

	res := &api.Route{
		Key:         obj.Namespace + "." + obj.Name + "." + strconv.Itoa(pos) + "." + strconv.Itoa(matchPos),
		RouteName:   obj.Namespace + "/" + obj.Name,
		ListenerKey: "",
		RuleName:    defaultString(r.Name, ""),
	}

	for _, match := range r.Matches {
		path, err := createADPPathMatch(match)
		if err != nil {
			return nil, err
		}
		headers, err := createADPHeadersMatch(match)
		if err != nil {
			return nil, err
		}
		//method, err := createADPMethodMatch(match)
		//if err != nil {
		//	return nil, err
		//}
		//query, err := createADPQueryMatch(match)
		//if err != nil {
		//	return nil, err
		//}
		res.Matches = append(res.Matches, &api.RouteMatch{
			Path:    path,
			Headers: headers,
			//Method:      method,
			//QueryParams: query,
		})
	}
	filters, err := buildADPFilters(ctx, obj.Namespace, r.Filters)
	if err != nil {
		return nil, err
	}
	res.Filters = filters

	if r.Timeouts != nil {
		res.TrafficPolicy = &api.TrafficPolicy{}
		if r.Timeouts.Request != nil {
			request, _ := time.ParseDuration(string(*r.Timeouts.Request))
			if request > 0 {
				res.TrafficPolicy.RequestTimeout = durationpb.New(request)
			}
		}
		if r.Timeouts.BackendRequest != nil {
			request, _ := time.ParseDuration(string(*r.Timeouts.BackendRequest))
			if request > 0 {
				res.TrafficPolicy.RequestTimeout = durationpb.New(request)
			}
		}
	}

	// Retry: todo
	route, backendErr, err := buildADPHTTPDestination(ctx, r.BackendRefs, obj.Namespace)
	if err != nil {
		return nil, err
	}
	res.Backends = route
	res.Hostnames = slices.Map(obj.Spec.Hostnames, func(e gwv1.Hostname) string {
		return string(e)
	})
	return res, backendErr
}

// parentKey holds info about a parentRef (eg route binding to a Gateway). This is a mirror of
// gwv1.ParentReference in a form that can be stored in a map
type parentKey struct {
	Kind schema.GroupVersionKind
	// Name is the original name of the resource (eg Kubernetes Gateway name)
	Name string
	// Namespace is the namespace of the resource
	Namespace string
}

func (p parentKey) String() string {
	return p.Kind.String() + "/" + p.Namespace + "/" + p.Name
}

// RouteParents holds information about things routes can reference as parents.
type RouteParents struct {
	gateways     krt.Collection[ir.Gateway]
	gatewayIndex krt.Index[parentKey, ir.Gateway]
}

// parentInfo holds info about a "parent" - something that can be referenced as a ParentRef in the API.
type parentInfo struct {
	// InternalName refers to the internal name we can reference it by. For example, "my-ns/my-gateway"
	InternalName string
	// AllowedKinds indicates which kinds can be admitted by this parent
	AllowedKinds []gwv1.RouteGroupKind
	// Hostnames is the hostnames that must be match to reference to the parent. For gateway this is listener hostname
	Hostnames []string

	SectionName gwv1.SectionName
	Port        gwv1.PortNumber
	Protocol    gwv1.ProtocolType
}

func (p RouteParents) fetch(ctx krt.HandlerContext, pk parentKey) []parentInfo {
	gateways := krt.Fetch(ctx, p.gateways, krt.FilterIndex(p.gatewayIndex, pk))
	var result []parentInfo

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
			pi := parentInfo{
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
				Kind:      o.Obj.GetObjectKind().GroupVersionKind(),
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

func createADPPathMatch(match gwv1.HTTPRouteMatch) (*api.PathMatch, *ConfigError) {
	tp := gwv1.PathMatchPathPrefix
	if match.Path.Type != nil {
		tp = *match.Path.Type
	}
	dest := "/"
	if match.Path.Value != nil {
		dest = *match.Path.Value
	}
	switch tp {
	case gwv1.PathMatchPathPrefix:
		// "When specified, a trailing `/` is ignored."
		if dest != "/" {
			dest = strings.TrimSuffix(dest, "/")
		}
		return &api.PathMatch{Kind: &api.PathMatch_PathPrefix{
			PathPrefix: dest,
		}}, nil
	case gwv1.PathMatchExact:
		return &api.PathMatch{Kind: &api.PathMatch_Exact{
			Exact: dest,
		}}, nil
	case gwv1.PathMatchRegularExpression:
		return &api.PathMatch{Kind: &api.PathMatch_Regex{
			Regex: dest,
		}}, nil
	default:
		// Should never happen, unless a new field is added
		return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported Path match type", tp)}
	}
}

func createADPHeadersMatch(match gwv1.HTTPRouteMatch) ([]*api.HeaderMatch, *ConfigError) {
	res := []*api.HeaderMatch{}
	for _, header := range match.Headers {
		tp := gwv1.HeaderMatchExact
		if header.Type != nil {
			tp = *header.Type
		}
		switch tp {
		case gwv1.HeaderMatchExact:
			res = append(res, &api.HeaderMatch{
				Name:  string(header.Name),
				Value: &api.HeaderMatch_Exact{Exact: header.Value},
			})
		case gwv1.HeaderMatchRegularExpression:
			res = append(res, &api.HeaderMatch{
				Name:  string(header.Name),
				Value: &api.HeaderMatch_Regex{Regex: header.Value},
			})
		default:
			// Should never happen, unless a new field is added
			return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported HeaderMatch type", tp)}
		}
	}

	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func createADPHeadersFilter(filter *gwv1.HTTPHeaderFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_RequestHeaderModifier{
			RequestHeaderModifier: &api.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPResponseHeadersFilter(filter *gwv1.HTTPHeaderFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_ResponseHeaderModifier{
			ResponseHeaderModifier: &api.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPRewriteFilter(filter *gwv1.HTTPURLRewriteFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	ff := &api.UrlRewrite{
		Host: string(ptr.OrEmpty(filter.Hostname)),
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case gwv1.PrefixMatchHTTPPathModifier:
			ff.Path = &api.UrlRewrite_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case gwv1.FullPathHTTPPathModifier:
			ff.Path = &api.UrlRewrite_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_UrlRewrite{
			UrlRewrite: ff,
		},
	}
}

func createADPMirrorFilter(
	ctx RouteContext,
	filter *gwv1.HTTPRequestMirrorFilter,
	ns string,
) (*api.RouteFilter, *ConfigError) {
	if filter == nil {
		return nil, nil
	}
	var weightOne int32 = 1
	dst, err := buildADPDestination(ctx, gwv1.HTTPBackendRef{
		BackendRef: gwv1.BackendRef{
			BackendObjectReference: filter.BackendRef,
			Weight:                 &weightOne,
		},
	}, ns)
	if err != nil {
		return nil, err
	}
	var percent float64
	if f := filter.Fraction; f != nil {
		percent = (100 * float64(f.Numerator)) / float64(ptr.OrDefault(f.Denominator, int32(100)))
	} else if p := filter.Percent; p != nil {
		percent = float64(*p)
	} else {
		percent = 100
	}
	if percent == 0 {
		return nil, nil
	}
	rm := &api.RequestMirror{
		Kind:       nil,
		Percentage: percent,
		Port:       dst.Port,
	}
	switch dk := dst.Kind.(type) {
	case *api.RouteBackend_Service:
		rm.Kind = &api.RequestMirror_Service{
			Service: dk.Service,
		}
	}
	return &api.RouteFilter{Kind: &api.RouteFilter_RequestMirror{RequestMirror: rm}}, nil
}

func createADPRedirectFilter(filter *gwv1.HTTPRequestRedirectFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	ff := &api.RequestRedirect{
		Scheme: ptr.OrEmpty(filter.Scheme),
		Host:   string(ptr.OrEmpty(filter.Hostname)),
		Port:   uint32(ptr.OrEmpty(filter.Port)),
		Status: uint32(ptr.OrEmpty(filter.StatusCode)),
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case gwv1.PrefixMatchHTTPPathModifier:
			ff.Path = &api.RequestRedirect_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case gwv1.FullPathHTTPPathModifier:
			ff.Path = &api.RequestRedirect_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_RequestRedirect{
			RequestRedirect: ff,
		},
	}
}

func headerListToADP(hl []gwv1.HTTPHeader) []*api.Header {
	return slices.Map(hl, func(hl gwv1.HTTPHeader) *api.Header {
		return &api.Header{
			Name:  string(hl.Name),
			Value: hl.Value,
		}
	})
}

func buildADPDestination(
	ctx RouteContext,
	to gwv1.HTTPBackendRef,
	ns string,
) (*api.RouteBackend, *ConfigError) {
	// TODO: enforce ref grant and check if the reference is allowed

	namespace := ptr.OrDefault((*string)(to.Namespace), ns)
	var invalidBackendErr *ConfigError
	var hostname string
	ref := normalizeReference(to.Group, to.Kind, wellknown.ServiceGVK)
	rb := &api.RouteBackend{
		Weight: ptr.OrDefault(to.Weight, 1),
	}
	var port *gwv1.PortNumber
	switch ref.GroupKind() {
	case wellknown.ServiceGVK.GroupKind():
		port = to.Port
		if strings.Contains(string(to.Name), ".") {
			return nil, &ConfigError{Reason: InvalidDestination, Message: "service name invalid; the name of the Service must be used, not the hostname."}
		}
		hostname = fmt.Sprintf("%s.%s.svc.%s", to.Name, namespace, ctx.DomainSuffix)
		key := namespace + "/" + string(to.Name)
		svc := ptr.Flatten(krt.FetchOne(ctx.Krt, ctx.Services, krt.FilterIndex(ctx.ServicesIndex, key)))
		if svc == nil {
			invalidBackendErr = &ConfigError{Reason: InvalidDestinationNotFound, Message: fmt.Sprintf("backend(%s) not found", hostname)}
		}
		rb.Kind = &api.RouteBackend_Service{Service: namespace + "/" + hostname}
	default:
		port = to.Port
		return nil, &ConfigError{
			Reason:  InvalidDestinationKind,
			Message: fmt.Sprintf("referencing unsupported backendRef: group %q kind %q", ptr.OrEmpty(to.Group), ptr.OrEmpty(to.Kind)),
		}
	}
	// All types currently require a Port, so we do this for everything; consider making this per-type if we have future types
	// that do not require port.
	if port == nil {
		// "Port is required when the referent is a Kubernetes Service."
		return nil, &ConfigError{Reason: InvalidDestination, Message: "port is required in backendRef"}
	}
	rb.Port = int32(*port)
	return rb, invalidBackendErr
}

func buildADPHTTPDestination(
	ctx RouteContext,
	forwardTo []gwv1.HTTPBackendRef,
	ns string,
) ([]*api.RouteBackend, *ConfigError, *ConfigError) {
	if forwardTo == nil {
		return nil, nil, nil
	}

	var invalidBackendErr *ConfigError
	var res []*api.RouteBackend
	for _, fwd := range forwardTo {
		dst, err := buildADPDestination(ctx, fwd, ns)
		if err != nil {
			logger.Error("error building destination", "err", err.Message)
			if isInvalidBackend(err) {
				invalidBackendErr = err
				// keep going, we will gracefully drop invalid backends
			} else {
				return nil, nil, err
			}
		}
		if dst != nil {
			filters, err := buildADPFilters(ctx, ns, fwd.Filters)
			if err != nil {
				return nil, nil, err
			}
			dst.Filters = filters
		}
		res = append(res, dst)
	}
	return res, invalidBackendErr, nil
}

func buildADPFilters(
	ctx RouteContext,
	ns string,
	inputFilters []gwv1.HTTPRouteFilter,
) ([]*api.RouteFilter, *ConfigError) {
	filters := []*api.RouteFilter{}
	var mirrorBackendErr *ConfigError
	for _, filter := range inputFilters {
		switch filter.Type {
		case gwv1.HTTPRouteFilterRequestHeaderModifier:
			h := createADPHeadersFilter(filter.RequestHeaderModifier)
			if h == nil {
				continue
			}
			filters = append(filters, h)
		case gwv1.HTTPRouteFilterResponseHeaderModifier:
			h := createADPResponseHeadersFilter(filter.ResponseHeaderModifier)
			if h == nil {
				continue
			}
			filters = append(filters, h)
		case gwv1.HTTPRouteFilterRequestRedirect:
			h := createADPRedirectFilter(filter.RequestRedirect)
			if h == nil {
				continue
			}
			filters = append(filters, h)
		case gwv1.HTTPRouteFilterRequestMirror:
			h, err := createADPMirrorFilter(ctx, filter.RequestMirror, ns)
			if err != nil {
				mirrorBackendErr = err
			} else {
				filters = append(filters, h)
			}
		case gwv1.HTTPRouteFilterURLRewrite:
			h := createADPRewriteFilter(filter.URLRewrite)
			if h == nil {
				continue
			}
			filters = append(filters, h)
		default:
			return nil, &ConfigError{
				Reason:  InvalidFilter,
				Message: fmt.Sprintf("unsupported filter type %q", filter.Type),
			}
		}
	}
	return filters, mirrorBackendErr
}

// https://github.com/kubernetes-sigs/gateway-api/blob/cea484e38e078a2c1997d8c7a62f410a1540f519/apis/v1beta1/httproute_types.go#L207-L212
func isInvalidBackend(err *ConfigError) bool {
	return err.Reason == InvalidDestinationPermit ||
		err.Reason == InvalidDestinationNotFound ||
		err.Reason == InvalidDestinationKind
}

func extractParentReferenceInfo(ctx RouteContext, parents RouteParents, obj controllers.Object) []routeParentReference {
	routeRefs, _, _ := getCommonRouteInfo(obj)
	localNamespace := obj.GetNamespace()
	var parentRefs []routeParentReference
	for _, ref := range routeRefs {
		pk, err := toInternalParentReference(ref, localNamespace)
		if err != nil {
			// Cannot handle the reference. Maybe it is for another controller, so we just ignore it
			continue
		}
		parentRef := parentReference{
			parentKey:   pk,
			SectionName: ptr.OrEmpty(ref.SectionName),
			Port:        ptr.OrEmpty(ref.Port),
		}
		currentParents := parents.fetch(ctx.Krt, pk)
		appendParent := func(pr parentInfo, parentRef parentReference) {
			// TODO: check reference allowed
			var hostname string
			if len(pr.Hostnames) > 0 {
				hostname = pr.Hostnames[0]
			}

			rpi := routeParentReference{
				InternalName:      pr.InternalName,
				InternalKind:      wellknown.GatewayGVK,
				Hostname:          hostname,
				OriginalReference: ref,
				ParentKey:         pk,
				ParentSection:     pr.SectionName,
			}
			parentRefs = append(parentRefs, rpi)
		}

		for _, parent := range currentParents {
			// Append all matches. Note we may be adding mismatch section or ports; this is handled later
			appendParent(parent, parentRef)
		}
	}
	// Ensure stable order
	slices.SortBy(parentRefs, func(a routeParentReference) string {
		return parentRefString(a.OriginalReference)
	})
	return parentRefs
}

func parentRefString(ref gwv1.ParentReference) string {
	return fmt.Sprintf("%s/%s/%s/%s/%d.%s",
		ptr.OrEmpty(ref.Group),
		ptr.OrEmpty(ref.Kind),
		ref.Name,
		ptr.OrEmpty(ref.SectionName),
		ptr.OrEmpty(ref.Port),
		ptr.OrEmpty(ref.Namespace))
}

func getCommonRouteInfo(spec any) ([]gwv1.ParentReference, []gwv1.Hostname, schema.GroupVersionKind) {
	switch t := spec.(type) {
	case *gwv1.HTTPRoute:
		return t.Spec.ParentRefs, t.Spec.Hostnames, wellknown.HTTPRouteGVK
	default:
		logger.Error("unknown type", "type", t)
		return nil, nil, schema.GroupVersionKind{}
	}
}

func toInternalParentReference(p gwv1.ParentReference, localNamespace string) (parentKey, error) {
	ref := normalizeReference(p.Group, p.Kind, wellknown.GatewayGVK)
	// TODO: check allowed parent references for GVK
	return parentKey{
		Kind: ref,
		Name: string(p.Name),
		// Unset namespace means "same namespace"
		Namespace: defaultString(p.Namespace, localNamespace),
	}, nil
}
