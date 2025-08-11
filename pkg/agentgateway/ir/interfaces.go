package ir

import (
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
)

// Re-export the types from pluginsdk for convenience in this package
type (
	AgentGatewayTranslationPass              = ir.AgentGatewayTranslationPass
	AgentGatewayRouteContext                 = ir.AgentGatewayRouteContext
	AgentGatewayTranslationBackendContext    = ir.AgentGatewayTranslationBackendContext
	UnimplementedAgentGatewayTranslationPass = ir.UnimplementedAgentGatewayTranslationPass
)
