package agentgateway

import (
	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// BuildExtauthIr pre-resolves extauth configuration
func BuildExtauthIr(krtctx krt.HandlerContext, tp *v1alpha1.TrafficPolicy, fetchGatewayExtension FetchAgwExtensionFunc) (*ExtauthIr, error) {
	var context map[string]string
	if tp.Spec.ExtAuth.ContextExtensions != nil {
		context = tp.Spec.ExtAuth.ContextExtensions
	}

	extRef := tp.Spec.ExtAuth.ExtensionRef
	extension, err := fetchGatewayExtension(krtctx, extRef, tp.GetNamespace())
	if err != nil {
		logger.Error("gateway extension not found for extauth", "error", err)
		return nil, err
	}

	extauthPolicy := &api.PolicySpec_ExtAuthz{
		ExtAuthz: &api.PolicySpec_ExternalAuth{
			Target:  extension.ExtAuth,
			Context: context,
		},
	}

	return &ExtauthIr{
		ExtAuthz: extauthPolicy,
	}, nil
}
