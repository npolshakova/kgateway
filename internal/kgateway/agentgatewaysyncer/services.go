package agentgatewaysyncer

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pilot/pkg/util/protoconv"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	corev1 "k8s.io/api/core/v1"
)

func buildServicesCollection(
	services krt.Collection[*corev1.Service],
	domainSuffix string,
) krt.Collection[ServiceInfo] {
	return krt.NewCollection(services, serviceServiceBuilder(domainSuffix))
}

func serviceServiceBuilder(
	domainSuffix string,
) krt.TransformationSingle[*corev1.Service, ServiceInfo] {
	return func(ctx krt.HandlerContext, s *corev1.Service) *ServiceInfo {
		if s.Spec.Type == corev1.ServiceTypeExternalName {
			// TODO: check if ExternalName services still work
			return nil
		}
		portNames := map[int32]ServicePortName{}
		for _, p := range s.Spec.Ports {
			portNames[p.Port] = ServicePortName{
				PortName:       p.Name,
				TargetPortName: p.TargetPort.StrVal,
			}
		}

		svc := constructService(ctx, s, domainSuffix)
		return precomputeServicePtr(&ServiceInfo{
			Service:       svc,
			PortNames:     portNames,
			LabelSelector: NewSelector(s.Spec.Selector),
			Source:        MakeSource(s),
		})
	}
}

// MakeSource is a helper to turn an Object into a TypedObject.
func MakeSource(o controllers.Object) TypedObject {
	return TypedObject{
		NamespacedName: config.NamespacedName(o),
		Kind:           o.GetObjectKind().GroupVersionKind().Kind,
	}
}

func NewSelector(l map[string]string) LabelSelector {
	return LabelSelector{l}
}

func constructService(
	ctx krt.HandlerContext,
	svc *corev1.Service,
	domainSuffix string,
) *api.Service {
	ports := make([]*api.Port, 0, len(svc.Spec.Ports))
	for _, p := range svc.Spec.Ports {
		var appProtocol api.AppProtocol
		if p.AppProtocol != nil {
			switch strings.ToLower(*p.AppProtocol) {
			case "grpc":
				appProtocol = api.AppProtocol_GRPC
			case "http":
				appProtocol = api.AppProtocol_HTTP11
			case "http2":
				appProtocol = api.AppProtocol_HTTP11
			}
		}

		ports = append(ports, &api.Port{
			ServicePort: uint32(p.Port),
			TargetPort:  uint32(p.TargetPort.IntVal),
			AppProtocol: appProtocol,
		})
	}

	addresses, err := slices.MapErr(getVIPs(svc), func(e string) (*api.NetworkAddress, error) {
		return toNetworkAddress(ctx, e)
	})
	if err != nil {
		logger.Error("fail to parse service", "svc", config.NamespacedName(svc), "error", err)
		return nil
	}

	var lb *api.LoadBalancing

	// First, use internal traffic policy if set.
	if itp := svc.Spec.InternalTrafficPolicy; itp != nil && *itp == corev1.ServiceInternalTrafficPolicyLocal {
		lb = &api.LoadBalancing{
			// Only allow endpoints on the same node.
			RoutingPreference: []api.LoadBalancing_Scope{
				api.LoadBalancing_NODE,
			},
			Mode: api.LoadBalancing_STRICT,
		}
	}
	// TODO: Check traffic distribution configuration

	if svc.Spec.PublishNotReadyAddresses {
		if lb == nil {
			lb = &api.LoadBalancing{}
		}
		lb.HealthPolicy = api.LoadBalancing_ALLOW_ALL
	}

	ipFamily := api.IPFamilies_AUTOMATIC
	if len(svc.Spec.IPFamilies) == 2 {
		ipFamily = api.IPFamilies_DUAL
	} else if len(svc.Spec.IPFamilies) == 1 {
		family := svc.Spec.IPFamilies[0]
		if family == corev1.IPv4Protocol {
			ipFamily = api.IPFamilies_IPV4_ONLY
		} else {
			ipFamily = api.IPFamilies_IPV6_ONLY
		}
	}
	// This is only checking one cluster - we'll merge later in the nested join to make sure
	// we get service VIPs from other clusters
	return &api.Service{
		Name:          svc.Name,
		Namespace:     svc.Namespace,
		Hostname:      string(kube.ServiceHostname(svc.Name, svc.Namespace, domainSuffix)),
		Addresses:     addresses,
		Ports:         ports,
		LoadBalancing: lb,
		IpFamilies:    ipFamily,
	}
}

func getVIPs(svc *corev1.Service) []string {
	res := []string{}
	cips := svc.Spec.ClusterIPs
	if len(cips) == 0 {
		cips = []string{svc.Spec.ClusterIP}
	}
	for _, cip := range cips {
		if cip != "" && cip != corev1.ClusterIPNone {
			res = append(res, cip)
		}
	}
	return res
}

func precomputeServicePtr(w *ServiceInfo) *ServiceInfo {
	return ptr.Of(precomputeService(*w))
}

func precomputeService(w ServiceInfo) ServiceInfo {
	addr := serviceToAddress(w.Service)
	w.MarshaledAddress = protoconv.MessageToAny(addr)
	w.AsAddress = AddressInfo{
		Address:   addr,
		Marshaled: w.MarshaledAddress,
	}
	return w
}

func serviceToAddress(s *api.Service) *api.Address {
	return &api.Address{
		Type: &api.Address_Service{
			Service: s,
		},
	}
}

func toNetworkAddress(ctx krt.HandlerContext, vip string) (*api.NetworkAddress, error) {
	ip, err := netip.ParseAddr(vip)
	if err != nil {
		return nil, fmt.Errorf("parse %v: %v", vip, err)
	}
	// TODO: support network id
	return &api.NetworkAddress{
		Address: ip.AsSlice(),
	}, nil
}
