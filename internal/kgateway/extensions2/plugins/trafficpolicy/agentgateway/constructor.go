package agentgateway

import (
	"context"
	"fmt"

	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
)

// FetchGatewayExtensionFunc defines the signature for fetching gateway extensions
type FetchGatewayExtensionFunc func(krtctx krt.HandlerContext, extensionRef *corev1.LocalObjectReference, ns string) (*TrafficPolicyAgentGatewayExtensionIR, error)

type AgwTrafficPolicyConstructor struct {
	commoncol         *common.CommonCollections
	gatewayExtensions krt.Collection[TrafficPolicyAgentGatewayExtensionIR]
	extBuilder        func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgentGatewayExtensionIR
}

func NewAgwTrafficPolicyConstructor(
	ctx context.Context,
	commoncol *common.CommonCollections,
) *AgwTrafficPolicyConstructor {
	extBuilder := TranslateGatewayExtensionBuilder(commoncol)
	defaultExtBuilder := func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgentGatewayExtensionIR {
		return extBuilder(krtctx, gExt)
	}
	gatewayExtensions := krt.NewCollection(commoncol.GatewayExtensions, defaultExtBuilder)
	return &AgwTrafficPolicyConstructor{
		commoncol:         commoncol,
		gatewayExtensions: gatewayExtensions,
		extBuilder:        extBuilder,
	}
}

func (c *AgwTrafficPolicyConstructor) ConstructIR(
	krtctx krt.HandlerContext,
	policyCR *v1alpha1.TrafficPolicy,
) (*AgwTrafficPolicyIr, []error) {
	policyIr := AgwTrafficPolicyIr{
		ct: policyCR.CreationTimestamp.Time,
	}
	outSpec := agwTrafficPolicySpecIr{}

	var errors []error
	// Construct extproc specific IR
	if err := constructExtAuth(krtctx, policyCR, c.FetchGatewayExtension, &outSpec); err != nil {
		errors = append(errors, err)
	}

	for _, err := range errors {
		logger.Error("error translating gateway extension", "namespace", policyCR.GetNamespace(), "name", policyCR.GetName(), "error", err)
	}
	policyIr.Spec = outSpec

	return &policyIr, errors
}

func (c *AgwTrafficPolicyConstructor) FetchGatewayExtension(krtctx krt.HandlerContext, extensionRef *corev1.LocalObjectReference, ns string) (*TrafficPolicyAgentGatewayExtensionIR, error) {
	if extensionRef == nil {
		return nil, fmt.Errorf("gateway extension ref is nil")
	}

	gwExtNN := types.NamespacedName{Name: extensionRef.Name, Namespace: ns}
	gatewayExtension := krt.FetchOne(krtctx, c.gatewayExtensions, krt.FilterObjectName(gwExtNN))
	if gatewayExtension == nil {
		return nil, fmt.Errorf("gateway extension %s not found", gwExtNN.String())
	}
	if gatewayExtension.Err != nil {
		return gatewayExtension, gatewayExtension.Err
	}
	return gatewayExtension, nil
}

func (c *AgwTrafficPolicyConstructor) HasSynced() bool {
	return c.gatewayExtensions.HasSynced()
}
