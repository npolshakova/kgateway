package plugins

import (
	"fmt"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

const (
	a2aProtocol   = "kgateway.dev/a2a"
	a2aPluginName = "a2a-policy-plugin"
)

// A2APolicyIr converts an A2A service to an agentgateway policy
type A2APolicyIr struct {
	policies []ADPPolicy
	ct       time.Time
}

func (p *A2APolicyIr) CreationTime() time.Time {
	return p.ct
}

func (p *A2APolicyIr) Equals(in any) bool {
	p2, ok := in.(*A2APolicyIr)
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

// NewA2APlugin creates a new A2A policy plugin
func NewA2APlugin(agw *AgwCollections) AgentgatewayPlugin {
	gk := wellknown.ServiceGVK.GroupKind()
	policyCol := krt.NewCollection(agw.Services, func(krtctx krt.HandlerContext, svc *corev1.Service) *PolicyWrapper {
		objSrc := ir.ObjectSource{
			Group:     gk.Group,
			Kind:      gk.Kind,
			Namespace: svc.Namespace,
			Name:      svc.Name,
		}

		policyIR, errors := translatePoliciesForService(svc)
		return &PolicyWrapper{
			ObjectSource: objSrc,
			Policy:       svc,
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
func (p *A2APolicyIr) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: wellknown.ServiceGVK.GroupKind().Group,
		Kind:  wellknown.ServiceGVK.GroupKind().Kind,
	}
}

// Name returns the name of this plugin
func (p *A2APolicyIr) Name() string {
	return a2aPluginName
}

// ApplyPolicies applies agentgateway policies for A2A services
func (p *A2APolicyIr) ApplyPolicies() []ADPPolicy {
	return p.policies
}

// translatePoliciesForService generates A2A policies for a single service
func translatePoliciesForService(svc *corev1.Service) (*A2APolicyIr, []error) {
	logger := logging.New("agentgateway/plugins/a2a")
	var a2aPolicies []ADPPolicy

	for _, port := range svc.Spec.Ports {
		if port.AppProtocol != nil && *port.AppProtocol == a2aProtocol {
			logger.Debug("found A2A service", "service", svc.Name, "namespace", svc.Namespace, "port", port.Port)

			svcRef := fmt.Sprintf("%v/%v", svc.Namespace, svc.Name)
			policy := &api.Policy{
				Name:   fmt.Sprintf("a2a/%s/%s/%d", svc.Namespace, svc.Name, port.Port),
				Target: &api.PolicyTarget{Kind: &api.PolicyTarget_Backend{Backend: svcRef}},
				Spec: &api.PolicySpec{Kind: &api.PolicySpec_A2A_{
					A2A: &api.PolicySpec_A2A{},
				}},
			}

			a2aPolicies = append(a2aPolicies, ADPPolicy{Policy: policy})
		}
	}

	return &A2APolicyIr{
		policies: a2aPolicies,
	}, nil
}

// Verify that A2APolicyIr implements the required interfaces
var _ PolicyPluginPass = (*A2APolicyIr)(nil)
