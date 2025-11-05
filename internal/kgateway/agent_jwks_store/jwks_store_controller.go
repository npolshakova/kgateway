package agent_jwksstore

import (
	"context"
	"time"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/jwks"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const JwksStoreConfigMapName = "jwks-store"

type JwksStoreController struct {
	client       client.Client
	podNamespace string
	jwksStore    *jwks.JwksStore
}

func NewJwksStoreController(mgr manager.Manager, podNamespace string) (*JwksStoreController, error) {
	controller := &JwksStoreController{
		client:       mgr.GetClient(),
		podNamespace: podNamespace,
	}

	err := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AgentgatewayPolicy{}, builder.WithPredicates(predicate.Funcs{
			CreateFunc:  func(e event.CreateEvent) bool { return true },
			DeleteFunc:  func(e event.DeleteEvent) bool { return true },
			UpdateFunc:  func(e event.UpdateEvent) bool { return true },
			GenericFunc: func(e event.GenericEvent) bool { return false },
		})).
		Named("jwksstore").
		Complete(controller)
	if err != nil {
		return nil, err
	}

	return controller, nil
}

func (r *JwksStoreController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := contextutils.LoggerFrom(ctx)
	log.Info("reconciling jwks store")

	policies := &v1alpha1.AgentgatewayPolicyList{}
	err := r.client.List(ctx, policies, &client.ListOptions{})
	if err != nil {
		return ctrl.Result{RequeueAfter: 1 * time.Second}, err
	}

	jwksSources := make([]jwks.JwksSource, 0)
	for _, policy := range policies.Items {
		if policy.Spec.Traffic == nil || policy.Spec.Traffic.JWTAuthentication == nil {
			continue
		}

		for _, provider := range policy.Spec.Traffic.JWTAuthentication.Providers {
			if provider.JWKS.Remote == nil {
				continue
			}
			jwksSources = append(jwksSources, jwks.JwksSource{JwksURL: provider.JWKS.Remote.JwksUri, Ttl: provider.JWKS.Remote.CacheDuration.Duration})
		}
	}

	err = r.jwksStore.UpdateJwksSources(ctx, jwksSources)
	return ctrl.Result{}, err
}
