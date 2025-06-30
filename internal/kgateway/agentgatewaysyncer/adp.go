package agentgatewaysyncer

import (
	"fmt"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/slices"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8s "sigs.k8s.io/gateway-api/apis/v1"
)

func createADPMethodMatch(match k8s.HTTPRouteMatch) (*api.MethodMatch, *ConfigError) {
	if match.Method == nil {
		return nil, nil
	}
	return &api.MethodMatch{
		Exact: string(*match.Method),
	}, nil
}

func createADPQueryMatch(match k8s.HTTPRouteMatch) ([]*api.QueryMatch, *ConfigError) {
	res := []*api.QueryMatch{}
	for _, header := range match.QueryParams {
		tp := k8s.QueryParamMatchExact
		if header.Type != nil {
			tp = *header.Type
		}
		switch tp {
		case k8s.QueryParamMatchExact:
			res = append(res, &api.QueryMatch{
				Name:  string(header.Name),
				Value: &api.QueryMatch_Exact{Exact: header.Value},
			})
		case k8s.QueryParamMatchRegularExpression:
			res = append(res, &api.QueryMatch{
				Name:  string(header.Name),
				Value: &api.QueryMatch_Regex{Regex: header.Value},
			})
		default:
			// Should never happen, unless a new field is added
			return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported QueryMatch type", tp)}
		}
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func createADPPathMatch(match k8s.HTTPRouteMatch) (*api.PathMatch, *ConfigError) {
	tp := k8s.PathMatchPathPrefix
	if match.Path.Type != nil {
		tp = *match.Path.Type
	}
	dest := "/"
	if match.Path.Value != nil {
		dest = *match.Path.Value
	}
	switch tp {
	case k8s.PathMatchPathPrefix:
		// "When specified, a trailing `/` is ignored."
		if dest != "/" {
			dest = strings.TrimSuffix(dest, "/")
		}
		return &api.PathMatch{Kind: &api.PathMatch_PathPrefix{
			PathPrefix: dest,
		}}, nil
	case k8s.PathMatchExact:
		return &api.PathMatch{Kind: &api.PathMatch_Exact{
			Exact: dest,
		}}, nil
	case k8s.PathMatchRegularExpression:
		return &api.PathMatch{Kind: &api.PathMatch_Regex{
			Regex: dest,
		}}, nil
	default:
		// Should never happen, unless a new field is added
		return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported Path match type", tp)}
	}
}

func createADPHeadersMatch(match k8s.HTTPRouteMatch) ([]*api.HeaderMatch, *ConfigError) {
	res := []*api.HeaderMatch{}
	for _, header := range match.Headers {
		tp := k8s.HeaderMatchExact
		if header.Type != nil {
			tp = *header.Type
		}
		switch tp {
		case k8s.HeaderMatchExact:
			res = append(res, &api.HeaderMatch{
				Name:  string(header.Name),
				Value: &api.HeaderMatch_Exact{Exact: header.Value},
			})
		case k8s.HeaderMatchRegularExpression:
			res = append(res, &api.HeaderMatch{
				Name:  string(header.Name),
				Value: &api.HeaderMatch_Regex{Regex: header.Value},
			})
		default:
			// Should never happen, unless a new field is added
			return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported HeaderMatch type", tp)}
		}
	}

	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func createADPHeadersFilter(filter *k8s.HTTPHeaderFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_RequestHeaderModifier{
			RequestHeaderModifier: &api.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPResponseHeadersFilter(filter *k8s.HTTPHeaderFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_ResponseHeaderModifier{
			ResponseHeaderModifier: &api.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPRewriteFilter(filter *k8s.HTTPURLRewriteFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}

	var hostname string
	if filter.Hostname != nil {
		hostname = string(*filter.Hostname)
	}
	ff := &api.UrlRewrite{
		Host: hostname,
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case k8s.PrefixMatchHTTPPathModifier:
			ff.Path = &api.UrlRewrite_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case k8s.FullPathHTTPPathModifier:
			ff.Path = &api.UrlRewrite_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_UrlRewrite{
			UrlRewrite: ff,
		},
	}
}

func createADPMirrorFilter(
	ctx RouteContext,
	filter *k8s.HTTPRequestMirrorFilter,
	ns string,
	k schema.GroupVersionKind,
) (*api.RouteFilter, *ConfigError) {
	if filter == nil {
		return nil, nil
	}
	var weightOne int32 = 1
	dst, err := buildADPDestination(ctx, k8s.HTTPBackendRef{
		BackendRef: k8s.BackendRef{
			BackendObjectReference: filter.BackendRef,
			Weight:                 &weightOne,
		},
	}, ns, k)
	if err != nil {
		return nil, err
	}
	var percent float64
	if f := filter.Fraction; f != nil {
		denominator := float64(100)
		if f.Denominator != nil {
			denominator = float64(*f.Denominator)
		}
		percent = (100 * float64(f.Numerator)) / denominator
	} else if p := filter.Percent; p != nil {
		percent = float64(*p)
	} else {
		percent = 100
	}
	if percent == 0 {
		return nil, nil
	}
	rm := &api.RequestMirror{
		Kind:       nil,
		Percentage: percent,
		Port:       dst.GetPort(),
	}
	switch dk := dst.GetKind().(type) {
	case *api.RouteBackend_Service:
		rm.Kind = &api.RequestMirror_Service{
			Service: dk.Service,
		}
	}
	return &api.RouteFilter{Kind: &api.RouteFilter_RequestMirror{RequestMirror: rm}}, nil
}

func createADPRedirectFilter(filter *k8s.HTTPRequestRedirectFilter) *api.RouteFilter {
	if filter == nil {
		return nil
	}
	var scheme, host string
	var port, statusCode uint32
	if filter.Scheme != nil {
		scheme = *filter.Scheme
	}
	if filter.Hostname != nil {
		host = string(*filter.Hostname)
	}
	if filter.Port != nil {
		port = uint32(*filter.Port)
	}
	if filter.StatusCode != nil {
		statusCode = uint32(*filter.StatusCode)
	}

	ff := &api.RequestRedirect{
		Scheme: scheme,
		Host:   host,
		Port:   port,
		Status: statusCode,
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case k8s.PrefixMatchHTTPPathModifier:
			ff.Path = &api.RequestRedirect_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case k8s.FullPathHTTPPathModifier:
			ff.Path = &api.RequestRedirect_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &api.RouteFilter{
		Kind: &api.RouteFilter_RequestRedirect{
			RequestRedirect: ff,
		},
	}
}

func headerListToADP(hl []k8s.HTTPHeader) []*api.Header {
	return slices.Map(hl, func(hl k8s.HTTPHeader) *api.Header {
		return &api.Header{
			Name:  string(hl.Name),
			Value: hl.Value,
		}
	})
}
