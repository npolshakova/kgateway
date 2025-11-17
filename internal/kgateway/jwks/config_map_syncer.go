package jwks

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const JwksStoreName = "jwks-store"

// ConfigMapSyncer is used for writing to/reading from a backing ConfigMap of jwks store data
// The format used to store jwks data is key-values of jwks-url:serialized jwks
// This is done to skip an additional serialization step during policy translation
type ConfigMapSyncer struct {
	Client              client.Client
	DeploymentNamespace string
}

func JwksFromConfigMap(cm *corev1.ConfigMap) (map[string]string, error) {
	jwksStore := cm.Data[JwksStoreName]
	allJwks := make(map[string]string)
	err := json.Unmarshal(([]byte)(jwksStore), &allJwks)
	if err != nil {
		return nil, err
	}
	return allJwks, nil
}

func (cs *ConfigMapSyncer) WriteJwksToConfigMap(ctx context.Context, updated map[string]string) error {
	log := log.FromContext(ctx)

	serializedUpdate, err := json.Marshal(updated)
	if err != nil {
		log.Error(err, "error serialiazing jwks store update")
		return err
	}

	existing := &corev1.ConfigMap{}
	err = cs.Client.Get(ctx, client.ObjectKey{Namespace: cs.DeploymentNamespace, Name: JwksStoreName}, existing)
	if client.IgnoreNotFound(err) != nil {
		log.Error(err, "error retrieving jwks ConfigMap store")
		return err
	}

	if err != nil { // not found
		cm := cs.newJwksStoreConfigMap()
		cm.Data[JwksStoreName] = string(serializedUpdate)
		err = cs.Client.Create(ctx, cm, &client.CreateOptions{})
		if err != nil {
			log.Error(err, "error creating jwks ConfigMap store")
			return err
		}
		return nil
	}

	existing.Data[JwksStoreName] = string(serializedUpdate)
	err = cs.Client.Update(ctx, existing, &client.UpdateOptions{})
	if err != nil {
		log.Error(err, "error updating jwks ConfigMap store")
		return err
	}

	return nil
}

func (cs *ConfigMapSyncer) LoadJwksFromConfigMap(ctx context.Context) (map[string]string, error) {
	log := log.FromContext(ctx)

	jwksStoreConfigMap := &corev1.ConfigMap{}
	err := cs.Client.Get(ctx, client.ObjectKey{Namespace: cs.DeploymentNamespace, Name: JwksStoreName}, jwksStoreConfigMap)

	if apierrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		log.Error(err, "error retrieving jwks ConfigMap store")
		return nil, err
	}

	return JwksFromConfigMap(jwksStoreConfigMap)
}

func (cs *ConfigMapSyncer) newJwksStoreConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JwksStoreName,
			Namespace: cs.DeploymentNamespace,
		},
		Data: make(map[string]string),
	}
}
