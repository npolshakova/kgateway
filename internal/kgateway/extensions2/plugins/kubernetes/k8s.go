package kubernetes

import (
	"context"
	"fmt"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/network"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/settings"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

const BackendClusterPrefix = "kube"

func NewPlugin(ctx context.Context, commoncol *common.CommonCollections) extensionsplug.Plugin {
	epSliceClient := kclient.New[*discoveryv1.EndpointSlice](commoncol.Client)
	endpointSlices := krt.WrapClient(epSliceClient, commoncol.KrtOpts.ToOptions("EndpointSlices")...)
	return NewPluginFromCollections(ctx, commoncol.KrtOpts, commoncol.Pods, commoncol.Services, endpointSlices, commoncol.Settings)
}

func NewPluginFromCollections(
	ctx context.Context,
	krtOpts krtutil.KrtOptions,
	pods krt.Collection[krtcollections.LocalityPod],
	services krt.Collection[*corev1.Service],
	endpointSlices krt.Collection[*discoveryv1.EndpointSlice],
	stngs settings.Settings,
) extensionsplug.Plugin {
	k8sServiceBackends := krt.NewManyCollection(services, func(kctx krt.HandlerContext, svc *corev1.Service) []ir.BackendObjectIR {
		uss := []ir.BackendObjectIR{}
		for _, port := range svc.Spec.Ports {
			uss = append(uss, BuildServiceBackendObjectIR(svc, port.Port, ptr.OrDefault(port.AppProtocol, port.Name)))
		}
		return uss
	}, krtOpts.ToOptions("KubernetesServiceBackends")...)

	inputs := krtcollections.NewGlooK8sEndpointInputs(stngs, krtOpts, endpointSlices, pods, k8sServiceBackends)
	k8sServiceEndpoints := krtcollections.NewK8sEndpoints(ctx, inputs)

	return extensionsplug.Plugin{
		ContributesBackends: map[schema.GroupKind]extensionsplug.BackendPlugin{
			wellknown.ServiceGVK.GroupKind(): {
				BackendInit: ir.BackendInit{
					InitBackend: processBackend,
				},
				Endpoints: k8sServiceEndpoints,
				Backends:  k8sServiceBackends,
			},
		},
		// TODO consider ContibutesPolicies allowing backendRef by networking.istio.io/Hostname
		// wellknown.ServiceGCK.GroupKind(): extensionsplug.PolicyPlugin{
		// 	GetBackendForRef: getBackendForHostnameRef,
		// },
	}
}

func BuildServiceBackendObjectIR(svc *corev1.Service, svcPort int32, svcProtocol string) ir.BackendObjectIR {
	var hostname string
	if svc.Spec.Type == corev1.ServiceTypeExternalName {
		hostname = svc.Spec.ExternalName
	} else {
		hostname = fmt.Sprintf("%s.%s.svc.%s", svc.Name, svc.Namespace, network.GetClusterDomainName())
	}

	return ir.BackendObjectIR{
		ObjectSource: ir.ObjectSource{
			Kind:      wellknown.ServiceGVK.Kind,
			Group:     wellknown.ServiceGVK.Group,
			Namespace: svc.Namespace,
			Name:      svc.Name,
		},
		Obj: svc,
		// TODO: fill in ObjIR
		Port:        svcPort,
		AppProtocol: ir.ParseAppProtocol(&svcProtocol),
		GvPrefix:    BackendClusterPrefix,
		// TODO: reevaluate knative dep, dedupe with pkg/utils/kubeutils/dns.go
		CanonicalHostname: hostname,
	}
}

func processBackend(ctx context.Context, in ir.BackendObjectIR, out *envoy_config_cluster_v3.Cluster) {
	svc, ok := in.Obj.(*corev1.Service)
	if ok && svc.Spec.Type == corev1.ServiceTypeExternalName {
		// external name service needs dns
		out.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
		}

		// need to set the LoadAssignment to make sure Envoy resolves the right DNS name + port combination,
		// otherwise it will use the cluster name
		out.LoadAssignment = &envoy_config_endpoint_v3.ClusterLoadAssignment{
			ClusterName: out.GetName(),
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{{
				LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{{
					HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
						Endpoint: &envoy_config_endpoint_v3.Endpoint{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{
										Address: in.CanonicalHostname,
										PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
											PortValue: uint32(in.Port),
										},
									},
								},
							},
						},
					},
				}},
			}},
		}
	} else {
		// Default: EDS
		out.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}

		out.EdsClusterConfig = &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
			EdsConfig: &envoy_config_core_v3.ConfigSource{
				ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
			},
		}
	}
}
