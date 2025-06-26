package agentgatewaysyncer

import (
	"context"
	"testing"

	"github.com/agentgateway/agentgateway/go/api"
	envoytypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dumpXDSCacheState is a helper function that dump the current state of the XDS cache for the agentgateway cache
func dumpXDSCacheState(ctx context.Context, cache envoycache.SnapshotCache) {
	logger.Info("current XDS cache state:")

	// Get all snapshot IDs from cache
	for _, nodeID := range cache.GetStatusKeys() {
		logger.Info("snapshot has node", "node_id", nodeID)

		snapshot, err := cache.GetSnapshot(nodeID)
		if err != nil {
			logger.Info("error getting snapshot", "error", err.Error())
			continue
		}

		// Check for Resource targets
		logger.Info("Resource targets version", "snapshot", snapshot.GetVersion(TargetTypeResourceUrl)) //nolint:sloglint // ignore msg-type
		resources := snapshot.GetResources(TargetTypeResourceUrl)
		for name := range resources {
			logger.Info("snapshot has resources", "name", name)
		}
	}
}

// TestXDSCacheState checks that the xds cache has resources properly set
func TestXDSCacheState(t *testing.T) {
	ctx := context.Background()
	cache := envoycache.NewSnapshotCache(false, envoycache.IDHash{}, nil)

	// Create test resources using the new API
	bindResource := &api.Resource{
		Kind: &api.Resource_Bind{
			Bind: &api.Bind{
				Key:  "test-bind",
				Port: 8080,
			},
		},
	}

	listenerResource := &api.Resource{
		Kind: &api.Resource_Listener{
			Listener: &api.Listener{
				Key:         "test-listener",
				Name:        "default",
				BindKey:     "test-bind",
				GatewayName: "test-gateway",
				Protocol:    api.Protocol_HTTP,
			},
		},
	}

	routeResource := &api.Resource{
		Kind: &api.Resource_Route{
			Route: &api.Route{
				Key:         "test-route",
				ListenerKey: "test-listener",
				RuleName:    "test-rule",
				RouteName:   "test-route",
				Matches: []*api.RouteMatch{
					{
						Path: &api.PathMatch{
							Kind: &api.PathMatch_PathPrefix{
								PathPrefix: "/test",
							},
						},
					},
				},
				Backends: []*api.RouteBackend{
					{
						Kind: &api.RouteBackend_Service{
							Service: "test-service",
						},
						Weight: 1,
						Port:   8080,
					},
				},
			},
		},
	}

	snapshot := &agentGwSnapshot{
		Config: envoycache.NewResources("v1", []envoytypes.Resource{
			&envoyResourceWithCustomName{
				Message: bindResource,
				Name:    "test-bind",
				version: 1,
			},
			&envoyResourceWithCustomName{
				Message: listenerResource,
				Name:    "test-listener",
				version: 2,
			},
			&envoyResourceWithCustomName{
				Message: routeResource,
				Name:    "test-route",
				version: 3,
			},
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

	// Verify Resource resources
	resources := retrievedSnapshot.GetResources(TargetTypeResourceUrl)
	assert.NotNil(t, resources)
	assert.Len(t, resources, 3)

	// Verify bind resource
	bindWrapper := resources["test-bind"].(*envoyResourceWithCustomName)
	bindRes := bindWrapper.Message.(*api.Resource)
	assert.NotNil(t, bindRes.GetBind())
	assert.Equal(t, "test-bind", bindRes.GetBind().Key)
	assert.Equal(t, uint32(8080), bindRes.GetBind().Port)

	// Verify listener resource
	listenerWrapper := resources["test-listener"].(*envoyResourceWithCustomName)
	listenerRes := listenerWrapper.Message.(*api.Resource)
	assert.NotNil(t, listenerRes.GetListener())
	assert.Equal(t, "test-listener", listenerRes.GetListener().Key)
	assert.Equal(t, "default", listenerRes.GetListener().Name)
	assert.Equal(t, api.Protocol_HTTP, listenerRes.GetListener().Protocol)

	// Verify route resource
	routeWrapper := resources["test-route"].(*envoyResourceWithCustomName)
	routeRes := routeWrapper.Message.(*api.Resource)
	assert.NotNil(t, routeRes.GetRoute())
	assert.Equal(t, "test-route", routeRes.GetRoute().Key)
	assert.Equal(t, "test-rule", routeRes.GetRoute().RuleName)
	assert.Len(t, routeRes.GetRoute().Matches, 1)
	assert.Equal(t, "/test", routeRes.GetRoute().Matches[0].GetPath().GetPathPrefix())
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
	bindResource := &api.Resource{
		Kind: &api.Resource_Bind{
			Bind: &api.Bind{
				Key:  "test-bind",
				Port: 8080,
			},
		},
	}

	listenerResource := &api.Resource{
		Kind: &api.Resource_Listener{
			Listener: &api.Listener{
				Key:         "test-listener",
				Name:        "default",
				BindKey:     "test-bind",
				GatewayName: "test-gateway",
				Protocol:    api.Protocol_HTTP,
			},
		},
	}

	routeResource := &api.Resource{
		Kind: &api.Resource_Route{
			Route: &api.Route{
				Key:         "test-route",
				ListenerKey: "test-listener",
				RuleName:    "test-rule",
				RouteName:   "test-route",
				Matches: []*api.RouteMatch{
					{
						Path: &api.PathMatch{
							Kind: &api.PathMatch_PathPrefix{
								PathPrefix: "/test",
							},
						},
					},
				},
				Backends: []*api.RouteBackend{
					{
						Kind: &api.RouteBackend_Service{
							Service: "test-service",
						},
						Weight: 1,
						Port:   8080,
					},
				},
			},
		},
	}

	snapshot := &agentGwSnapshot{
		Config: envoycache.NewResources("v1", []envoytypes.Resource{
			&envoyResourceWithCustomName{
				Message: bindResource,
				Name:    "test-bind",
				version: 1,
			},
			&envoyResourceWithCustomName{
				Message: listenerResource,
				Name:    "test-listener",
				version: 2,
			},
			&envoyResourceWithCustomName{
				Message: routeResource,
				Name:    "test-route",
				version: 3,
			},
		}),
	}

	// Test GetVersion
	assert.Equal(t, "v1", snapshot.GetVersion(TargetTypeResourceUrl))

	// Test GetResources
	resources := snapshot.GetResources(TargetTypeResourceUrl)
	assert.NotNil(t, resources)
	assert.Len(t, resources, 3)

	// Test GetVersionMap
	err := snapshot.ConstructVersionMap()
	require.NoError(t, err)

	a2aVersionMap := snapshot.GetVersionMap(TargetTypeResourceUrl)
	assert.NotNil(t, a2aVersionMap)
	assert.Len(t, a2aVersionMap, 3)

	// Verify specific resources exist
	assert.Contains(t, a2aVersionMap, "test-bind")
	assert.Contains(t, a2aVersionMap, "test-listener")
	assert.Contains(t, a2aVersionMap, "test-route")
}
