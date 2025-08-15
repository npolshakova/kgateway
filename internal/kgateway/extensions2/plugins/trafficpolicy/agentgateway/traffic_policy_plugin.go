package agentgateway

import (
	"context"
	"fmt"

	"github.com/agentgateway/agentgateway/go/api"
	skubeclient "istio.io/istio/pkg/config/schema/kubeclient"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	reportssdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	agwir "github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/client/clientset/versioned"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
	pluginsdkutils "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
)

var (
	logger = logging.New("plugin/agentgateway-trafficpolicy")
)

// PolicySubIR documents the expected interface that all policy sub-IRs should implement.
type PolicySubIR interface {
	// Equals compares this policy with another policy
	Equals(other PolicySubIR) bool

	// Validate performs PGV validation on the policy
	Validate() error

	// TODO: Merge. Just awkward as we won't be using the actual method type.
}

type ProviderNeededMap struct {
	// map filter_chain name -> provider name -> provider
	Providers map[string]map[string]*TrafficPolicyAgwExtensionIR
}

type agwTrafficPolicyPluginGwPass struct {
	reporter reports.Reporter
	agwir.UnimplementedAgentGatewayTranslationPass
}

func (a agwTrafficPolicyPluginGwPass) ApplyForRoute(pCtx *agwir.AgentGatewayRouteContext, out *api.Route, policies *[]*api.Policy) error {
	if pCtx.AttachedPolicies.Policies != nil {
		gk := schema.GroupKind{Group: wellknown.TrafficPolicyGVK.Group, Kind: wellknown.TrafficPolicyGVK.Kind}
		attachedPolicies, ok := pCtx.AttachedPolicies.Policies[gk]
		if !ok {
			return nil
		}
		if attachedPolicies == nil {
			return nil
		}
		for _, p := range attachedPolicies {
			rtPolicy, ok := p.PolicyIr.(AgentGatewayTrafficPolicyIr)
			if !ok {
				continue
			}
			// TODO: does the policy name need to be unique? (per target)
			var policyName string
			if p.PolicyRef != nil {
				policyName = fmt.Sprintf("trafficpolicy/%s/%s", p.PolicyRef.Namespace, p.PolicyRef.Name)
			} else {
				policyName = "trafficpolicy/unknown"
			}
			if rtPolicy.Extauth != nil {
				*policies = append(*policies, &api.Policy{
					Name: policyName + ":extauth",
					Target: &api.PolicyTarget{
						Kind: &api.PolicyTarget_Route{Route: out.GetRouteName()},
					},
					Spec: &api.PolicySpec{Kind: rtPolicy.Extauth.ExtAuthz},
				})
			}
		}
	}
	return nil
}

var _ agwir.AgentGatewayTranslationPass = &agwTrafficPolicyPluginGwPass{}

func registerTypes(ourCli versioned.Interface) {
	skubeclient.Register[*v1alpha1.TrafficPolicy](
		wellknown.TrafficPolicyGVR,
		wellknown.TrafficPolicyGVK,
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (runtime.Object, error) {
			return ourCli.GatewayV1alpha1().TrafficPolicies(namespace).List(context.Background(), o)
		},
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (watch.Interface, error) {
			return ourCli.GatewayV1alpha1().TrafficPolicies(namespace).Watch(context.Background(), o)
		},
	)
}

func NewPlugin(ctx context.Context, commoncol *common.CommonCollections) extensionsplug.Plugin {
	registerTypes(commoncol.OurClient)

	col := krt.WrapClient(kclient.NewFiltered[*v1alpha1.TrafficPolicy](
		commoncol.Client,
		kclient.Filter{ObjectFilter: commoncol.Client.ObjectFilter()},
	), commoncol.KrtOpts.ToOptions("TrafficPolicy")...)
	gk := wellknown.TrafficPolicyGVK.GroupKind()

	var errors []error
	translator := NewTrafficPolicyConstructor(ctx, commoncol)

	// TrafficPolicy IR will have TypedConfig -> implement backendroute method to add prompt guard, etc.
	policyCol := krt.NewCollection(col, func(krtctx krt.HandlerContext, policyCR *v1alpha1.TrafficPolicy) *ir.PolicyWrapper {
		objSrc := ir.ObjectSource{
			Group:     gk.Group,
			Kind:      gk.Kind,
			Namespace: policyCR.Namespace,
			Name:      policyCR.Name,
		}

		agwPolicyIr := translator.ConstructAgentGatewayIR(krtctx, policyCR)

		pol := &ir.PolicyWrapper{
			ObjectSource: objSrc,
			Policy:       policyCR,
			PolicyIR:     agwPolicyIr,
			TargetRefs:   pluginsdkutils.TargetRefsToPolicyRefsWithSectionName(policyCR.Spec.TargetRefs, policyCR.Spec.TargetSelectors),
			Errors:       errors,
		}
		return pol
	})

	return extensionsplug.Plugin{
		ContributesPolicies: map[schema.GroupKind]extensionsplug.PolicyPlugin{
			wellknown.TrafficPolicyGVK.GroupKind(): {
				NewAgentGatewayPass: func(reporter reportssdk.Reporter) agwir.AgentGatewayTranslationPass {
					return NewAgentGatewayTranslationPass(reporter)
				},
				Policies: policyCol,
				MergePolicies: func(pols []ir.PolicyAtt) ir.PolicyAtt {
					return policy.MergePolicies(pols, MergeTrafficPolicies)
				},
				// TODO: add policy status for agentgateway TrafficPolicy
				//GetPolicyStatus:   getPolicyStatusFn(commoncol.CrudClient),
				//PatchPolicyStatus: patchPolicyStatusFn(commoncol.CrudClient),
			},
		},
		ExtraHasSynced: translator.HasSynced,
	}
}

func NewAgentGatewayTranslationPass(reporter reports.Reporter) agwir.AgentGatewayTranslationPass {
	return &agwTrafficPolicyPluginGwPass{
		reporter: reporter,
	}
}
