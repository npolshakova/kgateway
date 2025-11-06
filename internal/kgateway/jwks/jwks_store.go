package jwks

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type JwksStore struct {
	mgr             manager.Manager
	jwksCache       *jwksCache
	jwksFetcher     *JwksFetcher
	configMapSyncer *ConfigMapSyncer
	updates         <-chan string
}

func BuildJwksStore(ctx context.Context, mgr manager.Manager, deploymentNamespace string) *JwksStore {
	log := log.Log.WithName("jwks store setup")
	log.Info("creating jwks store")

	jwksCache := NewJwksCache()
	jwksStore := &JwksStore{
		mgr:             mgr,
		jwksCache:       jwksCache,
		jwksFetcher:     NewJwksFetcher(jwksCache),
		configMapSyncer: &ConfigMapSyncer{Client: mgr.GetClient(), DeploymentNamespace: deploymentNamespace},
	}
	jwksStore.updates = jwksStore.jwksFetcher.SubscribeToUpdates()

	return jwksStore
}

func (s *JwksStore) Start(ctx context.Context) error {
	log := log.FromContext(ctx)

	if !s.mgr.GetCache().WaitForCacheSync(ctx) {
		return fmt.Errorf("failed waiting for caches to sync")
	}

	storedJwks, err := s.configMapSyncer.LoadJwksFromConfigMap(ctx)
	if err != nil {
		log.Error(err, "error loading jwks store from a ConfigMap")
	}
	err = s.jwksCache.LoadfromJson(storedJwks)
	if err != nil {
		log.Error(err, "error deserializing jwks store state")
	}

	go s.syncToConfigMap(ctx)
	go s.jwksFetcher.Run(ctx)

	<-ctx.Done()
	return nil
}

func (s *JwksStore) UpdateJwksSources(ctx context.Context, jwks []JwksSource) error {
	return s.jwksFetcher.UpdateJwksSources(ctx, jwks)
}

func (s *JwksStore) syncToConfigMap(ctx context.Context) {
	log := log.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-s.updates:
			log.Info("received an update")
			err := s.configMapSyncer.WriteJwksToConfigMap(ctx, update)
			if err != nil {
				log.Error(err, "error syncing jwks store state to ConfigMap")
			}
		}
	}
}

// NeedLeaderElection returns true to ensure that the JwksStore runs only on the leader
func (r *JwksStore) NeedLeaderElection() bool {
	return true
}
