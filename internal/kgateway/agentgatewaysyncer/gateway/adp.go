package gateway

import (
	"fmt"
	"strings"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/config"
	k8s "sigs.k8s.io/gateway-api/apis/v1"

	"istio.io/istio/pkg/ptr"
	"istio.io/istio/pkg/slices"
	"istio.io/istio/pkg/workloadapi"
)

func createADPMethodMatch(match k8s.HTTPRouteMatch) (*api.MethodMatch, *ConfigError) {
	if match.Method == nil {
		return nil, nil
	}
	return &workloadapi.MethodMatch{
		Exact: string(*match.Method),
	}, nil
}

func createADPQueryMatch(match k8s.HTTPRouteMatch) ([]*workloadapi.QueryMatch, *ConfigError) {
	res := []*workloadapi.QueryMatch{}
	for _, header := range match.QueryParams {
		tp := k8s.QueryParamMatchExact
		if header.Type != nil {
			tp = *header.Type
		}
		switch tp {
		case k8s.QueryParamMatchExact:
			res = append(res, &workloadapi.QueryMatch{
				Name:  string(header.Name),
				Value: &workloadapi.QueryMatch_Exact{Exact: header.Value},
			})
		case k8s.QueryParamMatchRegularExpression:
			res = append(res, &workloadapi.QueryMatch{
				Name:  string(header.Name),
				Value: &workloadapi.QueryMatch_Regex{Regex: header.Value},
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

func createADPPathMatch(match k8s.HTTPRouteMatch) (*workloadapi.PathMatch, *ConfigError) {
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
		return &workloadapi.PathMatch{Kind: &workloadapi.PathMatch_PathPrefix{
			PathPrefix: dest,
		}}, nil
	case k8s.PathMatchExact:
		return &workloadapi.PathMatch{Kind: &workloadapi.PathMatch_Exact{
			Exact: dest,
		}}, nil
	case k8s.PathMatchRegularExpression:
		return &workloadapi.PathMatch{Kind: &workloadapi.PathMatch_Regex{
			Regex: dest,
		}}, nil
	default:
		// Should never happen, unless a new field is added
		return nil, &ConfigError{Reason: InvalidConfiguration, Message: fmt.Sprintf("unknown type: %q is not supported Path match type", tp)}
	}
}

func createADPHeadersMatch(match k8s.HTTPRouteMatch) ([]*workloadapi.HeaderMatch, *ConfigError) {
	res := []*workloadapi.HeaderMatch{}
	for _, header := range match.Headers {
		tp := k8s.HeaderMatchExact
		if header.Type != nil {
			tp = *header.Type
		}
		switch tp {
		case k8s.HeaderMatchExact:
			res = append(res, &workloadapi.HeaderMatch{
				Name:  string(header.Name),
				Value: &workloadapi.HeaderMatch_Exact{Exact: header.Value},
			})
		case k8s.HeaderMatchRegularExpression:
			res = append(res, &workloadapi.HeaderMatch{
				Name:  string(header.Name),
				Value: &workloadapi.HeaderMatch_Regex{Regex: header.Value},
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

func createADPHeadersFilter(filter *k8s.HTTPHeaderFilter) *workloadapi.RouteFilter {
	if filter == nil {
		return nil
	}
	return &workloadapi.RouteFilter{
		Kind: &workloadapi.RouteFilter_RequestHeaderModifier{
			RequestHeaderModifier: &workloadapi.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPResponseHeadersFilter(filter *k8s.HTTPHeaderFilter) *workloadapi.RouteFilter {
	if filter == nil {
		return nil
	}
	return &workloadapi.RouteFilter{
		Kind: &workloadapi.RouteFilter_ResponseHeaderModifier{
			ResponseHeaderModifier: &workloadapi.HeaderModifier{
				Add:    headerListToADP(filter.Add),
				Set:    headerListToADP(filter.Set),
				Remove: filter.Remove,
			},
		},
	}
}

func createADPRewriteFilter(filter *k8s.HTTPURLRewriteFilter) *workloadapi.RouteFilter {
	if filter == nil {
		return nil
	}
	ff := &workloadapi.UrlRewrite{
		Host: string(ptr.OrEmpty(filter.Hostname)),
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case k8s.PrefixMatchHTTPPathModifier:
			ff.Path = &workloadapi.UrlRewrite_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case k8s.FullPathHTTPPathModifier:
			ff.Path = &workloadapi.UrlRewrite_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &workloadapi.RouteFilter{
		Kind: &workloadapi.RouteFilter_UrlRewrite{
			UrlRewrite: ff,
		},
	}
}

func createADPMirrorFilter(
	ctx RouteContext,
	filter *k8s.HTTPRequestMirrorFilter,
	ns string,
	enforceRefGrant bool,
	k config.GroupVersionKind,
) (*workloadapi.RouteFilter, *ConfigError) {
	if filter == nil {
		return nil, nil
	}
	var weightOne int32 = 1
	dst, err := buildADPDestination(ctx, k8s.HTTPBackendRef{
		BackendRef: k8s.BackendRef{
			BackendObjectReference: filter.BackendRef,
			Weight:                 &weightOne,
		},
	}, ns, enforceRefGrant, k)
	if err != nil {
		return nil, err
	}
	var percent float64
	if f := filter.Fraction; f != nil {
		percent = (100 * float64(f.Numerator)) / float64(ptr.OrDefault(f.Denominator, int32(100)))
	} else if p := filter.Percent; p != nil {
		percent = float64(*p)
	} else {
		percent = 100
	}
	if percent == 0 {
		return nil, nil
	}
	rm := &workloadapi.RequestMirror{
		Kind:       nil,
		Percentage: percent,
		Port:       dst.Port,
	}
	switch dk := dst.Kind.(type) {
	case *workloadapi.RouteBackend_Service:
		rm.Kind = &workloadapi.RequestMirror_Service{
			Service: dk.Service,
		}
	}
	return &workloadapi.RouteFilter{Kind: &workloadapi.RouteFilter_RequestMirror{RequestMirror: rm}}, nil
}

func createADPRedirectFilter(filter *k8s.HTTPRequestRedirectFilter) *workloadapi.RouteFilter {
	if filter == nil {
		return nil
	}
	ff := &workloadapi.RequestRedirect{
		Scheme: ptr.OrEmpty(filter.Scheme),
		Host:   string(ptr.OrEmpty(filter.Hostname)),
		Port:   uint32(ptr.OrEmpty(filter.Port)),
		Status: uint32(ptr.OrEmpty(filter.StatusCode)),
	}
	if filter.Path != nil {
		switch filter.Path.Type {
		case k8s.PrefixMatchHTTPPathModifier:
			ff.Path = &workloadapi.RequestRedirect_Prefix{Prefix: strings.TrimSuffix(*filter.Path.ReplacePrefixMatch, "/")}
		case k8s.FullPathHTTPPathModifier:
			ff.Path = &workloadapi.RequestRedirect_Full{Full: strings.TrimSuffix(*filter.Path.ReplaceFullPath, "/")}
		}
	}
	return &workloadapi.RouteFilter{
		Kind: &workloadapi.RouteFilter_RequestRedirect{
			RequestRedirect: ff,
		},
	}
}

func headerListToADP(hl []k8s.HTTPHeader) []*workloadapi.Header {
	return slices.Map(hl, func(hl k8s.HTTPHeader) *workloadapi.Header {
		return &workloadapi.Header{
			Name:  string(hl.Name),
			Value: hl.Value,
		}
	})
}
