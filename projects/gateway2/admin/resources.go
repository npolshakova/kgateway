package admin

import (
	"slices"

	"github.com/kgateway-dev/kgateway/projects/gateway2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/projects/gateway2/wellknown"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: these need to be updated
var (
	KubernetesCoreGVKs = []schema.GroupVersionKind{
		wellknown.SecretGVK,
		wellknown.ConfigMapGVK,
	}

	// PolicyGVKs is the set of GVKs that are used by a classic Gloo Gateway installation.
	// This is the common set of GVKs that are available when Edge Gateway APIs are being
	// used. See KubernetesGatewayIntegrationPolicyGVKs for the set of GVKs that are added
	// when the Kubernetes Gateway API is enabled
	// PolicyGVKs = []schema.GroupVersionKind{
	// 	gatewayv1.VirtualHostOptionGVK,
	// 	gatewayv1.RouteOptionGVK,
	// }

	KubernetesGatewayGVKs = []schema.GroupVersionKind{
		wellknown.GatewayClassGVK,
		wellknown.GatewayGVK,
		wellknown.HTTPRouteGVK,
		wellknown.ReferenceGrantGVK,
	}

	KubernetesGatewayIntegrationPolicyGVKs = []schema.GroupVersionKind{
		v1alpha1.GatewayParametersGVK,

		// While these are in fact Policy APIs, they are only enabled if the Kubernetes Gateway Integration is turned on
		// gatewayv1.ListenerOptionGVK,
		// gatewayv1.HttpListenerOptionGVK,
	}

	// CompleteInputSnapshotGVKs is the list of GVKs that will be returned by the InputSnapshot API
	CompleteInputSnapshotGVKs = slices.Concat(
		KubernetesCoreGVKs,
		KubernetesGatewayGVKs,
		KubernetesGatewayIntegrationPolicyGVKs,
	)
)
