package agentgatewaysyncer

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
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

	trafficPolicyCol := krt.NewManyCollection(inputs.TrafficPolicies, func(ctx krt.HandlerContext, trafficPolicy *v1alpha1.TrafficPolicy) []ADPPolicy {
		var trafficPolicies []ADPPolicy
		for _, target := range trafficPolicy.Spec.TargetRefs {
			switch string(target.Kind) {
			case wellknown.GatewayKind:
				gwTarget := &api.PolicyTarget{
					Kind: &api.PolicyTarget_Gateway{
						Gateway: trafficPolicy.Namespace + "/" + string(target.Name),
					},
				}
				if target.SectionName != nil {
					gwTarget = &api.PolicyTarget{
						Kind: &api.PolicyTarget_Listener{
							Listener: InternalGatewayName(trafficPolicy.Namespace, string(target.Name), string(*target.SectionName)),
						},
					}
				}
				translatedPolicies := translateTrafficPolicyToADP(ctx, inputs.GatewayExtensions, trafficPolicy, string(target.Name), gwTarget)
				trafficPolicies = append(trafficPolicies, translatedPolicies...)
			case wellknown.HTTPRouteKind:
				gwTarget := &api.PolicyTarget{
					Kind: &api.PolicyTarget_Route{
						Route: trafficPolicy.Namespace + "/" + string(target.Name),
					},
				}
				if target.SectionName != nil {
					gwTarget = &api.PolicyTarget{
						Kind: &api.PolicyTarget_RouteRule{
							RouteRule: trafficPolicy.Namespace + "/" + string(target.Name) + "/" + string(*target.SectionName),
						},
					}
				}
				translatedPolicies := translateTrafficPolicyToADP(ctx, inputs.GatewayExtensions, trafficPolicy, string(target.Name), gwTarget)
				trafficPolicies = append(trafficPolicies, translatedPolicies...)
			}
		}
		return trafficPolicies
	}, krtopts.ToOptions("A2APolicies")...)

	// For now, we apply all policies to all gateways. In the future, we can more precisely bind them to only relevant ones
	policiesByGateway := krt.NewCollection(binds, func(ctx krt.HandlerContext, i ADPResourcesForGateway) *ADPResourcesForGateway {
		var allResources []*api.Resource

		// Add inference policies
		inferences := slices.Map(krt.Fetch(ctx, inference), func(e ADPPolicy) *api.Resource {
			return toADPResource(e)
		})
		allResources = append(allResources, inferences...)

		// Add A2A policies
		a2aPolicies := slices.Map(krt.Fetch(ctx, a2a), func(e ADPPolicy) *api.Resource {
			return toADPResource(e)
		})
		allResources = append(allResources, a2aPolicies...)

		// Add TrafficPolicy policies
		trafficPolicies := slices.Map(krt.Fetch(ctx, trafficPolicyCol), func(e ADPPolicy) *api.Resource {
			return toADPResource(e)
		})
		allResources = append(allResources, trafficPolicies...)

		return &ADPResourcesForGateway{
			Resources: allResources,
			Gateway:   i.Gateway,
		}
	})

	return policiesByGateway
}

// translateTrafficPolicyToADP converts a TrafficPolicy to ADP Policy resources
func translateTrafficPolicyToADP(ctx krt.HandlerContext, gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension], trafficPolicy *v1alpha1.TrafficPolicy, policyTargetName string, policyTarget *api.PolicyTarget) []ADPPolicy {
	adpPolicies := make([]ADPPolicy, 0)

	// Generate a base policy name from the TrafficPolicy reference
	policyName := fmt.Sprintf("trafficpolicy/%s/%s/%s", trafficPolicy.Namespace, trafficPolicy.Name, policyTargetName)

	// Convert ExtAuth policy if present
	if trafficPolicy.Spec.ExtAuth != nil && trafficPolicy.Spec.ExtAuth.ExtensionRef != nil {
		// Look up the GatewayExtension referenced by the ExtAuth policy
		extensionName := trafficPolicy.Spec.ExtAuth.ExtensionRef.Name
		gwExtKey := fmt.Sprintf("%s/%s", trafficPolicy.Namespace, extensionName)
		gwExt := krt.FetchOne(ctx, gatewayExtensions, krt.FilterKey(gwExtKey))

		if gwExt != nil && (*gwExt).Spec.Type == v1alpha1.GatewayExtensionTypeExtAuth && (*gwExt).Spec.ExtAuth != nil {
			// Extract service target from GatewayExtension's ExtAuth configuration
			var extauthSvcTarget *api.BackendReference
			if (*gwExt).Spec.ExtAuth.GrpcService != nil && (*gwExt).Spec.ExtAuth.GrpcService.BackendRef != nil {
				backendRef := (*gwExt).Spec.ExtAuth.GrpcService.BackendRef
				serviceName := string(backendRef.Name)
				port := uint32(80) // default port
				if backendRef.Port != nil {
					port = uint32(*backendRef.Port)
				}
				// use trafficPolicy namespace as default
				namespace := trafficPolicy.Namespace
				if backendRef.Namespace != nil {
					namespace = string(*backendRef.Namespace)
				}
				serviceHost := kubeutils.ServiceFQDN(metav1.ObjectMeta{Namespace: namespace, Name: serviceName})
				extauthSvcTarget = &api.BackendReference{
					Kind: &api.BackendReference_Service{Service: namespace + "/" + serviceHost},
					Port: port,
				}
			}

			if extauthSvcTarget != nil {
				extauthPolicy := &api.Policy{
					Name:   policyName + ":extauth",
					Target: policyTarget,
					Spec: &api.PolicySpec{
						Kind: &api.PolicySpec_ExtAuthz{
							ExtAuthz: &api.PolicySpec_ExternalAuth{
								Target:  extauthSvcTarget,
								Context: trafficPolicy.Spec.ExtAuth.ContextExtensions,
							},
						},
					},
				}
				adpPolicies = append(adpPolicies, ADPPolicy{Policy: extauthPolicy})
			} else {
				// log warning, but continue processing other policies
				logger.Warn("failed to translate traffic policy", "policy", trafficPolicy.Name, "target", policyTargetName, "error", "missing service target")
			}
		} else {
			// log warning, but continue processing other policies
			logger.Error("gateway extension not of type ExtAuth", "extension", gwExtKey)
		}
	}

	return adpPolicies
}
