package agentgateway

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

type TrafficPolicyAgentGatewayExtensionIR struct {
	Name    string
	ExtType v1alpha1.GatewayExtensionType
	ExtAuth *api.PolicySpec_ExtAuthz
	Err     error
}

// ResourceName returns the unique name for this extension.
func (e TrafficPolicyAgentGatewayExtensionIR) ResourceName() string {
	return e.Name
}

func (e TrafficPolicyAgentGatewayExtensionIR) Equals(other TrafficPolicyAgentGatewayExtensionIR) bool {
	if e.ExtType != other.ExtType {
		return false
	}

	if !proto.Equal(e.ExtAuth.ExtAuthz, other.ExtAuth.ExtAuthz) {
		return false
	}

	if e.Err == nil && other.Err != nil {
		return false
	}
	if e.Err != nil && other.Err == nil {
		return false
	}
	if (e.Err != nil && other.Err != nil) && e.Err.Error() != other.Err.Error() {
		return false
	}

	return true
}

// Validate performs PGV-based validation on the gateway extension components
func (e TrafficPolicyAgentGatewayExtensionIR) Validate() error {
	// TODO: implement checks on agentgateway api
	return nil
}

func TranslateGatewayExtensionBuilder(commoncol *common.CommonCollections) func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgentGatewayExtensionIR {
	return func(krtctx krt.HandlerContext, gExt ir.GatewayExtension) *TrafficPolicyAgentGatewayExtensionIR {
		p := &TrafficPolicyAgentGatewayExtensionIR{
			Name:    krt.Named{Name: gExt.Name, Namespace: gExt.Namespace}.ResourceName(),
			ExtType: gExt.Type,
		}

		switch gExt.Type {
		case v1alpha1.GatewayExtensionTypeExtAuth:
			grpcService, err := ResolveExtGrpcService(krtctx, commoncol.BackendIndex, false, gExt.ObjectSource, gExt.ExtAuth.GrpcService)
			if err != nil {
				// TODO: should this be a warning
				p.Err = fmt.Errorf("failed to resolve ExtAuth backend: %w", err)
				return p
			}

			p.ExtAuth = &api.PolicySpec_ExtAuthz{
				ExtAuthz: &api.PolicySpec_ExternalAuth{
					Target: &api.BackendReference{
						Kind: grpcService,
					},
					Context: map[string]string{
						"": "",
					},
				},
			}
		default:
			slog.Warn("agentgateway does not support GatewayExtension type for type", "type", gExt.Type)
		}
		return p
	}
}

func ResolveExtGrpcService(krtctx krt.HandlerContext, backends *krtcollections.BackendIndex, disableExtensionRefValidation bool, objectSource ir.ObjectSource, grpcService *v1alpha1.ExtGrpcService) (*api.BackendReference_Service, error) {
	var backend *ir.BackendObjectIR
	if grpcService != nil {
		if grpcService.BackendRef == nil {
			return nil, errors.New("backend not provided")
		}
		backendRef := grpcService.BackendRef.BackendObjectReference

		var err error
		if disableExtensionRefValidation {
			backend, err = backends.GetBackendFromRefWithoutRefGrantValidation(krtctx, objectSource, backendRef)
		} else {
			backend, err = backends.GetBackendFromRef(krtctx, objectSource, backendRef)
		}
		if err != nil {
			return nil, err
		}
	}
	if backend == nil {
		return nil, errors.New("backend not found")
	}
	namespace := backend.GetNamespace()
	hostname := kubeutils.GetServiceHostname(backend.GetName(), namespace)
	envoyGrpcService := &api.BackendReference_Service{
		Service: namespace + "/" + hostname,
	}
	return envoyGrpcService, nil
}
