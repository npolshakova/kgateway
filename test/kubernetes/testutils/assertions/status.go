package assertions

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	inf "sigs.k8s.io/gateway-api-inference-extension/api/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
	"github.com/kgateway-dev/kgateway/v2/test/helpers"
)

// // Checks GetNamespacedStatuses status for gloo installation namespace
// func (p *Provider) EventuallyResourceStatusMatchesWarningReasons(getter helpers.InputResourceGetter, desiredStatusReasons []string, desiredReporter string, timeout ...time.Duration) {
// 	ginkgo.GinkgoHelper()

// 	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
// 	gomega.Eventually(func(g gomega.Gomega) {
// 		statusWarningsMatcher := matchers.MatchStatusInNamespace(
// 			p.installContext.InstallNamespace,
// 			gomega.And(matchers.HaveWarningStateWithReasonSubstrings(desiredStatusReasons...), matchers.HaveReportedBy(desiredReporter)),
// 		)

// 		status, err := getResourceNamespacedStatus(getter)
// 		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
// 		g.Expect(status).ToNot(gomega.BeNil())
// 		g.Expect(status).To(gomega.HaveValue(statusWarningsMatcher))
// 	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
// }

// func (p *Provider) EventuallyResourceStatusMatchesRejectedReasons(getter helpers.InputResourceGetter, desiredStatusReasons []string, desiredReporter string, timeout ...time.Duration) {
// 	ginkgo.GinkgoHelper()

// 	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
// 	gomega.Eventually(func(g gomega.Gomega) {
// 		statusRejectionsMatcher := matchers.MatchStatusInNamespace(
// 			p.installContext.InstallNamespace,
// 			gomega.And(matchers.HaveRejectedStateWithReasonSubstrings(desiredStatusReasons...), matchers.HaveReportedBy(desiredReporter)),
// 		)

// 		status, err := getResourceNamespacedStatus(getter)
// 		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
// 		g.Expect(status).ToNot(gomega.BeNil())
// 		g.Expect(status).To(gomega.HaveValue(statusRejectionsMatcher))
// 	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
// }

// func (p *Provider) EventuallyResourceStatusMatchesState(
// 	getter helpers.InputResourceGetter,
// 	desiredState core.Status_State,
// 	desiredReporter string,
// 	timeout ...time.Duration,
// ) {
// 	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
// 	p.Gomega.Eventually(func(g gomega.Gomega) {
// 		statusStateMatcher := matchers.MatchStatusInNamespace(
// 			p.installContext.InstallNamespace,
// 			gomega.And(matchers.HaveState(desiredState), matchers.HaveReportedBy(desiredReporter)),
// 		)
// 		status, err := getResourceNamespacedStatus(getter)
// 		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
// 		g.Expect(status).ToNot(gomega.BeNil())
// 		g.Expect(status).To(gomega.HaveValue(statusStateMatcher))
// 	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
// }

// func (p *Provider) EventuallyResourceStatusMatchesSubResource(
// 	getter helpers.InputResourceGetter,
// 	desiredSubresourceName string,
// 	desiredSubresource matchers.SoloKitSubresourceStatus,
// 	timeout ...time.Duration,
// ) {
// 	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
// 	p.Gomega.Eventually(func(g gomega.Gomega) {
// 		subResourceStatusMatcher := matchers.HaveSubResourceStatusState(desiredSubresourceName, desiredSubresource)
// 		status, err := getResourceNamespacedStatus(getter)
// 		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
// 		g.Expect(status).ToNot(gomega.BeNil())
// 		g.Expect(status).To(gomega.HaveValue(subResourceStatusMatcher))
// 	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
// }

// func getResourceNamespacedStatus(getter helpers.InputResourceGetter) (*core.NamespacedStatuses, error) {
// 	resource, err := getter()
// 	if err != nil {
// 		return &core.NamespacedStatuses{}, errors.Wrapf(err, "failed to get resource")
// 	}

// 	namespacedStatuses := resource.GetNamespacedStatuses()

// 	// In newer versions of kgateway we provide a default "empty" status, which allows us to patch it to perform updates
// 	// As a result, a nil check isn't enough to determine that that status hasn't been reported
// 	if namespacedStatuses == nil || namespacedStatuses.GetStatuses() == nil {
// 		return &core.NamespacedStatuses{}, errors.Wrapf(err, "waiting for %v status to be non-empty", resource.GetMetadata().GetName())
// 	}

// 	return namespacedStatuses, nil
// }

// EventuallyHTTPRouteStatusContainsMessage asserts that eventually at least one of the HTTPRoute's route parent statuses contains
// the given message substring.
func (p *Provider) EventuallyHTTPRouteStatusContainsMessage(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	message string,
	timeout ...time.Duration) {
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		matcher := matchers.HaveKubeGatewayRouteStatus(&matchers.KubeGatewayRouteStatus{
			Custom: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Parents": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Conditions": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Message": matchers.ContainSubstrings([]string{message}),
					})),
				})),
			}),
		})

		route := &gwv1.HTTPRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "can get httproute")
		g.Expect(route.Status.RouteStatus).To(gomega.HaveValue(matcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyHTTPRouteStatusContainsReason asserts that eventually at least one of the HTTPRoute's route parent statuses contains
// the given reason substring.
func (p *Provider) EventuallyHTTPRouteStatusContainsReason(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	reason string,
	timeout ...time.Duration,
) {
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		matcher := matchers.HaveKubeGatewayRouteStatus(&matchers.KubeGatewayRouteStatus{
			Custom: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Parents": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Conditions": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Reason": matchers.ContainSubstrings([]string{reason}),
					})),
				})),
			}),
		})

		route := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routeName,
				Namespace: routeNamespace,
			},
		}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "can get httproute")
		g.Expect(route.Status.RouteStatus).To(gomega.HaveValue(matcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGatewayCondition checks the provided Gateway condition is set to expect.
func (p *Provider) EventuallyGatewayCondition(
	ctx context.Context,
	gatewayName string,
	gatewayNamespace string,
	cond gwv1.GatewayConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		gateway := &gwv1.Gateway{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gatewayNamespace}, gateway)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get Gateway %s/%s", gatewayNamespace, gatewayName))

		condition := getConditionByType(gateway.Status.Conditions, string(cond))
		g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for Gateway %s/%s", cond, gatewayNamespace, gatewayName))
		g.Expect(condition.Status).To(gomega.Equal(expect), fmt.Sprintf("%v condition is not %v for Gateway %s/%s",
			cond, expect, gatewayNamespace, gatewayName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGatewayListenerAttachedRoutes checks the provided Gateway contains the expected attached routes for the listener.
func (p *Provider) EventuallyGatewayListenerAttachedRoutes(
	ctx context.Context,
	gatewayName string,
	gatewayNamespace string,
	listener gwv1.SectionName,
	routes int32,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		gateway := &gwv1.Gateway{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gatewayNamespace}, gateway)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get Gateway %s/%s", gatewayNamespace, gatewayName))

		found := false
		for _, l := range gateway.Status.Listeners {
			if l.Name == listener {
				found = true
				g.Expect(l.AttachedRoutes).To(gomega.Equal(routes), fmt.Sprintf("%v listener does not contain %d attached routes for Gateway %s/%s",
					l, routes, gatewayNamespace, gatewayName))
			}
		}
		g.Expect(found).To(gomega.BeTrue(), fmt.Sprintf("%v listener not found for Gateway %s/%s", listener, gatewayNamespace, gatewayName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func (p *Provider) EventuallyGatewayStatus(
	ctx context.Context,
	name string,
	namespace string,
	status gwv1.GatewayStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		gw := &gwv1.Gateway{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, gw)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get gateway %s/%s", namespace, name))

		for _, expected := range status.Conditions {
			condition := getConditionByType(gw.Status.Conditions, expected.Type)
			g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for gateway %s/%s", expected.Type, namespace, name))
			g.Expect(condition.Status).To(gomega.Equal(expected.Status), fmt.Sprintf("%v status is not %v for gateway %s/%s", expected, expected.Status, namespace, name))
			if expected.Reason != "" {
				g.Expect(condition.Reason).To(gomega.Equal(expected.Reason), fmt.Sprintf("%v reason is not %v for gateway %s/%s", expected, expected.Reason, namespace, name))
			}
		}

		for _, expectedListener := range status.Listeners {
			listenerStatus := getListenerStatus(gw.Status.Listeners, string(expectedListener.Name))
			g.Expect(listenerStatus).NotTo(gomega.BeNil(), fmt.Sprintf("%v listener status not found for listener %s", expectedListener.Name, expectedListener.Name))
			if expectedListener.AttachedRoutes != 0 {
				g.Expect(listenerStatus.AttachedRoutes).To(gomega.Equal(expectedListener.AttachedRoutes), fmt.Sprintf("%v condition is not %v for listener %s", expectedListener, expectedListener.AttachedRoutes, expectedListener.Name))
			}
			if expectedListener.SupportedKinds != nil {
				g.Expect(listenerStatus.SupportedKinds).To(gomega.ContainElements(expectedListener.SupportedKinds), fmt.Sprintf("%v condition is not %v for listener %s", expectedListener, expectedListener.SupportedKinds, expectedListener.Name))
			}

			for _, expected := range expectedListener.Conditions {
				condition := getConditionByType(listenerStatus.Conditions, expected.Type)
				g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for listener %s", expected, expectedListener.Name))
				g.Expect(condition.Status).To(gomega.Equal(expected.Status), fmt.Sprintf("%v condition is not %v for listener %s", expected, expected.Status, expectedListener.Name))
				if expected.Reason != "" {
					g.Expect(condition.Reason).To(gomega.Equal(expected.Reason), fmt.Sprintf("%v condition is not %v for listener %s", expected, expected.Reason, expectedListener.Name))
				}
			}
		}
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyHTTPRouteCondition checks that provided HTTPRoute condition is set to expect.
func (p *Provider) EventuallyHTTPRouteCondition(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	cond gwv1.RouteConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		route := &gwv1.HTTPRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get HTTPRoute %s/%s", routeNamespace, routeName)

		var conditionFound bool
		for _, parentStatus := range route.Status.Parents {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of HTTPRoute %s/%s",
			cond, expect, routeNamespace, routeName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyTCPRouteCondition checks that provided TCPRoute condition is set to expect.
func (p *Provider) EventuallyTCPRouteCondition(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	cond gwv1.RouteConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		route := &gwv1a2.TCPRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get TCPRoute %s/%s", routeNamespace, routeName)

		var conditionFound bool
		for _, parentStatus := range route.Status.Parents {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of TCPRoute %s/%s",
			cond, expect, routeNamespace, routeName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyTLSRouteCondition checks that provided TLSRoute condition is set to expect.
func (p *Provider) EventuallyTLSRouteCondition(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	cond gwv1.RouteConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		route := &gwv1a2.TLSRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get TLSRoute %s/%s", routeNamespace, routeName)

		var conditionFound bool
		for _, parentStatus := range route.Status.Parents {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of TLSRoute %s/%s",
			cond, expect, routeNamespace, routeName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGRPCRouteCondition checks that provided GRPCRoute condition is set to expect.
func (p *Provider) EventuallyGRPCRouteCondition(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	cond gwv1.RouteConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		route := &gwv1.GRPCRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get GRPCRoute %s/%s", routeNamespace, routeName)

		var conditionFound bool
		for _, parentStatus := range route.Status.Parents {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of GRPCRoute %s/%s",
			cond, expect, routeNamespace, routeName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyInferencePoolCondition checks that the specified InferencePool condition
// eventually has the desired status on any parent managed by Kgateway.
func (p *Provider) EventuallyInferencePoolCondition(
	ctx context.Context,
	poolName string,
	poolNamespace string,
	cond inf.InferencePoolConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()

	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		pool := &inf.InferencePool{}
		err := p.clusterContext.Client.Get(
			ctx,
			types.NamespacedName{Name: poolName, Namespace: poolNamespace},
			pool,
		)
		g.Expect(err).NotTo(gomega.HaveOccurred(),
			"failed to get InferencePool %s/%s", poolNamespace, poolName)

		var conditionFound bool
		for _, parent := range pool.Status.Parents {
			// Look for the first matching condition on any parent.
			if c := getConditionByType(parent.Conditions, string(cond)); c != nil && c.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(),
			fmt.Sprintf("%v condition is not %v for any parent of InferencePool %s/%s",
				cond, expect, poolNamespace, poolName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// Helper function to retrieve a condition by type from a list of conditions.
func getConditionByType(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}
	return nil
}

func (p *Provider) EventuallyListenerSetStatus(
	ctx context.Context,
	name string,
	namespace string,
	status gwxv1a1.ListenerSetStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		ls := &gwxv1a1.XListenerSet{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, ls)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get listenerset %s/%s", namespace, name))

		for _, expected := range status.Conditions {
			condition := getConditionByType(ls.Status.Conditions, expected.Type)
			g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for listenerset %s/%s", expected.Type, namespace, name))
			g.Expect(condition.Status).To(gomega.Equal(expected.Status), fmt.Sprintf("%v status is not %v for listenerset %s/%s", expected, expected.Status, namespace, name))
			if expected.Reason != "" {
				g.Expect(condition.Reason).To(gomega.Equal(expected.Reason), fmt.Sprintf("%v reason is not %v for listenerset %s/%s", expected, expected.Reason, namespace, name))
			}
		}

		for _, expectedListener := range status.Listeners {
			listenerStatus := getListenerEntryStatus(ls.Status.Listeners, string(expectedListener.Name))
			g.Expect(listenerStatus).NotTo(gomega.BeNil(), fmt.Sprintf("%v listener status not found for listener %s", expectedListener.Name, expectedListener.Name))
			if expectedListener.Port != 0 {
				g.Expect(listenerStatus.Port).To(gomega.Equal(expectedListener.Port), fmt.Sprintf("%v listener condition is not %v for listener %s", expectedListener, expectedListener.Port, expectedListener.Name))
			}
			if expectedListener.AttachedRoutes != 0 {
				g.Expect(listenerStatus.AttachedRoutes).To(gomega.Equal(expectedListener.AttachedRoutes), fmt.Sprintf("%v condition is not %v for listener %s", expectedListener, expectedListener.AttachedRoutes, expectedListener.Name))
			}
			if expectedListener.SupportedKinds != nil {
				g.Expect(listenerStatus.SupportedKinds).To(gomega.ContainElements(expectedListener.SupportedKinds), fmt.Sprintf("%v condition is not %v for listener %s", expectedListener, expectedListener.SupportedKinds, expectedListener.Name))
			}

			for _, expected := range expectedListener.Conditions {
				condition := getConditionByType(listenerStatus.Conditions, expected.Type)
				g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for listener %s", expected, expectedListener.Name))
				g.Expect(condition.Status).To(gomega.Equal(expected.Status), fmt.Sprintf("%v condition is not %v for listener %s", expected, expected.Status, expectedListener.Name))
				if expected.Reason != "" {
					g.Expect(condition.Reason).To(gomega.Equal(expected.Reason), fmt.Sprintf("%v condition is not %v for listener %s", expected, expected.Reason, expectedListener.Name))
				}
			}
		}
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func (p *Provider) EventuallyListenerSetAttachedRoutes(
	ctx context.Context,
	name string,
	namespace string,
	listener gwv1.SectionName,
	routes int32,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		ls := &gwxv1a1.XListenerSet{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, ls)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get listenerset %s/%s", namespace, name))

		for _, expectedListener := range ls.Status.Listeners {
			listenerStatus := getListenerEntryStatus(ls.Status.Listeners, string(expectedListener.Name))
			g.Expect(listenerStatus).NotTo(gomega.BeNil(), fmt.Sprintf("%v listener status not found for listener %s", expectedListener.Name, expectedListener.Name))
			g.Expect(listenerStatus.AttachedRoutes).To(gomega.Equal(expectedListener.AttachedRoutes), fmt.Sprintf("%v AttachedRoutes is not %v for listener %s", expectedListener, expectedListener.AttachedRoutes, expectedListener.Name))
		}
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func getListenerEntryStatus(listeners []gwxv1a1.ListenerEntryStatus, name string) *gwxv1a1.ListenerEntryStatus {
	for _, listener := range listeners {
		if string(listener.Name) == name {
			return &listener
		}
	}
	return nil
}

func getListenerStatus(listeners []gwv1.ListenerStatus, name string) *gwv1.ListenerStatus {
	for _, listener := range listeners {
		if string(listener.Name) == name {
			return &listener
		}
	}
	return nil
}

// EventuallyHTTPListenerPolicyCondition checks that provided HTTPListenerPolicy condition is set to expect.
func (p *Provider) EventuallyHTTPListenerPolicyCondition(
	ctx context.Context,
	name string,
	namespace string,
	cond gwv1.GatewayConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helpers.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		hlp := &v1alpha1.HTTPListenerPolicy{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, hlp)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get HTTPListenerPolicy %s/%s", namespace, name)

		var conditionFound bool
		for _, parentStatus := range hlp.Status.Ancestors {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of HTTPListenerPolicy %s/%s",
			cond, expect, namespace, name))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyTrafficPolicyAccepted validates that a TrafficPolicy has "Accepted" status.
// TrafficPolicy uses PolicyStatus with ancestors, so we check the first ancestor's conditions.
func (p *Provider) EventuallyTrafficPolicyAccepted(ctx context.Context, policy *v1alpha1.TrafficPolicy) {
	currentTimeout, pollingInterval := helpers.GetTimeouts()
	p.Gomega.Eventually(func(g gomega.Gomega) {
		obj := &v1alpha1.TrafficPolicy{}
		objKey := client.ObjectKeyFromObject(policy)
		err := p.clusterContext.Client.Get(ctx, objKey, obj)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get TrafficPolicy %s", objKey)

		g.Expect(obj.Status.Ancestors).ToNot(gomega.BeEmpty(), "TrafficPolicy should have ancestors")

		// Check first ancestor for Accepted condition
		ancestorStatus := obj.Status.Ancestors[0]
		cond := meta.FindStatusCondition(ancestorStatus.Conditions, string(v1alpha1.PolicyConditionAccepted))
		g.Expect(cond).NotTo(gomega.BeNil(), "TrafficPolicy should have Accepted condition")
		g.Expect(cond.Status).To(gomega.Equal(metav1.ConditionTrue), "TrafficPolicy should be accepted")
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyBackendAccepted validates that a Backend has "Accepted" status.
// Backend uses direct conditions on its status.
func (p *Provider) EventuallyBackendAccepted(ctx context.Context, backend *v1alpha1.Backend) {
	currentTimeout, pollingInterval := helpers.GetTimeouts()
	p.Gomega.Eventually(func(g gomega.Gomega) {
		obj := &v1alpha1.Backend{}
		objKey := client.ObjectKeyFromObject(backend)
		err := p.clusterContext.Client.Get(ctx, objKey, obj)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get Backend %s", objKey)

		cond := meta.FindStatusCondition(obj.Status.Conditions, string(v1alpha1.PolicyConditionAccepted))
		g.Expect(cond).NotTo(gomega.BeNil(), "Backend should have Accepted condition")
		g.Expect(cond.Status).To(gomega.Equal(metav1.ConditionTrue), "Backend should be accepted")
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyBackendConfigPolicyAccepted validates that a BackendConfigPolicy has "Accepted" status.
// BackendConfigPolicy uses PolicyStatus with ancestors, so we check the first ancestor's conditions.
func (p *Provider) EventuallyBackendConfigPolicyAccepted(ctx context.Context, policy *v1alpha1.BackendConfigPolicy) {
	currentTimeout, pollingInterval := helpers.GetTimeouts()
	p.Gomega.Eventually(func(g gomega.Gomega) {
		obj := &v1alpha1.BackendConfigPolicy{}
		objKey := client.ObjectKeyFromObject(policy)
		err := p.clusterContext.Client.Get(ctx, objKey, obj)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get BackendConfigPolicy %s", objKey)

		g.Expect(obj.Status.Ancestors).ToNot(gomega.BeEmpty(), "BackendConfigPolicy should have ancestors")

		// Check first ancestor for Accepted condition
		ancestorStatus := obj.Status.Ancestors[0]
		cond := meta.FindStatusCondition(ancestorStatus.Conditions, string(v1alpha1.PolicyConditionAccepted))
		g.Expect(cond).NotTo(gomega.BeNil(), "BackendConfigPolicy should have Accepted condition")
		g.Expect(cond.Status).To(gomega.Equal(metav1.ConditionTrue), "BackendConfigPolicy should be accepted")
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGatewayExtensionAccepted validates that a GatewayExtension has "Accepted" status.
// GatewayExtension uses direct conditions on its status.
func (p *Provider) EventuallyGatewayExtensionAccepted(ctx context.Context, extension *v1alpha1.GatewayExtension) {
	currentTimeout, pollingInterval := helpers.GetTimeouts()
	p.Gomega.Eventually(func(g gomega.Gomega) {
		obj := &v1alpha1.GatewayExtension{}
		objKey := client.ObjectKeyFromObject(extension)
		err := p.clusterContext.Client.Get(ctx, objKey, obj)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get GatewayExtension %s", objKey)

		cond := meta.FindStatusCondition(obj.Status.Conditions, string(v1alpha1.PolicyConditionAccepted))
		g.Expect(cond).NotTo(gomega.BeNil(), "GatewayExtension should have Accepted condition")
		g.Expect(cond.Status).To(gomega.Equal(metav1.ConditionTrue), "GatewayExtension should be accepted")
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}
