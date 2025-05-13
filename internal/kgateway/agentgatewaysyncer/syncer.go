package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"strings"

	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/xds"
)

type AgentGwSyncer struct {
	commonCols     *common.CommonCollections
	translator     *AgentGwTranslator
	controllerName string
	xDS            krt.Collection[AgentGwXdsResources]
	xdsCache       envoycache.SnapshotCache
	istioClient    kube.Client

	waitForSync []cache.InformerSynced
}

func NewAgentGwSyncer(
	ctx context.Context,
	controllerName string,
	mgr manager.Manager,
	client kube.Client,
	commonCols *common.CommonCollections,
	xdsCache envoycache.SnapshotCache,
) *AgentGwSyncer {
	// TODO: register types (auth, policy, etc.) if necessary

	return &AgentGwSyncer{
		commonCols:     commonCols,
		translator:     NewTranslator(ctx, commonCols),
		controllerName: controllerName,
		xdsCache:       xdsCache,
		//mgr:            mgr,
		istioClient: client,
	}
}

type AgentGwXdsResources struct {
	types.NamespacedName

	reports    reports.ReportMap
	A2ATargets envoycache.Resources
	MCPTargets envoycache.Resources
	Listeners  envoycache.Resources
}

func (r AgentGwXdsResources) ResourceName() string {
	return xds.OwnerNamespaceNameID(OwnerNodeId, r.Namespace, r.Name)
}

func (r AgentGwXdsResources) Equals(in AgentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.A2ATargets.Version == in.A2ATargets.Version &&
		r.MCPTargets.Version == in.MCPTargets.Version
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

type agentGwService struct {
	krt.Named
	ip       string
	port     int
	path     string
	protocol string // currently only A2A or MCP
	// The listeners which are allowed to connect to the target.
	allowedListener []string
}

func (r agentGwService) Equals(in agentGwService) bool {
	return r.ip == in.ip && r.port == in.port && r.path == in.path
}

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

func (s *AgentGwSyncer) Init(ctx context.Context, krtopts krtutil.KrtOptions) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Debugf("Init %s Agent Gateway Syncer", s.controllerName)

	// TODO: convert auth to rbac json config for agent gateways

	s.xDS = krt.NewCollection(s.commonCols.GatewayIndex.Gateways, func(kctx krt.HandlerContext, gw ir.Gateway) *AgentGwXdsResources {
		// skip agentgateway proxies as they are not envoy-based gateways
		if gw.Obj.Spec.GatewayClassName != wellknown.AgentGatewayClassName {
			logger.Debugf("skipping agentgateway proxy sync for %s.%s", gw.Obj.Name, gw.Obj.Namespace)
			return nil
		}

		logger.Debugf("building proxy for kube gw %s version %s", client.ObjectKeyFromObject(gw.Obj), gw.Obj.GetResourceVersion())

		xdsSnap, rm := s.translator.TranslateGateway(kctx, ctx, gw)
		if xdsSnap == nil {
			return nil
		}

		return toResources(gw, *xdsSnap, rm)
	}, krtopts.ToOptions("agentgateway-xds")...)

	s.waitForSync = []cache.InformerSynced{
		s.commonCols.HasSynced,
		s.translator.HasSynced,
		s.xDS.HasSynced,
	}
}

func (s *AgentGwSyncer) Start(ctx context.Context) error {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("starting %s Agent Gateway Syncer", s.controllerName)
	logger.Infof("waiting for Agent Gateway cache to sync")
	kube.WaitForCacheSync("Agent Gateway syncer", ctx.Done(), s.waitForSync...)

	s.xDS.RegisterBatch(func(events []krt.Event[AgentGwXdsResources], _ bool) {
		for _, e := range events {
			if e.Event == controllers.EventDelete {
				// TODO: do we need to handle deletes?
				continue
			}
			r := e.Latest()
			snapshot := &agentGwSnapshot{
				A2ATargets: r.A2ATargets,
				MCPTargets: r.MCPTargets,
				Listeners:  r.Listeners,
			}
			logger.Debugf("setting xds snapshot for %s", r.ResourceName())
			err := s.xdsCache.SetSnapshot(ctx, r.ResourceName(), snapshot)
			if err != nil {
				logger.Errorf("failed to set xds snapshot for %s: %v", r.ResourceName(), err)
				continue
			}
		}
	}, true)

	return nil
}

type agentGwSnapshot struct {
	A2ATargets envoycache.Resources
	MCPTargets envoycache.Resources
	Listeners  envoycache.Resources
	VersionMap map[string]map[string]string
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
	case TargetTypeA2AUrl:
		return m.A2ATargets.Items
	case TargetTypeMcpUrl:
		return m.MCPTargets.Items
	case TargetTypeListenerUrl:
		return m.Listeners.Items
	default:
		return nil
	}
}

func (m *agentGwSnapshot) GetVersion(typeURL string) string {
	switch typeURL {
	case TargetTypeA2AUrl:
		return m.A2ATargets.Version
	case TargetTypeMcpUrl:
		return m.MCPTargets.Version
	case TargetTypeListenerUrl:
		return m.Listeners.Version
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
		TargetTypeA2AUrl: m.A2ATargets.Items,
		TargetTypeMcpUrl: m.MCPTargets.Items,
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

func translateAgentService(svc *corev1.Service, allowedListeners []string) []agentGwService {
	var svcs []agentGwService
	for _, port := range svc.Spec.Ports {
		if port.AppProtocol == nil {
			continue
		}
		appProtocol := *port.AppProtocol
		if svc.Spec.ClusterIP == "" && svc.Spec.ExternalName == "" {
			continue
		}
		addr := svc.Spec.ClusterIP
		if addr == "" {
			addr = svc.Spec.ExternalName
		}
		path := ""
		if appProtocol == A2AProtocol {
			path = svc.Annotations[A2APathAnnotation]
		} else if appProtocol == MCPProtocol {
			path = svc.Annotations[MCPPathAnnotation]
		}

		svcs = append(svcs, agentGwService{
			Named: krt.Named{
				Name:      svc.Name,
				Namespace: svc.Namespace,
			},
			ip:              addr,
			port:            int(port.Port),
			path:            path,
			protocol:        appProtocol,
			allowedListener: allowedListeners,
		})
	}
	return svcs
}

func toResources(gw ir.Gateway, xdsSnap AgentGatewayTranslationResult, r reports.ReportMap) *AgentGwXdsResources {
	return &AgentGwXdsResources{
		NamespacedName: types.NamespacedName{
			Namespace: gw.Obj.GetNamespace(),
			Name:      gw.Obj.GetName(),
		},
		reports:    r,
		Listeners:  sliceToResources(xdsSnap.Listeners),
		A2ATargets: sliceToResources(xdsSnap.A2ATargets),
		MCPTargets: sliceToResources(xdsSnap.McpTargets),
	}
}

// TODO: move to shared utils
func sliceToResourcesHash[T proto.Message](slice []T) ([]envoytypes.ResourceWithTTL, uint64) {
	var slicePb []envoytypes.ResourceWithTTL
	var resourcesHash uint64
	for _, r := range slice {
		var m proto.Message = r
		hash := utils.HashProto(r)
		slicePb = append(slicePb, envoytypes.ResourceWithTTL{Resource: m})
		resourcesHash ^= hash
	}

	return slicePb, resourcesHash
}

func sliceToResources[T proto.Message](slice []T) envoycache.Resources {
	r, h := sliceToResourcesHash(slice)
	return envoycache.NewResourcesWithTTL(fmt.Sprintf("%d", h), r)
}
