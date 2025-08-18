package plugins

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/proto"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

const (
	inferencePluginName = "inference-pool-policy-plugin"
)

// InferencePlugin converts an inference pool to an agentgateway inference policy
type InferencePlugin struct {
	InferencePolicyIr InferencePolicyIr
}

type InferencePolicyIr struct {
	policies []ADPPolicy
}

func (p *InferencePolicyIr) Equals(in InferencePolicyIr) bool {
	// Compare policies slice
	if len(p.policies) != len(in.policies) {
		return false
	}
	for i, policy := range p.policies {
		if !proto.Equal(policy.Policy, in.policies[i].Policy) {
			return false
		}
	}

	return true
}

// NewInferencePlugin creates a new InferencePool policy plugin
func NewInferencePlugin(agw *AgwCollections) *InferencePlugin {
	logger := logging.New("agentgateway/plugins/inference")

	inferencePools := agw.InferencePools
	if inferencePools == nil {
		logger.Debug("inference pools collection is nil, skipping inference policy generation")
		return nil
	}

	domainSuffix := kubeutils.GetClusterDomainName()

	policyCol := krt.NewManyCollection(inferencePools, func(krtctx krt.HandlerContext, policyCR *inf.InferencePool) []ADPPolicy {
		return translateInferencePoolPolicies(krtctx, inferencePools, domainSuffix)
	})

	return &InferencePlugin{
		InferencePolicyIr: InferencePolicyIr{policies: policyCol.List()},
	}
}

// GroupKind returns the GroupKind of the policy this plugin handles
func (p *InferencePlugin) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: wellknown.InferencePoolGVK.GroupKind().Group,
		Kind:  wellknown.InferencePoolGVK.GroupKind().Kind,
	}
}

// Name returns the name of this plugin
func (p *InferencePlugin) Name() string {
	return inferencePluginName
}

// ApplyPolicies applies agentgateway policies for inference pools
func (p *InferencePlugin) ApplyPolicies() []ADPPolicy {
	return p.InferencePolicyIr.policies
}

// translateInferencePoolPolicies generates policies for inference pools
func translateInferencePoolPolicies(ctx krt.HandlerContext, inferencePools krt.Collection[*inf.InferencePool], domainSuffix string) []ADPPolicy {
	logger := logging.New("agentgateway/plugins/inference")
	logger.Debug("generating inference pool policies")

	var inferencePolicies []ADPPolicy

	// Fetch all inference pools and process them
	allInferencePools := krt.Fetch(ctx, inferencePools)

	for _, pool := range allInferencePools {
		policies := translatePoliciesForInferencePool(pool, domainSuffix)
		inferencePolicies = append(inferencePolicies, policies...)
	}

	logger.Info("generated inference pool policies", "count", len(inferencePolicies))
	return inferencePolicies
}

// translatePoliciesForInferencePool generates policies for a single inference pool
func translatePoliciesForInferencePool(pool *inf.InferencePool, domainSuffix string) []ADPPolicy {
	logger := logging.New("agentgateway/plugins/inference")

	// 'service/{namespace}/{hostname}:{port}'
	svc := fmt.Sprintf("service/%v/%v.%v.inference.%v:%v",
		pool.Namespace, pool.Name, pool.Namespace, domainSuffix, pool.Spec.TargetPortNumber)

	er := pool.Spec.ExtensionRef
	if er == nil {
		logger.Debug("inference pool has no extension ref", "pool", pool.Name)
		return nil
	}

	erf := er.ExtensionReference
	if erf.Group != nil && *erf.Group != "" {
		logger.Debug("inference pool extension ref has non-empty group, skipping", "pool", pool.Name, "group", *erf.Group)
		return nil
	}

	if erf.Kind != nil && *erf.Kind != "Service" {
		logger.Debug("inference pool extension ref is not a Service, skipping", "pool", pool.Name, "kind", *erf.Kind)
		return nil
	}

	eppPort := ptr.OrDefault(erf.PortNumber, 9002)

	eppSvc := fmt.Sprintf("%v/%v.%v.svc.%v",
		pool.Namespace, erf.Name, pool.Namespace, domainSuffix)
	eppPolicyTarget := fmt.Sprintf("service/%v:%v",
		eppSvc, eppPort)

	failureMode := api.PolicySpec_InferenceRouting_FAIL_CLOSED
	if er.FailureMode == nil || *er.FailureMode == inf.FailOpen {
		failureMode = api.PolicySpec_InferenceRouting_FAIL_OPEN
	}

	// Create the inference routing policy
	inferencePolicy := &api.Policy{
		Name:   pool.Namespace + "/" + pool.Name + ":inference",
		Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: svc}},
		Spec: &api.PolicySpec{
			Kind: &api.PolicySpec_InferenceRouting_{
				InferenceRouting: &api.PolicySpec_InferenceRouting{
					EndpointPicker: &api.BackendReference{
						Kind: &api.BackendReference_Service{Service: eppSvc},
						Port: uint32(eppPort),
					},
					FailureMode: failureMode,
				},
			},
		},
	}

	// Create the TLS policy for the endpoint picker
	// TODO: we would want some way if they explicitly set a BackendTLSPolicy for the EPP to respect that
	inferencePolicyTLS := &api.Policy{
		Name:   pool.Namespace + "/" + pool.Name + ":inferencetls",
		Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: eppPolicyTarget}},
		Spec: &api.PolicySpec{
			Kind: &api.PolicySpec_BackendTls{
				BackendTls: &api.PolicySpec_BackendTLS{
					// The spec mandates this :vomit:
					Insecure: wrappers.Bool(true),
				},
			},
		},
	}

	logger.Debug("generated inference pool policies",
		"pool", pool.Name,
		"namespace", pool.Namespace,
		"inference_policy", inferencePolicy.Name,
		"tls_policy", inferencePolicyTLS.Name)

	return []ADPPolicy{
		{Policy: inferencePolicy},
		{Policy: inferencePolicyTLS},
	}
}

// Verify that InferencePlugin implements the required interfaces
var _ PolicyPlugin = (*InferencePlugin)(nil)
var _ AgentgatewayPlugin = (*InferencePlugin)(nil)
