package agentgatewaysyncer

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	udpa "github.com/cncf/xds/go/udpa/type/v1"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gogoproto "github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"istio.io/istio/pkg/config/schema/kind"
	"istio.io/istio/pkg/util/hash"
	"istio.io/istio/pkg/util/protomarshal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
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

type ADPCacheResource struct {
	Gateway types.NamespacedName `json:"gateway"`
	reports reports.ReportMap

	Resources envoycache.Resources

	VersionMap map[string]map[string]string
}

func (r ADPCacheResource) ResourceName() string {
	return fmt.Sprintf("%s~%s", r.Gateway.Namespace, r.Gateway.Name)
}

func (r ADPCacheResource) Equals(in ADPCacheResource) bool {
	return r.Gateway == in.Gateway &&
		report{r.reports}.Equals(report{in.reports}) &&
		r.Resources.Version == in.Resources.Version
}

type ADPCacheAddress struct {
	NamespacedName types.NamespacedName
	Address        envoycache.Resources `json:"address"`
	proxyKey       string

	reports    reports.ReportMap
	VersionMap map[string]map[string]string
}

func (r ADPCacheAddress) ResourceName() string {
	return r.proxyKey
}

func (r ADPCacheAddress) Equals(in ADPCacheAddress) bool {
	return report{r.reports}.Equals(report{in.reports}) &&
		r.Address.Version == in.Address.Version
}

type ADPResource struct {
	Resource *api.Resource        `json:"resource"`
	Gateway  types.NamespacedName `json:"gateway"`

	reports reports.ReportMap
}

func (g ADPResource) ResourceName() string {
	switch t := g.Resource.GetKind().(type) {
	case *api.Resource_Bind:
		return "bind/" + t.Bind.GetKey()
	case *api.Resource_Listener:
		return "listener/" + t.Listener.GetKey()
	case *api.Resource_Route:
		return "route/" + t.Route.GetKey()
	}
	panic("unknown resource kind")
}

func (g ADPResource) Equals(other ADPResource) bool {
	return proto.Equal(g.Resource, other.Resource) && g.Gateway == other.Gateway
}

// Meta is metadata attached to each configuration unit.
// The revision is optional, and if provided, identifies the
// last update operation on the object.
type Meta struct {
	// GroupVersionKind is a short configuration name that matches the content message type
	// (e.g. "route-rule")
	GroupVersionKind schema.GroupVersionKind `json:"type,omitempty"`

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

type TypedResource struct {
	Kind schema.GroupVersionKind
	Name types.NamespacedName
}

// Spec defines the spec for the  In order to use below helper methods,
// this must be one of:
// * golang/protobuf Message
// * gogo/protobuf Message
// * Able to marshal/unmarshal using json
type Spec any

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
		UID:               types.UID(meta.UID),
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

type deepCopier interface {
	DeepCopyInterface() any
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

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetNamespace() string {
	return c.Namespace
}

func (c Config) GetCreationTimestamp() time.Time {
	return c.CreationTimestamp
}

func (c Config) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: c.Namespace,
		Name:      c.Name,
	}
}

type Index[K comparable, O any] interface {
	Lookup(k K) []O
	// AsCollection(opts ...CollectionOption) Collection[IndexObject[K, O]]
	objectHasKey(obj O, k K) bool
	extractKeys(o O) []K
	LookupCount(k K) int
}

type IndexObject[K comparable, O any] struct {
	Key     K
	Objects []O
}

func (i IndexObject[K, O]) ResourceName() string {
	return toString(i.Key)
}

func toString(rk any) string {
	tk, ok := rk.(string)
	if !ok {
		return rk.(fmt.Stringer).String()
	}
	return tk
}
