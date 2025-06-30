package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"

	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
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
	httprouteCol   krt.Collection[*gwv1.HTTPRoute] // TODO: index this

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
	httprouteCol krt.Collection[*gwv1.HTTPRoute],
) *AgentGwSyncer {
	// TODO: register types (auth, policy, etc.) if necessary
	return &AgentGwSyncer{
		commonCols:     commonCols,
		controllerName: controllerName,
		xdsCache:       xdsCache,
		// mgr:            mgr,
		istioClient:  client,
		clusterId:    clusterId,
		httprouteCol: httprouteCol,
	}
}

type agentGwXdsResources struct {
	types.NamespacedName

	reports reports.ReportMap
	// Resources config for gw (Bind, Listener. Route)
	ResourceConfig envoycache.Resources
	// Address config (Services, Workloads)
	AddressConfig envoycache.Resources
}

// Needs to match agentgateway role configured in client.rs (https://github.com/agentgateway/agentgateway/blob/main/crates/agentgateway/src/xds/client.rs)
func (r agentGwXdsResources) ResourceName() string {
	return fmt.Sprintf("%s~%s", r.Namespace, r.Name)
}

func (r agentGwXdsResources) Equals(in agentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.ResourceConfig.Version == in.ResourceConfig.Version &&
		r.AddressConfig.Version == in.AddressConfig.Version
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

	domainSuffix := "cluster.local" // todo make configurable

	// translate gateways to xds
	s.xDS = s.translateGateway(krtopts, domainSuffix)

	s.waitForSync = append(s.waitForSync, []cache.InformerSynced{
		s.commonCols.HasSynced,
		s.xDS.HasSynced,
	}...)
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
				ResourceConfig: r.ResourceConfig,
				AddressConfig:  r.AddressConfig,
			}
			logger.Debug("setting xds snapshot", "resourceName", r.ResourceName())
			logger.Debug("snapshot config", "resourceSnapshot", snapshot.ResourceConfig, "workloadSnapshot", snapshot.AddressConfig)
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
	AddressConfig  envoycache.Resources
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
		return m.AddressConfig.Items
	default:
		return nil
	}
}

func (m *agentGwSnapshot) GetVersion(typeURL string) string {
	switch typeURL {
	case TargetTypeResourceUrl:
		return m.ResourceConfig.Version
	case TargetTypeAddressUrl:
		return m.AddressConfig.Version
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
		TargetTypeAddressUrl:  m.AddressConfig.Items,
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
