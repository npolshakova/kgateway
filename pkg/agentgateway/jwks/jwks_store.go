package jwks

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type JwksStore struct {
	jwksCache       *jwksCache
	jwksFetcher     *JwksFetcher
	configMapSyncer *ConfigMapSyncer
	updates         <-chan string
}

func BuildJwksStore(ctx context.Context, deploymentNamespace string, client client.Client) *JwksStore {
	log := log.Log.WithName("jwks store setup")
	log.Info("creating jwks store")

	jwksCache := NewJwksCache()
	jwksStore := &JwksStore{
		jwksCache:       jwksCache,
		jwksFetcher:     NewJwksFetcher(jwksCache),
		configMapSyncer: &ConfigMapSyncer{Client: client, DeploymentNamespace: deploymentNamespace},
	}

	storedJwks, err := jwksStore.configMapSyncer.LoadJwksFromConfigMap(ctx)
	if err != nil {
		log.Error(err, "error loading jwks store from a ConfigMap")
	}
	err = jwksCache.LoadfromJson(storedJwks)
	if err != nil {
		log.Error(err, "error deserializing jwks store state")
	}

	jwksStore.updates = jwksStore.jwksFetcher.SubscribeToUpdates()

	return jwksStore
}

func (s *JwksStore) Start(ctx context.Context) {
	go s.syncToConfigMap(ctx)
	s.jwksFetcher.Run(ctx)
}

func (s *JwksStore) syncToConfigMap(ctx context.Context) {
	log := log.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-s.updates:
			err := s.configMapSyncer.WriteJwksToConfigMap(ctx, update)
			if err != nil {
				log.Error(err, "error syncing jwks store state to ConfigMap")
			}
		}
	}
}
