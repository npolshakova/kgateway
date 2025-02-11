package routepolicy

import (
	"context"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/kgateway-dev/kgateway/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/internal/gateway2/ir"
)

const (
	aiMetadataNamespace        = "envoy.filters.ai.solo.io"
	transformationFilterName   = "ai.transformation.solo.io"
	transformationEEFilterName = "ai.transformation_ee.solo.io"
	extProcFilterName          = "ai.extproc.solo.io"
	modSecurityFilterName      = "ai.modsecurity.solo.io"
	setMetadataFilterName      = "envoy.filters.http.set_filter_state"
	extProcUDSClusterName      = "ai_ext_proc_uds_cluster"
	extProcUDSSocketPath       = "@gloo-ai-sock"
)

func processAIRoutePolicy(ctx context.Context, aiConfig *v1alpha1.AIRoutePolicy, outputRoute *envoy_config_route_v3.Route, pCtx *ir.RouteContext) error {
	if outputRoute.GetTypedPerFilterConfig() == nil {
		outputRoute.TypedPerFilterConfig = make(map[string]*anypb.Any)
	}

	//// TODO:
	//// If it's not an AI route we want to disable our ext-proc filter just in case.
	//// This will have no effect if we don't add the listener filter
	//var aiUpstreams []*v1alpha1.AIUpstream
	////backends := pCtx.In.Backends
	////for _, backend := range backends {
	////	aiUpstreams = append(aiUpstreams, backend.Backend.Upstream)
	////}
	//
	//var llmModel string
	//byType := map[string]struct{}{}
	//for _, us := range aiUpstreams {
	//	if us.MultiPool != nil {
	//		llmModel = us.MultiPool.GetPriorities()[0].GetPool()[0].GetEndpoint()
	//		for _, pool := range us.GetAi().GetMulti().GetPriorities() {
	//			for _, ep := range pool.GetPool() {
	//				switch ep.GetLlm().(type) {
	//				case *ai.UpstreamSpec_MultiPool_Backend_Openai:
	//					byType["openai"] = struct{}{}
	//					llmModel = ep.GetOpenai().GetModel()
	//				case *ai.UpstreamSpec_MultiPool_Backend_Mistral:
	//					byType["mistral"] = struct{}{}
	//					llmModel = ep.GetMistral().GetModel()
	//				case *ai.UpstreamSpec_MultiPool_Backend_Anthropic:
	//					byType["anthropic"] = struct{}{}
	//					llmModel = ep.GetAnthropic().GetModel()
	//				case *ai.UpstreamSpec_MultiPool_Backend_AzureOpenai:
	//					byType["azure_openai"] = struct{}{}
	//					llmModel = ep.GetAzureOpenai().GetDeploymentName()
	//				case *ai.UpstreamSpec_MultiPool_Backend_Gemini:
	//					byType["gemini"] = struct{}{}
	//					llmModel = ep.GetGemini().GetModel()
	//				case *ai.UpstreamSpec_MultiPool_Backend_VertexAi:
	//					byType["vertex-ai"] = struct{}{}
	//					llmModel = ep.GetVertexAi().GetModel()
	//				}
	//			}
	//		}
	//		break
	//	} else if us.LLM != nil {
	//		if us.LLM.OpenAI != nil {
	//			byType["openai"] = struct{}{}
	//			llmModel = us.GetAi().GetOpenai().GetModel()
	//		} else if us.LLM.Mistral != nil {
	//			byType["mistral"] = struct{}{}
	//			llmModel = us.GetAi().GetMistral().GetModel()
	//		} else if us.LLM.Anthropic != nil {
	//			byType["anthropic"] = struct{}{}
	//			llmModel = us.GetAi().GetAnthropic().GetModel()
	//		} else if us.LLM.AzureOpenAI != nil {
	//			byType["azure_openai"] = struct{}{}
	//			llmModel = us.GetAi().GetAzureOpenai().GetDeploymentName()
	//		} else if us.LLM.Gemini != nil {
	//			byType["gemini"] = struct{}{}
	//			llmModel = us.GetAi().GetGemini().GetModel()
	//		} else if us.LLM.VertexAI != nil {
	//			byType["vertex-ai"] = struct{}{}
	//			llmModel = us.GetAi().GetVertexAi().GetModel()
	//		} else {
	//			// TODO: handle error
	//		}
	//	}
	//	break
	//}
	//
	//if len(byType) != 1 {
	//	return eris.Errorf("multiple AI backend types found for single ai route %+v", byType)
	//}
	//
	//// This is only len(1)
	//var llmProvider string
	//for k := range byType {
	//	llmProvider = k
	//}
	//
	//// Setup ext-proc route filter config, we will conditionally modify it based on certain route options.
	//// A heavily used part of this config is the `GrpcInitialMetadata`.
	//// This is used to add headers to the ext-proc request.
	//// These headers are used to configure the AI server on a per-request basis.
	//// This was the best available way to pass per-route configuration to the AI server.
	//extProcRouteSettings := &envoy_ext_proc_v3.ExtProcPerRoute{
	//	Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
	//		Overrides: &envoy_ext_proc_v3.ExtProcOverrides{},
	//	},
	//}
	//
	//// Add things which require basic AI upstream.
	//outputRoute.GetRoute().HostRewriteSpecifier = &envoy_config_route_v3.RouteAction_AutoHostRewrite{
	//	AutoHostRewrite: wrapperspb.Bool(true),
	//}
	//
	//// TODO:
	//// We only want to add the transformation filter if we have a single AI backend
	//// Otherwise we already have the transformation filter added by the weighted destination
	//
	//extProcRouteSettings.GetOverrides().GrpcInitialMetadata = append(extProcRouteSettings.GetOverrides().GrpcInitialMetadata,
	//	&envoy_config_core_v3.HeaderValue{
	//		Key:   "x-llm-provider",
	//		Value: llmProvider,
	//	},
	//)
	//// If the Upstream specifies a model, add a header to the ext-proc request
	//if llmModel != "" {
	//	extProcRouteSettings.GetOverrides().GrpcInitialMetadata = append(extProcRouteSettings.GetOverrides().GrpcInitialMetadata,
	//		&envoy_config_core_v3.HeaderValue{
	//			Key:   "x-llm-model",
	//			Value: llmModel,
	//		})
	//}
	//// If the route options specify this as a chat streaming route, add a header to the ext-proc request
	//if aiConfig.RouteType == v1alpha1.CHAT_STREAMING {
	//	// append streaming header if it's a streaming route
	//	extProcRouteSettings.GetOverrides().GrpcInitialMetadata = append(extProcRouteSettings.GetOverrides().GrpcInitialMetadata, &envoy_config_core_v3.HeaderValue{
	//		Key:   "x-chat-streaming",
	//		Value: "true",
	//	})
	//}
	//
	//// Add the x-request-id header to the ext-proc request.
	//// This is an optimization to allow us to not have to wait for the headers request to
	//// Initialize our logger/handler classes.
	//extProcRouteSettings.GetOverrides().GrpcInitialMetadata = append(extProcRouteSettings.GetOverrides().GrpcInitialMetadata,
	//	&envoy_config_core_v3.HeaderValue{
	//		Key:   "x-request-id",
	//		Value: "%REQ(X-REQUEST-ID)%",
	//	},
	//)
	//
	//if err := handleAiRouteOptions(in, p.transformationsByRoute[in], extProcRouteSettings); err != nil {
	//	return err
	//}
	//
	//// Setup transformation routeConfig/output type
	//for _, val := range transformationsByRoute[in] {
	//	transformation := &envoytransformation.RouteTransformations_RouteTransformation{
	//		Match: &envoytransformation.RouteTransformations_RouteTransformation_RequestMatch_{
	//			RequestMatch: &envoytransformation.RouteTransformations_RouteTransformation_RequestMatch{
	//				RequestTransformation: &envoytransformation.Transformation{
	//					// Set this env var to true to log the request/response info for each transformation
	//					LogRequestResponseInfo: wrapperspb.Bool(os.Getenv("AI_PLUGIN_DEBUG_TRANSFORMATIONS") == "true"),
	//					TransformationType: &envoytransformation.Transformation_TransformationTemplate{
	//						TransformationTemplate: val.transformation,
	//					},
	//				},
	//			},
	//		},
	//	}
	//	marshaled, err := utils.MessageToAny(&envoytransformation.RouteTransformations{
	//		Transformations: []*envoytransformation.RouteTransformations_RouteTransformation{transformation},
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	val.perFilterConfig[transformationFilterName] = marshaled
	//}
	//
	//marshaled, err := utils.MessageToAny(extProcRouteSettings)
	//if err != nil {
	//	return err
	//}
	//outputRoute.TypedPerFilterConfig[extProcFilterName] = marshaled
	//
	//listenerExtProcFilter[params.HttpListener] = struct{}{}
	return nil

}
