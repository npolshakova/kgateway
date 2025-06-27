package agentgatewaysyncer

import (
	"context"
	"testing"

	"github.com/agentgateway/agentgateway/go/api"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/krt/krttest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/stretchr/testify/assert"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
)

// TestAgentRouteResource tests the AgentRouteResource struct and its methods
func TestAgentRouteResource(t *testing.T) {
	tests := []struct {
		name     string
		resource AgentRouteResource
		other    AgentRouteResource
		expected bool
	}{
		{
			name: "identical resources",
			resource: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Route: &api.Route{
					Key:         "test-route",
					ListenerKey: "test-listener",
					RuleName:    "test-rule",
					RouteName:   "test-route",
				},
				Valid: true,
			},
			other: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Route: &api.Route{
					Key:         "test-route",
					ListenerKey: "test-listener",
					RuleName:    "test-rule",
					RouteName:   "test-route",
				},
				Valid: true,
			},
			expected: true,
		},
		{
			name: "different namespaces",
			resource: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Valid:          true,
			},
			other: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "other", Name: "test-route"},
				Valid:          true,
			},
			expected: false,
		},
		{
			name: "different validity",
			resource: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Valid:          true,
			},
			other: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Valid:          false,
			},
			expected: false,
		},
		{
			name: "one has route, other doesn't",
			resource: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Route:          &api.Route{Key: "test-route"},
				Valid:          true,
			},
			other: AgentRouteResource{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
				Valid:          true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resource.Equals(tt.other)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAgentRouteResourceResourceName tests the ResourceName method
func TestAgentRouteResourceResourceName(t *testing.T) {
	resource := AgentRouteResource{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
		Valid:          true,
	}

	expected := "default/test-route"
	result := resource.ResourceName()
	assert.Equal(t, expected, result)
}

// TestCanRouteAttachToGateway tests the canRouteAttachToGateway function
func TestCanRouteAttachToGateway(t *testing.T) {
	tests := []struct {
		name    string
		route   *gwv1.HTTPRoute
		gateway AgentGatewayResource
		expect  bool
	}{
		{
			name: "route attaches to gateway",
			route: &gwv1.HTTPRoute{
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name: "test-gateway",
							},
						},
					},
				},
			},
			gateway: AgentGatewayResource{
				NamespacedName: types.NamespacedName{Name: "test-gateway"},
				Valid:          true,
			},
			expect: true,
		},
		{
			name: "route doesn't attach to gateway - different name",
			route: &gwv1.HTTPRoute{
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name: "other-gateway",
							},
						},
					},
				},
			},
			gateway: AgentGatewayResource{
				NamespacedName: types.NamespacedName{Name: "test-gateway"},
				Valid:          true,
			},
			expect: false,
		},
		{
			name: "route doesn't attach to gateway - different namespace",
			route: &gwv1.HTTPRoute{
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name:      "test-gateway",
								Namespace: ptr.To(gwv1.Namespace("other-namespace")),
							},
						},
					},
				},
			},
			gateway: AgentGatewayResource{
				NamespacedName: types.NamespacedName{Name: "test-gateway", Namespace: "default"},
				Valid:          true,
			},
			expect: false,
		},
		{
			name: "route doesn't attach to gateway - section name specified",
			route: &gwv1.HTTPRoute{
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							{
								Name:        "test-gateway",
								SectionName: ptr.To(gwv1.SectionName("listener-1")),
							},
						},
					},
				},
			},
			gateway: AgentGatewayResource{
				NamespacedName: types.NamespacedName{Name: "test-gateway"},
				Valid:          true,
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canRouteAttachToGateway(tt.route, tt.gateway)
			assert.Equal(t, tt.expect, result)
		})
	}
}

// TestConvertPathMatch tests the convertPathMatch function
func TestConvertPathMatch(t *testing.T) {
	tests := []struct {
		name     string
		path     *gwv1.HTTPPathMatch
		expected *api.RouteMatch
	}{
		{
			name: "exact path match",
			path: &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchExact),
				Value: ptr.To("/api/v1"),
			},
			expected: &api.RouteMatch{
				Path: &api.PathMatch{
					Kind: &api.PathMatch_Exact{
						Exact: "/api/v1",
					},
				},
			},
		},
		{
			name: "path prefix match",
			path: &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchPathPrefix),
				Value: ptr.To("/api"),
			},
			expected: &api.RouteMatch{
				Path: &api.PathMatch{
					Kind: &api.PathMatch_PathPrefix{
						PathPrefix: "/api",
					},
				},
			},
		},
		{
			name: "regex path match",
			path: &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchRegularExpression),
				Value: ptr.To("/api/.*"),
			},
			expected: &api.RouteMatch{
				Path: &api.PathMatch{
					Kind: &api.PathMatch_Regex{
						Regex: "/api/.*",
					},
				},
			},
		},
		{
			name:     "nil path match",
			path:     nil,
			expected: nil,
		},
		{
			name: "path match without value",
			path: &gwv1.HTTPPathMatch{
				Type: ptr.To(gwv1.PathMatchExact),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertPathMatch(tt.path)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Path.GetExact(), result.Path.GetExact())
				assert.Equal(t, tt.expected.Path.GetPathPrefix(), result.Path.GetPathPrefix())
				assert.Equal(t, tt.expected.Path.GetRegex(), result.Path.GetRegex())
			}
		})
	}
}

// TestConvertHeaderMatch tests the convertHeaderMatch function
func TestConvertHeaderMatch(t *testing.T) {
	header := gwv1.HTTPHeaderMatch{
		Name:  "Content-Type",
		Value: "application/json",
	}

	// Currently returns nil as header matching is not implemented
	result := convertHeaderMatch(header)
	assert.Nil(t, result)
}

// TestConvertBackendRef tests the convertBackendRef function
func TestConvertBackendRef(t *testing.T) {
	// Create test agentGwAddress
	testAgentGwAddress := ServiceInfo{
		Service: &api.Service{
			Name:      "test-service",
			Namespace: "default",
			Hostname:  "test-service.default.cluster.local",
			Ports: []*api.Port{
				{
					ServicePort: 8080,
					TargetPort:  8080,
					AppProtocol: api.AppProtocol_HTTP11,
				},
			},
		},
	}

	// Initialize mock with test data
	var allInputs []any
	allInputs = append(allInputs, testAgentGwAddress)
	mock := krttest.NewMock(t, allInputs)
	services := krttest.GetMockCollection[ServiceInfo](mock)

	ctx := RouteContext{
		Krt:      krt.TestingDummyContext{},
		Services: services,
	}

	tests := []struct {
		name           string
		backend        gwv1.HTTPBackendRef
		routeNamespace string
		expectNil      bool
		expectedRoutes []*api.RouteBackend
	}{
		{
			name: "valid backend ref",
			backend: gwv1.HTTPBackendRef{
				BackendRef: gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name: "test-service",
					},
					Weight: ptr.To(int32(100)),
				},
			},
			routeNamespace: "default",
			expectNil:      false,
			expectedRoutes: []*api.RouteBackend{
				{
					Kind: &api.RouteBackend_Service{
						Service: "default/test-service.default.cluster.local",
					},
					Weight: 100,
					Port:   8080,
				},
			},
		},
		{
			name: "backend ref with explicit port",
			backend: gwv1.HTTPBackendRef{
				BackendRef: gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name: "test-service",
						Port: ptr.To(gwv1.PortNumber(9090)),
					},
				},
			},
			routeNamespace: "default",
			expectNil:      false,
			expectedRoutes: []*api.RouteBackend{
				{
					Kind: &api.RouteBackend_Service{
						Service: "default/test-service.default.cluster.local",
					},
					Weight: 1,
					Port:   9090,
				},
			},
		},
		{
			name: "backend ref with different namespace",
			backend: gwv1.HTTPBackendRef{
				BackendRef: gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      "test-service",
						Namespace: ptr.To(gwv1.Namespace("other-namespace")),
					},
				},
			},
			routeNamespace: "default",
			expectNil:      true, // Service not found in other namespace
			expectedRoutes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := convertBackendRef(ctx, tt.backend, tt.routeNamespace)
			if tt.expectNil {
				assert.Nil(t, results)
			} else {
				assert.NotEmpty(t, results)
				assert.Equal(t, len(tt.expectedRoutes), len(results))
				for idx, expectedRoute := range tt.expectedRoutes {
					resultRoute := results[idx]
					assert.Equal(t, expectedRoute.Kind, resultRoute.Kind)
					assert.Equal(t, expectedRoute.Weight, resultRoute.Weight)
					assert.Equal(t, expectedRoute.Port, resultRoute.Port)
				}
			}
		})
	}
}

// TestConvertHTTPFilter tests the convertHTTPFilter function
func TestConvertHTTPFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   gwv1.HTTPRouteFilter
		expected *api.RouteFilter
	}{
		{
			name: "URL rewrite filter",
			filter: gwv1.HTTPRouteFilter{
				Type: gwv1.HTTPRouteFilterURLRewrite,
				URLRewrite: &gwv1.HTTPURLRewriteFilter{
					Path: &gwv1.HTTPPathModifier{
						ReplaceFullPath: ptr.To("/new-path"),
					},
				},
			},
			expected: &api.RouteFilter{
				Kind: &api.RouteFilter_UrlRewrite{
					UrlRewrite: &api.UrlRewrite{
						Path: &api.UrlRewrite_Full{
							Full: "/new-path",
						},
					},
				},
			},
		},
		{
			name: "URL rewrite filter with prefix",
			filter: gwv1.HTTPRouteFilter{
				Type: gwv1.HTTPRouteFilterURLRewrite,
				URLRewrite: &gwv1.HTTPURLRewriteFilter{
					Path: &gwv1.HTTPPathModifier{
						ReplacePrefixMatch: ptr.To("/new-prefix"),
					},
				},
			},
			expected: &api.RouteFilter{
				Kind: &api.RouteFilter_UrlRewrite{
					UrlRewrite: &api.UrlRewrite{
						Path: &api.UrlRewrite_Prefix{
							Prefix: "/new-prefix",
						},
					},
				},
			},
		},
		{
			name: "unsupported filter type",
			filter: gwv1.HTTPRouteFilter{
				Type: "UnsupportedFilter",
			},
			expected: nil,
		},
		{
			name: "redirect filter (not implemented)",
			filter: gwv1.HTTPRouteFilter{
				Type:            gwv1.HTTPRouteFilterRequestRedirect,
				RequestRedirect: &gwv1.HTTPRequestRedirectFilter{},
			},
			expected: nil,
		},
		{
			name: "header modifier filter (not implemented)",
			filter: gwv1.HTTPRouteFilter{
				Type:                  gwv1.HTTPRouteFilterRequestHeaderModifier,
				RequestHeaderModifier: &gwv1.HTTPHeaderFilter{},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertHTTPFilter(tt.filter)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// Check that the filter type matches
				if urlRewrite := result.GetUrlRewrite(); urlRewrite != nil {
					expectedRewrite := tt.expected.GetUrlRewrite()
					assert.NotNil(t, expectedRewrite)
					if full := urlRewrite.GetFull(); full != "" {
						assert.Equal(t, expectedRewrite.GetFull(), full)
					}
					if prefix := urlRewrite.GetPrefix(); prefix != "" {
						assert.Equal(t, expectedRewrite.GetPrefix(), prefix)
					}
				}
			}
		})
	}
}

// TestConvertURLRewriteFilter tests the convertURLRewriteFilter function
func TestConvertURLRewriteFilter(t *testing.T) {
	tests := []struct {
		name     string
		rewrite  *gwv1.HTTPURLRewriteFilter
		expected *api.RouteFilter
	}{
		{
			name: "full path rewrite",
			rewrite: &gwv1.HTTPURLRewriteFilter{
				Path: &gwv1.HTTPPathModifier{
					ReplaceFullPath: ptr.To("/new-full-path"),
				},
			},
			expected: &api.RouteFilter{
				Kind: &api.RouteFilter_UrlRewrite{
					UrlRewrite: &api.UrlRewrite{
						Path: &api.UrlRewrite_Full{
							Full: "/new-full-path",
						},
					},
				},
			},
		},
		{
			name: "prefix rewrite",
			rewrite: &gwv1.HTTPURLRewriteFilter{
				Path: &gwv1.HTTPPathModifier{
					ReplacePrefixMatch: ptr.To("/new-prefix"),
				},
			},
			expected: &api.RouteFilter{
				Kind: &api.RouteFilter_UrlRewrite{
					UrlRewrite: &api.UrlRewrite{
						Path: &api.UrlRewrite_Prefix{
							Prefix: "/new-prefix",
						},
					},
				},
			},
		},
		{
			name:     "nil rewrite",
			rewrite:  nil,
			expected: nil,
		},
		{
			name:     "rewrite without path",
			rewrite:  &gwv1.HTTPURLRewriteFilter{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertURLRewriteFilter(tt.rewrite)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				urlRewrite := result.GetUrlRewrite()
				assert.NotNil(t, urlRewrite)
				expectedRewrite := tt.expected.GetUrlRewrite()
				if full := urlRewrite.GetFull(); full != "" {
					assert.Equal(t, expectedRewrite.GetFull(), full)
				}
				if prefix := urlRewrite.GetPrefix(); prefix != "" {
					assert.Equal(t, expectedRewrite.GetPrefix(), prefix)
				}
			}
		})
	}
}

// TestConvertHTTPRouteToAgentRoute tests the convertHTTPRouteToAgentRoute function
func TestConvertHTTPRouteToAgentRoute(t *testing.T) {
	// Create test agentGwAddress
	testAgentGwAddress := ServiceInfo{
		Service: &api.Service{
			Name:      "test-service",
			Namespace: "default",
			Hostname:  "test-service.default.cluster.local",
			Ports: []*api.Port{
				{
					ServicePort: 8080,
					TargetPort:  8080,
					AppProtocol: api.AppProtocol_HTTP11,
				},
			},
		},
	}

	// Initialize mock with test data
	var allInputs []any
	allInputs = append(allInputs, testAgentGwAddress)
	mock := krttest.NewMock(t, allInputs)
	services := krttest.GetMockCollection[ServiceInfo](mock)

	ctx := RouteContext{
		Krt:      krt.TestingDummyContext{},
		Services: services,
	}

	// Create test gateway
	gateway := AgentGatewayResource{
		NamespacedName: types.NamespacedName{Name: "test-gateway"},
		Listener: &Listener{
			Listener: &api.Listener{
				Key: "test-listener",
			},
		},
		Valid: true,
	}

	// Create test route
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-route",
			Namespace: "default",
		},
		Spec: gwv1.HTTPRouteSpec{
			Rules: []gwv1.HTTPRouteRule{
				{
					Matches: []gwv1.HTTPRouteMatch{
						{
							Path: &gwv1.HTTPPathMatch{
								Type:  ptr.To(gwv1.PathMatchPathPrefix),
								Value: ptr.To("/api"),
							},
						},
					},
					BackendRefs: []gwv1.HTTPBackendRef{
						{
							BackendRef: gwv1.BackendRef{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name: "test-service",
								},
								Weight: ptr.To(int32(100)),
							},
						},
					},
				},
			},
		},
	}

	rule := route.Spec.Rules[0]
	match := &route.Spec.Rules[0].Matches[0]

	result := convertHTTPRouteToAgentRoute(ctx, rule, route, 0, 0, match, gateway)

	assert.NotNil(t, result)
	assert.Equal(t, "default", result.Namespace)
	assert.Equal(t, "test-route", result.Name)
	assert.True(t, result.Valid)
	assert.NotNil(t, result.Route)

	// Check route details
	agentRoute := result.Route
	assert.Equal(t, "route-test-route-0-0-test-gateway", agentRoute.Key)
	assert.Equal(t, "test-listener", agentRoute.ListenerKey)
	assert.Equal(t, "rule-0", agentRoute.RuleName)
	assert.Equal(t, "route-0", agentRoute.RouteName)

	// Check matches
	assert.Len(t, agentRoute.Matches, 1)
	pathMatch := agentRoute.Matches[0].GetPath()
	assert.NotNil(t, pathMatch)
	assert.Equal(t, "/api", pathMatch.GetPathPrefix())

	// Check backends
	assert.Len(t, agentRoute.Backends, 1)
	backend := agentRoute.Backends[0]
	assert.Equal(t, int32(100), backend.Weight)
}

// TestAgentHTTPRouteCollection tests the AgentHTTPRouteCollection function
func TestAgentHTTPRouteCollection(t *testing.T) {
	// Create test HTTPRoute
	httpRoute := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-route",
			Namespace: "default",
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{
					{Name: "test-gateway"},
				},
			},
			Rules: []gwv1.HTTPRouteRule{
				{
					Matches: []gwv1.HTTPRouteMatch{
						{
							Path: &gwv1.HTTPPathMatch{
								Type:  ptr.To(gwv1.PathMatchPathPrefix),
								Value: ptr.To("/api"),
							},
						},
					},
					BackendRefs: []gwv1.HTTPBackendRef{
						{
							BackendRef: gwv1.BackendRef{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name: "test-service",
								},
							},
						},
					},
				},
			},
		},
	}

	// Create test gateway
	gateway := AgentGatewayResource{
		NamespacedName: types.NamespacedName{Name: "test-gateway"},
		Listener: &Listener{
			Listener: &api.Listener{
				Key: "test-listener",
			},
		},
		Valid: true,
	}

	// Create test service
	testAgentGwAddress := ServiceInfo{
		Service: &api.Service{
			Name:      "test-service",
			Namespace: "default",
			Hostname:  "test-service.default.cluster.local",
			Ports: []*api.Port{
				{
					ServicePort: 8080,
					TargetPort:  8080,
					AppProtocol: api.AppProtocol_HTTP11,
				},
			},
		},
	}

	// Initialize mock with test data
	var allInputs []any
	allInputs = append(allInputs, httpRoute)
	allInputs = append(allInputs, gateway)
	allInputs = append(allInputs, testAgentGwAddress)
	mock := krttest.NewMock(t, allInputs)

	httpRoutes := krttest.GetMockCollection[*gwv1.HTTPRoute](mock)
	gateways := krttest.GetMockCollection[AgentGatewayResource](mock)
	services := krttest.GetMockCollection[ServiceInfo](mock)

	inputs := RouteContextInputs{
		AgentGatewayResource: gateways,
		Services:             services,
	}

	krtopts := krtutil.KrtOptions{}

	result := AgentHTTPRouteCollection(httpRoutes, inputs, krtopts)

	assert.NotNil(t, result)
	assert.Equal(t, httpRoutes, result.Input)

	// Wait for collections to sync
	httpRoutes.WaitUntilSynced(context.Background().Done())
	gateways.WaitUntilSynced(context.Background().Done())
	result.Routes.WaitUntilSynced(context.Background().Done())

	// Verify that routes were created
	routes := result.Routes.List()
	assert.Len(t, routes, 1)

	route := routes[0]
	assert.Equal(t, "default", route.Namespace)
	assert.Equal(t, "test-route", route.Name)
	assert.True(t, route.Valid)
	assert.NotNil(t, route.Route)
}

// TestRouteResult tests the RouteResult struct
func TestRouteResult(t *testing.T) {
	// Create mock collections
	mock := krttest.NewMock(t, []any{})
	routes := krttest.GetMockCollection[AgentRouteResource](mock)
	input := krttest.GetMockCollection[*gwv1.HTTPRoute](mock)

	result := RouteResult[*gwv1.HTTPRoute]{
		Routes: routes,
		Input:  input,
	}

	assert.Equal(t, routes, result.Routes)
	assert.Equal(t, input, result.Input)
}
