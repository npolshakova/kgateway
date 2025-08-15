package agentgateway

import (
	"errors"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

// AgentGatewayTrafficPolicyIr is the internal representation of an agent gateway traffic policy.
// This mirrors the envoy TrafficPolicyIr pattern by pre-resolving all dependencies
// during collection building rather than at translation time.
type AgentGatewayTrafficPolicyIr struct {
	Extauth *ExtauthIr
	Errors  []error
	Ct      time.Time
}

func (a AgentGatewayTrafficPolicyIr) CreationTime() time.Time {
	return a.Ct
}

func (a AgentGatewayTrafficPolicyIr) Equals(other any) bool {
	otherBackend, ok := other.(*AgentGatewayTrafficPolicyIr)
	if !ok {
		return false
	}

	// Compare Static IR
	if !a.Extauth.Equals(otherBackend.Extauth) {
		return false
	}

	// Compare Errors - simple string comparison
	if len(a.Errors) != len(otherBackend.Errors) {
		return false
	}
	for i, err := range a.Errors {
		if err.Error() != otherBackend.Errors[i].Error() {
			return false
		}
	}

	return true
}

// Validate performs PGV-based validation on the traffic policy Ir
func (a AgentGatewayTrafficPolicyIr) Validate() error {
	// TODO: implement validation checks on agentgateway api
	return nil
}

// ExtauthIr contains pre-resolved data for external auth
type ExtauthIr struct {
	// Pre-resolved extauth configuration
	ExtAuthz *api.PolicySpec_ExtAuthz
}

func (s *ExtauthIr) Equals(other *ExtauthIr) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil || other == nil {
		return false
	}

	// Use protobuf equality for api.Extauth
	if !(s.ExtAuthz == nil && other.ExtAuthz == nil) {
		if !proto.Equal(s.ExtAuthz.ExtAuthz, other.ExtAuthz.ExtAuthz) {
			return false
		}
	}
	return true
}

type TrafficPolicyAgwExtensionIR struct {
	Name    string
	ExtType v1alpha1.GatewayExtensionType
	ExtAuth *api.BackendReference
	Err     error
}

// ResourceName returns the unique name for this extension.
func (e TrafficPolicyAgwExtensionIR) ResourceName() string {
	return e.Name
}

func (e TrafficPolicyAgwExtensionIR) Equals(other TrafficPolicyAgwExtensionIR) bool {
	if e.ExtType != other.ExtType {
		return false
	}

	if !proto.Equal(e.ExtAuth, other.ExtAuth) {
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
func (e TrafficPolicyAgwExtensionIR) Validate() error {
	// TODO: implement checks on agentgateway api
	return nil
}

func ResolveExtGrpcService(krtctx krt.HandlerContext, backends *krtcollections.BackendIndex, disableExtensionRefValidation bool, objectSource ir.ObjectSource, grpcService *v1alpha1.ExtGrpcService) (*api.BackendReference, error) {
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
	extauthSvc := &api.BackendReference{
		Kind: &api.BackendReference_Service{
			Service: namespace + "/" + hostname,
		},
		Port: uint32(backend.Port),
	}
	return extauthSvc, nil
}
