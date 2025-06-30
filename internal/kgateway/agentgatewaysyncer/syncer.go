package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/gateway"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"google.golang.org/protobuf/proto"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pkg/config/schema/gvr"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayalpha "sigs.k8s.io/gateway-api/apis/v1alpha2"
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
		controllerName: controllerName,
		xdsCache:       xdsCache,
		// mgr:            mgr,
		istioClient: client,
	}
}

type agentGwXdsResources struct {
	types.NamespacedName

	reports   reports.ReportMap
	Resources envoycache.Resources
	Addresses envoycache.Resources
}

func (r agentGwXdsResources) ResourceName() string {
	return fmt.Sprintf("%s~%s", r.Namespace, r.Name)
}

func (r agentGwXdsResources) Equals(in agentGwXdsResources) bool {
	return r.NamespacedName == in.NamespacedName &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.Resources.Version == in.Resources.Version &&
		r.Addresses.Version == in.Addresses.Version
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
	protocol string // currently only A2A and MCP
	// The listeners which are allowed to connect to the target.
	allowedListeners []string
}

func (r agentGwService) Equals(in agentGwService) bool {
	return r.ip == in.ip && r.port == in.port && r.path == in.path && r.protocol == in.protocol && slices.Equal(r.allowedListeners, in.allowedListeners)
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

func (s *AgentGwSyncer) Init(krtopts krtutil.KrtOptions) {
	logger.Debug("init agentgateway Syncer", "controllername", s.controllerName)

	inputs := gateway.Inputs{
		Namespaces: krt.NewInformer[*corev1.Namespace](s.istioClient),
		Secrets: krt.WrapClient[*corev1.Secret](
			kclient.NewFiltered[*corev1.Secret](s.istioClient, kubetypes.Filter{
				//FieldSelector: kubesecrets.SecretsFieldSelector,
				ObjectFilter: s.istioClient.ObjectFilter(),
			}),
		),
		Services: krt.WrapClient[*corev1.Service](
			kclient.NewFiltered[*corev1.Service](s.istioClient, kubetypes.Filter{ObjectFilter: s.istioClient.ObjectFilter()}),
		),
		GatewayClasses: buildClient[*gateway.GatewayClass](c, kc, gvr.GatewayClass, opts, "informer/GatewayClasses"),
		Gateways:       buildClient[*gateway.Gateway](c, kc, gvr.KubernetesGateway, opts, "informer/Gateways"),
		HTTPRoutes:     buildClient[*gateway.HTTPRoute](c, kc, gvr.HTTPRoute, opts, "informer/HTTPRoutes"),
		GRPCRoutes:     buildClient[*gatewayv1.GRPCRoute](c, kc, gvr.GRPCRoute, opts, "informer/GRPCRoutes"),

		ReferenceGrants: buildClient[*gateway.ReferenceGrant](c, kc, gvr.ReferenceGrant, opts, "informer/ReferenceGrants"),
		ServiceEntries:  buildClient[*networkingclient.ServiceEntry](c, kc, gvr.ServiceEntry, opts, "informer/ServiceEntries"),
		//InferencePools:  buildClient[*inf.InferencePool](c, kc, gvr.InferencePool, opts, "informer/InferencePools"),
	}
	if features.EnableAlphaGatewayAPI {
		inputs.TCPRoutes = buildClient[*gatewayalpha.TCPRoute](c, kc, gvr.TCPRoute, opts, "informer/TCPRoutes")
		inputs.TLSRoutes = buildClient[*gatewayalpha.TLSRoute](c, kc, gvr.TLSRoute, opts, "informer/TLSRoutes")
	} else {
		// If disabled, still build a collection but make it always empty
		inputs.TCPRoutes = krt.NewStaticCollection[*gatewayalpha.TCPRoute](nil, opts.WithName("disable/TCPRoutes")...)
		inputs.TLSRoutes = krt.NewStaticCollection[*gatewayalpha.TLSRoute](nil, opts.WithName("disable/TLSRoutes")...)
	}

	GatewayClasses := gateway.GatewayClassesCollection(inputs.GatewayClasses, opts)

	RefGrants := BuildReferenceGrants(ReferenceGrantsCollection(inputs.ReferenceGrants, opts))

	// Note: not fully complete until its join with route attachments to report attachedRoutes.
	// Do not register yet.
	Gateways := GatewayCollection(
		inputs.Gateways,
		GatewayClasses,
		inputs.Namespaces,
		RefGrants,
		inputs.Secrets,
		options.DomainSuffix,
		opts,
	)
	ports := krt.NewCollection(Gateways, func(ctx krt.HandlerContext, obj Gateway) *IndexObject[string, Gateway] {
		port := fmt.Sprint(obj.parentInfo.Port)
		return &IndexObject[string, Gateway]{
			Key:     port,
			Objects: []Gateway{obj},
		}
	}, opts.WithName("ports")...)

	Binds := krt.NewManyCollection(ports, func(ctx krt.HandlerContext, object IndexObject[string, Gateway]) []ADPResource {
		port, _ := strconv.Atoi(object.Key)
		uniq := sets.New[types.NamespacedName]()
		for _, gw := range object.Objects {
			uniq.Insert(types.NamespacedName{
				Namespace: gw.parent.Namespace,
				Name:      gw.parent.Name,
			})
		}
		return slices.Map(uniq.UnsortedList(), func(e types.NamespacedName) ADPResource {
			bind := Bind{
				Bind: &api.Bind{
					Key:  object.Key + "/" + e.String(),
					Port: uint32(port),
				},
			}
			return toResource(e, bind)
		})
	}, opts.WithName("Binds")...)

	Listeners := krt.NewCollection(Gateways, func(ctx krt.HandlerContext, obj gateway.Gateway) *gateway.ADPResource {
		l := &api.Listener{
			Key:         obj.ResourceName(),
			Name:        string(obj.parentInfo.SectionName),
			BindKey:     fmt.Sprint(obj.parentInfo.Port) + "/" + obj.parent.Namespace + "/" + obj.parent.Name,
			GatewayName: obj.parent.Namespace + "/" + obj.parent.Name,
			Hostname:    obj.parentInfo.OriginalHostname,
		}

		switch obj.parentInfo.Protocol {
		case gatewayv1.HTTPProtocolType:
			l.Protocol = api.Protocol_HTTP
		case gatewayv1.HTTPSProtocolType:
			l.Protocol = api.Protocol_HTTPS
			if obj.TLSInfo == nil {
				return nil
			}
			l.Tls = &api.TLSConfig{
				Cert:       obj.TLSInfo.Cert,
				PrivateKey: obj.TLSInfo.Key,
			}
		case gatewayv1.TLSProtocolType:
			l.Protocol = api.Protocol_TLS
			if obj.TLSInfo == nil {
				return nil
			}
			l.Tls = &api.TLSConfig{
				Cert:       obj.TLSInfo.Cert,
				PrivateKey: obj.TLSInfo.Key,
			}
		case gatewayv1.TCPProtocolType:
			l.Protocol = api.Protocol_TCP
		default:
			return nil
		}
		return toResourcep(types.NamespacedName{
			Namespace: obj.parent.Namespace,
			Name:      obj.parent.Name,
		}, gateway.ADPListener{l})
	}, krtopts.WithName("Listeners")...)

	routeParents := gateway.BuildRouteParents(Gateways)

	routeInputs := gateway.RouteContextInputs{
		Grants:         RefGrants,
		RouteParents:   routeParents,
		DomainSuffix:   options.DomainSuffix,
		Services:       inputs.Services,
		Namespaces:     inputs.Namespaces,
		ServiceEntries: inputs.ServiceEntries,
		InferencePools: inputs.InferencePools,
	}
	ADPRoutes := gateway.ADPRouteCollection(
		inputs.HTTPRoutes,
		routeInputs,
		opts,
	)

	s.waitForSync = []cache.InformerSynced{
		s.commonCols.HasSynced,
		gatewaysCol.HasSynced,
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
				Resources: r.Resources,
				Addresses: r.Addresses,
			}
			logger.Debug("setting xds snapshot", "resourcename", r.ResourceName())
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
	Resources  envoycache.Resources
	Addresses  envoycache.Resources
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
	case TargetTypeResourceUrl:
		return m.Resources.Items
	case TargetTypeAddressUrl:
		return m.Addresses.Items
	default:
		return nil
	}
}

func (m *agentGwSnapshot) GetVersion(typeURL string) string {
	switch typeURL {
	case TargetTypeResourceUrl:
		return m.Resources.Version
	case TargetTypeAddressUrl:
		return m.Resources.Version
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
		TargetTypeResourceUrl: m.Resources.Items,
		TargetTypeAddressUrl:  m.Addresses.Items,
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
