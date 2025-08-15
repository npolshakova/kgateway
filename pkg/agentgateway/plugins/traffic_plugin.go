package plugins

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	trafficPluginName = "traffic-policy-plugin"
	trafficPriority   = 10 // Highest priority
)

// TrafficPlugin implements the TrafficPolicyPlugin interface
type TrafficPlugin struct{}

// NewTrafficPlugin creates a new Traffic policy plugin
func NewTrafficPlugin() *TrafficPlugin {
	return &TrafficPlugin{}
}

// PolicyType returns the policy type for this plugin
func (p *TrafficPlugin) PolicyType() PolicyType {
	return PolicyTypeTraffic
}

// Name returns the name of this plugin
func (p *TrafficPlugin) Name() string {
	return trafficPluginName
}

// Priority returns the priority of this plugin
func (p *TrafficPlugin) Priority() int {
	return trafficPriority
}

// GeneratePolicies generates ADP policies for traffic policies
func (p *TrafficPlugin) GeneratePolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error) {
	logger := logging.New("agentgateway/plugins/traffic")

	if inputs.TrafficPolicies == nil {
		logger.Warn("traffic policies collection is nil, skipping traffic policy generation")
		return nil, nil
	}

	return p.GenerateTrafficPolicies(ctx, inputs.TrafficPolicies, inputs.GatewayExtensions)
}

// GenerateTrafficPolicies generates policies for traffic policies
func (p *TrafficPlugin) GenerateTrafficPolicies(ctx krt.HandlerContext, trafficPolicies krt.Collection[*v1alpha1.TrafficPolicy], gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension]) ([]ADPPolicy, error) {
	logger := logging.New("agentgateway/plugins/traffic")
	logger.Debug("generating traffic policies")

	var trafficPoliciesResult []ADPPolicy

	// Fetch all traffic policies and process them
	allTrafficPolicies := krt.Fetch(ctx, trafficPolicies)

	for _, trafficPolicy := range allTrafficPolicies {
		policies := p.generatePoliciesForTrafficPolicy(ctx, gatewayExtensions, trafficPolicy)
		trafficPoliciesResult = append(trafficPoliciesResult, policies...)
	}

	logger.Info("generated traffic policies", "count", len(trafficPoliciesResult))
	return trafficPoliciesResult, nil
}

// generatePoliciesForTrafficPolicy generates policies for a single traffic policy
func (p *TrafficPlugin) generatePoliciesForTrafficPolicy(ctx krt.HandlerContext, gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension], trafficPolicy *v1alpha1.TrafficPolicy) []ADPPolicy {
	logger := logging.New("agentgateway/plugins/traffic")
	var adpPolicies []ADPPolicy

	for _, target := range trafficPolicy.Spec.TargetRefs {
		var policyTarget *api.PolicyTarget

		switch string(target.Kind) {
		case wellknown.GatewayKind:
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Gateway{
					Gateway: trafficPolicy.Namespace + "/" + string(target.Name),
				},
			}
			// TODO(npolshak): add listener support once https://github.com/agentgateway/agentgateway/pull/323 goes in
			//if target.SectionName != nil {
			//	policyTarget = &api.PolicyTarget{
			//		Kind: &api.PolicyTarget_Listener{
			//			Listener: InternalGatewayName(trafficPolicy.Namespace, string(target.Name), string(*target.SectionName)),
			//		},
			//	}
			//}

		case wellknown.HTTPRouteKind:
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Route{
					Route: trafficPolicy.Namespace + "/" + string(target.Name),
				},
			}
			// TODO(npolshak): add route rule support once https://github.com/agentgateway/agentgateway/pull/323 goes in
			//if target.SectionName != nil {
			//	policyTarget = &api.PolicyTarget{
			//		Kind: &api.PolicyTarget_RouteRule{
			//			RouteRule: trafficPolicy.Namespace + "/" + string(target.Name) + "/" + string(*target.SectionName),
			//		},
			//	}
			//}

		default:
			logger.Warn("unsupported target kind", "kind", target.Kind, "policy", trafficPolicy.Name)
			continue
		}

		if policyTarget != nil {
			translatedPolicies := p.translateTrafficPolicyToADP(ctx, gatewayExtensions, trafficPolicy, string(target.Name), policyTarget)
			adpPolicies = append(adpPolicies, translatedPolicies...)
		}
	}

	return adpPolicies
}

// translateTrafficPolicyToADP converts a TrafficPolicy to ADP Policy resources
func (p *TrafficPlugin) translateTrafficPolicyToADP(ctx krt.HandlerContext, gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension], trafficPolicy *v1alpha1.TrafficPolicy, policyTargetName string, policyTarget *api.PolicyTarget) []ADPPolicy {
	adpPolicies := make([]ADPPolicy, 0)

	// Generate a base policy name from the TrafficPolicy reference
	policyName := fmt.Sprintf("trafficpolicy/%s/%s/%s", trafficPolicy.Namespace, trafficPolicy.Name, policyTargetName)

	// Convert ExtAuth policy if present
	if trafficPolicy.Spec.ExtAuth != nil && trafficPolicy.Spec.ExtAuth.ExtensionRef != nil {
		extAuthPolicies := p.processExtAuthPolicy(ctx, gatewayExtensions, trafficPolicy, policyName, policyTarget)
		adpPolicies = append(adpPolicies, extAuthPolicies...)
	}

	// TODO: Add support for other policy types as needed:
	// - RateLimit
	// - Transformation
	// - ExtProc
	// - AI policies
	// etc.

	return adpPolicies
}

// processExtAuthPolicy processes ExtAuth configuration and creates corresponding ADP policies
func (p *TrafficPlugin) processExtAuthPolicy(ctx krt.HandlerContext, gatewayExtensions krt.Collection[*v1alpha1.GatewayExtension], trafficPolicy *v1alpha1.TrafficPolicy, policyName string, policyTarget *api.PolicyTarget) []ADPPolicy {
	logger := logging.New("agentgateway/plugins/traffic")

	// Look up the GatewayExtension referenced by the ExtAuth policy
	extensionName := trafficPolicy.Spec.ExtAuth.ExtensionRef.Name
	gwExtKey := fmt.Sprintf("%s/%s", trafficPolicy.Namespace, extensionName)
	gwExt := krt.FetchOne(ctx, gatewayExtensions, krt.FilterKey(gwExtKey))

	if gwExt == nil || (*gwExt).Spec.Type != v1alpha1.GatewayExtensionTypeExtAuth || (*gwExt).Spec.ExtAuth == nil {
		logger.Error("gateway extension not found or not of type ExtAuth", "extension", gwExtKey)
		return nil
	}

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

	if extauthSvcTarget == nil {
		logger.Warn("failed to translate traffic policy", "policy", trafficPolicy.Name, "target", policyTarget, "error", "missing service target")
		return nil
	}

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

	logger.Debug("generated ExtAuth policy",
		"policy", trafficPolicy.Name,
		"adp_policy", extauthPolicy.Name,
		"target", extauthSvcTarget)

	return []ADPPolicy{{Policy: extauthPolicy}}
}

// Verify that TrafficPlugin implements the required interfaces
var _ TrafficPolicyPlugin = (*TrafficPlugin)(nil)
var _ ADPPolicyPlugin = (*TrafficPlugin)(nil)
var _ PolicyPlugin = (*TrafficPlugin)(nil)
