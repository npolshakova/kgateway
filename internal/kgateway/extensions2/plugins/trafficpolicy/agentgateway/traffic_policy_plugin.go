package agentgateway

import (
	"time"

	"github.com/agentgateway/agentgateway/go/api"

	agwir "github.com/kgateway-dev/kgateway/v2/pkg/agentgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	ir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
)

var (
	logger = logging.New("agentgateway/plugin/trafficpolicy")
)

// PolicySubIR documents the expected interface that all policy sub-IRs should implement.
type PolicySubIR interface {
	// Equals compares this policy with another policy
	Equals(other PolicySubIR) bool

	// Validate performs PGV validation on the policy
	Validate() error

	// TODO: Merge. Just awkward as we won't be using the actual method type.
}

type TrafficPolicy struct {
	ct   time.Time
	Spec trafficPolicySpecIr
}

type trafficPolicySpecIr struct {
	ExtAuth *extAuthIR
}

func (d *TrafficPolicy) CreationTime() time.Time {
	return d.ct
}

func (d *TrafficPolicy) Equals(in any) bool {
	d2, ok := in.(*TrafficPolicy)
	if !ok {
		return false
	}
	if d.ct != d2.ct {
		return false
	}

	if !d.Spec.ExtAuth.Equals(d2.Spec.ExtAuth) {
		return false
	}
	return true
}

// Validate performs PGV (protobuf-generated validation) validation by delegating
// to each policy sub-IR's Validate() method. This follows the exact same pattern as the Equals() method.
// PGV validation is always performed regardless of route replacement mode.
func (p *TrafficPolicy) Validate() error {
	var validators []func() error
	validators = append(validators, p.Spec.ExtAuth.Validate)
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

type trafficPolicyPluginGwPass struct {
	reporter reports.Reporter
	agwir.UnimplementedAgentGatewayTranslationPass

	extAuthPerProvider ProviderNeededMap
}

var _ agwir.AgentGatewayTranslationPass = &trafficPolicyPluginGwPass{}

func NewAgentGatewayPass(reporter reports.Reporter) agwir.AgentGatewayTranslationPass {
	return &trafficPolicyPluginGwPass{
		reporter: reporter,
	}
}

func (p *TrafficPolicy) Name() string {
	return "trafficpolicies"
}

// ApplyForRoute processes route-level configuration for agent gateway
func (p *trafficPolicyPluginGwPass) ApplyForRoute(pCtx *agwir.AgentGatewayRouteContext, out *api.Route) error {
	// TODO: Implement route-level policy application for agent gateway
	return nil
}

// ApplyForBackend processes backend-level configuration for each backend referenced in routes
func (p *trafficPolicyPluginGwPass) ApplyForBackend(pCtx *agwir.AgentGatewayTranslationBackendContext, out *api.Backend) error {
	// TODO: Implement backend-level policy application for agent gateway
	return nil
}

// ApplyForRouteBackend processes route-specific backend configuration
func (p *trafficPolicyPluginGwPass) ApplyForRouteBackend(policy ir.PolicyIR, pCtx *agwir.AgentGatewayTranslationBackendContext) error {
	// TODO: Implement route-specific backend policy application for agent gateway
	return nil
}

func (p *trafficPolicyPluginGwPass) handlePolicies(fcn string, spec trafficPolicySpecIr) {
	p.handleExtAuth(fcn, spec.ExtAuth)
}

func (p *trafficPolicyPluginGwPass) SupportsPolicyMerge() bool {
	return true
}

// MergeTrafficPolicies merges two TrafficPolicy IRs, returning a map that contains information
// about the origin policy reference for each merged field.
func MergeTrafficPolicies(
	p1, p2 *TrafficPolicy,
	p2Ref *ir.AttachedPolicyRef,
	p2MergeOrigins ir.MergeOrigins,
	mergeOpts policy.MergeOptions,
	mergeOrigins ir.MergeOrigins,
) {
	if p1 == nil || p2 == nil {
		return
	}

	mergeFuncs := []func(*TrafficPolicy, *TrafficPolicy, *ir.AttachedPolicyRef, ir.MergeOrigins, policy.MergeOptions, ir.MergeOrigins){
		mergeExtAuth,
		mergeExtProc,
	}

	for _, mergeFunc := range mergeFuncs {
		mergeFunc(p1, p2, p2Ref, p2MergeOrigins, mergeOpts, mergeOrigins)
	}
}

// mergeExtAuth merges ExtAuth policies
func mergeExtAuth(
	p1, p2 *TrafficPolicy,
	p2Ref *ir.AttachedPolicyRef,
	p2MergeOrigins ir.MergeOrigins,
	mergeOpts policy.MergeOptions,
	mergeOrigins ir.MergeOrigins,
) {
	// TODO: Implement ExtAuth merging
}

// mergeExtProc merges ExtProc policies
func mergeExtProc(
	p1, p2 *TrafficPolicy,
	p2Ref *ir.AttachedPolicyRef,
	p2MergeOrigins ir.MergeOrigins,
	mergeOpts policy.MergeOptions,
	mergeOrigins ir.MergeOrigins,
) {
	// TODO: Implement ExtProc merging
}
