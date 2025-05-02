package agentgateway

import (
	"path/filepath"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

const (
	a2aPort = 9090
	mcpPort = 8080
)

var (
	// Agent Gateway deployment
	agentgatewayManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "agentgateway.yaml")

	// Test A2A Agent
	a2aAgentManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "a2a.yaml")

	// Test MCP Server
	mcpManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "mcp.yaml")

	// Self-managed Gateway to configure the Agent Gateway
	gatewayManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "gateway-mcp-a2a.yaml")
)
