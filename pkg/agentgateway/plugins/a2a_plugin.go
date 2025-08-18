package plugins

import (
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	"google.golang.org/protobuf/proto"
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

// A2APlugin converts an a2a annotated service to an agentgateway a2a policy
type A2APlugin struct {
	A2APolicyIr A2APolicyIr
}

type A2APolicyIr struct {
	policies []ADPPolicy
}

func (p *A2APolicyIr) Equals(in A2APolicyIr) bool {
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

// NewA2APlugin creates a new A2A policy plugin
func NewA2APlugin(agw *AgwCollections) *A2APlugin {
	logger := logging.New("agentgateway/plugins/a2a")

	services := agw.Services
	if services == nil {
		logger.Warn("services collection is nil, skipping A2A policy generation")
		return nil
	}

	policies := translateA2APolicies(services)
	return &AgentgatewayPlugin{
		ContributesPolicies: map[schema.GroupKind]PolicyPlugin{wellknown.ServiceGVK.GroupKind(): &A2APlugin{A2APolicyIr: A2APolicyIr{policies: policies}})},
	}
}

// GroupKind returns the GroupKind of the policy this plugin handles
func (p *A2APlugin) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: wellknown.ServiceGVK.GroupKind().Group,
		Kind:  wellknown.ServiceGVK.GroupKind().Kind,
	}
}

// Name returns the name of this plugin
func (p *A2APlugin) Name() string {
	return a2aPluginName
}

// ApplyPolicies applies agentgateway policies for inference pools
func (p *A2APlugin) ApplyPolicies() []ADPPolicy {
	return p.A2APolicyIr.policies
}

// translateA2APolicies generates agentgateway policies for services with a2a protocol
func translateA2APolicies(services krt.Collection[*corev1.Service]) []ADPPolicy {
	logger := logging.New("agentgateway/plugins/a2a")
	logger.Debug("generating A2A policies")

	if services == nil {
		logger.Debug("services collection is nil, skipping A2A policy generation")
		return nil
	}

	var a2aPolicies []ADPPolicy
	policyCol := krt.NewManyCollection(services, func(krtctx krt.HandlerContext, svc *corev1.Service) []ADPPolicy {
		return translatePoliciesForService(svc)
	})
	a2aPolicies = policyCol.List()
	logger.Info("generated A2A policies", "count", len(a2aPolicies))
	return a2aPolicies
}

// translatePoliciesForService generates A2A policies for a single service
func translatePoliciesForService(svc *corev1.Service) []ADPPolicy {
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

	return a2aPolicies
}

// Verify that A2APlugin implements the required interfaces
var _ PolicyPlugin = (*A2APlugin)(nil)
var _ AgentgatewayPlugin = (*A2APlugin)(nil)
