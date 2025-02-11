// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"

	apiv1alpha1 "github.com/kgateway-dev/kgateway/api/applyconfiguration/api/v1alpha1"
	internal "github.com/kgateway-dev/kgateway/api/applyconfiguration/internal"
	v1alpha1 "github.com/kgateway-dev/kgateway/api/v1alpha1"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=gateway.kgateway.dev, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("AiExtension"):
		return &apiv1alpha1.AiExtensionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AiExtensionStats"):
		return &apiv1alpha1.AiExtensionStatsApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AIPromptEnrichment"):
		return &apiv1alpha1.AIPromptEnrichmentApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AIPromptGuard"):
		return &apiv1alpha1.AIPromptGuardApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AIRoutePolicy"):
		return &apiv1alpha1.AIRoutePolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AIUpstream"):
		return &apiv1alpha1.AIUpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AnthropicConfig"):
		return &apiv1alpha1.AnthropicConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AwsUpstream"):
		return &apiv1alpha1.AwsUpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AzureOpenAIConfig"):
		return &apiv1alpha1.AzureOpenAIConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CustomLabel"):
		return &apiv1alpha1.CustomLabelApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CustomResponse"):
		return &apiv1alpha1.CustomResponseApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DirectResponse"):
		return &apiv1alpha1.DirectResponseApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DirectResponseSpec"):
		return &apiv1alpha1.DirectResponseSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EnvoyBootstrap"):
		return &apiv1alpha1.EnvoyBootstrapApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EnvoyContainer"):
		return &apiv1alpha1.EnvoyContainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FieldDefault"):
		return &apiv1alpha1.FieldDefaultApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GatewayParameters"):
		return &apiv1alpha1.GatewayParametersApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GatewayParametersSpec"):
		return &apiv1alpha1.GatewayParametersSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GeminiConfig"):
		return &apiv1alpha1.GeminiConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GracefulShutdownSpec"):
		return &apiv1alpha1.GracefulShutdownSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Host"):
		return &apiv1alpha1.HostApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("HttpListenerPolicy"):
		return &apiv1alpha1.HttpListenerPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("HttpListenerPolicySpec"):
		return &apiv1alpha1.HttpListenerPolicySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Image"):
		return &apiv1alpha1.ImageApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("IstioContainer"):
		return &apiv1alpha1.IstioContainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("IstioIntegration"):
		return &apiv1alpha1.IstioIntegrationApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("KubernetesProxyConfig"):
		return &apiv1alpha1.KubernetesProxyConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ListenerPolicy"):
		return &apiv1alpha1.ListenerPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ListenerPolicySpec"):
		return &apiv1alpha1.ListenerPolicySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LLMProviders"):
		return &apiv1alpha1.LLMProvidersApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LocalPolicyTargetReference"):
		return &apiv1alpha1.LocalPolicyTargetReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Message"):
		return &apiv1alpha1.MessageApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("MistralConfig"):
		return &apiv1alpha1.MistralConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Moderation"):
		return &apiv1alpha1.ModerationApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("MultiPoolConfig"):
		return &apiv1alpha1.MultiPoolConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("OpenAIConfig"):
		return &apiv1alpha1.OpenAIConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("OpenAIModeration"):
		return &apiv1alpha1.OpenAIModerationApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Pod"):
		return &apiv1alpha1.PodApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PolicyAncestorStatus"):
		return &apiv1alpha1.PolicyAncestorStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PolicyStatus"):
		return &apiv1alpha1.PolicyStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Priority"):
		return &apiv1alpha1.PriorityApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PromptguardRequest"):
		return &apiv1alpha1.PromptguardRequestApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PromptguardResponse"):
		return &apiv1alpha1.PromptguardResponseApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ProxyDeployment"):
		return &apiv1alpha1.ProxyDeploymentApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Regex"):
		return &apiv1alpha1.RegexApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RegexMatch"):
		return &apiv1alpha1.RegexMatchApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RoutePolicy"):
		return &apiv1alpha1.RoutePolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RoutePolicySpec"):
		return &apiv1alpha1.RoutePolicySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SdsBootstrap"):
		return &apiv1alpha1.SdsBootstrapApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SdsContainer"):
		return &apiv1alpha1.SdsContainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Service"):
		return &apiv1alpha1.ServiceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ServiceAccount"):
		return &apiv1alpha1.ServiceAccountApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SingleAuthToken"):
		return &apiv1alpha1.SingleAuthTokenApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("StaticUpstream"):
		return &apiv1alpha1.StaticUpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("StatsConfig"):
		return &apiv1alpha1.StatsConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Upstream"):
		return &apiv1alpha1.UpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("UpstreamSpec"):
		return &apiv1alpha1.UpstreamSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("UpstreamStatus"):
		return &apiv1alpha1.UpstreamStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("VertexAIConfig"):
		return &apiv1alpha1.VertexAIConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Webhook"):
		return &apiv1alpha1.WebhookApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
