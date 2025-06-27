package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

var logger = logging.New("agentgateway/syncer")

// AgentGwSyncer synchronizes Kubernetes Gateway API resources with xDS for agentgateway proxies.
// It watches Gateway resources with the agentgateway class and translates them to agentgateway configuration.
type AgentGwSyncer struct {
	commonCols     *common.CommonCollections
	controllerName string
	xDS            krt.Collection[agentGwXdsResources]
	xdsCache       envoycache.SnapshotCache
	istioClient    kube.Client
	clusterId      string

	waitForSync []cache.InformerSynced
}

func NewAgentGwSyncer(
	ctx context.Context,
	controllerName string,
	mgr manager.Manager,
	client kube.Client,
	commonCols *common.CommonCollections,
	xdsCache envoycache.SnapshotCache,
	clusterId string,
) *AgentGwSyncer {
	// TODO: register types (auth, policy, etc.) if necessary
	return &AgentGwSyncer{
		commonCols:     commonCols,
		controllerName: controllerName,
		xdsCache:       xdsCache,
		// mgr:            mgr,
		istioClient: client,
		clusterId:   clusterId,
	}
}

type agentGwXdsResources struct {
	types.NamespacedName

	reports reports.ReportMap
	// ResourcesConfig (Bind, Listener, Route)
	ResourcesConfig envoycache.Resources
	// WorkloadConfig (Services, Workloads)
	WorkloadConfig envoycache.Resources
}

// Needs to match agentgateway role configured in client.rs (https://github.com/agentgateway/agentgateway/blob/main/crates/agentgateway/src/xds/client.rs)
func (r agentGwXdsResources) ResourceName() string {
	return fmt.Sprintf("%s~%s", r.Namespace, r.Name)
}

func (r agentGwXdsResources) Equals(in agentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.ResourcesConfig.Version == in.ResourcesConfig.Version &&
		r.WorkloadConfig.Version == in.WorkloadConfig.Version
}

type envoyResourceWithName struct {
	inner   envoytypes.ResourceWithName
	version uint64
}

func (r envoyResourceWithName) ResourceName() string {
	return r.inner.GetName()
}

func (r envoyResourceWithName) Equals(in envoyResourceWithName) bool {
	return r.version == in.version
}

type envoyResourceWithCustomName struct {
	proto.Message
	Name    string
	version uint64
}

func (r envoyResourceWithCustomName) ResourceName() string {
	return r.Name
}

func (r envoyResourceWithCustomName) GetName() string {
	return r.Name
}

func (r envoyResourceWithCustomName) Equals(in envoyResourceWithCustomName) bool {
	return r.version == in.version
}

var _ envoytypes.ResourceWithName = envoyResourceWithCustomName{}

type report struct {
	// lower case so krt doesn't error in debug handler
	reportMap reports.ReportMap
}

func (r report) ResourceName() string {
	return "report"
}

func (r report) Equals(in report) bool {
	return maps.Equal(r.reportMap.Gateways, in.reportMap.Gateways) &&
		maps.Equal(r.reportMap.HTTPRoutes, in.reportMap.HTTPRoutes) &&
		maps.Equal(r.reportMap.TCPRoutes, in.reportMap.TCPRoutes)
}

func (s *AgentGwSyncer) Init(krtopts krtutil.KrtOptions) {
	logger.Debug("init agentgateway Syncer", "controllername", s.controllerName)

	gatewaysCol := krt.NewCollection(s.commonCols.GatewayIndex.Gateways, func(kctx krt.HandlerContext, gw ir.Gateway) *ir.Gateway {
		if gw.Obj.Spec.GatewayClassName != wellknown.AgentGatewayClassName {
			return nil
		}
		return &gw
	}, krtopts.ToOptions("agentgateway")...)

	// these are workloadapi-style services combined from kube services (todo: support services entries, backend types, etc.)
	domainSuffix := "cluster.local" // todo make configurable
	// k8s service -> service info ir
	agentGwServices := buildServicesCollection(s.commonCols.Services, domainSuffix)
	// Gateway resources -> agentgateway resources
	agentGwResources := AgentGatewayCollection(gatewaysCol, agentGwServices, krtopts)

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

	// Create raw HTTPRoute collection
	filter := kclient.Filter{ObjectFilter: s.istioClient.ObjectFilter()}
	httpRoutes := krt.WrapClient(kclient.NewFiltered[*gwv1.HTTPRoute](s.istioClient, filter), krtopts.ToOptions("HTTPRoute")...)

	// Create route context inputs
	routeContextInputs := RouteContextInputs{
		AgentGatewayResource: agentGwResources,
		Services:             agentGwServices,
		Workloads:            agentGwWorkloads,
	}

	// Call AgentHTTPRouteCollection to get the routes
	httpRouteResult := AgentHTTPRouteCollection(httpRoutes, routeContextInputs, krtopts)

	// translate gateways to xds
	s.xDS = krt.NewCollection(agentGwResources, func(kctx krt.HandlerContext, gwResource AgentGatewayResource) *agentGwXdsResources {
		if !gwResource.Valid {
			// todo: error?
			return &agentGwXdsResources{
				NamespacedName:  gwResource.NamespacedName,
				ResourcesConfig: envoycache.NewResources("0", []envoytypes.Resource{}),
				WorkloadConfig:  envoycache.NewResources("0", []envoytypes.Resource{}),
			}
		}

		// Create resources for the new API format
		var resources []envoytypes.Resource
		var serviceResources []envoytypes.Resource
		var workloadResources []envoytypes.Resource

		// Add bind resource
		if gwResource.Bind != nil {
			bindResource := &api.Resource{
				Kind: &api.Resource_Bind{
					Bind: gwResource.Bind.Bind,
				},
			}
			resources = append(resources, &envoyResourceWithCustomName{
				Message: bindResource,
				Name:    gwResource.Bind.Key,
				version: utils.HashProto(bindResource),
			})
		}

		// Add listener resource
		if gwResource.Listener != nil {
			listenerResource := &api.Resource{
				Kind: &api.Resource_Listener{
					Listener: gwResource.Listener.Listener,
				},
			}
			resources = append(resources, &envoyResourceWithCustomName{
				Message: listenerResource,
				Name:    gwResource.Listener.Key,
				version: utils.HashProto(listenerResource),
			})
		}

		httproutes := krt.Fetch(kctx, httpRouteResult.Routes)
		// Add route resources
		for _, route := range httproutes {
			if !route.Valid {
				// todo: error
				continue
			}

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

		// Create workload address resources for agentgateway workload API
		addresses := krt.Fetch(kctx, agentGwServices)
		for _, addr := range addresses {
			// Create Service resource
			serviceResource := addr.Service

			// Create Address resource wrapping the Service
			addressResource := &api.Address{
				Type: &api.Address_Service{
					Service: serviceResource,
				},
			}

			serviceResources = append(serviceResources, &envoyResourceWithCustomName{
				Message: addressResource,
				Name:    addr.ResourceName(),
				version: utils.HashProto(addressResource),
			})
		}

		// TODO: use index?
		//workloadServiceIndex := krt.NewIndex[string, WorkloadInfo](workloads, func(o WorkloadInfo) []string {
		//	var svcs []string
		//	for svcName := range o.Workload.Services {
		//		svcs = append(svcs, svcName)
		//	}
		//	return svcs
		//})

		workloads := krt.Fetch(kctx, agentGwWorkloads)
		for _, workload := range workloads {
			// Create Address resource wrapping the Service
			addressResource := &api.Address{
				Type: &api.Address_Workload{
					Workload: workload.Workload,
				},
			}

			workloadResources = append(workloadResources, &envoyResourceWithCustomName{
				Message: addressResource,
				Name:    workload.ResourceName(),
				version: utils.HashProto(addressResource),
			})
		}

		// Create the resource wrapper
		var version uint64
		for _, res := range resources {
			version ^= res.(*envoyResourceWithCustomName).version
		}

		var serviceVersion uint64
		for _, res := range serviceResources {
			serviceVersion ^= res.(*envoyResourceWithCustomName).version
		}

		var workloadVersion uint64
		for _, res := range workloadResources {
			workloadVersion ^= res.(*envoyResourceWithCustomName).version
		}

		// Combine service and workload versions for WorkloadConfig
		combinedServiceVersion := serviceVersion ^ workloadVersion

		result := &agentGwXdsResources{
			NamespacedName:  gwResource.NamespacedName,
			ResourcesConfig: envoycache.NewResources(fmt.Sprintf("%d", version), resources),
			WorkloadConfig:  envoycache.NewResources(fmt.Sprintf("%d", combinedServiceVersion), append(serviceResources, workloadResources...)),
		}
		logger.Debug("created XDS resources for with ID", "gwname", gwResource.Name, "resourceid", result.ResourceName())
		return result
	}, krtopts.ToOptions("agentgateway-xds")...)

	s.waitForSync = []cache.InformerSynced{
		s.commonCols.HasSynced,
		gatewaysCol.HasSynced,
		agentGwServices.HasSynced,
		agentGwWorkloads.HasSynced,
		agentGwResources.HasSynced,
		httpRouteResult.Routes.HasSynced,
		s.xDS.HasSynced,
	}
}

func (s *AgentGwSyncer) Start(ctx context.Context) error {
	logger.Info("starting agentgateway Syncer", "controllername", s.controllerName)
	logger.Info("waiting for agentgateway cache to sync")

	// Wait for cache to sync
	if !kube.WaitForCacheSync("agentgateway syncer", ctx.Done(), s.waitForSync...) {
		return fmt.Errorf("agentgateway syncer waiting for cache to sync failed")
	}

	s.xDS.RegisterBatch(func(events []krt.Event[agentGwXdsResources], _ bool) {
		for _, e := range events {
			r := e.Latest()
			if e.Event == controllers.EventDelete {
				s.xdsCache.ClearSnapshot(r.ResourceName())
				continue
			}
			snapshot := &agentGwSnapshot{
				ResourceConfig: r.ResourcesConfig,
				WorkloadConfig: r.WorkloadConfig,
			}
			logger.Debug("setting xds snapshot", "resourceName", r.ResourceName())
			logger.Debug("snapshot config", "resourceSnapshot", snapshot.ResourceConfig, "workloadSnapshot", snapshot.WorkloadConfig)
			err := s.xdsCache.SetSnapshot(ctx, r.ResourceName(), snapshot)
			if err != nil {
				logger.Error("failed to set xds snapshot", "resourcename", r.ResourceName(), "error", err.Error())
				continue
			}
		}
	}, true)

	return nil
}

type agentGwSnapshot struct {
	ResourceConfig envoycache.Resources
	WorkloadConfig envoycache.Resources
	VersionMap     map[string]map[string]string
}

func (m *agentGwSnapshot) GetResources(typeURL string) map[string]envoytypes.Resource {
	resources := m.GetResourcesAndTTL(typeURL)
	result := make(map[string]envoytypes.Resource, len(resources))
	for k, v := range resources {
		result[k] = v.Resource
	}
	return result
}

func (m *agentGwSnapshot) GetResourcesAndTTL(typeURL string) map[string]envoytypes.ResourceWithTTL {
	switch typeURL {
	case TargetTypeResourceUrl:
		return m.ResourceConfig.Items
	case TargetTypeAddressUrl:
		return m.WorkloadConfig.Items
	default:
		return nil
	}
}

func (m *agentGwSnapshot) GetVersion(typeURL string) string {
	switch typeURL {
	case TargetTypeResourceUrl:
		return m.ResourceConfig.Version
	case TargetTypeAddressUrl:
		return m.WorkloadConfig.Version
	default:
		return ""
	}
}

func (m *agentGwSnapshot) ConstructVersionMap() error {
	if m == nil {
		return fmt.Errorf("missing snapshot")
	}
	if m.VersionMap != nil {
		return nil
	}

	m.VersionMap = make(map[string]map[string]string)
	resources := map[string]map[string]envoytypes.ResourceWithTTL{
		TargetTypeResourceUrl: m.ResourceConfig.Items,
		TargetTypeAddressUrl:  m.WorkloadConfig.Items,
	}

	for typeUrl, items := range resources {
		inner := make(map[string]string, len(items))
		for _, r := range items {
			marshaled, err := envoycache.MarshalResource(r.Resource)
			if err != nil {
				return err
			}
			v := envoycache.HashResource(marshaled)
			if v == "" {
				return fmt.Errorf("failed to build resource version")
			}
			inner[envoycache.GetResourceName(r.Resource)] = v
		}
		m.VersionMap[typeUrl] = inner
	}
	return nil
}

func (m *agentGwSnapshot) GetVersionMap(typeURL string) map[string]string {
	return m.VersionMap[typeURL]
}

var _ envoycache.ResourceSnapshot = &agentGwSnapshot{}
