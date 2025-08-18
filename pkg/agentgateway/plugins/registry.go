package plugins

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type AgentgatewayPlugin struct {
	ContributesPolicies map[schema.GroupKind]PolicyPlugin
}

func MergePlugins(plug ...AgentgatewayPlugin) AgentgatewayPlugin {
	ret := AgentgatewayPlugin{
		ContributesPolicies: make(map[schema.GroupKind]PolicyPlugin),
	}
	for _, p := range plug {
		// Merge contributed policies
		for gk, policy := range p.ContributesPolicies {
			ret.ContributesPolicies[gk] = policy
		}
	}
	return ret
}

func mergeSynced(funcs []func() bool) func() bool {
	return func() bool {
		for _, f := range funcs {
			if !f() {
				return false
			}
		}
		return true
	}
}

// Plugins registers all built-in policy plugins
func Plugins(agw *AgwCollections) []AgentgatewayPlugin {
	return []AgentgatewayPlugin{
		NewTrafficPlugin(agw),
		NewInferencePlugin(agw),
		NewA2APlugin(agw),
	}
}
