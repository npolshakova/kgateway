package agentgateway

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
)

// extprocIR represents external processing policy configuration
type extAuthIR struct {
	provider            *TrafficPolicyAgentGatewayExtensionIR
	Extauth             *api.PolicySpec_ExtAuthz
	disableAllProviders bool
}

var _ PolicySubIR = &extAuthIR{}

// Equals compares two extAuthIR instances for equality
func (e *extAuthIR) Equals(other PolicySubIR) bool {
	otherExtProc, ok := other.(*extAuthIR)
	if !ok {
		return false
	}
	if e == nil || otherExtProc == nil {
		return e == nil && otherExtProc == nil
	}
	if e.disableAllProviders != otherExtProc.disableAllProviders {
		return false
	}
	// TODO: Add proper comparison for provider and perRoute
	return true
}

// Validate performs validation on the extprocIR
func (e *extAuthIR) Validate() error {
	if e == nil {
		return nil
	}
	// TODO: Implement validation
	return nil
}

// constructExtAuth constructs the external authentication policy IR from the policy specification.
func constructExtAuth(
	krtctx krt.HandlerContext,
	in *v1alpha1.TrafficPolicy,
	fetchGatewayExtension FetchGatewayExtensionFunc,
	out *trafficPolicySpecIr,
) error {
	spec := in.Spec.ExtAuth
	if spec == nil {
		return nil
	}

	if spec.Disable != nil {
		out.ExtAuth = &extAuthIR{
			disableAllProviders: true,
		}
		return nil
	}

	perRouteCfg := buildExtAuthPerRouteFilterConfig(spec)

	provider, err := fetchGatewayExtension(krtctx, spec.ExtensionRef, in.GetNamespace())
	if err != nil {
		return fmt.Errorf("extauthz: %w", err)
	}
	if provider.ExtType != v1alpha1.GatewayExtensionTypeExtAuth || provider.ExtAuth == nil {
		return pluginutils.ErrInvalidExtensionType(v1alpha1.GatewayExtensionTypeExtAuth, provider.ExtType)
	}

	out.ExtAuth = &extAuthIR{
		provider: provider,
		Extauth:  perRouteCfg,
	}
	return nil
}

func buildExtAuthPerRouteFilterConfig(
	spec *v1alpha1.ExtAuthPolicy,
) *api.PolicySpec_ExtAuthz {

	// TODO: add support for WithRequestBody

	if spec.ContextExtensions != nil {
		return &api.PolicySpec_ExtAuthz{
			ExtAuthz: &api.PolicySpec_ExternalAuth{
				Context: spec.ContextExtensions,
			},
		}
	}
	return nil
}

func (p *trafficPolicyPluginGwPass) handleExtAuth(fcn string, extAuth *extAuthIR) {
	if extAuth == nil {
		return
	}

	// TODO: handle global disable

	// TODO: translate
}
