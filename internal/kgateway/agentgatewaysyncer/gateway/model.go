// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gateway

import (
	"bytes"
	"cmp"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	udpa "github.com/cncf/xds/go/udpa/type/v1"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
	"istio.io/api/label"
	"istio.io/istio/pilot/pkg/credentials"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/trustbundle"
	"istio.io/istio/pilot/pkg/util/protoconv"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/config/mesh/meshwatcher"
	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	"istio.io/istio/pkg/network"
	"istio.io/istio/pkg/spiffe"
	"istio.io/istio/pkg/util/gogoprotomarshal"
	"istio.io/istio/pkg/util/identifier"
	"istio.io/istio/pkg/util/protomarshal"
	"istio.io/istio/pkg/xds"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/kind"
	"istio.io/istio/pkg/util/hash"
	netutil "istio.io/istio/pkg/util/net"
	"istio.io/istio/pkg/util/sets"
)

// Statically link protobuf descriptors from UDPA
var _ = udpa.TypedStruct{}

type ConfigHash uint64

// ConfigKey describe a specific config item.
// In most cases, the name is the config's name. However, for ServiceEntry it is service's FQDN.
type ConfigKey struct {
	Kind      kind.Kind
	Name      string
	Namespace string
}

func (key ConfigKey) HashCode() ConfigHash {
	h := hash.New()
	h.Write([]byte{byte(key.Kind)})
	// Add separator / to avoid collision.
	h.WriteString("/")
	h.WriteString(key.Namespace)
	h.WriteString("/")
	h.WriteString(key.Name)
	return ConfigHash(h.Sum64())
}

func (key ConfigKey) String() string {
	return key.Kind.String() + "/" + key.Namespace + "/" + key.Name
}

// ConfigsOfKind extracts configs of the specified kind.
func ConfigsOfKind(configs sets.Set[ConfigKey], kind kind.Kind) sets.Set[ConfigKey] {
	ret := make(sets.Set[ConfigKey])

	for conf := range configs {
		if conf.Kind == kind {
			ret.Insert(conf)
		}
	}

	return ret
}

// HasConfigsOfKind returns true if configs has changes of type kind
func HasConfigsOfKind(configs sets.Set[ConfigKey], kind kind.Kind) bool {
	for c := range configs {
		if c.Kind == kind {
			return true
		}
	}
	return false
}

// ConfigNamesOfKind extracts config names of the specified kind.
func ConfigNamesOfKind(configs sets.Set[ConfigKey], kind kind.Kind) sets.String {
	ret := sets.New[string]()

	for conf := range configs {
		if conf.Kind == kind {
			ret.Insert(conf.Name)
		}
	}

	return ret
}

// ConfigStore describes a set of platform agnostic APIs that must be supported
// by the underlying platform to store and retrieve Istio configuration.
//
// Configuration key is defined to be a combination of the type, name, and
// namespace of the configuration object. The configuration key is guaranteed
// to be unique in the store.
//
// The storage interface presented here assumes that the underlying storage
// layer supports _Get_ (list), _Update_ (update), _Create_ (create) and
// _Delete_ semantics but does not guarantee any transactional semantics.
//
// _Update_, _Create_, and _Delete_ are mutator operations. These operations
// are asynchronous, and you might not see the effect immediately (e.g. _Get_
// might not return the object by key immediately after you mutate the store.)
// Intermittent errors might occur even though the operation succeeds, so you
// should always check if the object store has been modified even if the
// mutating operation returns an error.  Objects should be created with
// _Create_ operation and updated with _Update_ operation.
//
// Resource versions record the last mutation operation on each object. If a
// mutation is applied to a different revision of an object than what the
// underlying storage expects as defined by pure equality, the operation is
// blocked.  The client of this interface should not make assumptions about the
// structure or ordering of the revision identifier.
//
// Object references supplied and returned from this interface should be
// treated as read-only. Modifying them violates thread-safety.
type ConfigStore interface {
	// Schemas exposes the configuration type schema known by the config store.
	// The type schema defines the bidirectional mapping between configuration
	// types and the protobuf encoding schema.
	Schemas() collection.Schemas

	// Get retrieves a configuration element by a type and a key
	Get(typ GroupVersionKind, name, namespace string) *Config

	// List returns objects by type and namespace.
	// Use "" for the namespace to list across namespaces.
	List(typ GroupVersionKind, namespace string) []Config

	// Create adds a new configuration object to the store. If an object with the
	// same name and namespace for the type already exists, the operation fails
	// with no side effects.
	Create(config Config) (revision string, err error)

	// Update modifies an existing configuration object in the store.  Update
	// requires that the object has been created.  Resource version prevents
	// overriding a value that has been changed between prior _Get_ and _Put_
	// operation to achieve optimistic concurrency. This method returns a new
	// revision if the operation succeeds.
	Update(config Config) (newRevision string, err error)
	UpdateStatus(config Config) (newRevision string, err error)

	// Patch applies only the modifications made in the PatchFunc rather than doing a full replace. Useful to avoid
	// read-modify-write conflicts when there are many concurrent-writers to the same resource.
	Patch(orig Config, patchFn PatchFunc) (string, error)

	// Delete removes an object from the store by key
	// For k8s, resourceVersion must be fulfilled before a deletion is carried out.
	// If not possible, a 409 Conflict status will be returned.
	Delete(typ GroupVersionKind, name, namespace string, resourceVersion *string) error
}

// ResolveShortnameToFQDN uses metadata information to resolve a reference
// to shortname of the service to FQDN
func ResolveShortnameToFQDN(hostname string, meta Meta) host.Name {
	if len(hostname) == 0 {
		// only happens when the gateway-api BackendRef is invalid
		return ""
	}
	out := hostname
	// Treat the wildcard hostname as fully qualified. Any other variant of a wildcard hostname will contain a `.` too,
	// and skip the next if, so we only need to check for the literal wildcard itself.
	if hostname == "*" {
		return host.Name(out)
	}

	// if the hostname is a valid ipv4 or ipv6 address, do not append domain or namespace
	if netutil.IsValidIPAddress(hostname) {
		return host.Name(out)
	}

	// if FQDN is specified, do not append domain or namespace to hostname
	if !strings.Contains(hostname, ".") {
		if meta.Namespace != "" {
			out = out + "." + meta.Namespace
		}

		// FIXME this is a gross hack to hardcode a service's domain name in kubernetes
		// BUG this will break non kubernetes environments if they use shortnames in the
		// rules.
		if meta.Domain != "" {
			out = out + ".svc." + meta.Domain
		}
	}

	return host.Name(out)
}

// resolveGatewayName uses metadata information to resolve a reference
// to shortname of the gateway to FQDN
func resolveGatewayName(gwname string, meta Meta) string {
	out := gwname

	// New way of binding to a gateway in remote namespace
	// is ns/name. Old way is either FQDN or short name
	if !strings.Contains(gwname, "/") {
		if !strings.Contains(gwname, ".") {
			// we have a short name. Resolve to a gateway in same namespace
			out = meta.Namespace + "/" + gwname
		} else {
			// parse namespace from FQDN. This is very hacky, but meant for backward compatibility only
			// This is a legacy FQDN format. Transform name.ns.svc.cluster.local -> ns/name
			i := strings.Index(gwname, ".")
			fqdn := strings.Index(gwname[i+1:], ".")
			if fqdn == -1 {
				out = gwname[i+1:] + "/" + gwname[:i]
			} else {
				out = gwname[i+1:i+1+fqdn] + "/" + gwname[:i]
			}
		}
	} else {
		// remove the . from ./gateway and substitute it with the namespace name
		i := strings.Index(gwname, "/")
		if gwname[:i] == "." {
			out = meta.Namespace + "/" + gwname[i+1:]
		}
	}
	return out
}

// MostSpecificHostMatch compares the maps of specific and wildcard hosts to the needle, and returns the longest element
// matching the needle and it's value, or false if no element in the maps matches the needle.
func MostSpecificHostMatch[V any](needle host.Name, specific map[host.Name]V, wildcard map[host.Name]V) (host.Name, V, bool) {
	if needle.IsWildCarded() {
		// exact match first
		if v, ok := wildcard[needle]; ok {
			return needle, v, true
		}

		return mostSpecificHostWildcardMatch(string(needle[1:]), wildcard)
	}

	// exact match first
	if v, ok := specific[needle]; ok {
		return needle, v, true
	}

	// check wildcard
	return mostSpecificHostWildcardMatch(string(needle), wildcard)
}

func mostSpecificHostWildcardMatch[V any](needle string, wildcard map[host.Name]V) (host.Name, V, bool) {
	found := false
	var matchHost host.Name
	var matchValue V

	for h, v := range wildcard {
		if strings.HasSuffix(needle, string(h[1:])) {
			if !found {
				matchHost = h
				matchValue = wildcard[h]
				found = true
			} else if host.MoreSpecific(h, matchHost) {
				matchHost = h
				matchValue = v
			}
		}
	}

	return matchHost, matchValue, found
}

// sortConfigByCreationTime sorts the list of config objects in ascending order by their creation time (if available)
func sortConfigByCreationTime(configs []Config) []Config {
	sort.Slice(configs, func(i, j int) bool {
		if r := configs[i].CreationTimestamp.Compare(configs[j].CreationTimestamp); r != 0 {
			return r == -1 // -1 means i is less than j, so return true
		}
		// If creation time is the same, then behavior is nondeterministic. In this case, we can
		// pick an arbitrary but consistent ordering based on name and namespace, which is unique.
		// CreationTimestamp is stored in seconds, so this is not uncommon.
		if r := cmp.Compare(configs[i].Name, configs[j].Name); r != 0 {
			return r == -1
		}
		return cmp.Compare(configs[i].Namespace, configs[j].Namespace) == -1
	})
	return configs
}

type (
	Node                    = pm.Node
	NodeMetadata            = pm.NodeMetadata
	NodeMetaProxyConfig     = pm.NodeMetaProxyConfig
	NodeType                = pm.NodeType
	BootstrapNodeMetadata   = pm.BootstrapNodeMetadata
	TrafficInterceptionMode = pm.TrafficInterceptionMode
	PodPort                 = pm.PodPort
	StringBool              = pm.StringBool
	IPMode                  = pm.IPMode
)

const (
	SidecarProxy = pm.SidecarProxy
	Router       = pm.Router
	Waypoint     = pm.Waypoint
	Ztunnel      = pm.Ztunnel

	IPv4 = pm.IPv4
	IPv6 = pm.IPv6
	Dual = pm.Dual
)

var _ mesh.Holder = &Environment{}

func NewEnvironment() *Environment {
	var cache XdsCache
	if features.EnableXDSCaching {
		cache = NewXdsCache()
	} else {
		cache = DisabledCache{}
	}
	return &Environment{
		pushContext:   NewPushContext(),
		Cache:         cache,
		EndpointIndex: NewEndpointIndex(cache),
	}
}

// Watcher is a type alias to keep the embedded type name stable.
type Watcher = meshwatcher.WatcherCollection

// Environment provides an aggregate environmental API for Pilot
type Environment struct {
	// Discovery interface for listing services and instances.
	ServiceDiscovery

	// Config interface for listing routing rules
	ConfigStore

	// Watcher is the watcher for the mesh config (to be merged into the config store)
	Watcher

	// NetworksWatcher (loaded from a config map) provides information about the
	// set of networks inside a mesh and how to route to endpoints in each
	// network. Each network provides information about the endpoints in a
	// routable L3 network. A single routable L3 network can have one or more
	// service registries.
	NetworksWatcher mesh.NetworksWatcher

	NetworkManager *NetworkManager

	// mutex used for protecting Environment.pushContext
	mutex sync.RWMutex
	// pushContext holds information during push generation. It is reset on config change, at the beginning
	// of the pushAll. It will hold all errors and stats and possibly caches needed during the entire cache computation.
	// DO NOT USE EXCEPT FOR TESTS AND HANDLING OF NEW CONNECTIONS.
	// ALL USE DURING A PUSH SHOULD USE THE ONE CREATED AT THE
	// START OF THE PUSH, THE GLOBAL ONE MAY CHANGE AND REFLECT A DIFFERENT
	// CONFIG AND PUSH
	pushContext *PushContext

	// DomainSuffix provides a default domain for the Istio server.
	DomainSuffix string

	// TrustBundle: List of Mesh TrustAnchors
	TrustBundle *trustbundle.TrustBundle

	clusterLocalServices ClusterLocalProvider

	CredentialsController credentials.MulticlusterController

	GatewayAPIController GatewayController

	// EndpointShards for a service. This is a global (per-server) list, built from
	// incremental updates. This is keyed by service and namespace
	EndpointIndex *EndpointIndex

	// Cache for XDS resources.
	Cache XdsCache
}

func (e *Environment) Mesh() *meshMeshConfig {
	if e != nil && e.Watcher != nil {
		return e.Watcher.Mesh()
	}
	return nil
}

func (e *Environment) MeshNetworks() *meshMeshNetworks {
	if e != nil && e.NetworksWatcher != nil {
		return e.NetworksWatcher.Networks()
	}
	return nil
}

// SetPushContext sets the push context with lock protected
func (e *Environment) SetPushContext(pc *PushContext) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.pushContext = pc
}

// PushContext returns the push context with lock protected
func (e *Environment) PushContext() *PushContext {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.pushContext
}

// GetDiscoveryAddress parses the DiscoveryAddress specified via Mesh
func (e *Environment) GetDiscoveryAddress() (host.Name, string, error) {
	proxyConfig := mesh.DefaultProxyConfig()
	if e.Mesh().DefaultConfig != nil {
		proxyConfig = e.Mesh().DefaultConfig
	}
	hostname, port, err := net.SplitHostPort(proxyDiscoveryAddress)
	if err != nil {
		return "", "", fmt.Errorf("invalid Istiod Address: %s, %v", proxyDiscoveryAddress, err)
	}
	if _, err := strconv.Atoi(port); err != nil {
		return "", "", fmt.Errorf("invalid Istiod Port: %s, %s, %v", port, proxyDiscoveryAddress, err)
	}
	return host.Name(hostname), port, nil
}

func (e *Environment) AddMeshHandler(h func()) {
	if e != nil && e.Watcher != nil {
		e.Watcher.AddMeshHandler(h)
	}
}

func (e *Environment) AddNetworksHandler(h func()) {
	if e != nil && e.NetworksWatcher != nil {
		e.NetworksWatcher.AddNetworksHandler(h)
	}
}

func (e *Environment) AddMetric(metric monitoring.Metric, key string, proxyID, msg string) {
	if e != nil {
		e.PushContext().AddMetric(metric, key, proxyID, msg)
	}
}

// Init initializes the Environment for use.
func (e *Environment) Init() {
	// Use a default DomainSuffix, if none was provided.
	if len(e.DomainSuffix) == 0 {
		e.DomainSuffix = constants.DefaultClusterLocalDomain
	}

	// Create the cluster-local service registry.
	e.clusterLocalServices = NewClusterLocalProvider(e)
}

func (e *Environment) InitNetworksManager(updater XDSUpdater) (err error) {
	e.NetworkManager, err = NewNetworkManager(e, updater)
	return
}

func (e *Environment) ClusterLocal() ClusterLocalProvider {
	return e.clusterLocalServices
}

func (e *Environment) GetProxyConfigOrDefault(ns string, labels, annotations map[string]string, meshConfig *meshMeshConfig) *meshProxyConfig {
	push := e.PushContext()
	if push != nil && push.ProxyConfigs != nil {
		if generatedProxyConfig := push.ProxyConfigs.EffectiveProxyConfig(
			&NodeMetadata{
				Namespace:   ns,
				Labels:      labels,
				Annotations: annotations,
			}, meshConfig); generatedProxyConfig != nil {
			return generatedProxyConfig
		}
	}
	return mesh.DefaultProxyConfig()
}

// Resources is an alias for array of marshaled resources.
type Resources = []*discovery.Resource

// DeletedResources is an alias for array of strings that represent removed resources in delta.
type DeletedResources = []string

func AnyToUnnamedResources(r []*anypb.Any) Resources {
	a := make(Resources, 0, len(r))
	for _, rr := range r {
		a = append(a, &discovery.Resource{Resource: rr})
	}
	return a
}

// XdsUpdates include information about the subset of updated resources.
// See for example EDS incremental updates.
type XdsUpdates = sets.Set[ConfigKey]

// XdsLogDetails contains additional metadata that is captured by Generators and used by xds processors
// like Ads and Delta to uniformly log.
type XdsLogDetails struct {
	Incremental    bool
	AdditionalInfo string
}

var DefaultXdsLogDetails = XdsLogDetails{}

// XdsResourceGenerator creates the response for a typeURL DiscoveryRequest or DeltaDiscoveryRequest. If no generator
// is associated with a Proxy, the default (a networking.core.ConfigGenerator instance) will be used.
// The server may associate a different generator based on client metadata. Different
// WatchedResources may use same or different Generator.
// Note: any errors returned will completely close the XDS stream. Use with caution; typically and empty
// or no response is preferred.
type XdsResourceGenerator interface {
	// Generate generates the Sotw resources for Xds.
	Generate(proxy *Proxy, w *WatchedResource, req *PushRequest) (Resources, XdsLogDetails, error)
}

// XdsDeltaResourceGenerator generates Sotw and delta resources.
type XdsDeltaResourceGenerator interface {
	XdsResourceGenerator
	// GenerateDeltas returns the changed and removed resources, along with whether or not delta was actually used.
	GenerateDeltas(proxy *Proxy, req *PushRequest, w *WatchedResource) (Resources, DeletedResources, XdsLogDetails, bool, error)
}

// Proxy contains information about an specific instance of a proxy (envoy sidecar, gateway,
// etc). The Proxy is initialized when a sidecar connects to Pilot, and populated from
// 'node' info in the protocol as well as data extracted from registries.
//
// In current Istio implementation nodes use a 4-parts '~' delimited ID.
// Type~IPAddress~ID~Domain
type Proxy struct {
	sync.RWMutex

	// Type specifies the node type. First part of the ID.
	Type NodeType

	// IPAddresses is the IP addresses of the proxy used to identify it and its
	// co-located service instances. Example: "10.60.1.6". In some cases, the host
	// where the proxy and service instances reside may have more than one IP address
	IPAddresses []string

	// ID is the unique platform-specific sidecar proxy ID. For k8s it is the pod ID and
	// namespace <podName.namespace>.
	ID string

	// Locality is the location of where Envoy proxy runs. This is extracted from
	// the registry where possible. If the registry doesn't provide a locality for the
	// proxy it will use the one sent via ADS that can be configured in the Envoy bootstrap
	Locality *core.Locality

	// DNSDomain defines the DNS domain suffix for short hostnames (e.g.
	// "default.svc.cluster.local")
	DNSDomain string

	// ConfigNamespace defines the namespace where this proxy resides
	// for the purposes of network scoping.
	// NOTE: DO NOT USE THIS FIELD TO CONSTRUCT DNS NAMES
	ConfigNamespace string

	// Labels specifies the set of workload instance (ex: k8s pod) labels associated with this node.
	// Labels can be different from that in Metadata because of pod labels update after startup,
	// while NodeMetadata.Labels are set during bootstrap.
	Labels map[string]string

	// Metadata key-value pairs extending the Node identifier
	Metadata *NodeMetadata

	// the sidecarScope associated with the proxy
	SidecarScope *SidecarScope

	// the sidecarScope associated with the proxy previously
	PrevSidecarScope *SidecarScope

	// The merged gateways associated with the proxy if this is a Router
	MergedGateway *MergedGateway

	// PrevMergedGateway contains information about merged gateway associated with the proxy previously
	PrevMergedGateway *PrevMergedGateway

	// ServiceTargets contains a list of all Services associated with the proxy, contextualized for this particular proxy.
	// These are unique to this proxy, as the port information is specific to it - while a ServicePort is shared with the
	// service, the target port may be distinct per-endpoint. So this maintains a view specific to this proxy.
	// ServiceTargets will maintain a list entry for each Service-port, so if we have 2 services each with 3 ports, we
	// would have 6 entries.
	ServiceTargets []ServiceTarget

	// Istio version associated with the Proxy
	IstioVersion *IstioVersion

	// VerifiedIdentity determines whether a proxy had its identity verified. This
	// generally occurs by JWT or mTLS authentication. This can be false when
	// connecting over plaintext. If this is set to true, we can verify the proxy has
	// access to ConfigNamespace namespace. However, other options such as node type
	// are not part of an Istio identity and thus are not verified.
	VerifiedIdentity *spiffe.Identity

	// IPMode of proxy.
	ipMode IPMode

	// GlobalUnicastIP stores the global unicast IP if available, otherwise nil
	GlobalUnicastIP string

	// XdsResourceGenerator is used to generate resources for the node, based on the PushContext.
	// If nil, the default networking/core v2 generator is used. This field can be set
	// at connect time, based on node metadata, to trigger generation of a different style
	// of configuration.
	XdsResourceGenerator XdsResourceGenerator

	// WatchedResources contains the list of watched resources for the proxy, keyed by the DiscoveryRequest TypeUrl.
	WatchedResources map[string]*WatchedResource

	// XdsNode is the xDS node identifier
	XdsNode *core.Node

	workloadEntryName        string
	workloadEntryAutoCreated bool

	// LastPushContext stores the most recent push context for this proxy. This will be monotonically
	// increasing in version. Requests should send config based on this context; not the global latest.
	// Historically, the latest was used which can cause problems when computing whether a push is
	// required, as the computed sidecar scope version would not monotonically increase.
	LastPushContext *PushContext
	// LastPushTime records the time of the last push. This is used in conjunction with
	// LastPushContext; the XDS cache depends on knowing the time of the PushContext to determine if a
	// key is stale or not.
	LastPushTime time.Time
}

type WatchedResource = xds.WatchedResource

// GetView returns a restricted view of the mesh for this proxy. The view can be
// restricted by network (via ISTIO_META_REQUESTED_NETWORK_VIEW).
// If not set, we assume that the proxy wants to see endpoints in any network.
func (node *Proxy) GetView() ProxyView {
	return newProxyView(node)
}

// InNetwork returns true if the proxy is on the given network, or if either
// the proxy's network or the given network is unspecified ("").
func (node *Proxy) InNetwork(network network.ID) bool {
	return node == nil || identifier.IsSameOrEmpty(network.String(), node.Metadata.Network.String())
}

// InCluster returns true if the proxy is in the given cluster, or if either
// the proxy's cluster id or the given cluster id is unspecified ("").
func (node *Proxy) InCluster(cluster cluster.ID) bool {
	return node == nil || identifier.IsSameOrEmpty(cluster.String(), node.Metadata.ClusterID.String())
}

// IsWaypointProxy returns true if the proxy is acting as a waypoint proxy in an ambient mesh.
func (node *Proxy) IsWaypointProxy() bool {
	return node.Type == Waypoint
}

// IsZTunnel returns true if the proxy is acting as a ztunnel in an ambient mesh.
func (node *Proxy) IsZTunnel() bool {
	return node.Type == Ztunnel
}

// IsAmbient returns true if the proxy is acting as either a ztunnel or a waypoint proxy in an ambient mesh.
func (node *Proxy) IsAmbient() bool {
	return node.IsWaypointProxy() || node.IsZTunnel()
}

var NodeTypes = [...]NodeType{SidecarProxy, Router, Waypoint, Ztunnel}

// SetGatewaysForProxy merges the Gateway objects associated with this
// proxy and caches the merged object in the proxy Node. This is a convenience hack so that
// callers can simply call push.MergedGateways(node) instead of having to
// fetch all the gateways and invoke the merge call in multiple places (lds/rds).
// Must be called after ServiceTargets are set
func (node *Proxy) SetGatewaysForProxy(ps *PushContext) {
	if node.Type != Router {
		return
	}
	var prevMergedGateway MergedGateway
	if node.MergedGateway != nil {
		prevMergedGateway = *node.MergedGateway
	}
	node.MergedGateway = ps.mergeGateways(node)
	node.PrevMergedGateway = &PrevMergedGateway{
		ContainsAutoPassthroughGateways: prevMergedGateway.ContainsAutoPassthroughGateways,
		AutoPassthroughSNIHosts:         prevMergedGateway.GetAutoPassthroughGatewaySNIHosts(),
	}
}

func (node *Proxy) ShouldUpdateServiceTargets(updates sets.Set[ConfigKey]) bool {
	// we only care for services which can actually select this proxy
	for config := range updates {
		if Kind == kind.ServiceEntry || Namespace == node.Metadata.Namespace {
			return true
		}
	}

	return false
}

func (node *Proxy) SetServiceTargets(serviceDiscovery ServiceDiscovery) {
	instances := serviceDiscovery.GetProxyServiceTargets(node)

	// Keep service instances in order of creation/hostname.
	sort.SliceStable(instances, func(i, j int) bool {
		if instances[i].Service != nil && instances[j].Service != nil {
			if !instances[i].Service.CreationTime.Equal(instances[j].Service.CreationTime) {
				return instances[i].Service.CreationTime.Before(instances[j].Service.CreationTime)
			}
			// Additionally, sort by hostname just in case services created automatically at the same second.
			return instances[i].Service.Hostname < instances[j].Service.Hostname
		}
		return true
	})

	node.ServiceTargets = instances
}

// SetWorkloadLabels will set the node.Labels.
// It merges both node meta labels and workload labels and give preference to workload labels.
func (node *Proxy) SetWorkloadLabels(env *Environment) {
	// If this is VM proxy, do not override labels at all, because in istio test we use pod to simulate VM.
	if node.IsVM() {
		node.Labels = node.Metadata.Labels
		return
	}
	labels := env.GetProxyWorkloadLabels(node)
	if labels != nil {
		node.Labels = make(map[string]string, len(labels)+len(node.Metadata.StaticLabels))
		// we can't just equate proxy workload labels to node meta labels as it may be customized by user
		// with `ISTIO_METAJSON_LABELS` env (pkg/bootstrap/go extractAttributesMetadata).
		// so, we fill the `ISTIO_METAJSON_LABELS` as well.
		for k, v := range node.Metadata.StaticLabels {
			node.Labels[k] = v
		}
		for k, v := range labels {
			node.Labels[k] = v
		}
	} else {
		// If could not find pod labels, fallback to use the node metadata labels.
		node.Labels = node.Metadata.Labels
	}
}

// DiscoverIPMode discovers the IP Versions supported by Proxy based on its IP addresses.
func (node *Proxy) DiscoverIPMode() {
	node.ipMode = pm.DiscoverIPMode(node.IPAddresses)
	node.GlobalUnicastIP = networkutil.GlobalUnicastIP(node.IPAddresses)
}

// SupportsIPv4 returns true if proxy supports IPv4 addresses.
func (node *Proxy) SupportsIPv4() bool {
	return node.ipMode == IPv4 || node.ipMode == Dual
}

// SupportsIPv6 returns true if proxy supports IPv6 addresses.
func (node *Proxy) SupportsIPv6() bool {
	return node.ipMode == IPv6 || node.ipMode == Dual
}

// IsIPv6 returns true if proxy only supports IPv6 addresses.
func (node *Proxy) IsIPv6() bool {
	return node.ipMode == IPv6
}

func (node *Proxy) IsDualStack() bool {
	return node.ipMode == Dual
}

// GetIPMode returns proxy's ipMode
func (node *Proxy) GetIPMode() IPMode {
	return node.ipMode
}

// SetIPMode set node's ip mode
// Note: Donot use this function directly in most cases, use DiscoverIPMode instead.
func (node *Proxy) SetIPMode(mode IPMode) {
	node.ipMode = mode
}

// ParseMetadata parses the opaque Metadata from an Envoy Node into string key-value pairs.
// Any non-string values are ignored.
func ParseMetadata(metadata *structpb.Struct) (*NodeMetadata, error) {
	if metadata == nil {
		return &NodeMetadata{}, nil
	}

	bootstrapNodeMeta, err := ParseBootstrapNodeMetadata(metadata)
	if err != nil {
		return nil, err
	}
	return &bootstrapNodeMeta.NodeMetadata, nil
}

// ParseBootstrapNodeMetadata parses the opaque Metadata from an Envoy Node into string key-value pairs.
func ParseBootstrapNodeMetadata(metadata *structpb.Struct) (*BootstrapNodeMetadata, error) {
	if metadata == nil {
		return &BootstrapNodeMetadata{}, nil
	}

	b, err := protomarshal.MarshalProtoNames(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to read node metadata %v: %v", metadata, err)
	}
	meta := &BootstrapNodeMetadata{}
	if err := json.Unmarshal(b, meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node metadata (%v): %v", string(b), err)
	}
	return meta, nil
}

// ParseServiceNodeWithMetadata parse the Envoy Node from the string generated by ServiceNode
// function and the metadata.
func ParseServiceNodeWithMetadata(nodeID string, metadata *NodeMetadata) (*Proxy, error) {
	parts := strings.Split(nodeID, serviceNodeSeparator)
	out := &Proxy{
		Metadata: metadata,
	}

	if len(parts) != 4 {
		return out, fmt.Errorf("missing parts in the service node %q", nodeID)
	}

	if !pm.IsApplicationNodeType(NodeType(parts[0])) {
		return out, fmt.Errorf("invalid node type (valid types: %v) in the service node %q", NodeTypes, nodeID)
	}
	out.Type = NodeType(parts[0])

	// Get all IP Addresses from Metadata
	if hasValidIPAddresses(metadata.InstanceIPs) {
		out.IPAddresses = metadata.InstanceIPs
	} else if netutil.IsValidIPAddress(parts[1]) {
		// Fall back, use IP from node id, it's only for backward-compatibility, IP should come from metadata
		out.IPAddresses = append(out.IPAddresses, parts[1])
	}

	// Does query from ingress or router have to carry valid IP address?
	if len(out.IPAddresses) == 0 {
		return out, fmt.Errorf("no valid IP address in the service node id or metadata")
	}

	out.ID = parts[2]
	out.DNSDomain = parts[3]
	if len(metadata.IstioVersion) == 0 {
		log.Warnf("Istio Version is not found in metadata for %v, which may have undesirable side effects", out.ID)
	}
	out.IstioVersion = ParseIstioVersion(metadata.IstioVersion)
	return out, nil
}

// GetOrDefault returns either the value, or the default if the value is empty. Useful when retrieving node metadata fields.
func GetOrDefault(s string, def string) string {
	return pm.GetOrDefault(s, def)
}

// GetProxyConfigNamespace extracts the namespace associated with the proxy
// from the proxy metadata or the proxy ID
func GetProxyConfigNamespace(proxy *Proxy) string {
	if proxy == nil {
		return ""
	}

	// First look for ISTIO_META_CONFIG_NAMESPACE
	// All newer proxies (from Istio 1.1 onwards) are supposed to supply this
	if len(proxy.Metadata.Namespace) > 0 {
		return proxy.Metadata.Namespace
	}

	// if not found, for backward compatibility, extract the namespace from
	// the proxy domain. this is a k8s specific hack and should be enabled
	parts := strings.Split(proxy.DNSDomain, ".")
	if len(parts) > 1 { // k8s will have namespace.<domain>
		return parts[0]
	}

	return ""
}

const (
	serviceNodeSeparator = "~"
)

// hasValidIPAddresses returns true if the input ips are all valid, otherwise returns false.
func hasValidIPAddresses(ipAddresses []string) bool {
	if len(ipAddresses) == 0 {
		return false
	}
	for _, ipAddress := range ipAddresses {
		if !netutil.IsValidIPAddress(ipAddress) {
			return false
		}
	}
	return true
}

const (
	// InterceptionNone indicates that the workload is not using IPtables for traffic interception
	InterceptionNone TrafficInterceptionMode = "NONE"

	// InterceptionTproxy implies traffic intercepted by IPtables with TPROXY mode
	InterceptionTproxy TrafficInterceptionMode = "TPROXY"

	// InterceptionRedirect implies traffic intercepted by IPtables with REDIRECT mode
	// This is our default mode
	InterceptionRedirect TrafficInterceptionMode = "REDIRECT"
)

// GetInterceptionMode extracts the interception mode associated with the proxy
// from the proxy metadata
func (node *Proxy) GetInterceptionMode() TrafficInterceptionMode {
	if node == nil {
		return InterceptionRedirect
	}

	switch node.Metadata.InterceptionMode {
	case "TPROXY":
		return InterceptionTproxy
	case "REDIRECT":
		return InterceptionRedirect
	case "NONE":
		return InterceptionNone
	}

	return InterceptionRedirect
}

// IsUnprivileged returns true if the proxy has explicitly indicated that it is
// unprivileged, i.e. it cannot bind to the privileged ports 1-1023.
func (node *Proxy) IsUnprivileged() bool {
	if node == nil || node.Metadata == nil {
		return false
	}
	// expect explicit "true" value
	unprivileged, _ := strconv.ParseBool(node.Metadata.UnprivilegedPod)
	return unprivileged
}

// CanBindToPort returns true if the proxy can bind to a given port.
// canbind indicates whether the proxy can bind to the port.
// knownlistener indicates whether the check failed if the proxy is trying to bind to a port that is reserved for a static listener or virtual listener.
func (node *Proxy) CanBindToPort(bindTo bool, proxy *Proxy, push *PushContext,
	bind string, port int, protocol protocol.Instance, wildcard string,
) (canbind bool, knownlistener bool) {
	if bindTo {
		if isPrivilegedPort(port) && node.IsUnprivileged() {
			return false, false
		}
	}
	if conflictWithReservedListener(proxy, push, bind, port, protocol, wildcard) {
		return false, true
	}
	return true, false
}

// conflictWithReservedListener checks whether the listener address bind:port conflicts with
// - static listener port：default is 15021 and 15090
// - virtual listener port: default is 15001 and 15006 (only need to check for outbound listener)
func conflictWithReservedListener(proxy *Proxy, push *PushContext, bind string, port int, protocol protocol.Instance, wildcard string) bool {
	if bind != "" {
		if bind != wildcard {
			return false
		}
	} else if !protocol.IsHTTP() {
		// if the protocol is HTTP and bind == "", the listener address will be 0.0.0.0:port
		return false
	}

	var conflictWithStaticListener, conflictWithVirtualListener bool

	// bind == wildcard
	// or bind unspecified, but protocol is HTTP
	if proxy.Metadata != nil {
		conflictWithStaticListener = proxy.Metadata.EnvoyStatusPort == port || proxy.Metadata.EnvoyPrometheusPort == port
	}
	if push != nil {
		conflictWithVirtualListener = int(push.Mesh.ProxyListenPort) == port || int(push.Mesh.ProxyInboundListenPort) == port
	}
	return conflictWithStaticListener || conflictWithVirtualListener
}

// isPrivilegedPort returns true if a given port is in the range 1-1023.
func isPrivilegedPort(port int) bool {
	// check for 0 is important because:
	// 1) technically, 0 is not a privileged port; any process can ask to bind to 0
	// 2) this function will be receiving 0 on input in the case of UDS listeners
	return 0 < port && port < 1024
}

func (node *Proxy) IsVM() bool {
	// TODO use node metadata to indicate that this is a VM instead of the TestVMLabel
	return node.Metadata.Labels[constants.TestVMLabel] != ""
}

func (node *Proxy) IsProxylessGrpc() bool {
	return node.Metadata != nil && node.Metadata.Generator == "grpc"
}

func (node *Proxy) GetNodeName() string {
	if node.Metadata != nil && len(node.Metadata.NodeName) > 0 {
		return node.Metadata.NodeName
	}
	return ""
}

func (node *Proxy) GetClusterID() cluster.ID {
	if node == nil || node.Metadata == nil {
		return ""
	}
	return node.Metadata.ClusterID
}

func (node *Proxy) GetNamespace() string {
	if node == nil || node.Metadata == nil {
		return ""
	}
	return node.Metadata.Namespace
}

func (node *Proxy) GetID() string {
	if node == nil {
		return ""
	}
	return node.ID
}

func (node *Proxy) FuzzValidate() bool {
	if node.Metadata == nil {
		return false
	}
	found := false
	for _, t := range NodeTypes {
		if node.Type == t {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	return len(node.IPAddresses) != 0
}

func (node *Proxy) EnableHBONEListen() bool {
	return node.IsAmbient() || (features.EnableSidecarHBONEListening && bool(node.Metadata.EnableHBONE))
}

func (node *Proxy) SetWorkloadEntry(name string, create bool) {
	node.Lock()
	defer node.Unlock()
	node.workloadEntryName = name
	node.workloadEntryAutoCreated = create
}

func (node *Proxy) WorkloadEntry() (string, bool) {
	node.RLock()
	defer node.RUnlock()
	return node.workloadEntryName, node.workloadEntryAutoCreated
}

// ShallowCloneWatchedResources clones the watched resources, both the keys and values are shallow copy.
func (node *Proxy) ShallowCloneWatchedResources() map[string]*WatchedResource {
	node.RLock()
	defer node.RUnlock()
	return maps.Clone(node.WatchedResources)
}

// DeepCloneWatchedResources clones the watched resources
func (node *Proxy) DeepCloneWatchedResources() map[string]WatchedResource {
	node.RLock()
	defer node.RUnlock()
	m := make(map[string]WatchedResource, len(node.WatchedResources))
	for k, v := range node.WatchedResources {
		m[k] = *v
	}
	return m
}

func (node *Proxy) GetWatchedResourceTypes() sets.String {
	node.RLock()
	defer node.RUnlock()

	ret := sets.NewWithLength[string](len(node.WatchedResources))
	for typeURL := range node.WatchedResources {
		ret.Insert(typeURL)
	}
	return ret
}

func (node *Proxy) GetWatchedResource(typeURL string) *WatchedResource {
	node.RLock()
	defer node.RUnlock()

	return node.WatchedResources[typeURL]
}

func (node *Proxy) NonceSent(typeURL string) string {
	node.RLock()
	defer node.RUnlock()

	wr := node.WatchedResources[typeURL]
	if wr != nil {
		return wr.NonceSent
	}
	return ""
}

func (node *Proxy) NonceAcked(typeURL string) string {
	node.RLock()
	defer node.RUnlock()

	wr := node.WatchedResources[typeURL]
	if wr != nil {
		return wr.NonceAcked
	}
	return ""
}

func (node *Proxy) Clusters() []string {
	node.RLock()
	defer node.RUnlock()
	wr := node.WatchedResources[v3.EndpointType]
	if wr != nil {
		return wr.ResourceNames.UnsortedList()
	}
	return nil
}

func (node *Proxy) NewWatchedResource(typeURL string, names []string) {
	node.Lock()
	defer node.Unlock()

	node.WatchedResources[typeURL] = &WatchedResource{TypeUrl: typeURL, ResourceNames: sets.New(names...)}
	// For all EDS requests that we have already responded with in the same stream let us
	// force the response. It is important to respond to those requests for Envoy to finish
	// warming of those resources(Clusters).
	// This can happen with the following sequence
	// 1. Envoy disconnects and reconnects to Istiod.
	// 2. Envoy sends EDS request and we respond with it.
	// 3. Envoy sends CDS request and we respond with clusters.
	// 4. Envoy detects a change in cluster state and tries to warm those clusters and send EDS request for them.
	// 5. We should respond to the EDS request with Endpoints to let Envoy finish cluster warming.
	// Refer to https://github.com/envoyproxy/envoy/issues/13009 for more details.
	for _, dependent := range WarmingDependencies(typeURL) {
		if dwr, exists := node.WatchedResources[dependent]; exists {
			dwr.AlwaysRespond = true
		}
	}
}

// WarmingDependencies returns the dependent typeURLs that need to be responded with
// for warming of this typeURL.
func WarmingDependencies(typeURL string) []string {
	switch typeURL {
	case v3.ClusterType:
		return []string{v3.EndpointType}
	default:
		return nil
	}
}

func (node *Proxy) AddOrUpdateWatchedResource(r *WatchedResource) {
	if r == nil {
		return
	}
	node.Lock()
	defer node.Unlock()
	node.WatchedResources[r.TypeUrl] = r
}

func (node *Proxy) UpdateWatchedResource(typeURL string, updateFn func(*WatchedResource) *WatchedResource) {
	node.Lock()
	defer node.Unlock()
	r := node.WatchedResources[typeURL]
	r = updateFn(r)
	if r != nil {
		node.WatchedResources[typeURL] = r
	} else {
		delete(node.WatchedResources, typeURL)
	}
}

func (node *Proxy) DeleteWatchedResource(typeURL string) {
	node.Lock()
	defer node.Unlock()

	delete(node.WatchedResources, typeURL)
}

type GatewayController interface {
	ConfigStoreController
	// Reconcile updates the internal state of the gateway controller for a given input. This should be
	// called before any List/Get calls if the state has changed
	Reconcile(ctx *PushContext)
	// SecretAllowed determines if a SDS credential is accessible to a given namespace.
	// For example, for resourceName of `kubernetes-gateway://ns-name/secret-name` and namespace of `ingress-ns`,
	// this would return true only if there was a policy allowing `ingress-ns` to access Secrets in the `ns-name` namespace.
	SecretAllowed(resourceName string, namespace string) bool
	Collection() krt.Collection[ADPResource]
}

// OutboundListenerClass is a helper to turn a NodeType for outbound to a ListenerClass.
func OutboundListenerClass(t NodeType) istionetworking.ListenerClass {
	if t == Router {
		return istionetworking.ListenerClassGateway
	}
	return istionetworking.ListenerClassSidecarOutbound
}

type ADPResource struct {
	Resource *api.Resource        `json:"resource"`
	Gateway  types.NamespacedName `json:"gateway"`
}

func (g ADPResource) ResourceName() string {
	switch t := g.Resource.Kind.(type) {
	case *api.Resource_Bind:
		return "bind/" + t.Bind.Key
	case *api.Resource_Listener:
		return "listener/" + t.Listener.Key
	case *api.Resource_Route:
		return "route/" + t.Route.Key
	}
	panic("unknown resource kind")
}

func (g ADPResource) Equals(other ADPResource) bool {
	return proto.Equal(g.Resource, other.Resource) && g.Gateway == other.Gateway
}

// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
"bytes"
"encoding/json"
"fmt"
"reflect"
"time"

gogojsonpb "github.com/gogo/protobuf/jsonpb" // nolint: depguard
gogoproto "github.com/gogo/protobuf/proto"   // nolint: depguard
gogotypes "github.com/gogo/protobuf/types"   // nolint: depguard
"google.golang.org/protobuf/proto"
"google.golang.org/protobuf/reflect/protoreflect"
"google.golang.org/protobuf/types/known/anypb"
"google.golang.org/protobuf/types/known/structpb"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/apimachinery/pkg/runtime/schema"
kubetypes "k8s.io/apimachinery/pkg/types"
"sigs.k8s.io/yaml"

"istio.io/api/label"
"istio.io/istio/pilot/pkg/util/protoconv"
"istio.io/istio/pkg/maps"
"istio.io/istio/pkg/ptr"
"istio.io/istio/pkg/slices"
"istio.io/istio/pkg/util/gogoprotomarshal"
"istio.io/istio/pkg/util/protomarshal"
"istio.io/istio/pkg/util/sets"
)

// Meta is metadata attached to each configuration unit.
// The revision is optional, and if provided, identifies the
// last update operation on the object.
type Meta struct {
	// GroupVersionKind is a short configuration name that matches the content message type
	// (e.g. "route-rule")
	GroupVersionKind GroupVersionKind `json:"type,omitempty"`

	// UID
	UID string `json:"uid,omitempty"`

	// Name is a unique immutable identifier in a namespace
	Name string `json:"name,omitempty"`

	// Namespace defines the space for names (optional for some types),
	// applications may choose to use namespaces for a variety of purposes
	// (security domains, fault domains, organizational domains)
	Namespace string `json:"namespace,omitempty"`

	// Domain defines the suffix of the fully qualified name past the namespace.
	// Domain is not a part of the unique key unlike name and namespace.
	Domain string `json:"domain,omitempty"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	Annotations map[string]string `json:"annotations,omitempty"`

	// ResourceVersion is an opaque identifier for tracking updates to the config registry.
	// The implementation may use a change index or a commit log for the revision.
	// The config client should not make any assumptions about revisions and rely only on
	// exact equality to implement optimistic concurrency of read-write operations.
	//
	// The lifetime of an object of a particular revision depends on the underlying data store.
	// The data store may compactify old revisions in the interest of storage optimization.
	//
	// An empty revision carries a special meaning that the associated object has
	// not been stored and assigned a revision.
	ResourceVersion string `json:"resourceVersion,omitempty"`

	// CreationTimestamp records the creation time
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`

	// OwnerReferences allows specifying in-namespace owning objects.
	OwnerReferences []metav1.OwnerReference `json:"ownerReferences,omitempty"`

	// A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
	Generation int64 `json:"generation,omitempty"`
}

// Config is a configuration unit consisting of the type of configuration, the
// key identifier that is unique per type, and the content represented as a
// protobuf message.
type Config struct {
	Meta

	// Spec holds the configuration object as a gogo protobuf message
	Spec Spec

	// Status holds long-running status.
	Status Status
}

func LabelsInRevision(lbls map[string]string, rev string) bool {
	configEnv, f := lbls[label.IoIstioRev.Name]
	if !f {
		// This is a global object, and always included
		return true
	}
	// If the revision is empty, this means we don't specify a revision, and
	// we should always include it
	if rev == "" {
		return true
	}
	// Otherwise, only return true if revisions equal
	return configEnv == rev
}

func LabelsInRevisionOrTags(lbls map[string]string, rev string, tags sets.Set[string]) bool {
	if LabelsInRevision(lbls, rev) {
		return true
	}
	configEnv := lbls[label.IoIstioRev.Name]
	// Otherwise, only return true if revisions equal
	return tags.Contains(configEnv)
}

func ObjectInRevision(o *Config, rev string) bool {
	return LabelsInRevision(o.Labels, rev)
}

// Spec defines the spec for the  In order to use below helper methods,
// this must be one of:
// * golang/protobuf Message
// * gogo/protobuf Message
// * Able to marshal/unmarshal using json
type Spec any

func ToProto(s Spec) (*anypb.Any, error) {
	// golang protobuf. Use protoreflect.ProtoMessage to distinguish from gogo
	// golang/protobuf 1.4+ will have this interface. Older golang/protobuf are gogo compatible
	// but also not used by Istio at all.
	if pb, ok := s.(protoreflect.ProtoMessage); ok {
		return protoconv.MessageToAnyWithError(pb)
	}

	// gogo protobuf
	if pb, ok := s.(gogoproto.Message); ok {
		gogoany, err := gogotypes.MarshalAny(pb)
		if err != nil {
			return nil, err
		}
		return &anypb.Any{
			TypeUrl: gogoany.TypeUrl,
			Value:   gogoany.Value,
		}, nil
	}

	js, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	pbs := &structpb.Struct{}
	if err := protomarshal.Unmarshal(js, pbs); err != nil {
		return nil, err
	}
	return protoconv.MessageToAnyWithError(pbs)
}

func ToMap(s Spec) (map[string]any, error) {
	js, err := ToJSON(s)
	if err != nil {
		return nil, err
	}

	// Unmarshal from json bytes to go map
	var data map[string]any
	err = json.Unmarshal(js, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ToRaw(s Spec) (json.RawMessage, error) {
	js, err := ToJSON(s)
	if err != nil {
		return nil, err
	}

	// Unmarshal from json bytes to go map
	return js, nil
}

func ToJSON(s Spec) ([]byte, error) {
	return toJSON(s, false)
}

func ToPrettyJSON(s Spec) ([]byte, error) {
	return toJSON(s, true)
}

func toJSON(s Spec, pretty bool) ([]byte, error) {
	indent := ""
	if pretty {
		indent = "    "
	}

	// golang protobuf. Use protoreflect.ProtoMessage to distinguish from gogo
	// golang/protobuf 1.4+ will have this interface. Older golang/protobuf are gogo compatible
	// but also not used by Istio at all.
	if _, ok := s.(protoreflect.ProtoMessage); ok {
		if pb, ok := s.(proto.Message); ok {
			b, err := protomarshal.MarshalIndent(pb, indent)
			return b, err
		}
	}

	b := &bytes.Buffer{}
	// gogo protobuf
	if pb, ok := s.(gogoproto.Message); ok {
		err := (&gogojsonpb.Marshaler{Indent: indent}).Marshal(b, pb)
		return b.Bytes(), err
	}
	if pretty {
		return json.MarshalIndent(s, "", indent)
	}
	return json.Marshal(s)
}

type deepCopier interface {
	DeepCopyInterface() any
}

func ApplyYAML(s Spec, yml string) error {
	js, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}
	return ApplyJSON(s, string(js))
}

func ApplyJSONStrict(s Spec, js string) error {
	// golang protobuf. Use protoreflect.ProtoMessage to distinguish from gogo
	// golang/protobuf 1.4+ will have this interface. Older golang/protobuf are gogo compatible
	// but also not used by Istio at all.
	if _, ok := s.(protoreflect.ProtoMessage); ok {
		if pb, ok := s.(proto.Message); ok {
			err := protomarshal.ApplyJSONStrict(js, pb)
			return err
		}
	}

	// gogo protobuf
	if pb, ok := s.(gogoproto.Message); ok {
		err := gogoprotomarshal.ApplyJSONStrict(js, pb)
		return err
	}

	d := json.NewDecoder(bytes.NewReader([]byte(js)))
	d.DisallowUnknownFields()
	return d.Decode(&s)
}

func ApplyJSON(s Spec, js string) error {
	// golang protobuf. Use protoreflect.ProtoMessage to distinguish from gogo
	// golang/protobuf 1.4+ will have this interface. Older golang/protobuf are gogo compatible
	// but also not used by Istio at all.
	if _, ok := s.(protoreflect.ProtoMessage); ok {
		if pb, ok := s.(proto.Message); ok {
			err := protomarshal.ApplyJSON(js, pb)
			return err
		}
	}

	// gogo protobuf
	if pb, ok := s.(gogoproto.Message); ok {
		err := gogoprotomarshal.ApplyJSON(js, pb)
		return err
	}

	return json.Unmarshal([]byte(js), &s)
}

func DeepCopy(s any) any {
	if s == nil {
		return nil
	}
	// If deep copy is defined, use that
	if dc, ok := s.(deepCopier); ok {
		return dc.DeepCopyInterface()
	}

	// golang protobuf. Use protoreflect.ProtoMessage to distinguish from gogo
	// golang/protobuf 1.4+ will have this interface. Older golang/protobuf are gogo compatible
	// but also not used by Istio at all.
	if _, ok := s.(protoreflect.ProtoMessage); ok {
		if pb, ok := s.(proto.Message); ok {
			return protomarshal.Clone(pb)
		}
	}

	// gogo protobuf
	if pb, ok := s.(gogoproto.Message); ok {
		return gogoproto.Clone(pb)
	}

	// If we don't have a deep copy method, we will have to do some reflection magic. Its not ideal,
	// but all Istio types have an efficient deep copy.
	js, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	data := reflect.New(reflect.TypeOf(s)).Interface()
	if err := json.Unmarshal(js, data); err != nil {
		return nil
	}
	data = reflect.ValueOf(data).Elem().Interface()
	return data
}

func (c *Config) Equals(other *Config) bool {
	am, bm := c.Meta, other.Meta
	if am.GroupVersionKind != bm.GroupVersionKind {
		return false
	}
	if am.UID != bm.UID {
		return false
	}
	if am.Name != bm.Name {
		return false
	}
	if am.Namespace != bm.Namespace {
		return false
	}
	if am.Domain != bm.Domain {
		return false
	}
	if !maps.Equal(am.Labels, bm.Labels) {
		return false
	}
	if !maps.Equal(am.Annotations, bm.Annotations) {
		return false
	}
	if am.ResourceVersion != bm.ResourceVersion {
		return false
	}
	if am.CreationTimestamp != bm.CreationTimestamp {
		return false
	}
	if !slices.EqualFunc(am.OwnerReferences, bm.OwnerReferences, func(a metav1.OwnerReference, b metav1.OwnerReference) bool {
		if a.APIVersion != b.APIVersion {
			return false
		}
		if a.Kind != b.Kind {
			return false
		}
		if a.Name != b.Name {
			return false
		}
		if a.UID != b.UID {
			return false
		}
		if !ptr.Equal(a.Controller, b.Controller) {
			return false
		}
		if !ptr.Equal(a.BlockOwnerDeletion, b.BlockOwnerDeletion) {
			return false
		}
		return true
	}) {
		return false
	}
	if am.Generation != bm.Generation {
		return false
	}

	if !equals(c.Spec, other.Spec) {
		return false
	}
	if !equals(c.Status, other.Status) {
		return false
	}
	return true
}

func equals(a any, b any) bool {
	if _, ok := a.(protoreflect.ProtoMessage); ok {
		if pb, ok := a.(proto.Message); ok {
			return proto.Equal(pb, b.(proto.Message))
		}
	}
	// We do NOT do gogo here. The reason is Kubernetes has hacked up almost-gogo types that do not allow Equals() calls

	return reflect.DeepEqual(a, b)
}

type Status any

// Key function for the configuration objects
func Key(grp, ver, typ, name, namespace string) string {
	return grp + "/" + ver + "/" + typ + "/" + namespace + "/" + name // Format: %s/%s/%s/%s/%s
}

// Key is the unique identifier for a configuration object
func (meta *Meta) Key() string {
	return Key(
		meta.GroupVersionKind.Group, meta.GroupVersionKind.Version, meta.GroupVersionKind.Kind,
		meta.Name, meta.Namespace)
}

func (meta *Meta) ToObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:              meta.Name,
		Namespace:         meta.Namespace,
		UID:               kubetypes.UID(meta.UID),
		ResourceVersion:   meta.ResourceVersion,
		Generation:        meta.Generation,
		CreationTimestamp: metav1.NewTime(meta.CreationTimestamp),
		Labels:            meta.Labels,
		Annotations:       meta.Annotations,
		OwnerReferences:   meta.OwnerReferences,
	}
}

func (c Config) DeepCopy() Config {
	var clone Config
	clone.Meta = c.Meta
	clone.Labels = maps.Clone(c.Labels)
	clone.Annotations = maps.Clone(clone.Annotations)
	clone.Spec = DeepCopy(c.Spec)
	if c.Status != nil {
		clone.Status = DeepCopy(c.Status)
	}
	return clone
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetNamespace() string {
	return c.Namespace
}

func (c Config) GetCreationTimestamp() time.Time {
	return c.CreationTimestamp
}

func (c Config) NamespacedName() kubetypes.NamespacedName {
	return kubetypes.NamespacedName{
		Namespace: c.Namespace,
		Name:      c.Name,
	}
}

var _ fmt.Stringer = GroupVersionKind{}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

func (g GroupVersionKind) String() string {
	return g.CanonicalGroup() + "/" + g.Version + "/" + g.Kind
}

// GroupVersion returns the group/version similar to what would be found in the apiVersion field of a Kubernetes resource.
func (g GroupVersionKind) GroupVersion() string {
	if g.Group == "" {
		return g.Version
	}
	return g.Group + "/" + g.Version
}

func FromKubernetesGVK(gvk schema.GroupVersionKind) GroupVersionKind {
	return GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
}

// Kubernetes returns the same GVK, using the Kubernetes object type
func (g GroupVersionKind) Kubernetes() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   g.Group,
		Version: g.Version,
		Kind:    g.Kind,
	}
}

func CanonicalGroup(group string) string {
	if group != "" {
		return group
	}
	return "core"
}

// CanonicalGroup returns the group with defaulting applied. This means an empty group will
// be treated as "core", following Kubernetes API standards
func (g GroupVersionKind) CanonicalGroup() string {
	return CanonicalGroup(g.Group)
}

// PatchFunc provides the cached config as a base for modification. Only diff the between the cfg
// parameter and the returned Config will be applied.
type PatchFunc func(cfg Config) (Config, kubetypes.PatchType)

type Namer interface {
	GetName() string
	GetNamespace() string
}

func NamespacedName[T Namer](o T) kubetypes.NamespacedName {
	return kubetypes.NamespacedName{
		Namespace: o.GetNamespace(),
		Name:      o.GetName(),
	}
}

