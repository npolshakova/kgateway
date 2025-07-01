package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"

	istio "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/util/protoconv"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
)

func toResourcep(gw types.NamespacedName, t any) *ADPResource {
	res := toResource(gw, t)
	return &res
}

func toResource(gw types.NamespacedName, t any) ADPResource {
	switch tt := t.(type) {
	case Bind:
		return ADPResource{Resource: &api.Resource{Kind: &api.Resource_Bind{tt.Bind}}, Gateway: gw}
	case ADPListener:
		return ADPResource{Resource: &api.Resource{Kind: &api.Resource_Listener{tt.Listener}}, Gateway: gw}
	case ADPRoute:
		return ADPResource{Resource: &api.Resource{Kind: &api.Resource_Route{tt.Route}}, Gateway: gw}
	}
	panic("unknown resource kind")
}

// TODO: we need some way to associate this to a specific instance of the proxy!!
type Bind struct {
	*api.Bind
}

func (g Bind) ResourceName() string {
	return g.Key
}

func (g Bind) Equals(other Bind) bool {
	return protoconv.Equals(g, other)
}

type ADPListener struct {
	*api.Listener
}

func (g ADPListener) ResourceName() string {
	return g.Key
}

func (g ADPListener) Equals(other ADPListener) bool {
	return protoconv.Equals(g, other)
}

type ADPRoute struct {
	*api.Route
}

func (g ADPRoute) ResourceName() string {
	return g.Key
}

func (g ADPRoute) Equals(other ADPRoute) bool {
	return protoconv.Equals(g, other)
}

type TLSInfo struct {
	Cert []byte
	Key  []byte `json:"-"`
}

type PortBindings struct {
	Gateway
	Port string
}

func (g PortBindings) ResourceName() string {
	return g.Gateway.Name
}

func (g PortBindings) Equals(other PortBindings) bool {
	return g.Gateway.Equals(other.Gateway) &&
		g.Port == other.Port
}

type Gateway struct {
	*Config
	parent     parentKey
	parentInfo parentInfo
	TLSInfo    *TLSInfo
	Valid      bool
}

func (g Gateway) ResourceName() string {
	return g.Config.Name
}

func (g Gateway) Equals(other Gateway) bool {
	return g.Config.Equals(other.Config) &&
		g.Valid == other.Valid // TODO: ok to ignore parent/parentInfo?
}

func GatewayCollection(
	agentGatewayClassName string,
	gateways krt.Collection[*gateway.Gateway],
	gatewayClasses krt.Collection[GatewayClass],
	namespaces krt.Collection[*corev1.Namespace],
	grants ReferenceGrants,
	secrets krt.Collection[*corev1.Secret],
	domainSuffix string,
	krtopts krtutil.KrtOptions,
) krt.Collection[Gateway] {
	gw := krt.NewManyCollection(gateways, func(ctx krt.HandlerContext, obj *gateway.Gateway) []Gateway {
		if string(obj.Spec.GatewayClassName) != agentGatewayClassName {
			return nil // ignore non agentgateway gws
		}

		var result []Gateway
		kgw := obj.Spec
		status := obj.Status.DeepCopy()
		class := fetchClass(ctx, gatewayClasses, kgw.GatewayClassName)
		if class == nil {
			return nil
		}
		controllerName := class.Controller
		var servers []*istio.Server

		// Extract the addresses. A gateway will bind to a specific Service
		gatewayServices, err := extractGatewayServices(domainSuffix, obj)
		if len(gatewayServices) == 0 && err != nil {
			// Short circuit if its a hard failure
			// TODO: log
			return nil
		}

		for i, l := range kgw.Listeners {
			server, tlsInfo, programmed := buildListener(ctx, secrets, grants, namespaces, obj, status, l, i, controllerName)

			servers = append(servers, server)
			meta := parentMeta(obj, &l.Name)
			// Each listener generates a Gateway with a single Server. This allows binding to a specific listener.
			gatewayConfig := Config{
				Meta: Meta{
					CreationTimestamp: obj.CreationTimestamp.Time,
					GroupVersionKind:  schema.GroupVersionKind{Group: wellknown.GatewayGroup, Kind: wellknown.GatewayKind},
					Name:              InternalGatewayName(obj.Name, string(l.Name)),
					Annotations:       meta,
					Namespace:         obj.Namespace,
					Domain:            domainSuffix,
				},
				// TODO: move away from istio gateway ir
				Spec: &istio.Gateway{
					Servers: []*istio.Server{server},
				},
			}

			allowed, _ := generateSupportedKinds(l)
			ref := parentKey{
				Kind:      wellknown.GatewayGVK,
				Name:      obj.Name,
				Namespace: obj.Namespace,
			}
			pri := parentInfo{
				InternalName:     obj.Namespace + "/" + gatewayConfig.Name,
				AllowedKinds:     allowed,
				Hostnames:        server.GetHosts(),
				OriginalHostname: string(ptr.OrEmpty(l.Hostname)),
				SectionName:      l.Name,
				Port:             l.Port,
				Protocol:         l.Protocol,
			}

			res := Gateway{
				Config:     &gatewayConfig,
				Valid:      programmed,
				TLSInfo:    tlsInfo,
				parent:     ref,
				parentInfo: pri,
			}
			result = append(result, res)
		}

		return result
	}, krtopts.ToOptions("KubernetesGateway")...)

	return gw
}

// RouteParents holds information about things routes can reference as parents.
type RouteParents struct {
	gateways     krt.Collection[Gateway]
	gatewayIndex krt.Index[parentKey, Gateway]
}

func (p RouteParents) fetch(ctx krt.HandlerContext, pk parentKey) []*parentInfo {
	return slices.Map(krt.Fetch(ctx, p.gateways, krt.FilterIndex(p.gatewayIndex, pk)), func(gw Gateway) *parentInfo {
		return &gw.parentInfo
	})
}

func BuildRouteParents(
	gateways krt.Collection[Gateway],
) RouteParents {
	idx := krt.NewIndex(gateways, func(o Gateway) []parentKey {
		return []parentKey{o.parent}
	})
	return RouteParents{
		gateways:     gateways,
		gatewayIndex: idx,
	}
}

// InternalGatewayName returns the name of the internal Istio Gateway corresponding to the
// specified gateway-api gateway and listener.
func InternalGatewayName(gwName, lName string) string {
	return fmt.Sprintf("%s-%s-%s", gwName, AgentgatewayName, lName)
}
