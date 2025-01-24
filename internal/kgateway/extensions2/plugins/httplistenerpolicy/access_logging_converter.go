package httplistenerpolicy

import (
	"context"
	"encoding/json"
	"fmt"

	envoyaccesslog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyalfile "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	cel "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/filters/cel/v3"
	envoygrpc "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/grpc/v3"
	envoy_metadata_formatter "github.com/envoyproxy/go-control-plane/envoy/extensions/formatter/metadata/v3"
	envoy_req_without_query "github.com/envoyproxy/go-control-plane/envoy/extensions/formatter/req_without_query/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	envoytype "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/runtime"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	kgateway_wellknown "github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

// convertAccessLogConfig transforms a list of AccessLog configurations into Envoy AccessLog configurations
func convertAccessLogConfig(ctx context.Context, configs []v1alpha1.AccessLog) ([]*envoyaccesslog.AccessLog, error) {
	if len(configs) == 0 {
		return nil, nil
	}

	logger := contextutils.LoggerFrom(ctx).Desugar()
	var results []*envoyaccesslog.AccessLog

	for _, logConfig := range configs {
		accessLogCfg, err := createAccessLogConfiguration(logger, logConfig)
		if err != nil {
			return nil, err
		}
		results = append(results, accessLogCfg)
	}

	return results, nil
}

// createAccessLogConfiguration creates an Envoy AccessLog configuration for a single log config
func createAccessLogConfiguration(logger *zap.Logger, logConfig v1alpha1.AccessLog) (*envoyaccesslog.AccessLog, error) {
	// Validate mutual exclusivity of sink types
	if logConfig.FileSink != nil && logConfig.GrpcService != nil {
		return nil, eris.New("access log config cannot have both file sink and grpc service")
	}

	var (
		accessLogCfg *envoyaccesslog.AccessLog
		err          error
	)

	switch {
	case logConfig.FileSink != nil:
		accessLogCfg, err = createFileAccessLog(logger, logConfig.FileSink)
	case logConfig.GrpcService != nil:
		accessLogCfg, err = createGrpcAccessLog(logger, logConfig.GrpcService)
	default:
		return nil, eris.New("no access log sink specified")
	}

	if err != nil {
		return nil, err
	}

	// Add filter if specified
	if logConfig.Filter != nil {
		if err := addAccessLogFilter(logger, accessLogCfg, logConfig.Filter); err != nil {
			return nil, err
		}
	}

	return accessLogCfg, nil
}

// createFileAccessLog generates a file-based access log configuration
func createFileAccessLog(logger *zap.Logger, fileSink *v1alpha1.FileSink) (*envoyaccesslog.AccessLog, error) {
	fileCfg := &envoyalfile.FileAccessLog{Path: fileSink.Path}

	// Validate format configuration
	if fileSink.StringFormat != "" && fileSink.JsonFormat != nil {
		return nil, eris.New("access log config cannot have both string format and json format")
	}

	formatterExtensions, err := getFormatterExtensions()
	if err != nil {
		return nil, err
	}

	switch {
	case fileSink.StringFormat != "":
		fileCfg.AccessLogFormat = &envoyalfile.FileAccessLog_LogFormat{
			LogFormat: &envoycore.SubstitutionFormatString{
				Format: &envoycore.SubstitutionFormatString_TextFormatSource{
					TextFormatSource: &envoycore.DataSource{
						Specifier: &envoycore.DataSource_InlineString{
							InlineString: fileSink.StringFormat,
						},
					},
				},
				Formatters: formatterExtensions,
			},
		}
	case fileSink.JsonFormat != nil:
		fileCfg.AccessLogFormat = &envoyalfile.FileAccessLog_LogFormat{
			LogFormat: &envoycore.SubstitutionFormatString{
				Format: &envoycore.SubstitutionFormatString_JsonFormat{
					JsonFormat: convertJsonFormat(fileSink.JsonFormat),
				},
				Formatters: formatterExtensions,
			},
		}
	}

	return newAccessLogWithConfig(wellknown.FileAccessLog, fileCfg)
}

// createGrpcAccessLog generates a gRPC-based access log configuration
func createGrpcAccessLog(logger *zap.Logger, grpcService *v1alpha1.GrpcService) (*envoyaccesslog.AccessLog, error) {
	var cfg envoygrpc.HttpGrpcAccessLogConfig
	if err := copyGrpcSettings(&cfg, grpcService); err != nil {
		logger.Error(fmt.Sprintf("error converting grpc access log config: %s", err.Error()))
		return nil, err
	}

	return newAccessLogWithConfig(wellknown.HTTPGRPCAccessLog, &cfg)
}

// addAccessLogFilter adds filtering logic to an access log configuration
func addAccessLogFilter(logger *zap.Logger, accessLogCfg *envoyaccesslog.AccessLog, filter *v1alpha1.AccessLogFilter) error {
	var (
		filters []*envoyaccesslog.AccessLogFilter
		err     error
	)

	switch {
	case filter.OrFilter != nil:
		filters, err = translateOrFilters(logger, filter.OrFilter)
		if err != nil {
			return err
		}
		accessLogCfg.GetFilter().FilterSpecifier = &envoyaccesslog.AccessLogFilter_OrFilter{
			OrFilter: &envoyaccesslog.OrFilter{Filters: filters},
		}
	case filter.AndFilter != nil:
		filters, err = translateOrFilters(logger, filter.AndFilter)
		if err != nil {
			return err
		}
		accessLogCfg.GetFilter().FilterSpecifier = &envoyaccesslog.AccessLogFilter_AndFilter{
			AndFilter: &envoyaccesslog.AndFilter{Filters: filters},
		}
	case filter.FilterType != nil:
		accessLogCfg.Filter, err = translateFilter(logger, filter.FilterType)
		if err != nil {
			return err
		}
	}

	return nil
}

// translateOrFilters translates a slice of filter types
func translateOrFilters(logger *zap.Logger, filters []*v1alpha1.FilterType) ([]*envoyaccesslog.AccessLogFilter, error) {
	result := make([]*envoyaccesslog.AccessLogFilter, 0, len(filters))
	for _, filter := range filters {
		cfg, err := translateFilter(logger, filter)
		if err != nil {
			return nil, err
		}
		result = append(result, cfg)
	}
	return result, nil
}

func translateFilter(logger *zap.Logger, filter *v1alpha1.FilterType) (*envoyaccesslog.AccessLogFilter, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}

	var alCfg *envoyaccesslog.AccessLogFilter
	switch {
	case filter.StatusCodeFilter != nil:
		op, err := toEnvoyComparisonOpType(filter.StatusCodeFilter.Comparison.Op)
		if err != nil {
			return nil, err
		}

		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_StatusCodeFilter{
				StatusCodeFilter: &envoyaccesslog.StatusCodeFilter{
					Comparison: &envoyaccesslog.ComparisonFilter{
						Op: op,
						Value: &envoycore.RuntimeUInt32{
							DefaultValue: filter.StatusCodeFilter.Comparison.Value.DefaultValue,
							RuntimeKey:   filter.StatusCodeFilter.Comparison.Value.RuntimeKey,
						},
					},
				},
			},
		}

	case filter.DurationFilter != nil:
		op, err := toEnvoyComparisonOpType(filter.DurationFilter.Comparison.Op)
		if err != nil {
			return nil, err
		}

		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_DurationFilter{
				DurationFilter: &envoyaccesslog.DurationFilter{
					Comparison: &envoyaccesslog.ComparisonFilter{
						Op: op,
						Value: &envoycore.RuntimeUInt32{
							DefaultValue: filter.DurationFilter.Comparison.Value.DefaultValue,
							RuntimeKey:   filter.DurationFilter.Comparison.Value.RuntimeKey,
						},
					},
				},
			},
		}

	case filter.NotHealthCheckFilter != nil:
		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_NotHealthCheckFilter{
				NotHealthCheckFilter: &envoyaccesslog.NotHealthCheckFilter{},
			},
		}

	case filter.TraceableFilter != nil:
		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_TraceableFilter{
				TraceableFilter: &envoyaccesslog.TraceableFilter{},
			},
		}

	case filter.RuntimeFilter != nil:
		denominator, err := toEnvoyDenominatorType(filter.RuntimeFilter.PercentSampled.Denominator)
		if err != nil {
			return nil, err
		}
		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_RuntimeFilter{
				RuntimeFilter: &envoyaccesslog.RuntimeFilter{
					RuntimeKey: filter.RuntimeFilter.RuntimeKey,
					PercentSampled: &envoytype.FractionalPercent{
						Numerator:   filter.RuntimeFilter.PercentSampled.Numerator,
						Denominator: denominator,
					},
					UseIndependentRandomness: filter.RuntimeFilter.UseIndependentRandomness,
				},
			},
		}

	case filter.HeaderFilter != nil:
		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_HeaderFilter{
				HeaderFilter: &envoyaccesslog.HeaderFilter{
					Header: &envoyroute.HeaderMatcher{
						Name:                 string(filter.HeaderFilter.Header.Name),
						HeaderMatchSpecifier: createHeaderMatchSpecifier(logger, filter.HeaderFilter.Header),
					},
				},
			},
		}

	case filter.ResponseFlagFilter != nil:
		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_ResponseFlagFilter{
				ResponseFlagFilter: &envoyaccesslog.ResponseFlagFilter{
					Flags: filter.ResponseFlagFilter.Flags,
				},
			},
		}

	case filter.GrpcStatusFilter != nil:
		statuses := make([]envoyaccesslog.GrpcStatusFilter_Status, len(filter.GrpcStatusFilter.Statuses))
		for i, status := range filter.GrpcStatusFilter.Statuses {
			envoyGrpcStatusType, err := toEnvoyGRPCStatusType(status)
			if err != nil {
				return nil, err
			}
			statuses[i] = envoyGrpcStatusType
		}

		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_GrpcStatusFilter{
				GrpcStatusFilter: &envoyaccesslog.GrpcStatusFilter{
					Statuses: statuses,
					Exclude:  filter.GrpcStatusFilter.Exclude,
				},
			},
		}

	case filter.CELFilter != nil:
		celExpressionFilter := &cel.ExpressionFilter{
			Expression: filter.CELFilter.Match,
		}
		celCfg, err := utils.MessageToAny(celExpressionFilter)
		if err != nil {
			logger.Error(fmt.Sprintf("error converting CEL filter: %s", err.Error()))
			return nil, err
		}

		alCfg = &envoyaccesslog.AccessLogFilter{
			FilterSpecifier: &envoyaccesslog.AccessLogFilter_ExtensionFilter{
				ExtensionFilter: &envoyaccesslog.ExtensionFilter{
					Name: kgateway_wellknown.CELExtensionFilter,
					ConfigType: &envoyaccesslog.ExtensionFilter_TypedConfig{
						TypedConfig: celCfg,
					},
				},
			},
		}

	default:
		return nil, fmt.Errorf("no valid filter type specified")
	}

	return alCfg, nil
}

// Helper function to create header match specifier
func createHeaderMatchSpecifier(logger *zap.Logger, header gwv1.HTTPHeaderMatch) *envoyroute.HeaderMatcher_StringMatch {
	switch *header.Type {
	case gwv1.HeaderMatchExact:
		return &envoyroute.HeaderMatcher_StringMatch{
			StringMatch: &envoymatcher.StringMatcher{
				IgnoreCase: false,
				MatchPattern: &envoymatcher.StringMatcher_Exact{
					Exact: header.Value,
				},
			},
		}
	case gwv1.HeaderMatchRegularExpression:
		return &envoyroute.HeaderMatcher_StringMatch{
			StringMatch: &envoymatcher.StringMatcher{
				IgnoreCase: false,
				MatchPattern: &envoymatcher.StringMatcher_SafeRegex{
					SafeRegex: &envoymatcher.RegexMatcher{
						Regex: header.Value,
					},
				},
			},
		}
	default:
		logger.Error(fmt.Sprintf("unsupported header match type: %s", *header.Type))
		return nil
	}
}

func validateFilter(filter *v1alpha1.FilterType) error {
	count := 0
	if filter.StatusCodeFilter != nil {
		count++
	}
	if filter.DurationFilter != nil {
		count++
	}
	if filter.NotHealthCheckFilter != nil {
		count++
	}
	if filter.TraceableFilter != nil {
		count++
	}
	if filter.RuntimeFilter != nil {
		count++
	}
	if filter.HeaderFilter != nil {
		count++
	}
	if filter.ResponseFlagFilter != nil {
		count++
	}
	if filter.GrpcStatusFilter != nil {
		count++
	}
	if filter.CELFilter != nil {
		count++
	}

	if count != 1 {
		return fmt.Errorf("exactly one AccessLogFilter type must be set")
	}

	return nil
}

func convertJsonFormat(jsonFormat *runtime.RawExtension) *structpb.Struct {
	if jsonFormat == nil {
		return nil
	}

	var formatMap map[string]interface{}
	if err := json.Unmarshal(jsonFormat.Raw, &formatMap); err != nil {
		return nil
	}

	structVal, err := structpb.NewStruct(formatMap)
	if err != nil {
		return nil
	}

	return structVal
}

func copyGrpcSettings(cfg *envoygrpc.HttpGrpcAccessLogConfig, grpcService *v1alpha1.GrpcService) error {
	if grpcService == nil {
		return eris.New("grpc service object cannot be nil")
	}

	svc := &envoycore.GrpcService{
		TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
			EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
				ClusterName: grpcService.StaticClusterName,
			},
		},
	}
	cfg.AdditionalRequestHeadersToLog = grpcService.AdditionalRequestHeadersToLog
	cfg.AdditionalResponseHeadersToLog = grpcService.AdditionalResponseHeadersToLog
	cfg.AdditionalResponseTrailersToLog = grpcService.AdditionalResponseTrailersToLog
	cfg.CommonConfig = &envoygrpc.CommonGrpcAccessLogConfig{
		LogName:             grpcService.LogName,
		GrpcService:         svc,
		TransportApiVersion: envoycore.ApiVersion_V3,
	}
	return cfg.Validate()
}

func getFormatterExtensions() ([]*envoycore.TypedExtensionConfig, error) {
	reqWithoutQueryFormatter := &envoy_req_without_query.ReqWithoutQuery{}
	reqWithoutQueryFormatterTc, err := utils.MessageToAny(reqWithoutQueryFormatter)
	if err != nil {
		return nil, err
	}

	mdFormatter := &envoy_metadata_formatter.Metadata{}
	mdFormatterTc, err := utils.MessageToAny(mdFormatter)
	if err != nil {
		return nil, err
	}

	return []*envoycore.TypedExtensionConfig{
		{
			Name:        "envoy.formatter.req_without_query",
			TypedConfig: reqWithoutQueryFormatterTc,
		},
		{
			Name:        "envoy.formatter.metadata",
			TypedConfig: mdFormatterTc,
		},
	}, nil

}

func newAccessLogWithConfig(name string, config proto.Message) (*envoyaccesslog.AccessLog, error) {
	s := &envoyaccesslog.AccessLog{
		Name: name,
	}

	if config != nil {
		marshalledConf, err := utils.MessageToAny(config)
		if err != nil {
			// this should NEVER HAPPEN!
			return nil, err
		}

		s.ConfigType = &envoyaccesslog.AccessLog_TypedConfig{
			TypedConfig: marshalledConf,
		}
	}

	return s, nil
}

// String provides a string representation for the Op enum.
func toEnvoyComparisonOpType(op v1alpha1.Op) (envoyaccesslog.ComparisonFilter_Op, error) {
	switch op {
	case v1alpha1.EQ:
		return envoyaccesslog.ComparisonFilter_EQ, nil
	case v1alpha1.GE:
		return envoyaccesslog.ComparisonFilter_EQ, nil
	case v1alpha1.LE:
		return envoyaccesslog.ComparisonFilter_EQ, nil
	default:
		return 0, eris.Errorf("Unknown OP (%s)", op)
	}
}

func toEnvoyDenominatorType(denominatorType v1alpha1.DenominatorType) (envoytype.FractionalPercent_DenominatorType, error) {
	switch denominatorType {
	case v1alpha1.HUNDRED:
		return envoytype.FractionalPercent_HUNDRED, nil
	case v1alpha1.TEN_THOUSAND:
		return envoytype.FractionalPercent_TEN_THOUSAND, nil
	case v1alpha1.MILLION:
		return envoytype.FractionalPercent_MILLION, nil
	default:
		return 0, eris.Errorf("Unknown DenominatorType (%s)", denominatorType)
	}
}

func toEnvoyGRPCStatusType(grpcStatus v1alpha1.GrpcStatus) (envoyaccesslog.GrpcStatusFilter_Status, error) {
	switch grpcStatus {
	case v1alpha1.OK:
		return envoyaccesslog.GrpcStatusFilter_OK, nil
	case v1alpha1.CANCELED:
		return envoyaccesslog.GrpcStatusFilter_CANCELED, nil
	case v1alpha1.UNKNOWN:
		return envoyaccesslog.GrpcStatusFilter_UNKNOWN, nil
	case v1alpha1.INVALID_ARGUMENT:
		return envoyaccesslog.GrpcStatusFilter_INVALID_ARGUMENT, nil
	case v1alpha1.DEADLINE_EXCEEDED:
		return envoyaccesslog.GrpcStatusFilter_DEADLINE_EXCEEDED, nil
	case v1alpha1.NOT_FOUND:
		return envoyaccesslog.GrpcStatusFilter_NOT_FOUND, nil
	case v1alpha1.ALREADY_EXISTS:
		return envoyaccesslog.GrpcStatusFilter_ALREADY_EXISTS, nil
	case v1alpha1.PERMISSION_DENIED:
		return envoyaccesslog.GrpcStatusFilter_PERMISSION_DENIED, nil
	case v1alpha1.RESOURCE_EXHAUSTED:
		return envoyaccesslog.GrpcStatusFilter_RESOURCE_EXHAUSTED, nil
	case v1alpha1.FAILED_PRECONDITION:
		return envoyaccesslog.GrpcStatusFilter_FAILED_PRECONDITION, nil
	case v1alpha1.ABORTED:
		return envoyaccesslog.GrpcStatusFilter_ABORTED, nil
	case v1alpha1.OUT_OF_RANGE:
		return envoyaccesslog.GrpcStatusFilter_OUT_OF_RANGE, nil
	case v1alpha1.UNIMPLEMENTED:
		return envoyaccesslog.GrpcStatusFilter_UNIMPLEMENTED, nil
	case v1alpha1.INTERNAL:
		return envoyaccesslog.GrpcStatusFilter_INTERNAL, nil
	case v1alpha1.UNAVAILABLE:
		return envoyaccesslog.GrpcStatusFilter_UNAVAILABLE, nil
	case v1alpha1.DATA_LOSS:
		return envoyaccesslog.GrpcStatusFilter_DATA_LOSS, nil
	case v1alpha1.UNAUTHENTICATED:
		return envoyaccesslog.GrpcStatusFilter_UNAUTHENTICATED, nil
	default:
		return 0, eris.Errorf("Unknown GRPCStatus (%s)", grpcStatus)
	}
}
