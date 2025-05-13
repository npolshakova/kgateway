package agentgatewaysyncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/a2a"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer/mcp"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type AgentGwTranslator struct {
	commonCols *common.CommonCollections

	waitForSync []cache.InformerSynced

	logger *zap.Logger
}

func NewTranslator(
	ctx context.Context,
	commonCols *common.CommonCollections,
) *AgentGwTranslator {
	return &AgentGwTranslator{
		commonCols:  commonCols,
		logger:      contextutils.LoggerFrom(ctx).Desugar().With(zap.String("component", "agentgateway_translator")),
		waitForSync: []cache.InformerSynced{},
	}
}

func (t *AgentGwTranslator) HasSynced() bool {
	for _, sync := range t.waitForSync {
		if !sync() {
			return false
		}
	}
	return true
}

type AgentGatewayTranslationResult struct {
	Listeners  []*Listener
	A2ATargets []*a2a.Target
	McpTargets []*mcp.Target
}

// ctx needed for logging; remove once we refactor logging.
func (t *AgentGwTranslator) TranslateGateway(kctx krt.HandlerContext, ctx context.Context, gw ir.Gateway) (*AgentGatewayTranslationResult, reports.ReportMap) {
	logger := contextutils.LoggerFrom(ctx)

	rm := reports.NewReportMap()
	r := reports.NewReporter(&rm)
	logger.Debugf("building agentgateway proxy for kube gw %s version %s", client.ObjectKeyFromObject(gw.Obj), gw.Obj.GetResourceVersion())

	result := &AgentGatewayTranslationResult{
		Listeners:  make([]*Listener, 0),
		A2ATargets: make([]*a2a.Target, 0),
		McpTargets: make([]*mcp.Target, 0),
	}

	// translate each listener in the gateway
	listErrs := make([]error, 0)
	for _, l := range gw.Listeners {
		// TODO: propagate errors so we can allow the retain last config mode
		l, err := t.ComputeListener(l)
		if err != nil {
			listErrs = append(listErrs, err)
			continue
		}
		result.Listeners = append(result.Listeners, l)
	}

	agentGwServices := krt.NewManyCollection(t.commonCols.Services, func(kctx krt.HandlerContext, s *corev1.Service) []agentGwService {
		var agentGwServices []agentGwService
		var allowedListeners []string
		for _, listener := range gw.Listeners {
			if listener.Protocol != A2AProtocol && listener.Protocol != MCPProtocol {
				continue
			}
			logger.Debugf("Found agent gateway service %s/%s", s.Namespace, s.Name)
			if listener.AllowedRoutes == nil {
				// only allow agent services in same namespace
				if s.Namespace == gw.Obj.Namespace {
					allowedListeners = append(allowedListeners, string(listener.Name))
				}
			} else if listener.AllowedRoutes.Namespaces.From != nil {
				switch *listener.AllowedRoutes.Namespaces.From {
				case gwv1.NamespacesFromAll:
					allowedListeners = append(allowedListeners, string(listener.Name))
				case gwv1.NamespacesFromSame:
					// only allow agent services in same namespace
					if s.Namespace == gw.Obj.Namespace {
						allowedListeners = append(allowedListeners, string(listener.Name))
					}
				case gwv1.NamespacesFromSelector:
					// TODO: implement namespace selectors
					contextutils.LoggerFrom(ctx).Errorf("namespace selectors not supported for agent gateways")
					continue
				}
			}
		}
		agentGwServices = translateAgentService(s, allowedListeners)
		for _, agentSvc := range agentGwServices {
			if agentSvc.protocol == A2AProtocol {
				a2aTarget := &a2a.Target{
					Name: getTargetName(agentSvc.ResourceName()),
					Host: agentSvc.ip,
					Port: uint32(agentSvc.port),
					Path: agentSvc.path,
				}
				result.A2ATargets = append(result.A2ATargets, a2aTarget)
			} else if agentSvc.protocol == MCPProtocol {
				mcpTarget := &mcp.Target{
					// Note: No slashes allowed here (must match ^[a-zA-Z0-9-]+$)
					Name: getTargetName(agentSvc.ResourceName()),
					Target: &mcp.Target_Sse{
						Sse: &mcp.Target_SseTarget{
							Host: agentSvc.ip,
							Port: uint32(agentSvc.port),
							Path: agentSvc.path,
						},
					},
				}
				result.McpTargets = append(result.McpTargets, mcpTarget)
			}
		}
		return agentGwServices
	})
	t.waitForSync = append(t.waitForSync, agentGwServices.HasSynced)

	if len(listErrs) > 0 {
		err := errors.Join(listErrs...)
		r.Gateway(gw.Obj).SetCondition(reporter.GatewayCondition{
			Type:    gwv1.GatewayConditionProgrammed,
			Reason:  gwv1.GatewayReasonInvalid,
			Status:  metav1.ConditionFalse,
			Message: "Error processing listeners: " + err.Error(),
		})
	} else {
		r.Gateway(gw.Obj).SetCondition(reporter.GatewayCondition{
			Type:    gwv1.GatewayConditionProgrammed,
			Status:  metav1.ConditionFalse,
			Reason:  gwv1.GatewayReasonInvalid,
			Message: "Translated with errors",
		})
	}

	// we are recomputing xds snapshots as proxies have changed, signal that we need to sync xds with these new snapshots
	return result, rm
}

func (t *AgentGwTranslator) ComputeListener(
	lis ir.Listener,
) (*Listener, error) {
	var listener *Listener

	if string(lis.Protocol) == A2AProtocol {
		listener = &Listener{
			Name:     string(lis.Name),
			Protocol: Listener_A2A,
			// TODO: Add suppose for stdio listener
			Listener: &Listener_Sse{
				Sse: &SseListener{
					Address: "[::]",
					Port:    uint32(lis.Port),
				},
			},
		}
	} else if string(lis.Protocol) == MCPProtocol {
		listener = &Listener{
			Name:     string(lis.Name),
			Protocol: Listener_MCP,
			// TODO: Add suppose for stdio listener
			Listener: &Listener_Sse{
				Sse: &SseListener{
					Address: "[::]",
					Port:    uint32(lis.Port),
				},
			},
		}
	} else {
		// Not a valid protocol for Agent Gateway
		return nil, fmt.Errorf("unsupported protocol %s for listener %s", lis.Protocol, lis.Name)
	}

	// TODO: add TLS configuration to listener

	return listener, nil
}
