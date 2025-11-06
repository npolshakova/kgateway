package jwks

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const JwksStoreName = "jwks-store"

type ConfigMapSyncer struct {
	Client              client.Client
	DeploymentNamespace string
}

func (cs *ConfigMapSyncer) WriteJwksToConfigMap(ctx context.Context, cmUpdate string) error {
	log := log.FromContext(ctx)

	existing := &corev1.ConfigMap{}
	err := cs.Client.Get(ctx, client.ObjectKey{Namespace: cs.DeploymentNamespace, Name: JwksStoreName}, existing)
	if client.IgnoreNotFound(err) != nil {
		log.Error(err, "error retrieiving jwks ConfigMap store")
		return err
	}

	if err != nil { // not found
		cm := cs.newJwksStoreConfigMap()
		cm.Data = map[string]string{"jwks": cmUpdate}
		err = cs.Client.Create(ctx, cm, &client.CreateOptions{})
		if err != nil {
			log.Error(err, "error creating jwks ConfigMap store")
			return err
		}
		return nil
	}

	existing.Data["jwks"] = cmUpdate
	err = cs.Client.Update(ctx, existing, &client.UpdateOptions{})
	if err != nil {
		log.Error(err, "error updating jwks ConfigMap store")
		return err
	}

	return nil
}

func (cs *ConfigMapSyncer) LoadJwksFromConfigMap(ctx context.Context) (string, error) {
	log := log.FromContext(ctx)

	jwksStoreConfigMap := &corev1.ConfigMap{}
	err := cs.Client.Get(ctx, client.ObjectKey{Namespace: cs.DeploymentNamespace, Name: JwksStoreName}, jwksStoreConfigMap)

	if apierrors.IsNotFound(err) {
		return "", nil
	}
	if err != nil {
		log.Error(err, "error retrieving jwks ConfigMap store")
		return "", err
	}

	return jwksStoreConfigMap.Data["jwks"], nil
}

func (cs *ConfigMapSyncer) newJwksStoreConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JwksStoreName,
			Namespace: cs.DeploymentNamespace,
		},
	}
}
