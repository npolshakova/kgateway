package agentgatewaysyncer

import (
	"context"
	"testing"

	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/a2a"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/mcp"
)

// dumpXDSCacheState is a helper function that dump the current state of the XDS cache for the agent gateway cache
func dumpXDSCacheState(ctx context.Context, cache envoycache.SnapshotCache) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("Current XDS cache state:")

	// Get all snapshot IDs from cache
	for _, nodeID := range cache.GetStatusKeys() {
		logger.Infof("Node ID: %s", nodeID)

		snapshot, err := cache.GetSnapshot(nodeID)
		if err != nil {
			logger.Infof("  Error getting snapshot: %v", err)
			continue
		}

		// Check for A2A targets
		logger.Infof("  A2A Targets (version %s):", snapshot.GetVersion(TargetTypeA2AUrl))
		resources := snapshot.GetResources(TargetTypeA2AUrl)
		for name := range resources {
			logger.Infof("    - %s", name)
		}

		// Check for MCP targets
		logger.Infof("  MCP Targets (version %s):", snapshot.GetVersion(TargetTypeMcpUrl))
		resources = snapshot.GetResources(TargetTypeMcpUrl)
		for name := range resources {
			logger.Infof("    - %s", name)
		}
	}
}

// TestXDSCacheState checks that the xds cache has targets and listeners properly set
func TestXDSCacheState(t *testing.T) {
	ctx := context.Background()
	cache := envoycache.NewSnapshotCache(false, envoycache.IDHash{}, nil)

	a2aTarget := &a2a.Target{
		Name:      "test-a2a-service",
		Host:      "10.0.0.1",
		Port:      8080,
		Path:      "/api",
		Listeners: []string{"a2a-listener"},
	}
	mcpTarget := &mcp.Target{
		Name: "test-mcp-service",
		Target: &mcp.Target_Sse{
			Sse: &mcp.Target_SseTarget{
				Host: "10.0.0.2",
				Port: 8081,
				Path: "/events",
			},
		},
		Listeners: []string{"mcp-listener"},
	}
	listener := &Listener{
		Name:     "test-listener",
		Protocol: Listener_A2A,
		Listener: &Listener_Sse{
			Sse: &SseListener{
				Address: "[::]",
				Port:    8080,
			},
		},
	}

	snapshot := &agentGwSnapshot{
		AgentGwA2AServices: envoycache.NewResources("v1", []envoytypes.Resource{
			a2aTarget,
		}),
		AgentGwMcpServices: envoycache.NewResources("v1", []envoytypes.Resource{
			mcpTarget,
		}),
		Listeners: envoycache.NewResources("v1", []envoytypes.Resource{
			listener,
		}),
	}

	// Set the snapshot in the cache
	err := cache.SetSnapshot(ctx, "test-node", snapshot)
	require.NoError(t, err)

	// Test dumping the cache state
	dumpXDSCacheState(ctx, cache)

	// Verify the resources were properly set
	retrievedSnapshot, err := cache.GetSnapshot("test-node")
	require.NoError(t, err)

	// Verify A2A resources
	a2aResources := retrievedSnapshot.GetResources(TargetTypeA2AUrl)
	assert.NotNil(t, a2aResources)
	assert.Contains(t, a2aResources, "test-a2a-service")
	retrievedA2A := a2aResources["test-a2a-service"].(*a2a.Target)
	assert.Equal(t, "10.0.0.1", retrievedA2A.Host)
	assert.Equal(t, uint32(8080), retrievedA2A.Port)
	assert.Equal(t, "/api", retrievedA2A.Path)

	// Verify MCP resources
	mcpResources := retrievedSnapshot.GetResources(TargetTypeMcpUrl)
	assert.NotNil(t, mcpResources)
	assert.Contains(t, mcpResources, "test-mcp-service")
	retrievedMCP := mcpResources["test-mcp-service"].(*mcp.Target)
	assert.Equal(t, "10.0.0.2", retrievedMCP.GetSse().Host)
	assert.Equal(t, uint32(8081), retrievedMCP.GetSse().Port)
	assert.Equal(t, "/events", retrievedMCP.GetSse().Path)

	// Verify Listener resources
	listenerResources := retrievedSnapshot.GetResources(TargetTypeListenerUrl)
	assert.NotNil(t, listenerResources)
	assert.Contains(t, listenerResources, "test-listener")
	retrievedListener := listenerResources["test-listener"].(*Listener)
	assert.Equal(t, Listener_A2A, retrievedListener.Protocol)
	assert.Equal(t, uint32(8080), retrievedListener.GetSse().Port)
}

// TestGetTargetName checks that the getTargetName function correctly formats target names
func TestGetTargetName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "test-service",
			expected: "test-service",
		},
		{
			name:     "name with slashes",
			input:    "namespace/service",
			expected: "namespace-service",
		},
		{
			name:     "name with invalid characters",
			input:    "test@service#123",
			expected: "test-service-123",
		},
		{
			name:     "name with multiple consecutive dashes",
			input:    "test--service",
			expected: "test-service",
		},
		{
			name:     "name with leading/trailing dashes",
			input:    "-test-service-",
			expected: "test-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTargetName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAgentGwSnapshot checks that the snapshot GetVersion and GetResources methods work as expected
func TestAgentGwSnapshot(t *testing.T) {
	a2aTarget := &a2a.Target{
		Name:      "test-a2a-service",
		Host:      "10.0.0.1",
		Port:      8080,
		Path:      "/api",
		Listeners: []string{"a2a-listener"},
	}
	mcpTarget := &mcp.Target{
		Name: "test-mcp-service",
		Target: &mcp.Target_Sse{
			Sse: &mcp.Target_SseTarget{
				Host: "10.0.0.2",
				Port: 8081,
				Path: "/events",
			},
		},
		Listeners: []string{"mcp-listener"},
	}
	listener := &Listener{
		Name:     "test-listener",
		Protocol: Listener_A2A,
		Listener: &Listener_Sse{
			Sse: &SseListener{
				Address: "[::]",
				Port:    8080,
			},
		},
	}

	// manually build the snapshot
	snapshot := &agentGwSnapshot{
		AgentGwA2AServices: envoycache.NewResources("v1", []envoytypes.Resource{
			a2aTarget,
		}),
		AgentGwMcpServices: envoycache.NewResources("v1", []envoytypes.Resource{
			mcpTarget,
		}),
		Listeners: envoycache.NewResources("v1", []envoytypes.Resource{
			listener,
		}),
	}

	// Construct the version map based on the snapshot
	err := snapshot.ConstructVersionMap()
	assert.NoError(t, err)

	assert.Equal(t, "v1", snapshot.GetVersion(TargetTypeA2AUrl))
	assert.Equal(t, "v1", snapshot.GetVersion(TargetTypeMcpUrl))
	assert.Equal(t, "v1", snapshot.GetVersion(TargetTypeListenerUrl))
	assert.Equal(t, "", snapshot.GetVersion("invalid-type"))

	a2aResources := snapshot.GetResources(TargetTypeA2AUrl)
	assert.NotNil(t, a2aResources)
	assert.Len(t, a2aResources, 1)
	a2aVersionMap := snapshot.GetVersionMap(TargetTypeA2AUrl)
	assert.NotNil(t, a2aVersionMap)

	mcpResources := snapshot.GetResources(TargetTypeMcpUrl)
	assert.NotNil(t, mcpResources)
	assert.Len(t, mcpResources, 1)
	mcpVersionMap := snapshot.GetVersionMap(TargetTypeMcpUrl)
	assert.NotNil(t, mcpVersionMap)

	listenerResources := snapshot.GetResources(TargetTypeListenerUrl)
	assert.NotNil(t, listenerResources)
	assert.Len(t, listenerResources, 1)
	listenerVersionMap := snapshot.GetVersionMap(TargetTypeMcpUrl)
	assert.NotNil(t, listenerVersionMap)

	err = snapshot.ConstructVersionMap()
	assert.NoError(t, err)
	assert.NotNil(t, snapshot.VersionMap)
}
