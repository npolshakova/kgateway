package agentgatewaysyncer

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"google.golang.org/protobuf/types/known/anypb"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pilot/pkg/util/protoconv"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/config/schema/gvk"
	kubeutil "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	kubelabels "istio.io/istio/pkg/kube/labels"
	"istio.io/istio/pkg/network"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/util/sets"
	corev1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
)

// Local tunnel protocol constants to avoid protobuf namespace conflict
const (
	TunnelProtocolNone  = 0
	TunnelProtocolHBone = 1
)

// WorkloadConfig represents the configuration for creating a workload
type WorkloadConfig struct {
	UID            string
	Name           string
	Namespace      string
	ServiceAccount string
	Addresses      []string
	TunnelProtocol int32 // Use int32 instead of api.TunnelProtocol to avoid conflict
	Node           string
	Network        string
	WorkloadType   api.WorkloadType
	Status         api.WorkloadStatus
	NetworkMode    api.NetworkMode
	Services       map[string]*ServicePortConfig
}

// ServicePortConfig represents the configuration for a service port
type ServicePortConfig struct {
	ServicePort uint32
	TargetPort  uint32
	AppProtocol api.AppProtocol
}

// generatePodUID generates a UID for a pod following the Istio pattern
func generatePodUID(clusterID string, namespace, name string) string {
	return clusterID + "//" + "Pod/" + namespace + "/" + name
}

// generatePodUIDFromPod generates a UID for a pod following the Istio pattern
func generatePodUIDFromPod(clusterID string, p *corev1.Pod) string {
	return generatePodUID(clusterID, p.Namespace, p.Name)
}

// getWorkloadStatus determines the health status of a workload
func getWorkloadStatus(ready bool, deletionTimestamp *metav1.Time) api.WorkloadStatus {
	if !ready || deletionTimestamp != nil {
		return api.WorkloadStatus_UNHEALTHY
	}
	return api.WorkloadStatus_HEALTHY
}

// createWorkloadResource creates a single workload resource from configuration
func createWorkloadResource(config WorkloadConfig) envoytypes.Resource {
	fmt.Printf("DEBUG: createWorkloadResource called with config %+v\n", config)

	// Convert service configurations to PortList
	services := make(map[string]*api.PortList)
	for serviceKey, portConfig := range config.Services {
		services[serviceKey] = &api.PortList{
			Ports: []*api.Port{
				{
					ServicePort: portConfig.ServicePort,
					TargetPort:  portConfig.TargetPort,
					AppProtocol: portConfig.AppProtocol,
				},
			},
		}
	}

	workload := &api.Workload{
		Uid:            config.UID,
		Name:           config.Name,
		Namespace:      config.Namespace,
		ServiceAccount: config.ServiceAccount,
		// Addresses:      addresses, // Temporarily removed to fix parsing error
		TunnelProtocol: api.TunnelProtocol(config.TunnelProtocol),
		Node:           config.Node,
		Network:        config.Network,
		WorkloadType:   config.WorkloadType,
		Status:         config.Status,
		NetworkMode:    config.NetworkMode,
		Services:       services,
		ClusterId:      "cluster1",      // TODO: make configurable
		TrustDomain:    "cluster.local", // TODO: make configurable
	}

	addressResource := &api.Address{
		Type: &api.Address_Workload{
			Workload: workload,
		},
	}

	result := &envoyResourceWithCustomName{
		Message: addressResource,
		Name:    config.UID,
		version: utils.HashProto(addressResource),
	}

	fmt.Printf("DEBUG: Successfully created workload resource with name %s\n", result.Name)
	return result
}

// selectedWorkload adds the following to LocalityPod:
// * fields specific to workload entry (portMapping, network, weight)
// Usable with FilterSelect
type selectedWorkload struct {
	// the workload that is selected
	krtcollections.LocalityPod

	// workload entry has workload-level port mappings
	portMapping map[string]uint32
	// network id (istio concept of network)
	network string
	// weight from workloadentry
	weight uint32
}

// TargetRef is a subset of the Kubernetes ObjectReference which has some fields we don't care about
type TargetRef struct {
	Kind      string
	Namespace string
	Name      string
	UID       types.UID
}

func (t TargetRef) String() string {
	return t.Kind + "/" + t.Namespace + "/" + t.Name + "/" + string(t.UID)
}

type Node struct {
	Name     string
	Locality *api.Locality
}

// NetworkGateway is the gateway of a network
type NetworkGateway struct {
	// Network is the ID of the network where this Gateway resides.
	Network network.ID
	// Cluster is the ID of the k8s cluster where this Gateway resides.
	Cluster cluster.ID
	// gateway ip address
	Addr string
	// gateway port
	Port uint32
	// HBONEPort if non-zero indicates that the gateway supports HBONE
	HBONEPort uint32
	// ServiceAccount the gateway runs as
	ServiceAccount types.NamespacedName
}

// buildWorkloadsCollection builds out the core Workload object type used in agentgateway mode.
// A Workload represents a single addressable unit of compute -- typically a Pod or a VM.
// Workloads can come from a variety of sources; these are joined together to build one complete `Collection[WorkloadInfo]`.
func buildWorkloadsCollection(
	pods krt.Collection[*corev1.Pod],
	workloadServices krt.Collection[ServiceInfo],
	endpointSlices krt.Collection[*discovery.EndpointSlice],
	domainSuffix, clusterId string,
) krt.Collection[WorkloadInfo] {
	endpointSlicesAddressIndex := endpointSliceAddressIndex(endpointSlices)
	// Workloads coming from pods. There should be one workload for each (running) Pod.
	return krt.NewCollection(
		pods,
		podWorkloadBuilder(
			workloadServices,
			endpointSlices,
			endpointSlicesAddressIndex,
			domainSuffix, clusterId,
		),
	)
}

// endpointSliceAddressIndex builds an index from IP Address
func endpointSliceAddressIndex(EndpointSlices krt.Collection[*discovery.EndpointSlice]) krt.Index[TargetRef, *discovery.EndpointSlice] {
	return krt.NewIndex(EndpointSlices, func(es *discovery.EndpointSlice) []TargetRef {
		if es.AddressType == discovery.AddressTypeFQDN {
			// Currently we do not support FQDN.
			return nil
		}
		_, f := es.Labels[discovery.LabelServiceName]
		if !f {
			// Not for a service; we don't care about it.
			return nil
		}
		res := make([]TargetRef, 0, len(es.Endpoints))
		for _, ep := range es.Endpoints {
			if ep.TargetRef == nil || ep.TargetRef.Kind != gvk.Pod.Kind {
				// We only want pods here
				continue
			}
			tr := TargetRef{
				Kind:      ep.TargetRef.Kind,
				Namespace: ep.TargetRef.Namespace,
				Name:      ep.TargetRef.Name,
				UID:       ep.TargetRef.UID,
			}
			res = append(res, tr)
		}
		return res
	})
}

func podWorkloadBuilder(
	workloadServices krt.Collection[ServiceInfo],
	endpointSlices krt.Collection[*discovery.EndpointSlice],
	endpointSlicesAddressIndex krt.Index[TargetRef, *discovery.EndpointSlice],
	domainSuffix, clusterId string,
) krt.TransformationSingle[*corev1.Pod, WorkloadInfo] {
	return func(ctx krt.HandlerContext, p *corev1.Pod) *WorkloadInfo {
		// TODO: Pod Is Pending but have a pod IP should be a valid workload, we should build it
		// See https://github.com/istio/istio/issues/48854

		k8sPodIPs := getPodIPs(p)
		if len(k8sPodIPs) == 0 {
			return nil
		}
		podIPs, err := slices.MapErr(k8sPodIPs, func(e corev1.PodIP) ([]byte, error) {
			n, err := netip.ParseAddr(e.IP)
			if err != nil {
				return nil, err
			}
			return n.AsSlice(), nil
		})
		if err != nil {
			// Is this possible? Probably not in typical case, but anyone could put garbage there.
			return nil
		}
		services := krt.Fetch(ctx, workloadServices)
		services = append(services, matchingServicesWithoutSelectors(ctx, p, services, workloadServices, endpointSlices, endpointSlicesAddressIndex, domainSuffix)...)

		// Logic from https://github.com/kubernetes/kubernetes/blob/7c873327b679a70337288da62b96dd610858181d/staging/src/k8s.io/endpointslice/utils.go#L37
		// Kubernetes has Ready, Serving, and Terminating. We only have a boolean, which is sufficient for our cases
		status := api.WorkloadStatus_HEALTHY
		if !IsPodReady(p) || p.DeletionTimestamp != nil {
			status = api.WorkloadStatus_UNHEALTHY
		}

		w := &api.Workload{
			Uid:            generatePodUID(clusterId, p.Name, p.Namespace),
			Name:           p.Name,
			Namespace:      p.Namespace,
			ClusterId:      clusterId,
			Addresses:      podIPs,
			ServiceAccount: p.Spec.ServiceAccountName,
			Services:       constructServices(p, services),
			Status:         status,
			// TODO: support other fields (Locality, Node, Network, etc.)
		}

		if p.Spec.HostNetwork {
			w.NetworkMode = api.NetworkMode_HOST_NETWORK
		}

		w.WorkloadName = workloadNameForPod(p)
		w.WorkloadType = api.WorkloadType_POD // backwards compatibility
		w.CanonicalName, w.CanonicalRevision = kubelabels.CanonicalService(p.Labels, w.WorkloadName)

		return precomputeWorkloadPtr(&WorkloadInfo{
			Workload:     w,
			Labels:       p.Labels,
			Source:       wellknown.PodGVK.Kind,
			CreationTime: p.CreationTimestamp.Time,
		})
	}
}

func precomputeWorkloadPtr(w *WorkloadInfo) *WorkloadInfo {
	return ptr.Of(precomputeWorkload(*w))
}

func precomputeWorkload(w WorkloadInfo) WorkloadInfo {
	addr := workloadToAddress(w.Workload)
	w.MarshaledAddress = protoconv.MessageToAny(addr)
	w.AsAddress = AddressInfo{
		Address:   addr,
		Marshaled: w.MarshaledAddress,
	}
	return w
}

func workloadToAddress(w *api.Workload) *api.Address {
	return &api.Address{
		Type: &api.Address_Workload{
			Workload: w,
		},
	}
}

func workloadNameForPod(pod *corev1.Pod) string {
	objMeta, _ := kubeutil.GetWorkloadMetaFromPod(pod)
	return objMeta.Name
}

type AddressInfo struct {
	*api.Address
	Marshaled *anypb.Any
}

type ServiceInfo struct {
	Service *api.Service
	// LabelSelectors for the Service. Note these are only used internally, not sent over XDS
	LabelSelector LabelSelector
	// PortNames provides a mapping of ServicePort -> port names. Note these are only used internally, not sent over XDS
	PortNames map[int32]ServicePortName
	// Source is the type that introduced this service.
	Source TypedObject
	// MarshaledAddress contains the pre-marshaled representation.
	// Note: this is an Address -- not a Service.
	MarshaledAddress *anypb.Any
	// AsAddress contains a pre-created AddressInfo representation. This ensures we do not need repeated conversions on
	// the hotpath
	AsAddress AddressInfo
}

func (s ServiceInfo) GetNamespace() string {
	return s.Service.Namespace
}

func (s ServiceInfo) ResourceName() string {
	return namespacedHostname(s.Service.Namespace, s.Service.Hostname)
}

var _ krt.ResourceNamer = ServiceInfo{}

type LabelSelector struct {
	Labels map[string]string
}

type TypedObject struct {
	types.NamespacedName
	Kind string
}

type ServicePortName struct {
	PortName       string
	TargetPortName string
}

type WorkloadInfo struct {
	Workload *api.Workload
	// Labels for the workload. Note these are only used internally, not sent over XDS
	Labels map[string]string
	// Source is the kind type that introduced this workload.
	Source string
	// CreationTime is the time when the workload was created. Note this is used internally only.
	CreationTime time.Time
	// MarshaledAddress contains the pre-marshaled representation.
	// Note: this is an Address -- not a Workload.
	MarshaledAddress *anypb.Any
	// AsAddress contains a pre-created AddressInfo representation. This ensures we do not need repeated conversions on
	// the hotpath
	AsAddress AddressInfo
}

func (w WorkloadInfo) ResourceName() string {
	return fmt.Sprintf("%s/%s", w.Workload.Namespace, w.Workload.Name)
}

var _ krt.ResourceNamer = WorkloadInfo{}

func getPodIPs(p *corev1.Pod) []corev1.PodIP {
	k8sPodIPs := p.Status.PodIPs
	if len(k8sPodIPs) == 0 && p.Status.PodIP != "" {
		k8sPodIPs = []corev1.PodIP{{IP: p.Status.PodIP}}
	}
	return k8sPodIPs
}

// IsPodReady is copied from kubernetes/pkg/api/v1/pod/utils.go
func IsPodReady(pod *corev1.Pod) bool {
	return IsPodReadyConditionTrue(pod.Status)
}

// IsPodReadyConditionTrue returns true if a pod is ready; false otherwise.
func IsPodReadyConditionTrue(status corev1.PodStatus) bool {
	condition := GetPodReadyCondition(status)
	return condition != nil && condition.Status == corev1.ConditionTrue
}

func GetPodReadyCondition(status corev1.PodStatus) *corev1.PodCondition {
	_, condition := GetPodCondition(&status, corev1.PodReady)
	return condition
}

func GetPodCondition(status *corev1.PodStatus, conditionType corev1.PodConditionType) (int, *corev1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	return GetPodConditionFromList(status.Conditions, conditionType)
}

// GetPodConditionFromList extracts the provided condition from the given list of condition and
// returns the index of the condition and the condition. Returns -1 and nil if the condition is not present.
func GetPodConditionFromList(conditions []corev1.PodCondition, conditionType corev1.PodConditionType) (int, *corev1.PodCondition) {
	if conditions == nil {
		return -1, nil
	}
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return i, &conditions[i]
		}
	}
	return -1, nil
}

// matchingServicesWithoutSelectors finds all Services that match a given pod that do not use selectors.
// See https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors for more info.
// For selector service, we query by the selector elsewhere, so this only handles the services that are NOT already found
// by a selector.
// For EndpointSlices that happen to point to the same IP as the pod, but are not directly bound to the pod (via TargetRef),
// we ignore them here. These will produce a Workload directly from the EndpointSlice, but with limited information;
// we do not implicitly merge a Pod with an EndpointSlice just based on IP.
func matchingServicesWithoutSelectors(
	ctx krt.HandlerContext,
	p *corev1.Pod,
	alreadyMatchingServices []ServiceInfo,
	workloadServices krt.Collection[ServiceInfo],
	endpointSlices krt.Collection[*discovery.EndpointSlice],
	endpointSlicesAddressIndex krt.Index[TargetRef, *discovery.EndpointSlice],
	domainSuffix string,
) []ServiceInfo {
	var res []ServiceInfo
	// Build out our set of already-matched services to avoid double-selecting a service
	seen := sets.NewWithLength[string](len(alreadyMatchingServices))
	for _, s := range alreadyMatchingServices {
		seen.Insert(s.Service.Hostname)
	}
	tr := TargetRef{
		Kind:      gvk.Pod.Kind,
		Namespace: p.Namespace,
		Name:      p.Name,
		UID:       p.UID,
	}
	// For each IP, find any endpointSlices referencing it.
	matchedSlices := krt.Fetch(ctx, endpointSlices, krt.FilterIndex(endpointSlicesAddressIndex, tr))
	for _, es := range matchedSlices {
		serviceName, f := es.Labels[discovery.LabelServiceName]
		if !f {
			// Not for a service; we don't care about it.
			continue
		}
		hostname := string(kube.ServiceHostname(serviceName, es.Namespace, domainSuffix))
		if seen.Contains(hostname) {
			// We already know about this service
			continue
		}
		// This pod is included in the EndpointSlice. We need to fetch the Service object for it, by key.
		serviceKey := es.Namespace + "/" + hostname
		svcs := krt.Fetch(ctx, workloadServices, krt.FilterKey(serviceKey), krt.FilterGeneric(func(a any) bool {
			// Only find Service, not Service Entry
			return a.(ServiceInfo).Source.Kind == wellknown.ServiceGVK.Kind
		}))
		if len(svcs) == 0 {
			// no service found
			continue
		}
		// There SHOULD only be one. This is only for `Service` which has unique hostnames.
		svc := svcs[0]
		res = append(res, svc)
	}
	return res
}

func constructServices(p *corev1.Pod, services []ServiceInfo) map[string]*api.PortList {
	res := map[string]*api.PortList{}
	for _, svc := range services {
		n := namespacedHostname(svc.Service.Namespace, svc.Service.Hostname)
		pl := &api.PortList{
			Ports: make([]*api.Port, 0, len(svc.Service.Ports)),
		}
		res[n] = pl
		for _, port := range svc.Service.Ports {
			targetPort := port.TargetPort
			// The svc.Ports represents the api.Service, which drops the port name info and just has numeric target Port.
			// TargetPort can be 0 which indicates its a named port. Check if its a named port and replace with the real targetPort if so.
			if named, f := svc.PortNames[int32(port.ServicePort)]; f && named.TargetPortName != "" {
				// Pods only match on TargetPort names
				tp, ok := FindPortName(p, named.TargetPortName)
				if !ok {
					// Port not present for this workload. Exclude the port entirely
					continue
				}
				targetPort = uint32(tp)
			}

			pl.Ports = append(pl.Ports, &api.Port{
				ServicePort: port.ServicePort,
				TargetPort:  targetPort,
			})
		}
	}
	return res
}

func FindPortName(pod *corev1.Pod, name string) (int32, bool) {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.Name == name {
				return port.ContainerPort, true
			}
		}
	}
	return 0, false
}

func namespacedHostname(namespace, hostname string) string {
	return namespace + "/" + hostname
}
