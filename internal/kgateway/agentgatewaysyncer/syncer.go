package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"strings"

	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/a2a"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/mcp"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/query"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	gwtranslator "github.com/kgateway-dev/kgateway/v2/internal/kgateway/translator/gateway"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/xds"
)

type AgentGwSyncer struct {
	commonCols     *common.CommonCollections
	translator     *agentGwTranslator
	controllerName string
	xDS            krt.Collection[agentGwXdsResources]
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
		translator:     newTranslator(ctx, commonCols),
		controllerName: controllerName,
		xdsCache:       xdsCache,
		//mgr:            mgr,
		istioClient: client,
	}
}

type agentGwXdsResources struct {
	types.NamespacedName

	reports            reports.ReportMap
	AgentGwA2AServices envoycache.Resources
	AgentGwMcpServices envoycache.Resources
	Listeners          envoycache.Resources
}

func (r agentGwXdsResources) ResourceName() string {
	return xds.OwnerNamespaceNameID(OwnerNodeId, r.Namespace, r.Name)
}

func (r agentGwXdsResources) Equals(in agentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.AgentGwA2AServices.Version == in.AgentGwA2AServices.Version &&
		r.AgentGwMcpServices.Version == in.AgentGwMcpServices.Version
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
	ip   string
	port int
	path string
}

func (r agentGwService) Equals(in agentGwService) bool {
	return r.ip == in.ip && r.port == in.port && r.path == in.path
}

type agentGwTranslator struct {
	gwtranslator extensionsplug.KGwTranslator
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

	gateways := krt.NewCollection(s.commonCols.GatewayIndex.Gateways, func(kctx krt.HandlerContext, gw ir.Gateway) *ir.Gateway {
		if gw.Obj.Spec.GatewayClassName != wellknown.AgentGatewayClassName {
			return nil
		}
		return &gw
	}, krtopts.ToOptions("agentgateway")...)

	a2aServices := krt.NewManyCollection(s.commonCols.Services, func(kctx krt.HandlerContext, s *corev1.Service) []agentGwService {
		var result []agentGwService
		for _, port := range s.Spec.Ports {
			if port.AppProtocol == nil {
				continue
			}
			appProtocol := *port.AppProtocol
			if appProtocol != "kgateway.dev/a2a" {
				continue
			}
			if s.Spec.ClusterIP == "" && s.Spec.ExternalName == "" {
				continue
			}
			addr := s.Spec.ClusterIP
			if addr == "" {
				addr = s.Spec.ExternalName
			}
			path := ""
			if appProtocol == A2AProtocol {
				path = s.Annotations[A2APathAnnotation]
			}
			logger.Debugf("Found a2a agent gateway service %s/%s", s.Namespace, s.Name)
			result = append(result, agentGwService{
				Named: krt.Named{
					Name:      s.Name,
					Namespace: s.Namespace,
				},
				ip:   addr,
				port: int(port.Port),
				path: path,
			})
		}
		return result
	}, krtopts.ToOptions("agentGwA2AService")...)

	mcpServices := krt.NewManyCollection(s.commonCols.Services, func(kctx krt.HandlerContext, s *corev1.Service) []agentGwService {
		var result []agentGwService
		for _, port := range s.Spec.Ports {
			if port.AppProtocol == nil {
				continue
			}
			appProtocol := *port.AppProtocol
			if appProtocol != "kgateway.dev/mcp" {
				continue
			}
			if s.Spec.ClusterIP == "" && s.Spec.ExternalName == "" {
				continue
			}
			addr := s.Spec.ClusterIP
			if addr == "" {
				addr = s.Spec.ExternalName
			}
			path := ""
			if appProtocol == MCPProtocol {
				path = s.Annotations[MCPPathAnnotation]
			}
			logger.Debugf("Found mcp agent gateway service %s/%s", s.Namespace, s.Name)
			result = append(result, agentGwService{
				Named: krt.Named{
					Name:      s.Name,
					Namespace: s.Namespace,
				},
				ip:   addr,
				port: int(port.Port),
				path: path,
			})
		}
		return result
	}, krtopts.ToOptions("agentGwMcpService")...)

	xdsA2AServices := krt.NewCollection(a2aServices, func(kctx krt.HandlerContext, s agentGwService) *envoyResourceWithName {
		t := &a2a.Target{
			Name: getTargetName(s.ResourceName()),
			Host: s.ip,
			Port: uint32(s.port),
			Path: s.path,
		}
		return &envoyResourceWithName{inner: t, version: utils.HashProto(t)}
	}, krtopts.ToOptions("a2a-target-xds")...)

	xdsMcpServices := krt.NewCollection(mcpServices, func(kctx krt.HandlerContext, s agentGwService) *envoyResourceWithName {
		t := &mcp.Target{
			// Note: No slashes allowed here (must match ^[a-zA-Z0-9-]+$)
			Name: getTargetName(s.ResourceName()),
			Target: &mcp.Target_Sse{
				Sse: &mcp.Target_SseTarget{
					Host: s.ip,
					Port: uint32(s.port),
					Path: s.path,
				},
			},
		}
		return &envoyResourceWithName{inner: t, version: utils.HashProto(t)}
	}, krtopts.ToOptions("mcp-target-xds")...)

	// translate gateways to xds
	s.xDS = krt.NewCollection(gateways, func(kctx krt.HandlerContext, gw ir.Gateway) *agentGwXdsResources {
		// listeners for the agent gateway
		agwListeners := make([]envoytypes.Resource, 0, len(gw.Listeners))
		var listenerVersion uint64
		var listener *Listener
		// TODO: Use AllowedRoute to filter namespace -> listener names
		var a2aListeners, mcpListeners []string
		for _, gwListener := range gw.Listeners {
			if string(gwListener.Protocol) == A2AProtocol {
				listener = &Listener{
					Name:     string(gwListener.Name),
					Protocol: Listener_A2A,
					// TODO: Add suppose for stdio listener
					Listener: &Listener_Sse{
						Sse: &SseListener{
							Address: "[::]",
							Port:    uint32(gwListener.Port),
						},
					},
				}

				a2aListeners = append(a2aListeners, string(gwListener.Name))
			} else if string(gwListener.Protocol) == MCPProtocol {
				listener = &Listener{
					Name:     string(gwListener.Name),
					Protocol: Listener_MCP,
					// TODO: Add suppose for stdio listener
					Listener: &Listener_Sse{
						Sse: &SseListener{
							Address: "[::]",
							Port:    uint32(gwListener.Port),
						},
					},
				}
				mcpListeners = append(mcpListeners, string(gwListener.Name))
			} else {
				// Not a valid protocol for Agent Gateway
				continue
			}
			// Update listenerVersion to be the result
			listenerVersion ^= utils.HashProto(listener)
			agwListeners = append(agwListeners, listener)
		}

		// a2a services
		a2aServiceResources := krt.Fetch(kctx, xdsA2AServices)
		logger.Debugf("Found %d A2A resources for gateway %s/%s", len(a2aServiceResources), gw.Namespace, gw.Name)
		a2aResources := make([]envoytypes.Resource, len(a2aServiceResources))
		var a2aVersion uint64
		for i, res := range a2aServiceResources {
			a2aVersion ^= res.version
			target := res.inner.(*a2a.Target)
			target.Listeners = a2aListeners
			a2aResources[i] = target
		}
		// mcp services
		mcpServiceResources := krt.Fetch(kctx, xdsMcpServices)
		logger.Debugf("Found %d MCP resources for gateway %s/%s", len(mcpServiceResources), gw.Namespace, gw.Name)
		mcpResources := make([]envoytypes.Resource, len(mcpServiceResources))
		var mcpVersion uint64
		for i, res := range mcpServiceResources {
			mcpVersion ^= res.version
			target := res.inner.(*mcp.Target)
			target.Listeners = mcpListeners
			mcpResources[i] = target
		}
		result := &agentGwXdsResources{
			NamespacedName:     types.NamespacedName{Namespace: gw.Namespace, Name: gw.Name},
			AgentGwA2AServices: envoycache.NewResources(fmt.Sprintf("%d", a2aVersion), a2aResources),
			AgentGwMcpServices: envoycache.NewResources(fmt.Sprintf("%d", mcpVersion), mcpResources),
			Listeners:          envoycache.NewResources(fmt.Sprintf("%d", listenerVersion), agwListeners),
		}
		logger.Debugf("Created XDS resources for %s with ID %s", gw.Name, result.ResourceName())
		return result
	}, krtopts.ToOptions("agentgateway-xds")...)

	s.waitForSync = []cache.InformerSynced{
		s.commonCols.HasSynced,
		xdsA2AServices.HasSynced,
		xdsMcpServices.HasSynced,
		gateways.HasSynced,
		a2aServices.HasSynced,
		mcpServices.HasSynced,
		s.xDS.HasSynced,
	}
}

func (s *AgentGwSyncer) Start(ctx context.Context) error {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("starting %s Agent Gateway Syncer", s.controllerName)
	logger.Infof("waiting for Agent Gateway cache to sync")
	kube.WaitForCacheSync("Agent Gateway syncer", ctx.Done(), s.waitForSync...)

	s.xDS.RegisterBatch(func(events []krt.Event[agentGwXdsResources], _ bool) {
		for _, e := range events {
			if e.Event == controllers.EventDelete {
				// TODO: do we need to handle deletes?
				continue
			}
			r := e.Latest()
			snapshot := &agentGwSnapshot{
				AgentGwA2AServices: r.AgentGwA2AServices,
				AgentGwMcpServices: r.AgentGwMcpServices,
				Listeners:          r.Listeners,
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
	AgentGwA2AServices envoycache.Resources
	AgentGwMcpServices envoycache.Resources
	Listeners          envoycache.Resources
	VersionMap         map[string]map[string]string
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
		return m.AgentGwA2AServices.Items
	case TargetTypeMcpUrl:
		return m.AgentGwMcpServices.Items
	case TargetTypeListenerUrl:
		return m.Listeners.Items
	default:
		return nil
	}
}

func (m *agentGwSnapshot) GetVersion(typeURL string) string {
	switch typeURL {
	case TargetTypeA2AUrl:
		return m.AgentGwA2AServices.Version
	case TargetTypeMcpUrl:
		return m.AgentGwMcpServices.Version
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
		TargetTypeA2AUrl: m.AgentGwA2AServices.Items,
		TargetTypeMcpUrl: m.AgentGwMcpServices.Items,
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

func newTranslator(
	ctx context.Context,
	commonCols *common.CommonCollections,
) *agentGwTranslator {
	return &agentGwTranslator{
		gwtranslator: gwtranslator.NewTranslator(query.NewData(commonCols)),
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
