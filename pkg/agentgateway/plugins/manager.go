package plugins

import (
	"fmt"
	"sort"

	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

// DefaultPolicyManager implements the PolicyManager interface
type DefaultPolicyManager struct {
	plugins             []PolicyPlugin
	adpPlugins          []ADPPolicyPlugin
	agentGatewayPlugins []AgentGatewayPolicyPlugin
	pluginsByType       map[PolicyType][]PolicyPlugin
}

// NewPolicyManager creates a new DefaultPolicyManager
func NewPolicyManager() *DefaultPolicyManager {
	return &DefaultPolicyManager{
		plugins:             make([]PolicyPlugin, 0),
		adpPlugins:          make([]ADPPolicyPlugin, 0),
		agentGatewayPlugins: make([]AgentGatewayPolicyPlugin, 0),
		pluginsByType:       make(map[PolicyType][]PolicyPlugin),
	}
}

// RegisterPlugin registers a policy plugin
func (m *DefaultPolicyManager) RegisterPlugin(plugin PolicyPlugin) error {
	if plugin == nil {
		return fmt.Errorf("cannot register nil plugin")
	}

	managerLogger := logging.New("agentgateway/plugins/manager")
	managerLogger.Info("registering policy plugin", "name", plugin.Name(), "type", plugin.PolicyType())

	// Add to general plugins list
	m.plugins = append(m.plugins, plugin)

	// Add to type-specific map
	policyType := plugin.PolicyType()
	if m.pluginsByType[policyType] == nil {
		m.pluginsByType[policyType] = make([]PolicyPlugin, 0)
	}
	m.pluginsByType[policyType] = append(m.pluginsByType[policyType], plugin)

	// Sort by priority (lower numbers have higher priority)
	sort.Slice(m.pluginsByType[policyType], func(i, j int) bool {
		return m.pluginsByType[policyType][i].Priority() < m.pluginsByType[policyType][j].Priority()
	})

	// Add to specialized plugin lists if applicable
	if adpPlugin, ok := plugin.(ADPPolicyPlugin); ok {
		m.adpPlugins = append(m.adpPlugins, adpPlugin)
		// Sort ADP plugins by priority
		sort.Slice(m.adpPlugins, func(i, j int) bool {
			return m.adpPlugins[i].Priority() < m.adpPlugins[j].Priority()
		})
	}

	if agwPlugin, ok := plugin.(AgentGatewayPolicyPlugin); ok {
		m.agentGatewayPlugins = append(m.agentGatewayPlugins, agwPlugin)
		// Sort agent gateway plugins by priority
		sort.Slice(m.agentGatewayPlugins, func(i, j int) bool {
			return m.agentGatewayPlugins[i].Priority() < m.agentGatewayPlugins[j].Priority()
		})
	}

	return nil
}

// GetPluginsByType returns all plugins of a specific type
func (m *DefaultPolicyManager) GetPluginsByType(policyType PolicyType) []PolicyPlugin {
	plugins, exists := m.pluginsByType[policyType]
	if !exists {
		return make([]PolicyPlugin, 0)
	}
	// Return a copy to prevent external modification
	result := make([]PolicyPlugin, len(plugins))
	copy(result, plugins)
	return result
}

// GetADPPlugins returns all ADP policy plugins
func (m *DefaultPolicyManager) GetADPPlugins() []ADPPolicyPlugin {
	// Return a copy to prevent external modification
	result := make([]ADPPolicyPlugin, len(m.adpPlugins))
	copy(result, m.adpPlugins)
	return result
}

// GetAgentGatewayPlugins returns all agent gateway policy plugins
func (m *DefaultPolicyManager) GetAgentGatewayPlugins() []AgentGatewayPolicyPlugin {
	// Return a copy to prevent external modification
	result := make([]AgentGatewayPolicyPlugin, len(m.agentGatewayPlugins))
	copy(result, m.agentGatewayPlugins)
	return result
}

// GenerateAllPolicies generates policies from all registered ADP plugins
func (m *DefaultPolicyManager) GenerateAllPolicies(ctx krt.HandlerContext, inputs *PolicyInputs) ([]ADPPolicy, error) {
	var allPolicies []ADPPolicy
	var allErrors []error

	for _, plugin := range m.adpPlugins {
		managerLogger := logging.New("agentgateway/plugins/manager")
		managerLogger.Debug("generating policies", "plugin", plugin.Name(), "type", plugin.PolicyType())

		policies, err := plugin.GeneratePolicies(ctx, inputs)
		if err != nil {
			managerLogger.Error("failed to generate policies", "plugin", plugin.Name(), "error", err)
			allErrors = append(allErrors, fmt.Errorf("plugin %s failed: %w", plugin.Name(), err))
			continue
		}

		allPolicies = append(allPolicies, policies...)
		managerLogger.Debug("generated policies", "plugin", plugin.Name(), "count", len(policies))
	}

	managerLogger := logging.New("agentgateway/plugins/manager")
	if len(allErrors) > 0 {
		// Log errors but don't fail completely - return partial results
		for _, err := range allErrors {
			managerLogger.Error("policy generation error", "error", err)
		}
	}

	managerLogger.Info("generated all policies", "total_policies", len(allPolicies), "errors", len(allErrors))
	return allPolicies, nil
}

// GetRegisteredPlugins returns all registered plugins (for debugging/introspection)
func (m *DefaultPolicyManager) GetRegisteredPlugins() []PolicyPlugin {
	// Return a copy to prevent external modification
	result := make([]PolicyPlugin, len(m.plugins))
	copy(result, m.plugins)
	return result
}
