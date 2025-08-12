package agentgateway

import (
	"log/slog"

	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func mergeExtAuth(
	p1, p2 *AgentGatewayTrafficPolicyIr,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	p2MergeOrigins pluginsdkir.MergeOrigins,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.Extauth, p2.Extauth, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.Extauth != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.Extauth = p2.Extauth
		mergeOrigins.SetOne("extAuth", p2Ref, p2MergeOrigins)

	default:
		slog.Warn("unsupported merge strategy for extAuth policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
}
