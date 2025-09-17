package translator

import (
	"fmt"

	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/maps"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
)

const (
	// Resource name format strings
	resourceNameFormat = "%s~%s"
)

// AgentGwXdsResources represents XDS resources for a single agent gateway
type AgentGwXdsResources struct {
	types.NamespacedName

	// Status Reports for this gateway
	Reports        reports.ReportMap
	AttachedRoutes map[string]uint

	// Resources config for gateway (Bind, Listener, Route)
	ResourceConfig envoycache.Resources

	// Address config (Services, Workloads)
	AddressConfig envoycache.Resources
}

// ResourceName needs to match agentgateway role configured in agentgateway
func (r AgentGwXdsResources) ResourceName() string {
	return fmt.Sprintf(resourceNameFormat, r.Namespace, r.Name)
}

func (r AgentGwXdsResources) Equals(in AgentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{reportMap: r.Reports, attachedRoutes: r.AttachedRoutes}.Equals(report{reportMap: in.Reports, attachedRoutes: in.AttachedRoutes}) &&
		r.ResourceConfig.Version == in.ResourceConfig.Version &&
		r.AddressConfig.Version == in.AddressConfig.Version
}

// ResourceWithCustomName is a wrapper for the envoy ResourceWithName type
type ResourceWithCustomName struct {
	proto.Message
	Name    string
	Version uint64
}

func (r ResourceWithCustomName) ResourceName() string {
	return r.Name
}

func (r ResourceWithCustomName) GetName() string {
	return r.Name
}

func (r ResourceWithCustomName) Equals(in ResourceWithCustomName) bool {
	return r.Version == in.Version
}

var _ envoytypes.ResourceWithName = ResourceWithCustomName{}

type report struct {
	// lower case so krt doesn't error in debug handler
	reportMap      reports.ReportMap
	attachedRoutes map[string]uint
}

// RouteReports contains all route-related Reports
type RouteReports struct {
	HTTPRoutes map[types.NamespacedName]*reports.RouteReport
	GRPCRoutes map[types.NamespacedName]*reports.RouteReport
	TCPRoutes  map[types.NamespacedName]*reports.RouteReport
	TLSRoutes  map[types.NamespacedName]*reports.RouteReport
}

func (r RouteReports) ResourceName() string {
	return "route-Reports"
}

func (r RouteReports) Equals(in RouteReports) bool {
	return maps.Equal(r.HTTPRoutes, in.HTTPRoutes) &&
		maps.Equal(r.GRPCRoutes, in.GRPCRoutes) &&
		maps.Equal(r.TCPRoutes, in.TCPRoutes) &&
		maps.Equal(r.TLSRoutes, in.TLSRoutes)
}

// ListenerSetReports contains all listener set Reports
type ListenerSetReports struct {
	Reports map[types.NamespacedName]*reports.ListenerSetReport
}

func (l ListenerSetReports) ResourceName() string {
	return "listenerset-Reports"
}

func (l ListenerSetReports) Equals(in ListenerSetReports) bool {
	return maps.Equal(l.Reports, in.Reports)
}

// GatewayReports contains gateway Reports along with attached routes information
type GatewayReports struct {
	Reports        map[types.NamespacedName]*reports.GatewayReport
	AttachedRoutes map[types.NamespacedName]map[string]uint
}

func (g GatewayReports) ResourceName() string {
	return "gateway-Reports"
}

func (g GatewayReports) Equals(in GatewayReports) bool {
	if !maps.Equal(g.Reports, in.Reports) {
		return false
	}

	// Compare AttachedRoutes manually since it contains nested maps
	if len(g.AttachedRoutes) != len(in.AttachedRoutes) {
		return false
	}
	for key, gRoutes := range g.AttachedRoutes {
		inRoutes, exists := in.AttachedRoutes[key]
		if !exists {
			return false
		}
		if !maps.Equal(gRoutes, inRoutes) {
			return false
		}
	}

	return true
}

func (r report) ResourceName() string {
	return "Report"
}

func (r report) Equals(in report) bool {
	if !maps.Equal(r.reportMap.Gateways, in.reportMap.Gateways) {
		return false
	}
	if !maps.Equal(r.reportMap.ListenerSets, in.reportMap.ListenerSets) {
		return false
	}
	if !maps.Equal(r.reportMap.HTTPRoutes, in.reportMap.HTTPRoutes) {
		return false
	}
	if !maps.Equal(r.reportMap.TCPRoutes, in.reportMap.TCPRoutes) {
		return false
	}
	if !maps.Equal(r.reportMap.TLSRoutes, in.reportMap.TLSRoutes) {
		return false
	}
	if !maps.Equal(r.reportMap.Policies, in.reportMap.Policies) {
		return false
	}
	if !maps.Equal(r.attachedRoutes, in.attachedRoutes) {
		return false
	}
	return true
}
