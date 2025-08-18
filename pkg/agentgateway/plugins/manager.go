package plugins

import (
	"fmt"

	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

// DefaultPolicyManager implements the PolicyManager interface
type DefaultPolicyManager struct {
	plugins       []PolicyPlugin
	pluginsByType map[PolicyType][]PolicyPlugin
}

// NewPolicyManager creates a new DefaultPolicyManager
func NewPolicyManager() *DefaultPolicyManager {
	return &DefaultPolicyManager{
		plugins:       make([]PolicyPlugin, 0),
		pluginsByType: make(map[PolicyType][]PolicyPlugin),
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

// GenerateAllPolicies generates policies from all registered ADP plugins
func (m *DefaultPolicyManager) GenerateAllPolicies(ctx krt.HandlerContext, agw *AgwCollections, policyInput PolicyInput) ([]ADPPolicy, error) {
	var allPolicies []ADPPolicy
	var allErrors []error

	for _, plugin := range m.plugins {
		managerLogger := logging.New("agentgateway/plugins/manager")
		managerLogger.Debug("generating policies", "plugin", plugin.Name(), "type", plugin.PolicyType())

		policies, err := plugin.GeneratePolicies(ctx, agw, policyInput)
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
