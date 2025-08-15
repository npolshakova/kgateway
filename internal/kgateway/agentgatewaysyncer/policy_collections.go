package agentgatewaysyncer

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/slices"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/plugins"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	a2aProtocol = "kgateway.dev/a2a"
)

func ADPPolicyCollection(inputs Inputs, binds krt.Collection[ADPResourcesForGateway], krtopts krtutil.KrtOptions, policyManager *plugins.DefaultPolicyManager) krt.Collection[ADPResourcesForGateway] {
	domainSuffix := kubeutils.GetClusterDomainName()

	// Create policy inputs for all plugins
	policyInputs := &plugins.PolicyInputs{
		Services:          inputs.Services,
		InferencePools:    inputs.InferencePools,
		TrafficPolicies:   inputs.TrafficPolicies,
		GatewayExtensions: inputs.GatewayExtensions,
		DomainSuffix:      domainSuffix,
	}

	// Generate all policies using the plugin system
	allPoliciesCol := krt.NewCollection(binds, func(ctx krt.HandlerContext, i ADPResourcesForGateway) *ADPResourcesForGateway {
		logger.Debug("generating policies for gateway", "gateway", i.Gateway)

		// Generate all policies from all registered plugins
		allPolicies, err := policyManager.GenerateAllPolicies(ctx, policyInputs)
		if err != nil {
			logger.Error("failed to generate policies", "error", err)
			// Return empty resources but don't fail completely
			return &ADPResourcesForGateway{
				Resources: []*api.Resource{},
				Gateway:   i.Gateway,
			}
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
