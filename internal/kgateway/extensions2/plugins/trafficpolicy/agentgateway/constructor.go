package agentgateway

import (
	"context"
	"fmt"
	"log/slog"

	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	common "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/collections"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

// FetchAgwExtensionFunc defines the signature for fetching gateway extensions
type FetchAgwExtensionFunc func(krtctx krt.HandlerContext, extensionRef *corev1.LocalObjectReference, ns string) (*TrafficPolicyAgwExtensionIR, error)

func TranslateGatewayExtensionBuilder(commoncol *common.CommonCollections) func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgwExtensionIR {
	return func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgwExtensionIR {
		p := &TrafficPolicyAgwExtensionIR{
			Name:    krt.Named{Name: gExt.Name, Namespace: gExt.Namespace}.ResourceName(),
			ExtType: gExt.Type,
		}

		switch gExt.Type {
		case v1alpha1.GatewayExtensionTypeExtAuth:
			extauthService, err := ResolveExtGrpcService(krtctx, commoncol.BackendIndex, false, gExt.ObjectSource, gExt.ExtAuth.GrpcService)
			if err != nil {
				p.Err = fmt.Errorf("failed to resolve ExtAuth backend: %w", err)
				return p
			}

			p.ExtAuth = extauthService
		default:
			slog.Warn("agentgateway does not support GatewayExtension type for type", "type", gExt.Type)
		}
		return p
	}
}

// MergeTrafficPolicies merges two AgentGatewayTrafficPolicyIr IRs, returning a map that contains information
// about the origin policy reference for each merged field.
func MergeTrafficPolicies(
	p1, p2 *AgentGatewayTrafficPolicyIr,
	p2Ref *ir.AttachedPolicyRef,
	p2MergeOrigins ir.MergeOrigins,
	mergeOpts policy.MergeOptions,
	mergeOrigins ir.MergeOrigins,
) {
	if p1 == nil || p2 == nil {
		return
	}

	mergeFuncs := []func(*AgentGatewayTrafficPolicyIr, *AgentGatewayTrafficPolicyIr, *ir.AttachedPolicyRef, ir.MergeOrigins, policy.MergeOptions, ir.MergeOrigins){
		mergeExtAuth,
	}

	for _, mergeFunc := range mergeFuncs {
		mergeFunc(p1, p2, p2Ref, p2MergeOrigins, mergeOpts, mergeOrigins)
	}
}

type TrafficPolicyConstructor struct {
	commoncol     *common.CommonCollections
	agwExtensions krt.Collection[TrafficPolicyAgwExtensionIR]
	extBuilder    func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgwExtensionIR
}

func NewTrafficPolicyConstructor(
	ctx context.Context,
	commoncol *common.CommonCollections,
) *TrafficPolicyConstructor {
	agwExtBuilder := TranslateGatewayExtensionBuilder(commoncol)
	agwExtFunc := func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgwExtensionIR {
		return agwExtBuilder(krtctx, gExt)
	}
	agwExtensions := krt.NewCollection(commoncol.GatewayExtensions, agwExtFunc)
	return &TrafficPolicyConstructor{
		commoncol:     commoncol,
		agwExtensions: agwExtensions,
		extBuilder:    agwExtBuilder,
	}
}

func (c *TrafficPolicyConstructor) ConstructAgentGatewayIR(
	krtctx krt.HandlerContext,
	policyCR *v1alpha1.TrafficPolicy,
) AgentGatewayTrafficPolicyIr {
	return buildAgentGatewayTrafficPolicyIr(krtctx, c.FetchAgwExtension, policyCR)
}

// buildAgentGatewayTrafficPolicyIr translates a TrafficPolicy to an AgentGatewayTrafficPolicyIr
func buildAgentGatewayTrafficPolicyIr(
	krtctx krt.HandlerContext,
	fetchGatewayExtension FetchAgwExtensionFunc,
	tp *v1alpha1.TrafficPolicy,
) AgentGatewayTrafficPolicyIr {
	tpIr := AgentGatewayTrafficPolicyIr{
		Ct: tp.CreationTimestamp.Time,
	}

	if tp.Spec.ExtAuth != nil {
		extauth, err := BuildExtauthIr(krtctx, tp, fetchGatewayExtension)
		if err != nil {
			tpIr.Errors = append(tpIr.Errors, err)
		}
		tpIr.Extauth = extauth
	}

	return tpIr
}

func (c *TrafficPolicyConstructor) FetchAgwExtension(krtctx krt.HandlerContext, extensionRef *corev1.LocalObjectReference, ns string) (*TrafficPolicyAgwExtensionIR, error) {
	if extensionRef == nil {
		return nil, fmt.Errorf("gateway extension ref is nil")
	}

	gwExtNN := types.NamespacedName{Name: extensionRef.Name, Namespace: ns}
	gatewayExtension := krt.FetchOne(krtctx, c.agwExtensions, krt.FilterObjectName(gwExtNN))
	if gatewayExtension == nil {
		return nil, fmt.Errorf("gateway extension %s not found", gwExtNN.String())
	}
	if gatewayExtension.Err != nil {
		return gatewayExtension, gatewayExtension.Err
	}
	return gatewayExtension, nil
}

func (c *TrafficPolicyConstructor) HasSynced() bool {
	return c.agwExtensions.HasSynced()
}
