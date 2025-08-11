package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	corev1 "k8s.io/api/core/v1"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugins/trafficpolicy/agentgateway"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	a2aProtocol = "kgateway.dev/a2a"
)

func ADPPolicyCollection(inputs Inputs, binds krt.Collection[ADPResourcesForGateway], krtopts krtutil.KrtOptions) krt.Collection[ADPResourcesForGateway] {
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
	policiesByGateway := krt.NewCollection(binds, func(ctx krt.HandlerContext, i ADPResourcesForGateway) *ADPResourcesForGateway {
		inferences := slices.Map(krt.Fetch(ctx, inference), func(e ADPPolicy) *api.Resource {
			return toADPResource(e)
		})
		a2aPolicies := slices.Map(krt.Fetch(ctx, a2a), func(e ADPPolicy) *api.Resource {
			return toADPResource(e)
		})

		// Look up backend targeting policies and convert them to ADP Policy resources
		var backendPolicyResources []*api.Resource

		// Get all backends with policies from the backend index
		for _, backendCol := range inputs.Backends.BackendsWithPolicy() {
			backends := krt.Fetch(ctx, backendCol)
			for _, backend := range backends {
				if backend == nil {
					continue
				}

				// Create ObjectSource for the backend
				targetRef := pluginsdkir.ObjectSource{
					Group:     backend.Group,
					Kind:      backend.Kind,
					Namespace: backend.Namespace,
					Name:      backend.Name,
				}

				inputs.Backends.BackendsWithPolicy()

				// Look up policies targeting this backend
				pols := inputs.Backends.PolicyIndex().LookupTargetingPolicies(ctx,
					pluginsdk.BackendAttachmentPoint,
					targetRef,
					"", // no section name for backends
					backend.GetObjectLabels())

				// Convert PolicyAtt objects to ADP Policy resources
				for _, policyAtt := range pols {
					if policyAtt.PolicyIr != nil {
						// Try to convert the policy IR to an ADP policy
						if adpPolicy := convertPolicyAttToADPPolicy(policyAtt, backend); adpPolicy != nil {
							backendPolicyResources = append(backendPolicyResources, toADPResource(ADPPolicy{adpPolicy}))
						}
					}
				}
			}
		}

		// Combine all policy resources
		allResources := append(inferences, a2aPolicies...)
		allResources = append(allResources, backendPolicyResources...)

		return &ADPResourcesForGateway{
			Resources: allResources,
			Gateway:   i.Gateway,
		}
	})

	return policiesByGateway
}

// convertPolicyAttToADPPolicy converts a PolicyAtt to an ADP Policy
func convertPolicyAttToADPPolicy(policyAtt pluginsdkir.PolicyAtt, backend *pluginsdkir.BackendObjectIR) *api.Policy {
	if policyAtt.PolicyIr == nil {
		return nil
	}

	// Only process agent gateway traffic policies
	agwPolicyIr, ok := policyAtt.PolicyIr.(*agentgateway.TrafficPolicy)
	if !ok {
		// This PolicyIR is not an agent gateway traffic policy
		// Return nil to skip processing
		return nil
	}

	// PolicyRef can be nil if the attachment was done via extension ref or if PolicyAtt is the result of MergePolicies
	if policyAtt.PolicyRef == nil {
		// Cannot create a meaningful policy name without PolicyRef, skip processing
		return nil
	}

	// Create a unique policy name using the policy reference and backend
	policyName := fmt.Sprintf("%s/%s->%s/%s",
		policyAtt.PolicyRef.Namespace, policyAtt.PolicyRef.Name,
		backend.Namespace, backend.Name)

	// Create the backend target reference
	backendTarget := fmt.Sprintf("%s/%s", backend.Namespace, backend.Name)

	var translatedSpec *api.PolicySpec
	if agwPolicyIr.Spec.ExtAuth != nil {
		translatedSpec = &api.PolicySpec{
			Kind: agwPolicyIr.Spec.ExtAuth.Extauth,
		}
	}

	// For now, we create a generic ADP policy that references the backend
	// In the future, this could be extended to convert specific policy types
	// to their corresponding ADP policy specifications
	adpPolicy := &api.Policy{
		Name:   policyName,
		Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: backendTarget}},
		// Note: The actual policy spec conversion would depend on the policy type
		// This is a placeholder that could be extended based on specific policy types
		Spec: translatedSpec,
	}

	return adpPolicy
}
