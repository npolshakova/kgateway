package jwks

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kgateway-dev/kgateway/v2/pkg/apiclient"
)

const JwksStoreName = "jwks-store"
const JwksStoreComponent = "app.kubernetes.io/component"

var JwksStoreLabel = map[string]string{JwksStoreComponent: JwksStoreName}

// ConfigMapSyncer is used for writing to/reading from a backing ConfigMap of jwks store data
// The format used to store jwks data is key-values of jwks-url:serialized jwks
// This is done to skip an additional serialization step during policy translation
type ConfigMapSyncer struct {
	Client              apiclient.Client
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

func (cs *ConfigMapSyncer) WriteJwksToConfigMap(ctx context.Context, updated []map[string]string) error {
	log := log.FromContext(ctx)

	serializedUpdates := make([]string, len(updated))
	for i, m := range updated {
		serializedUpdate, err := json.Marshal(m)
		if err != nil {
			log.Error(err, "error serialiazing jwks store update")
			return err
		}
		serializedUpdates[i] = string(serializedUpdate)
	}

	allExistingStores, err := cs.fetchExistingStores(ctx)
	if err != nil {
		log.Error(err, "error retrieving jwks ConfigMap stores")
		return err
	}

	for i, u := range serializedUpdates {
		existing, ok := allExistingStores[jwksStoreName(i)]
		if !ok {
			cm := cs.newJwksStoreConfigMap(jwksStoreName(i))
			cm.Data[JwksStoreName] = u
			cm, err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Create(ctx, cm, metav1.CreateOptions{})
			if err != nil {
				log.Error(err, "error creating jwks ConfigMap store")
				return err
			}
		} else {
			existing.Data[JwksStoreName] = u
			_, err = cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Update(ctx, existing, metav1.UpdateOptions{})
			if err != nil {
				log.Error(err, "error updating jwks ConfigMap store")
				return err
			}
		}
	}

	// do gc
	if len(allExistingStores) > len(updated) {
		for i := len(updated); i < len(allExistingStores); i++ {
			err = cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Delete(ctx, jwksStoreName(i), metav1.DeleteOptions{})
			if err != nil {
				log.Error(err, "error deleting unused jwks ConfigMap store")
				return err
			}
		}
	}

	return nil
}

func (cs *ConfigMapSyncer) LoadJwksFromConfigMaps(ctx context.Context) ([]map[string]string, error) {
	log := log.FromContext(ctx)

	allExistingStores, err := cs.fetchExistingStores(ctx)
	if err != nil {
		log.Error(err, "error retrieving jwks ConfigMap stores")
		return nil, err
	}

	if len(allExistingStores) == 0 {
		return nil, nil
	}

	toret := make([]map[string]string, len(allExistingStores))
	for i := range len(allExistingStores) {
		cm, ok := allExistingStores[jwksStoreName(i)]
		if !ok {
			return nil, fmt.Errorf("error loading jwks stores, store '%s' doesn't exist", jwksStoreName(i))
		}

		jwks, err := JwksFromConfigMap(cm)
		if err != nil {
			log.Error(err, "error deserializing jwks ConfigMap store", "store", jwksStoreName(i))
			return nil, err
		}

		toret[i] = jwks
	}

	return toret, nil
}

func (cs *ConfigMapSyncer) newJwksStoreConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cs.DeploymentNamespace,
			Labels:    JwksStoreLabel,
		},
		Data: make(map[string]string),
	}
}

func (cs *ConfigMapSyncer) fetchExistingStores(ctx context.Context) (map[string]*corev1.ConfigMap, error) {
	allExistingStores, err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: JwksStoreComponent + "=" + JwksStoreName,
	})
	if err != nil {
		return nil, err
	}

	toret := make(map[string]*corev1.ConfigMap)
	for _, s := range allExistingStores.Items {
		toret[s.Name] = &s
	}

	return toret, nil
}

func jwksStoreName(i int) string {
	return fmt.Sprintf("%s-%d", JwksStoreName, i)
}
