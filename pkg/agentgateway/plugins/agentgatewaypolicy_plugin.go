package plugins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
)

var logger = logging.New("agentgateway/plugins")

// convertStatusCollection converts the specific TrafficPolicy status collection
// to the generic controllers.Object status collection expected by the interface
func convertStatusCollection(col krt.Collection[krt.ObjectWithStatus[*v1alpha1.AgentgatewayPolicy, gwv1.PolicyStatus]]) krt.StatusCollection[controllers.Object, gwv1.PolicyStatus] {
	return krt.MapCollection(col, func(item krt.ObjectWithStatus[*v1alpha1.AgentgatewayPolicy, gwv1.PolicyStatus]) krt.ObjectWithStatus[controllers.Object, gwv1.PolicyStatus] {
		return krt.ObjectWithStatus[controllers.Object, gwv1.PolicyStatus]{
			Obj:    controllers.Object(item.Obj),
			Status: item.Status,
		}
	})
}

// NewAgentPlugin creates a new AgentgatewayPolicy plugin
func NewAgentPlugin(agw *AgwCollections) AgwPlugin {
	col := krt.WrapClient(kclient.NewFilteredDelayed[*v1alpha1.AgentgatewayPolicy](
		agw.Client,
		wellknown.AgentgatewayPolicyGVR,
		kclient.Filter{ObjectFilter: agw.Client.ObjectFilter()},
	), agw.KrtOpts.ToOptions("AgentgatewayPolicy")...)
	policyStatusCol, policyCol := krt.NewStatusManyCollection(col, func(krtctx krt.HandlerContext, policyCR *v1alpha1.AgentgatewayPolicy) (
		*gwv1.PolicyStatus,
		[]AgwPolicy,
	) {
		return TranslateAgentgatewayPolicy(krtctx, policyCR, agw)
	})

	return AgwPlugin{
		ContributesPolicies: map[schema.GroupKind]PolicyPlugin{
			wellknown.AgentgatewayPolicyGVK.GroupKind(): {
				Policies:       policyCol,
				PolicyStatuses: convertStatusCollection(policyStatusCol),
			},
		},
		ExtraHasSynced: func() bool {
			return policyCol.HasSynced() && policyStatusCol.HasSynced()
		},
	}
}

type PolicyCtx struct {
	Krt         krt.HandlerContext
	Collections *AgwCollections
}

type ResolvedTarget struct {
	AgentgatewayTarget *api.PolicyTarget
	GatewayTarget      gwv1.ParentReference
}

// TranslateAgentgatewayPolicy generates policies for a single traffic policy
func TranslateAgentgatewayPolicy(
	ctx krt.HandlerContext,
	policy *v1alpha1.AgentgatewayPolicy,
	agw *AgwCollections,
) (*gwv1.PolicyStatus, []AgwPolicy) {
	var agwPolicies []AgwPolicy

	pctx := PolicyCtx{Krt: ctx, Collections: agw}

	var policyTargets []ResolvedTarget
	// TODO: add selectors
	for _, target := range policy.Spec.TargetRefs {
		var policyTarget *api.PolicyTarget
		// Build a base ParentReference for status

		gk := schema.GroupKind{
			Group: string(target.Group),
			Kind:  string(target.Kind),
		}
		parentRef := gwv1.ParentReference{
			Name:      target.Name,
			Namespace: ptr.Of(gwv1.Namespace(policy.Namespace)),
			Group:     ptr.Of(gwv1.Group(gk.Group)),
			Kind:      ptr.Of(gwv1.Kind(gk.Kind)),
		}
		if target.SectionName != nil {
			parentRef.SectionName = target.SectionName
		}
		// TODO: add support for XListenerSet
		switch gk {
		case wellknown.GatewayGVK.GroupKind():
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Gateway{
					Gateway: utils.InternalGatewayName(policy.Namespace, string(target.Name), ""),
				},
			}
			if target.SectionName != nil {
				policyTarget = &api.PolicyTarget{
					Kind: &api.PolicyTarget_Listener{
						Listener: utils.InternalGatewayName(policy.Namespace, string(target.Name), string(*target.SectionName)),
					},
				}
			}

		case wellknown.HTTPRouteGVK.GroupKind():
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Route{
					Route: utils.InternalRouteRuleName(policy.Namespace, string(target.Name), ""),
				},
			}
			if target.SectionName != nil {
				policyTarget = &api.PolicyTarget{
					Kind: &api.PolicyTarget_RouteRule{
						RouteRule: utils.InternalRouteRuleName(policy.Namespace, string(target.Name), string(*target.SectionName)),
					},
				}
			}

		case wellknown.BackendGVK.GroupKind():
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Backend{
					Backend: utils.InternalBackendName(policy.Namespace, string(target.Name), ""),
				},
			}
			if target.SectionName != nil {
				policyTarget = &api.PolicyTarget{
					Kind: &api.PolicyTarget_SubBackend{
						SubBackend: utils.InternalBackendName(policy.Namespace, string(target.Name), string(*target.SectionName)),
					},
				}
			}

		case wellknown.ServiceGVK.GroupKind():
			hostname := kubeutils.GetServiceHostname(string(target.Name), policy.Namespace)
			policyTarget = &api.PolicyTarget{
				Kind: &api.PolicyTarget_Service{
					Service: policy.Namespace + "/" + hostname,
				},
			}
			if target.SectionName != nil {
				policyTarget = &api.PolicyTarget{
					Kind: &api.PolicyTarget_Backend{
						Backend: fmt.Sprintf("service/%s/%s:%s", policy.Namespace, hostname, *target.SectionName),
					},
				}
			}

			// TODO: inferencepool

		default:
			// TODO(npolshak): support attaching policies to k8s services, serviceentries, and other backends
			logger.Warn("unsupported target kind", "kind", target.Kind, "policy", policy.Name)
			continue
		}
		policyTargets = append(policyTargets, ResolvedTarget{
			AgentgatewayTarget: policyTarget,
			GatewayTarget:      parentRef,
		})
	}

	var ancestors []gwv1.PolicyAncestorStatus
	for _, policyTarget := range policyTargets {
		translatedPolicies, err := translatePolicyToAgw(pctx, policy, policyTarget.AgentgatewayTarget)
		agwPolicies = append(agwPolicies, translatedPolicies...)
		var conds []metav1.Condition
		if err != nil {
			// If we produced some policies alongside errors, treat as partial validity
			if len(translatedPolicies) > 0 {
				meta.SetStatusCondition(&conds, metav1.Condition{
					Type:    string(v1alpha1.PolicyConditionAccepted),
					Status:  metav1.ConditionTrue,
					Reason:  string(v1alpha1.PolicyReasonPartiallyValid),
					Message: err.Error(),
				})
			} else {
				// No policies produced and error present -> invalid
				meta.SetStatusCondition(&conds, metav1.Condition{
					Type:    string(v1alpha1.PolicyConditionAccepted),
					Status:  metav1.ConditionTrue,
					Reason:  string(v1alpha1.PolicyReasonInvalid),
					Message: err.Error(),
				})
				meta.SetStatusCondition(&conds, metav1.Condition{
					Type:    string(v1alpha1.PolicyConditionAttached),
					Status:  metav1.ConditionFalse,
					Reason:  string(v1alpha1.PolicyReasonPending),
					Message: "Policy is not attached due to invalid status",
				})
			}
		} else {
			// Check for partial validity
			// Build success conditions per ancestor
			meta.SetStatusCondition(&conds, metav1.Condition{
				Type:    string(v1alpha1.PolicyConditionAccepted),
				Status:  metav1.ConditionTrue,
				Reason:  string(v1alpha1.PolicyReasonValid),
				Message: reporter.PolicyAcceptedMsg,
			})
			meta.SetStatusCondition(&conds, metav1.Condition{
				Type:    string(v1alpha1.PolicyConditionAttached),
				Status:  metav1.ConditionTrue,
				Reason:  string(v1alpha1.PolicyReasonAttached),
				Message: reporter.PolicyAttachedMsg,
			})
		}
		// TODO: validate the target exists with dataplane https://github.com/kgateway-dev/kgateway/issues/12275
		// Ensure LastTransitionTime is set for all conditions
		for i := range conds {
			if conds[i].LastTransitionTime.IsZero() {
				conds[i].LastTransitionTime = metav1.Now()
			}
		}
		// Only append valid ancestors: require non-empty controllerName and parentRef name
		if agw.ControllerName != "" && string(policyTarget.GatewayTarget.Name) != "" {
			ancestors = append(ancestors, gwv1.PolicyAncestorStatus{
				AncestorRef:    policyTarget.GatewayTarget,
				ControllerName: v1alpha2.GatewayController(agw.ControllerName),
				Conditions:     conds,
			})
		}
	}

	// Build final status from accumulated ancestors
	status := gwv1.PolicyStatus{Ancestors: ancestors}

	if len(status.Ancestors) > 15 {
		ignored := status.Ancestors[15:]
		status.Ancestors = status.Ancestors[:15]
		status.Ancestors = append(status.Ancestors, gwv1.PolicyAncestorStatus{
			AncestorRef: gwv1.ParentReference{
				Group: ptr.Of(gwv1.Group("gateway.kgateway.dev")),
				Name:  "StatusSummary",
			},
			ControllerName: gwv1.GatewayController(agw.ControllerName),
			Conditions: []metav1.Condition{
				{
					Type:    "StatusSummarized",
					Status:  metav1.ConditionTrue,
					Reason:  "StatusSummary",
					Message: fmt.Sprintf("%d AncestorRefs ignored due to max status size", len(ignored)),
				},
			},
		})
	}

	// sort all parents for consistency with Equals and for Update
	// match sorting semantics of istio/istio, see:
	// https://github.com/istio/istio/blob/6dcaa0206bcaf20e3e3b4e45e9376f0f96365571/pilot/pkg/config/kube/gateway/conditions.go#L188-L193
	slices.SortStableFunc(status.Ancestors, func(a, b gwv1.PolicyAncestorStatus) int {
		return strings.Compare(reports.ParentString(a.AncestorRef), reports.ParentString(b.AncestorRef))
	})

	return &status, agwPolicies
}

// translateTrafficPolicyToAgw converts a TrafficPolicy to agentgateway Policy resources
func translatePolicyToAgw(
	ctx PolicyCtx,
	policy *v1alpha1.AgentgatewayPolicy,
	policyTarget *api.PolicyTarget,
) ([]AgwPolicy, error) {
	agwPolicies := make([]AgwPolicy, 0)
	var errs []error

	frontend, err := translateFrontendPolicyToAgw(policy, policyTarget)
	agwPolicies = append(agwPolicies, frontend...)
	if err != nil {
		errs = append(errs, err)
	}

	traffic, err := translateTrafficPolicyToAgw(ctx, policy, policyTarget)
	agwPolicies = append(agwPolicies, traffic...)
	if err != nil {
		errs = append(errs, err)
	}

	backend, err := translateBackendPolicyToAgw(ctx, policy, policyTarget)
	agwPolicies = append(agwPolicies, backend...)
	if err != nil {
		errs = append(errs, err)
	}

	return agwPolicies, errors.Join(errs...)
}
