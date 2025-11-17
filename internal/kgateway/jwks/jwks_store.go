package jwks

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/apiclient"
)

var JwksStoreNamespacedName = func() *types.NamespacedName {
	return nil
}

type JwksStore struct {
	mgr             manager.Manager
	jwksCache       *jwksCache
	jwksFetcher     *JwksFetcher
	configMapSyncer *ConfigMapSyncer
	updates         <-chan map[string]string
	latestJwksQueue utils.AsyncQueue[JwksSources]
}

func BuildJwksStore(ctx context.Context, cli apiclient.Client, jwksQueue utils.AsyncQueue[JwksSources], deploymentNamespace string) *JwksStore {
	log := log.Log.WithName("jwks store setup")
	log.Info("creating jwks store")

	jwksCache := NewJwksCache()
	jwksStore := &JwksStore{
		jwksCache:       jwksCache,
		latestJwksQueue: jwksQueue,
		jwksFetcher:     NewJwksFetcher(jwksCache),
		configMapSyncer: &ConfigMapSyncer{Client: cli, DeploymentNamespace: deploymentNamespace},
	}
	jwksStore.updates = jwksStore.jwksFetcher.SubscribeToUpdates()
	JwksStoreNamespacedName = func() *types.NamespacedName {
		return &types.NamespacedName{Namespace: deploymentNamespace, Name: JwksStoreName}
	}
	return jwksStore
}

func (s *JwksStore) Start(ctx context.Context) error {
	log := log.FromContext(ctx)

	storedJwks, err := s.configMapSyncer.LoadJwksFromConfigMap(ctx)
	if err != nil {
		log.Error(err, "error loading jwks store from a ConfigMap")
	}
	jwks, err := LoadJwksfromJson(storedJwks)
	if err != nil {
		log.Error(err, "error deserializing jwks store state")
	}
	s.jwksCache.Jwks = jwks

	go s.syncToConfigMap(ctx)
	go s.jwksFetcher.Run(ctx)
	go s.updateJwksSources(ctx)

	<-ctx.Done()
	return nil
}

func (s *JwksStore) updateJwksSources(ctx context.Context) {
	log := log.FromContext(ctx)
	for {
		latestJwks, err := s.latestJwksQueue.Dequeue(ctx)
		if err != nil {
			log.Error(err, "error dequeuing jwks update")
			return
		}
		s.jwksFetcher.UpdateJwksSources(ctx, latestJwks)
	}
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

// JwksStore runs only on the leader
func (r *JwksStore) NeedLeaderElection() bool {
	return true
}
