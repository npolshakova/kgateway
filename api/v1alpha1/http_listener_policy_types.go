package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=httplistenerpolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=httplistenerpolicies/status,verbs=get;update;patch

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=gateway,app.kubernetes.io/name=gateway}
// +kubebuilder:resource:categories=kgateway,shortName=hlp
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="gateway.networking.k8s.io/policy=Direct"
type HTTPListenerPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HTTPListenerPolicySpec `json:"spec,omitempty"`
	Status PolicyStatus           `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type HTTPListenerPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPListenerPolicy `json:"items"`
}

type HTTPListenerPolicySpec struct {
	TargetRef LocalPolicyTargetReference `json:"targetRef,omitempty"`
	Compress  bool                       `json:"compress,omitempty"`

	// AccessLoggingConfig contains various settings for Envoy's access logging service.
	// See here for more information: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/accesslog/v3/accesslog.proto
	// +kubebuilder:validation:Items={type=object}
	AccessLog []AccessLog `json:"accessLog,omitempty"`
}

// AccessLog represents the top-level access log configuration.
type AccessLog struct {
	// Output access logs to local file
	FileSink *FileSink `json:"fileSink,omitempty"`

	// Send access logs to gRPC service
	GrpcService *GrpcService `json:"grpcService,omitempty"`

	// Filter access logs configuration
	Filter *AccessLogFilter `json:"filter,omitempty"`
}

// FileSink represents the file sink configuration for access logs.
// +kubebuilder:validation:XValidation:message="only one of 'StringFormat' or 'JsonFormat' may be set",rule="(has(self.stringFormat) && !has(self.jsonFormat)) || (!has(self.stringFormat) && has(self.jsonFormat))"
type FileSink struct {
	// the file path to which the file access logging service will sink
	// +kubebuilder:validation:Required
	Path string `json:"path"`
	// the format string by which envoy will format the log lines
	// https://www.envoyproxy.io/docs/envoy/v1.14.1/configuration/observability/access_log#config-access-log-format-strings
	StringFormat string `json:"stringFormat,omitempty"`
	// the format object by which to envoy will emit the logs in a structured way.
	// https://www.envoyproxy.io/docs/envoy/v1.14.1/configuration/observability/access_log#format-dictionaries
	JsonFormat *runtime.RawExtension `json:"jsonFormat,omitempty"`
}

// GrpcService represents the gRPC service configuration for access logs.
type GrpcService struct {
	// name of log stream
	// +kubebuilder:validation:Required
	LogName string `json:"logName"`

	// The backend gRPC service. Can be any type of supported backed (Kubernetes Service, kgateway Upstream, etc..)
	// +kubebuilder:validation:Required
	BackendRef *gwv1.BackendRef `json:"backendRef"`

	// Additional request headers to log in the access log
	AdditionalRequestHeadersToLog []string `json:"additionalRequestHeadersToLog,omitempty"`

	// Additional response headers to log in the access log
	AdditionalResponseHeadersToLog []string `json:"additionalResponseHeadersToLog,omitempty"`

	// Additional response trailers to log in the access log
	AdditionalResponseTrailersToLog []string `json:"additionalResponseTrailersToLog,omitempty"`
}

// AccessLogFilter represents the top-level filter structure.
type AccessLogFilter struct {
	*FilterType `json:",inline"` // embedded to allow for validation
	// +kube:validation:MinItems=2
	AndFilter []FilterType `json:"andFilter,omitempty"`
	// +kube:validation:MinItems=2
	OrFilter []FilterType `json:"orFilter,omitempty"`
}

// FilterType represents the type of filter to apply (only one of these should be set).
type FilterType struct {
	StatusCodeFilter     *StatusCodeFilter   `json:"statusCodeFilter,omitempty"`
	DurationFilter       *DurationFilter     `json:"durationFilter,omitempty"`
	NotHealthCheckFilter bool                `json:"notHealthCheckFilter,omitempty"`
	TraceableFilter      bool                `json:"traceableFilter,omitempty"`
	HeaderFilter         *HeaderFilter       `json:"headerFilter,omitempty"`
	ResponseFlagFilter   *ResponseFlagFilter `json:"responseFlagFilter,omitempty"`
	GrpcStatusFilter     *GrpcStatusFilter   `json:"grpcStatusFilter,omitempty"`
	CELFilter            *CELFilter          `json:"celFilter,omitempty"`
}

// ComparisonFilter represents a filter based on a comparison.
type ComparisonFilter struct {
	// +kubebuilder:validation:Required
	Op Op `json:"op,omitempty"`

	// Value to compare against.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4294967295
	Value uint32 `json:"value,omitempty"`
}

// Op represents comparison operators.
// +kubebuilder:validation:Enum=EQ;GE;LE
type Op string

const (
	EQ Op = "EQ" // Equal
	GE Op = "GQ" // Greater or equal
	LE Op = "LE" // Less or equal
)

// StatusCodeFilter filters based on HTTP status code.
type StatusCodeFilter ComparisonFilter

// DurationFilter filters based on request duration.
type DurationFilter ComparisonFilter

// DenominatorType defines the fraction percentages support several fixed denominator values.
// +kubebuilder:validation:enum=HUNDRED,TEN_THOUSAND,MILLION
type DenominatorType string

const (
	// 100.
	//
	// **Example**: 1/100 = 1%.
	HUNDRED DenominatorType = "HUNDRED"
	// 10,000.
	//
	// **Example**: 1/10000 = 0.01%.
	TEN_THOUSAND DenominatorType = "TEN_THOUSAND"
	// 1,000,000.
	//
	// **Example**: 1/1000000 = 0.0001%.
	MILLION DenominatorType = "MILLION"
)

// HeaderFilter filters requests based on headers.
type HeaderFilter struct {
	// +kubebuilder:validation:Required
	Header gwv1.HTTPHeaderMatch `json:"header"`
}

// ResponseFlagFilter filters based on response flags.
type ResponseFlagFilter struct {
	// +kubebuilder:validation:MinItems=1
	Flags []string `json:"flags"`
}

// CELFilter filters requests based on Common Expression Language (CEL).
type CELFilter struct {
	// The CEL expressions to evaluate. AccessLogs are only emitted when the CEL expressions evaluates to true.
	// see: https://www.envoyproxy.io/docs/envoy/latest/xds/type/v3/cel.proto.html#common-expression-language-cel-proto
	Match string `json:"match"`
}

// GrpcStatusFilter filters gRPC requests based on their response status.
type GrpcStatusFilter struct {
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Items={type=object}
	Statuses []GrpcStatus `json:"statuses,omitempty"`
	Exclude  bool         `json:"exclude,omitempty"`
}

// GrpcStatus represents possible gRPC statuses.
// +kubebuilder:validation:Enum=OK;CANCELED;UNKNOWN;INVALID_ARGUMENT;DEADLINE_EXCEEDED;NOT_FOUND;ALREADY_EXISTS;PERMISSION_DENIED;RESOURCE_EXHAUSTED;FAILED_PRECONDITION;ABORTED;OUT_OF_RANGE;UNIMPLEMENTED;INTERNAL;UNAVAILABLE;DATA_LOSS;UNAUTHENTICATED
type GrpcStatus string

const (
	OK                  GrpcStatus = "OK"
	CANCELED            GrpcStatus = "CANCELED"
	UNKNOWN             GrpcStatus = "UNKNOWN"
	INVALID_ARGUMENT    GrpcStatus = "INVALID_ARGUMENT"
	DEADLINE_EXCEEDED   GrpcStatus = "DEADLINE_EXCEEDED"
	NOT_FOUND           GrpcStatus = "NOT_FOUND"
	ALREADY_EXISTS      GrpcStatus = "ALREADY_EXISTS"
	PERMISSION_DENIED   GrpcStatus = "PERMISSION_DENIED"
	RESOURCE_EXHAUSTED  GrpcStatus = "RESOURCE_EXHAUSTED"
	FAILED_PRECONDITION GrpcStatus = "FAILED_PRECONDITION"
	ABORTED             GrpcStatus = "ABORTED"
	OUT_OF_RANGE        GrpcStatus = "OUT_OF_RANGE"
	UNIMPLEMENTED       GrpcStatus = "UNIMPLEMENTED"
	INTERNAL            GrpcStatus = "INTERNAL"
	UNAVAILABLE         GrpcStatus = "UNAVAILABLE"
	DATA_LOSS           GrpcStatus = "DATA_LOSS"
	UNAUTHENTICATED     GrpcStatus = "UNAUTHENTICATED"
)
