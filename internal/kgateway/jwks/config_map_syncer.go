package jwks

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kgateway-dev/kgateway/v2/pkg/apiclient"
)

const jwksStorePrefix = "jwks-store"
const JwksStoreComponent = "app.kubernetes.io/component"

var JwksStoreLabel = map[string]string{JwksStoreComponent: jwksStorePrefix}

// ConfigMapSyncer is used for writing to/reading from backing ConfigMaps
type ConfigMapSyncer struct {
	Client              apiclient.Client
	DeploymentNamespace string
}

func JwksFromConfigMap(cm *corev1.ConfigMap) (map[string]string, error) {
	jwksStore := cm.Data[jwksStorePrefix]
	jwks := make(map[string]string)
	err := json.Unmarshal(([]byte)(jwksStore), &jwks)
	if err != nil {
		return nil, err
	}
	return jwks, nil
}

func (cs *ConfigMapSyncer) WriteJwksToConfigMaps(ctx context.Context, updates map[string]string) error {
	log := log.FromContext(ctx)
	errs := make([]error, 0)

	for uri, jwks := range updates {
		switch jwks {
		case "": // empty jwks == remove the underlying ConfigMap
			err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Delete(ctx, JwksConfigMapName(uri), metav1.DeleteOptions{})
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "error deleting jwks ConfigMap")
				errs = append(errs, err)
			}
		default:
			existing, err := cs.fetchPersistedJwks(ctx, uri)
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "error fetching persisted jwks")
			}

			cmData, err := json.Marshal(map[string]string{uri: jwks})
			if err != nil {
				log.Error(err, "error serialiazing jwks")
				errs = append(errs, err)
				continue
			}

			if existing == nil {
				cm := cs.newJwksStoreConfigMap(JwksConfigMapName(uri))

				cm.Data[jwksStorePrefix] = string(cmData)
				cm, err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Create(ctx, cm, metav1.CreateOptions{})
				if err != nil {
					log.Error(err, "error persisting jwks to ConfigMap")
					errs = append(errs, err)
					continue
				}
			} else {
				existing.Data[jwksStorePrefix] = string(cmData)
				_, err = cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Update(ctx, existing, metav1.UpdateOptions{})
				if err != nil {
					log.Error(err, "error updating jwks ConfigMap")
					errs = append(errs, err)
					continue
				}
			}
		}
	}

	return errors.Join(errs...)
}

func (cs *ConfigMapSyncer) LoadJwksFromConfigMaps(ctx context.Context) (map[string]string, error) {
	log := log.FromContext(ctx)

	allPersistedJwks, err := cs.fetchAllPersistedJwks(ctx)
	if err != nil {
		log.Error(err, "error retrieving jwks ConfigMaps")
		return nil, err
	}

	if len(allPersistedJwks) == 0 {
		return nil, nil
	}

	errs := make([]error, 0)
	toret := make(map[string]string)
	for _, cm := range allPersistedJwks {
		jwks, err := JwksFromConfigMap(cm)
		if err != nil {
			log.Error(err, "error deserializing jwks ConfigMap", "ConfigMap", cm.Name)
			errs = append(errs, err)
			continue
		}

		maps.Copy(toret, jwks)
	}

	return toret, errors.Join(errs...)
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

func (cs *ConfigMapSyncer) fetchPersistedJwks(ctx context.Context, jwksUri string) (*corev1.ConfigMap, error) {
	existingJwks, err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).Get(ctx, JwksConfigMapName(jwksUri), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return existingJwks, nil
}

func (cs *ConfigMapSyncer) fetchAllPersistedJwks(ctx context.Context) ([]*corev1.ConfigMap, error) {
	allExistingStores, err := cs.Client.Kube().CoreV1().ConfigMaps(cs.DeploymentNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: JwksStoreComponent + "=" + jwksStorePrefix,
	})
	if err != nil {
		return nil, err
	}

	toret := make([]*corev1.ConfigMap, len(allExistingStores.Items))
	for idx, s := range allExistingStores.Items {
		toret[idx] = &s
	}

	return toret, nil
}

func JwksConfigMapName(jwksUri string) string {
	hash := md5.Sum([]byte(jwksUri))
	return fmt.Sprintf("%s-%s", jwksStorePrefix, hex.EncodeToString(hash[:]))
}
