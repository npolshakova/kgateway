package plugins

import (
	"fmt"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

const (
	inferencePluginName = "inference-pool-policy-plugin"
)

// InferencePluginIr converts a InferencePool to an agentgateway policy
type InferencePluginIr struct {
	policies []ADPPolicy
	ct       time.Time
}

func (p *InferencePluginIr) CreationTime() time.Time {
	return p.ct
}

func (p *InferencePluginIr) Equals(in any) bool {
	p2, ok := in.(*InferencePluginIr)
	if !ok {
		return false
	}
	if len(p.policies) != len(p2.policies) {
		return false
	}
	for i, policy := range p.policies {
		if !policy.Equals(&p2.policies[i]) {
			return false
		}
	}
	return true
}

// NewInferencePlugin creates a new InferencePool policy plugin
func NewInferencePlugin(agw *AgwCollections) AgentgatewayPlugin {
	gk := wellknown.ServiceGVK.GroupKind()
	domainSuffix := kubeutils.GetClusterDomainName()
	policyCol := krt.NewCollection(agw.InferencePools, func(krtctx krt.HandlerContext, infPool *inf.InferencePool) *PolicyWrapper {
		objSrc := ir.ObjectSource{
			Group:     gk.Group,
			Kind:      gk.Kind,
			Namespace: infPool.Namespace,
			Name:      infPool.Name,
		}

		policyIR, errors := translatePoliciesForInferencePool(infPool, domainSuffix)
		return &PolicyWrapper{
			ObjectSource: objSrc,
			Policy:       infPool,
			PolicyIR:     policyIR,
			Errors:       errors,
		}
	})
	return AgentgatewayPlugin{
		ContributesPolicies: map[schema.GroupKind]PolicyPlugin{
			wellknown.ServiceGVK.GroupKind(): {
				Policies: policyCol,
			},
		},
	}
}

// GroupKind returns the GroupKind of the policy this plugin handles
func (p *InferencePluginIr) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: wellknown.InferencePoolGVK.GroupKind().Group,
		Kind:  wellknown.InferencePoolGVK.GroupKind().Kind,
	}
}

// Name returns the name of this plugin
func (p *InferencePluginIr) Name() string {
	return inferencePluginName
}

// ApplyPolicies applies agentgateway policies for inference pools
func (p *InferencePluginIr) ApplyPolicies() []ADPPolicy {
	return p.policies
}

// translatePoliciesForInferencePool generates policies for a single inference pool
func translatePoliciesForInferencePool(pool *inf.InferencePool, domainSuffix string) (*InferencePluginIr, []error) {
	logger := logging.New("agentgateway/plugins/inference")
	var errors []error

	// 'service/{namespace}/{hostname}:{port}'
	svc := fmt.Sprintf("service/%v/%v.%v.inference.%v:%v",
		pool.Namespace, pool.Name, pool.Namespace, domainSuffix, pool.Spec.TargetPortNumber)

	er := pool.Spec.ExtensionRef
	if er == nil {
		logger.Debug("inference pool has no extension ref", "pool", pool.Name)
		errors = append(errors, fmt.Errorf("inference pool has no extension ref: %v", pool.Name))
		return nil, errors
	}

	erf := er.ExtensionReference
	if erf.Group != nil && *erf.Group != "" {
		logger.Debug("inference pool extension ref has non-empty group, skipping", "pool", pool.Name, "group", *erf.Group)
		errors = append(errors, fmt.Errorf("inference pool extension ref has non-empty group: %v", pool.Name))
		return nil, errors
	}

	if erf.Kind != nil && *erf.Kind != "Service" {
		logger.Debug("inference pool extension ref is not a Service, skipping", "pool", pool.Name, "kind", *erf.Kind)
		errors = append(errors, fmt.Errorf("inference pool extension ref is not a Service: %v", pool.Name))
		return nil, errors
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

	return &InferencePluginIr{
		policies: []ADPPolicy{
			{Policy: inferencePolicy},
			{Policy: inferencePolicyTLS},
		},
	}, errors
}

// Verify that InferencePlugin implements the required interfaces
var _ PolicyPluginPass = (*InferencePluginIr)(nil)
