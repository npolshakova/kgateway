package gateway

import (
	"context"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
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

		logger.Info("Resource targets version", "snapshot", snapshot.GetVersion(TargetTypeResourceUrl))
		resources := snapshot.GetResources(TargetTypeResourceUrl)
		for name := range resources {
			logger.Info("snapshot has resources", "name", name)
		}

		logger.Info("Address targets version", "snapshot", snapshot.GetVersion(TargetTypeAddressUrl))
		resources = snapshot.GetResources(TargetTypeAddressUrl)
		for name := range resources {
			logger.Info("snapshot has resources", "name", name)
		}
	}
}
