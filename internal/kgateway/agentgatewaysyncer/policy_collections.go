package agentgatewaysyncer

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"

	"github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/plugins"
)

func ADPPolicyCollection(binds krt.Collection[ADPResourcesForGateway], agwPlugins plugins.AgentgatewayPlugin) krt.Collection[ADPResourcesForGateway] {
	// Generate all policies using the plugin system
	allPoliciesCol := krt.NewCollection(binds, func(ctx krt.HandlerContext, i ADPResourcesForGateway) *ADPResourcesForGateway {
		logger.Debug("generating policies for gateway", "gateway", i.Gateway)

		var allPolicies []plugins.ADPPolicy
		// Generate all policies from all registered plugins using contributed policies
		for _, plugin := range agwPlugins.ContributesPolicies {
			policies := plugin.ApplyPolicies()
			allPolicies = append(allPolicies, policies...)
		}

		// Convert all plugins.ADPPolicy structs to api.Resource structs
		allResources := slices.Map(allPolicies, func(policy plugins.ADPPolicy) *api.Resource {
			return toADPResource(ADPPolicy{policy.Policy})
		})

		logger.Info("generated policies for gateway", "gateway", i.Gateway, "policy_count", len(allResources))

		return &ADPResourcesForGateway{
			Resources: allResources,
			Gateway:   i.Gateway,
		}
	})

	return allPoliciesCol
}
