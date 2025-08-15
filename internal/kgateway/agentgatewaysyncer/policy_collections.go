package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugins/trafficpolicy/agentgateway"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	a2aProtocol = "kgateway.dev/a2a"
)

func ADPPolicyCollection(inputs Inputs, krtopts krtutil.KrtOptions, plugins pluginsdk.Plugin) krt.Collection[ADPResourcesForGateway] {
	domainSuffix := kubeutils.GetClusterDomainName()

	inference := krt.NewManyCollection(inputs.InferencePools, func(ctx krt.HandlerContext, i *inf.InferencePool) []ADPPolicy {
		// 'service/{namespace}/{hostname}:{port}'
		svc := fmt.Sprintf("service/%v/%v.%v.inference.%v:%v", i.Namespace, i.Name, i.Namespace, domainSuffix, i.Spec.TargetPortNumber)
		er := i.Spec.ExtensionRef
		if er == nil {
			return nil
		}
		erf := er.ExtensionReference
		if erf.Group != nil && *erf.Group != "" {
			return nil
		}

		if erf.Kind != nil && *erf.Kind != "Service" {
			return nil
		}
		eppPort := ptr.OrDefault(erf.PortNumber, 9002)

		eppSvc := fmt.Sprintf("%v/%v.%v.svc.%v",
			i.Namespace, erf.Name, i.Namespace, domainSuffix)
		eppPolicyTarget := fmt.Sprintf("service/%v:%v",
			eppSvc, eppPort)

		failureMode := api.PolicySpec_InferenceRouting_FAIL_CLOSED
		if er.FailureMode == nil || *er.FailureMode == inf.FailOpen {
			failureMode = api.PolicySpec_InferenceRouting_FAIL_OPEN
		}
		inferencePolicy := &api.Policy{
			Name:   i.Namespace + "/" + i.Name + ":inference",
			Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: svc}},
			Spec: &api.PolicySpec{
				Kind: &api.PolicySpec_InferenceRouting_{
					InferenceRouting: &api.PolicySpec_InferenceRouting{
						EndpointPicker: &api.BackendReference{
							Kind: &api.BackendReference_Service{Service: eppSvc},
							Port: uint32(eppPort),
						},
						FailureMode: failureMode,
					},
				},
			},
		}

		// TODO: we would want some way if they explicitly set a BackendTLSPolicy for the EPP to respect that
		inferencePolicyTLS := &api.Policy{
			Name:   i.Namespace + "/" + i.Name + ":inferencetls",
			Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: eppPolicyTarget}},
			Spec: &api.PolicySpec{
				Kind: &api.PolicySpec_BackendTls{
					BackendTls: &api.PolicySpec_BackendTLS{
						// The spec mandates this :vomit:
						Insecure: wrappers.Bool(true),
					},
				},
			},
		}

		return []ADPPolicy{{inferencePolicy}, {inferencePolicyTLS}}
	}, krtopts.ToOptions("InferencePoolPolicies")...)

	a2a := krt.NewManyCollection(inputs.Services, func(ctx krt.HandlerContext, svc *corev1.Service) []ADPPolicy {
		var a2aPolicies []ADPPolicy
		for _, port := range svc.Spec.Ports {
			if port.AppProtocol != nil && *port.AppProtocol == a2aProtocol {
				svcRef := fmt.Sprintf("%v/%v", svc.Namespace, svc.Name)
				a2aPolicies = append(a2aPolicies, ADPPolicy{&api.Policy{
					Name:   fmt.Sprintf("a2a/%s/%s/%d", svc.Namespace, svc.Name, port.Port),
					Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: svcRef}},
					Spec: &api.PolicySpec{Kind: &api.PolicySpec_A2A_{
						A2A: &api.PolicySpec_A2A{},
					}},
				}})
			}
		}
		return a2aPolicies
	}, krtopts.ToOptions("A2APolicies")...)

	// For now, we apply all policies to all gateways. In the future, we can more precisely bind them to only relevant ones
	policiesByGateway := krt.NewCollection(inputs.GatewayIndex.Gateways, func(ctx krt.HandlerContext, i ir.Gateway) *ADPResourcesForGateway {
		var allResources []*api.Resource

		// Add inference policies
		inferences := slices.Map(krt.Fetch(ctx, inference), func(e ADPPolicy) *api.Resource {
			return toADPResources(e)[0]
		})
		allResources = append(allResources, inferences...)

		// Add A2A policies
		a2aPolicies := slices.Map(krt.Fetch(ctx, a2a), func(e ADPPolicy) *api.Resource {
			return toADPResources(e)[0]
		})
		allResources = append(allResources, a2aPolicies...)

		// Add all policies attached to the Gateway
		var attachedGatewayPolicyResources []*api.Resource
		gwPolicyTarget := &api.PolicyTarget{
			Kind: &api.PolicyTarget_Gateway{
				Gateway: i.Namespace + "/" + i.Name,
			},
		}
		for gvk, policyAttList := range i.AttachedHttpPolicies.Policies {
			if gvk.Group == wellknown.TrafficPolicyGVK.Group && gvk.Kind == wellknown.TrafficPolicyGVK.Kind {
				for _, policy := range policyAttList {
					// Convert TrafficPolicy IR to ADP Policy resources
					adpPolicies := convertTrafficPolicyIRToADPPolicies(policy, gwPolicyTarget)
					attachedGatewayPolicyResources = append(attachedGatewayPolicyResources, adpPolicies...)
				}
			}
		}
		for gvk, policyAttList := range i.AttachedListenerPolicies.Policies {
			if gvk.Group == wellknown.TrafficPolicyGVK.Group && gvk.Kind == wellknown.TrafficPolicyGVK.Kind {
				for _, policy := range policyAttList {
					adpPolicies := convertTrafficPolicyIRToADPPolicies(policy, gwPolicyTarget)
					attachedGatewayPolicyResources = append(attachedGatewayPolicyResources, adpPolicies...)
				}
			}
		}
		allResources = append(allResources, attachedGatewayPolicyResources...)

		attachedListenerPolicyResources := make([]*api.Resource, 0)
		for _, listener := range i.Listeners {
			listenerPolicyTarget := &api.PolicyTarget{
				Kind: &api.PolicyTarget_Listener{
					// TODO: check this
					Listener: string(listener.Name),
				},
			}
			for gvk, policyAttList := range listener.AttachedPolicies.Policies {
				if gvk.Group == wellknown.TrafficPolicyGVK.Group && gvk.Kind == wellknown.TrafficPolicyGVK.Kind {
					for _, policy := range policyAttList {
						adpPolicies := convertTrafficPolicyIRToADPPolicies(policy, listenerPolicyTarget)
						attachedListenerPolicyResources = append(attachedListenerPolicyResources, adpPolicies...)
					}
				}
			}
		}
		allResources = append(allResources, attachedListenerPolicyResources...)

		return &ADPResourcesForGateway{
			Resources: allResources,
			Gateway:   types.NamespacedName{Namespace: i.Namespace, Name: i.Name},
		}
	})

	return policiesByGateway
}

// convertTrafficPolicyIRToADPPolicies converts a policy attachment to ADP Policy resources
func convertTrafficPolicyIRToADPPolicies(policyAtt ir.PolicyAtt, policyTarget *api.PolicyTarget) []*api.Resource {
	var resources []*api.Resource

	policy, ok := policyAtt.PolicyIr.(agentgateway.AgentGatewayTrafficPolicyIr)
	if !ok {
		return resources
	}

	// Generate policy name - use the original policy reference
	var policyName string
	if policyAtt.PolicyRef != nil {
		policyName = fmt.Sprintf("trafficpolicy/%s/%s", policyAtt.PolicyRef.Namespace, policyAtt.PolicyRef.Name)
	} else {
		policyName = "trafficpolicy/unknown"
	}

	// Convert ExtAuth policy if present
	if policy.Extauth != nil && policy.Extauth.ExtAuthz != nil {
		extAuthPolicy := &api.Policy{
			Name:   policyName + ":extauth",
			Target: policyTarget,
			Spec: &api.PolicySpec{
				Kind: policy.Extauth.ExtAuthz,
			},
		}
		resources = append(resources, &api.Resource{
			Kind: &api.Resource_Policy{Policy: extAuthPolicy},
		})
	}

	return resources
}
