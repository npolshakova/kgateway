package agentgatewaysyncer

import (
	"context"
	"fmt"
	"maps"
	"strconv"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
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
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayalpha "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gateway "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

var logger = logging.New("agentgateway/syncer")

// AgentGwSyncer synchronizes Kubernetes Gateway API resources with xDS for agentgateway proxies.
// It watches Gateway resources with the agentgateway class and translates them to agentgateway configuration.
type AgentGwSyncer struct {
	commonCols     *common.CommonCollections
	controllerName string
	resourcexDS    krt.Collection[ADPCacheResource]
	addressxDS     krt.Collection[ADPCacheAddress]
	xdsCache       envoycache.SnapshotCache
	client         kube.Client
	domainSuffix   string

	waitForSync []cache.InformerSynced
}

func NewAgentGwSyncer(
	ctx context.Context,
	controllerName string,
	mgr manager.Manager,
	client kube.Client,
	commonCols *common.CommonCollections,
	xdsCache envoycache.SnapshotCache,
	domainSuffix string,
) *AgentGwSyncer {
	// TODO: register types (auth, policy, etc.) if necessary
	return &AgentGwSyncer{
		commonCols:     commonCols,
		controllerName: controllerName,
		xdsCache:       xdsCache,
		// mgr:            mgr,
		client:       client,
		domainSuffix: domainSuffix,
	}
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

type Inputs struct {
	Namespaces krt.Collection[*corev1.Namespace]

	Services krt.Collection[*corev1.Service]
	Secrets  krt.Collection[*corev1.Secret]

	GatewayClasses  krt.Collection[*gateway.GatewayClass]
	Gateways        krt.Collection[*gateway.Gateway]
	HTTPRoutes      krt.Collection[*gateway.HTTPRoute]
	GRPCRoutes      krt.Collection[*gatewayv1.GRPCRoute]
	TCPRoutes       krt.Collection[*gatewayalpha.TCPRoute]
	TLSRoutes       krt.Collection[*gatewayalpha.TLSRoute]
	ReferenceGrants krt.Collection[*gateway.ReferenceGrant]
	ServiceEntries  krt.Collection[*networkingclient.ServiceEntry]
	InferencePools  krt.Collection[*inf.InferencePool]
}

func (s *AgentGwSyncer) Init(krtopts krtutil.KrtOptions) {
	logger.Debug("init agentgateway Syncer", "controllername", s.controllerName)

	inputs := Inputs{
		Namespaces: krt.NewInformer[*corev1.Namespace](s.client),
		Secrets: krt.WrapClient[*corev1.Secret](
			kclient.NewFiltered[*corev1.Secret](s.client, kubetypes.Filter{
				//FieldSelector: kubesecrets.SecretsFieldSelector,
				ObjectFilter: s.client.ObjectFilter(),
			}),
		),
		Services: krt.WrapClient[*corev1.Service](
			kclient.NewFiltered[*corev1.Service](s.client, kubetypes.Filter{ObjectFilter: s.client.ObjectFilter()}),
			krtopts.ToOptions("informer/Services")...),

		GatewayClasses: krt.WrapClient(kclient.New[*gateway.GatewayClass](s.client), krtopts.ToOptions("informer/GatewayClasses")...),
		Gateways:       krt.WrapClient(kclient.New[*gateway.Gateway](s.client), krtopts.ToOptions("informer/Gateways")...),
		HTTPRoutes:     krt.WrapClient(kclient.New[*gateway.HTTPRoute](s.client), krtopts.ToOptions("informer/HTTPRoutes")...),
		GRPCRoutes:     krt.WrapClient(kclient.New[*gatewayv1.GRPCRoute](s.client), krtopts.ToOptions("informer/GRPCRoutes")...),

		ReferenceGrants: krt.WrapClient(kclient.New[*gateway.ReferenceGrant](s.client), krtopts.ToOptions("informer/ReferenceGrants")...),
		ServiceEntries:  krt.WrapClient(kclient.New[*networkingclient.ServiceEntry](s.client), krtopts.ToOptions("informer/ServiceEntries")...),
		//InferencePools:  krt.WrapClient(kclient.New[*inf.InferencePool](s.client), krtopts.ToOptions("informer/InferencePools")...),
	}
	if features.EnableAlphaGatewayAPI {
		inputs.TCPRoutes = krt.WrapClient(kclient.New[*gatewayalpha.TCPRoute](s.client), krtopts.ToOptions("informer/TCPRoutes")...)
		inputs.TLSRoutes = krt.WrapClient(kclient.New[*gatewayalpha.TLSRoute](s.client), krtopts.ToOptions("informer/TLSRoutes")...)
	} else {
		// If disabled, still build a collection but make it always empty
		inputs.TCPRoutes = krt.NewStaticCollection[*gatewayalpha.TCPRoute](nil, krtopts.ToOptions("disable/TCPRoutes")...)
		inputs.TLSRoutes = krt.NewStaticCollection[*gatewayalpha.TLSRoute](nil, krtopts.ToOptions("disable/TLSRoutes")...)
	}

	GatewayClasses := GatewayClassesCollection(inputs.GatewayClasses, krtopts)

	RefGrants := BuildReferenceGrants(ReferenceGrantsCollection(inputs.ReferenceGrants, krtopts))

	// Note: not fully complete until its join with route attachments to report attachedRoutes.
	// Do not register yet.
	Gateways := GatewayCollection(
		inputs.Gateways,
		GatewayClasses,
		inputs.Namespaces,
		RefGrants,
		inputs.Secrets,
		s.domainSuffix,
		krtopts,
	)
	ports := krt.NewCollection(Gateways, func(ctx krt.HandlerContext, obj Gateway) *IndexObject[string, Gateway] {
		port := fmt.Sprint(obj.parentInfo.Port)
		return &IndexObject[string, Gateway]{
			Key:     port,
			Objects: []Gateway{obj},
		}
	}, krtopts.ToOptions("ports")...)

	Binds := krt.NewManyCollection(ports, func(ctx krt.HandlerContext, object IndexObject[string, Gateway]) []ADPResource {
		port, _ := strconv.Atoi(object.Key)
		uniq := sets.New[types.NamespacedName]()
		for _, gw := range object.Objects {
			uniq.Insert(types.NamespacedName{
				Namespace: gw.parent.Namespace,
				Name:      gw.parent.Name,
			})
		}
		var binds []ADPResource
		for _, obj := range uniq.UnsortedList() {
			bind := Bind{
				Bind: &api.Bind{
					Key:  object.Key + "/" + obj.String(),
					Port: uint32(port),
				},
			}
			binds = append(binds, toResource(obj, bind))
		}
		return binds
	}, krtopts.ToOptions("Binds")...)

	Listeners := krt.NewCollection(Gateways, func(ctx krt.HandlerContext, obj Gateway) *ADPResource {
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
		}, ADPListener{l})
	}, krtopts.ToOptions("Listeners")...)

	routeParents := BuildRouteParents(Gateways)

	routeInputs := RouteContextInputs{
		Grants:         RefGrants,
		RouteParents:   routeParents,
		DomainSuffix:   s.domainSuffix,
		Services:       inputs.Services,
		Namespaces:     inputs.Namespaces,
		ServiceEntries: inputs.ServiceEntries,
		InferencePools: inputs.InferencePools,
	}
	ADPRoutes := ADPRouteCollection(
		inputs.HTTPRoutes,
		routeInputs,
		krtopts,
	)

	// TODO: inference pool

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
	nsClient := kclient.NewFiltered[*corev1.Namespace](
		s.commonCols.Client,
		kclient.Filter{ObjectFilter: s.commonCols.Client.ObjectFilter()},
	)
	namespaces := krt.WrapClient(nsClient, s.commonCols.KrtOpts.ToOptions("Namespaces")...)

	seInformer := kclient.NewDelayedInformer[*networkingclient.ServiceEntry](
		s.client, gvr.ServiceEntry,
		kubetypes.StandardInformer, kclient.Filter{ObjectFilter: s.client.ObjectFilter()},
	)
	serviceEntries := krt.WrapClient(seInformer, krtopts.ToOptions("ServiceEntries")...)

	workloadIndex := index{
		services: servicesCollection{},
	}

	// these are agw api-style services combined from kube services and service entries
	WorkloadServices := workloadIndex.ServicesCollection(inputs.Services, serviceEntries, namespaces, krtopts)
	avcAddresses := krt.NewCollection(WorkloadServices, func(ctx krt.HandlerContext, obj ServiceInfo) *ADPCacheAddress {
		var cacheResources []envoytypes.Resource
		addrMessage := obj.AsAddress.Address
		cacheResources = append(cacheResources, &envoyResourceWithCustomName{
			Message: addrMessage,
			Name:    obj.ResourceName(),
			version: utils.HashProto(addrMessage),
		})

		// Create the resource wrappers
		var resourceVersion uint64
		for _, res := range cacheResources {
			resourceVersion ^= res.(*envoyResourceWithCustomName).version
		}

		result := &ADPCacheAddress{
			NamespacedName: types.NamespacedName{Name: obj.Service.GetName(), Namespace: obj.Service.GetNamespace()},
			Address:        envoycache.NewResources(fmt.Sprintf("%d", resourceVersion), cacheResources),
		}
		logger.Debug("created XDS resources for svc address with ID", "addr", fmt.Sprintf("%s,%s", obj.Service.GetName(), obj.Service.GetNamespace()), "resourceid", result.ResourceName())
		return result
	})

	Workloads := workloadIndex.WorkloadsCollection(
		pods,
		WorkloadServices,
		serviceEntries,
		endpointSlices,
		namespaces,
		krtopts,
	)

	proxyKey := "default~agent-gateway" // TODO: don't hard code, use s.perclientSnapCollection
	workloadAddresses := krt.NewCollection(Workloads, func(ctx krt.HandlerContext, obj WorkloadInfo) *ADPCacheAddress {
		var cacheResources []envoytypes.Resource
		addrMessage := obj.AsAddress.Address
		cacheResources = append(cacheResources, &envoyResourceWithCustomName{
			Message: addrMessage,
			Name:    obj.ResourceName(),
			version: utils.HashProto(addrMessage),
		})

		// Create the resource wrappers
		var resourceVersion uint64
		for _, res := range cacheResources {
			resourceVersion ^= res.(*envoyResourceWithCustomName).version
		}

		result := &ADPCacheAddress{
			NamespacedName: types.NamespacedName{Name: obj.Workload.GetName(), Namespace: obj.Workload.GetNamespace()},
			Address:        envoycache.NewResources(fmt.Sprintf("%d", resourceVersion), cacheResources),
			proxyKey:       proxyKey,
		}
		logger.Debug("created XDS resources for workload address with ID", "addr", fmt.Sprintf("%s,%s", obj.Workload.GetName(), obj.Workload.GetNamespace()), "resourceid", result.ResourceName())
		return result
	})
	s.addressxDS = krt.JoinCollection([]krt.Collection[ADPCacheAddress]{avcAddresses, workloadAddresses}, krtopts.ToOptions("ADPAddresses")...)

	resources := krt.JoinCollection([]krt.Collection[ADPResource]{Binds, Listeners, ADPRoutes}, krtopts.ToOptions("ADPResources")...)
	s.resourcexDS = krt.NewCollection(resources, func(ctx krt.HandlerContext, obj ADPResource) *ADPCacheResource {
		var cacheResources []envoytypes.Resource
		cacheResources = append(cacheResources, &envoyResourceWithCustomName{
			Message: obj.Resource,
			Name:    obj.ResourceName(),
			version: utils.HashProto(obj.Resource),
		})

		// Create the resource wrappers
		var resourceVersion uint64
		for _, res := range cacheResources {
			resourceVersion ^= res.(*envoyResourceWithCustomName).version
		}

		result := &ADPCacheResource{
			Gateway:   obj.Gateway,
			Resources: envoycache.NewResources(fmt.Sprintf("%d", resourceVersion), cacheResources),
		}
		logger.Debug("created XDS resources for gateway with ID", "gwname", fmt.Sprintf("%s,%s", obj.Gateway.Name, obj.Gateway.Namespace), "resourceid", result.ResourceName())
		return result
	})

	s.waitForSync = []cache.InformerSynced{
		s.commonCols.HasSynced,
		// resources
		Binds.HasSynced,
		Listeners.HasSynced,
		ADPRoutes.HasSynced,
		s.resourcexDS.HasSynced,
		// addresses
		serviceEntries.HasSynced,
		namespaces.HasSynced,
		pods.HasSynced,
		endpointSlices.HasSynced,
		WorkloadServices.HasSynced,
		Workloads.HasSynced,
		s.addressxDS.HasSynced,
	}
}

func (s *AgentGwSyncer) Start(ctx context.Context) error {
	logger.Info("starting agentgateway Syncer", "controllername", s.controllerName)
	logger.Info("waiting for agentgateway cache to sync")

	// Wait for cache to sync
	if !kube.WaitForCacheSync("agentgateway syncer", ctx.Done(), s.waitForSync...) {
		return fmt.Errorf("agentgateway syncer waiting for cache to sync failed")
	}

	s.resourcexDS.RegisterBatch(func(events []krt.Event[ADPCacheResource], _ bool) {
		for _, e := range events {
			r := e.Latest()
			if e.Event == controllers.EventDelete {
				s.xdsCache.ClearSnapshot(r.ResourceName())
				continue
			}
			snapshot := &agentGwSnapshot{
				Resources: r.Resources,
			}
			logger.Debug("setting xds snapshot", "resourcename", r.ResourceName())
			err := s.xdsCache.SetSnapshot(ctx, r.ResourceName(), snapshot)
			if err != nil {
				logger.Error("failed to set xds snapshot", "resourcename", r.ResourceName(), "error", err.Error())
				continue
			}
		}
	}, true)

	s.addressxDS.RegisterBatch(func(events []krt.Event[ADPCacheAddress], _ bool) {
		for _, e := range events {
			r := e.Latest()
			if e.Event == controllers.EventDelete {
				s.xdsCache.ClearSnapshot(r.ResourceName())
				continue
			}
			snapshot := &agentGwSnapshot{
				Addresses: r.Address,
			}
			logger.Debug("setting xds snapshot", "resourcename", r.ResourceName())
			snapWrap := e.Latest()
			err := s.xdsCache.SetSnapshot(ctx, snapWrap.proxyKey, snapshot)
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
