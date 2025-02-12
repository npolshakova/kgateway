// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"

	apiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/applyconfiguration/api/v1alpha1"
	internal "github.com/kgateway-dev/kgateway/v2/api/applyconfiguration/internal"
	v1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=gateway.kgateway.dev, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("AccessLog"):
		return &apiv1alpha1.AccessLogApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AccessLogFilter"):
		return &apiv1alpha1.AccessLogFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AiExtension"):
		return &apiv1alpha1.AiExtensionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AiExtensionStats"):
		return &apiv1alpha1.AiExtensionStatsApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AwsUpstream"):
		return &apiv1alpha1.AwsUpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CELFilter"):
		return &apiv1alpha1.CELFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CustomLabel"):
		return &apiv1alpha1.CustomLabelApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DirectResponse"):
		return &apiv1alpha1.DirectResponseApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DirectResponseSpec"):
		return &apiv1alpha1.DirectResponseSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DurationFilter"):
		return &apiv1alpha1.DurationFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EnvoyBootstrap"):
		return &apiv1alpha1.EnvoyBootstrapApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EnvoyContainer"):
		return &apiv1alpha1.EnvoyContainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FileSink"):
		return &apiv1alpha1.FileSinkApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FilterType"):
		return &apiv1alpha1.FilterTypeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FractionalPercent"):
		return &apiv1alpha1.FractionalPercentApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GatewayParameters"):
		return &apiv1alpha1.GatewayParametersApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GatewayParametersSpec"):
		return &apiv1alpha1.GatewayParametersSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GracefulShutdownSpec"):
		return &apiv1alpha1.GracefulShutdownSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GrpcService"):
		return &apiv1alpha1.GrpcServiceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("GrpcStatusFilter"):
		return &apiv1alpha1.GrpcStatusFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("HeaderFilter"):
		return &apiv1alpha1.HeaderFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Host"):
		return &apiv1alpha1.HostApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("HTTPListenerPolicy"):
		return &apiv1alpha1.HTTPListenerPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("HTTPListenerPolicySpec"):
		return &apiv1alpha1.HTTPListenerPolicySpecApplyConfiguration{}
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
	case v1alpha1.SchemeGroupVersion.WithKind("LocalPolicyTargetReference"):
		return &apiv1alpha1.LocalPolicyTargetReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Pod"):
		return &apiv1alpha1.PodApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PolicyAncestorStatus"):
		return &apiv1alpha1.PolicyAncestorStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PolicyStatus"):
		return &apiv1alpha1.PolicyStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ProxyDeployment"):
		return &apiv1alpha1.ProxyDeploymentApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ResponseFlagFilter"):
		return &apiv1alpha1.ResponseFlagFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RoutePolicy"):
		return &apiv1alpha1.RoutePolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RoutePolicySpec"):
		return &apiv1alpha1.RoutePolicySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RuntimeFilter"):
		return &apiv1alpha1.RuntimeFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SdsBootstrap"):
		return &apiv1alpha1.SdsBootstrapApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SdsContainer"):
		return &apiv1alpha1.SdsContainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Service"):
		return &apiv1alpha1.ServiceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ServiceAccount"):
		return &apiv1alpha1.ServiceAccountApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("StaticUpstream"):
		return &apiv1alpha1.StaticUpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("StatsConfig"):
		return &apiv1alpha1.StatsConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("StatusCodeFilter"):
		return &apiv1alpha1.StatusCodeFilterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Upstream"):
		return &apiv1alpha1.UpstreamApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("UpstreamSpec"):
		return &apiv1alpha1.UpstreamSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("UpstreamStatus"):
		return &apiv1alpha1.UpstreamStatusApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
