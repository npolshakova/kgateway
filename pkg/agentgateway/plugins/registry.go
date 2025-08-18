package plugins

import (
	sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type AgentgatewayPlugin struct {
	ContributesPolicies map[schema.GroupKind]PolicyPlugin
	// extra has sync beyond primary resources in the collections above
	ExtraHasSynced func() bool
}

func MergePlugins(plug ...AgentgatewayPlugin) AgentgatewayPlugin {
	ret := AgentgatewayPlugin{
		ContributesPolicies: make(map[schema.GroupKind]PolicyPlugin),
	}
	var funcs []sdk.GwTranslatorFactory
	var hasSynced []func() bool
	for _, p := range plug {
		if p.ExtraHasSynced != nil {
			hasSynced = append(hasSynced, p.ExtraHasSynced)
		}
	}
	ret.ExtraHasSynced = mergeSynced(hasSynced)
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
