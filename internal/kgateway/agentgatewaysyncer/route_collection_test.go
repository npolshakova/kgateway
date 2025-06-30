package agentgatewaysyncer

import (
	"context"
	"testing"

	"github.com/agentgateway/agentgateway/go/api"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/krt/krttest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/stretchr/testify/assert"
)

// TestAgentRouteResourceResourceName tests the ResourceName method
func TestAgentRouteResourceResourceName(t *testing.T) {
	resource := AgentRouteResource{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-route"},
	}

	expected := "default/test-route"
	result := resource.ResourceName()
	assert.Equal(t, expected, result)
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

// TestAgentGatewayRouteCollection tests the agentGatewayRouteCollection function
func TestAgentGatewayRouteCollection(t *testing.T) {
	tests := []struct {
		name           string
		httpRoutes     []*gwv1.HTTPRoute
		services       []*corev1.Service
		gateways       []*gwv1.Gateway
		expectedRoutes int
		expectedError  bool
	}{
		{
			name: "basic HTTP route with service backend",
			httpRoutes: []*gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-route",
						Namespace: "default",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{
								{
									Name:      "test-gateway",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Kind:      ptr.To(gwv1.Kind("Gateway")),
								},
							},
						},
						Hostnames: []gwv1.Hostname{"example.com"},
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
												Name:      "test-service",
												Namespace: ptr.To(gwv1.Namespace("default")),
												Port:      ptr.To(gwv1.PortNumber(8080)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 8080,
							},
						},
					},
				},
			},
			gateways: []*gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-gateway",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Gateway",
						APIVersion: "gateway.networking.k8s.io/v1",
					},
					Spec: gwv1.GatewaySpec{
						GatewayClassName: "agentgateway",
						Listeners: []gwv1.Listener{
							{
								Name:     "http",
								Port:     80,
								Protocol: gwv1.HTTPProtocolType,
							},
						},
					},
				},
			},
			expectedRoutes: 1,
			expectedError:  false,
		},
		{
			name: "HTTP route with multiple rules",
			httpRoutes: []*gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "multi-rule-route",
						Namespace: "default",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{
								{
									Name:      "test-gateway",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Kind:      ptr.To(gwv1.Kind("Gateway")),
								},
							},
						},
						Hostnames: []gwv1.Hostname{"example.com"},
						Rules: []gwv1.HTTPRouteRule{
							{
								Matches: []gwv1.HTTPRouteMatch{
									{
										Path: &gwv1.HTTPPathMatch{
											Type:  ptr.To(gwv1.PathMatchExact),
											Value: ptr.To("/api/v1"),
										},
									},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name:      "service1",
												Namespace: ptr.To(gwv1.Namespace("default")),
												Port:      ptr.To(gwv1.PortNumber(8080)),
											},
										},
									},
								},
							},
							{
								Matches: []gwv1.HTTPRouteMatch{
									{
										Path: &gwv1.HTTPPathMatch{
											Type:  ptr.To(gwv1.PathMatchPathPrefix),
											Value: ptr.To("/api/v2"),
										},
									},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name:      "service2",
												Namespace: ptr.To(gwv1.Namespace("default")),
												Port:      ptr.To(gwv1.PortNumber(8081)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 8080,
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service2",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 8081,
							},
						},
					},
				},
			},
			gateways: []*gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-gateway",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Gateway",
						APIVersion: "gateway.networking.k8s.io/v1",
					},
					Spec: gwv1.GatewaySpec{
						GatewayClassName: "agentgateway",
						Listeners: []gwv1.Listener{
							{
								Name:     "http",
								Port:     80,
								Protocol: gwv1.HTTPProtocolType,
							},
						},
					},
				},
			},
			expectedRoutes: 2,
			expectedError:  false,
		},
		{
			name: "HTTP route with filters",
			httpRoutes: []*gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "filtered-route",
						Namespace: "default",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{
								{
									Name:      "test-gateway",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Kind:      ptr.To(gwv1.Kind("Gateway")),
								},
							},
						},
						Hostnames: []gwv1.Hostname{"example.com"},
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
								Filters: []gwv1.HTTPRouteFilter{
									{
										Type: gwv1.HTTPRouteFilterRequestHeaderModifier,
										RequestHeaderModifier: &gwv1.HTTPHeaderFilter{
											Add: []gwv1.HTTPHeader{
												{
													Name:  "X-Custom-Header",
													Value: "custom-value",
												},
											},
										},
									},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name:      "test-service",
												Namespace: ptr.To(gwv1.Namespace("default")),
												Port:      ptr.To(gwv1.PortNumber(8080)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 8080,
							},
						},
					},
				},
			},
			gateways: []*gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-gateway",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Gateway",
						APIVersion: "gateway.networking.k8s.io/v1",
					},
					Spec: gwv1.GatewaySpec{
						GatewayClassName: "agentgateway",
						Listeners: []gwv1.Listener{
							{
								Name:     "http",
								Port:     80,
								Protocol: gwv1.HTTPProtocolType,
							},
						},
					},
				},
			},
			expectedRoutes: 1,
			expectedError:  false,
		},
		{
			name: "HTTP route with timeouts",
			httpRoutes: []*gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "timeout-route",
						Namespace: "default",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{
								{
									Name:      "test-gateway",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Kind:      ptr.To(gwv1.Kind("Gateway")),
								},
							},
						},
						Hostnames: []gwv1.Hostname{"example.com"},
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
								Timeouts: &gwv1.HTTPRouteTimeouts{
									Request:        ptr.To(gwv1.Duration("30s")),
									BackendRequest: ptr.To(gwv1.Duration("10s")),
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name:      "test-service",
												Namespace: ptr.To(gwv1.Namespace("default")),
												Port:      ptr.To(gwv1.PortNumber(8080)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-service",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 8080,
							},
						},
					},
				},
			},
			gateways: []*gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-gateway",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Gateway",
						APIVersion: "gateway.networking.k8s.io/v1",
					},
					Spec: gwv1.GatewaySpec{
						GatewayClassName: "agentgateway",
						Listeners: []gwv1.Listener{
							{
								Name:     "http",
								Port:     80,
								Protocol: gwv1.HTTPProtocolType,
							},
						},
					},
				},
			},
			expectedRoutes: 1,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock collections
			var allInputs []any
			for _, route := range tt.httpRoutes {
				allInputs = append(allInputs, route)
			}
			for _, svc := range tt.services {
				allInputs = append(allInputs, svc)
			}
			// Convert gwv1.Gateway to ir.Gateway for the mock
			for _, gw := range tt.gateways {
				irGateway := ir.Gateway{
					ObjectSource: ir.ObjectSource{
						Group:     gw.GetObjectKind().GroupVersionKind().Group,
						Kind:      gw.GetObjectKind().GroupVersionKind().Kind,
						Namespace: gw.Namespace,
						Name:      gw.Name,
					},
					Obj:       gw,
					Listeners: make([]ir.Listener, 0, len(gw.Spec.Listeners)),
				}
				for _, l := range gw.Spec.Listeners {
					irGateway.Listeners = append(irGateway.Listeners, ir.Listener{
						Listener: l,
						Parent:   gw,
					})
				}
				t.Logf("Debug: Created ir.Gateway: %+v", irGateway)
				t.Logf("Debug: ir.Gateway ObjectSource: %+v", irGateway.ObjectSource)
				allInputs = append(allInputs, irGateway)
			}

			mock := krttest.NewMock(t, allInputs)

			// Debug: Check what's in the mock
			t.Logf("Debug: Mock created with %d inputs", len(allInputs))
			for i, input := range allInputs {
				t.Logf("Debug: Input %d: %T", i, input)
			}

			// Setup collections
			httpRoutes := krttest.GetMockCollection[*gwv1.HTTPRoute](mock)
			services := krttest.GetMockCollection[*corev1.Service](mock)
			gateways := krttest.GetMockCollection[ir.Gateway](mock)

			// Create services index
			servicesIndex := krt.NewIndex(services, func(svc *corev1.Service) []string {
				return []string{svc.Namespace + "/" + svc.Name}
			})

			// Create route parents
			routeParents := BuildRouteParents(gateways)

			// Create route inputs
			routeInputs := RouteContextInputs{
				RouteParents:  routeParents,
				DomainSuffix:  "cluster.local",
				Services:      services,
				ServicesIndex: servicesIndex,
				Namespaces:    nil, // Not needed for this test
			}

			// Wait for collections to sync
			httpRoutes.WaitUntilSynced(context.Background().Done())
			services.WaitUntilSynced(context.Background().Done())
			gateways.WaitUntilSynced(context.Background().Done())

			// Call the function under test
			routeCollection := agentGatewayRouteCollection(httpRoutes, routeInputs)

			// Add debugging to understand what's happening
			allRoutes := krt.Fetch(krt.TestingDummyContext{}, routeCollection.routes)
			t.Logf("Debug: Found %d routes", len(allRoutes))

			// Debug: Check if HTTP routes are being processed
			allHTTPRoutes := krt.Fetch(krt.TestingDummyContext{}, httpRoutes)
			t.Logf("Debug: Found %d HTTP routes", len(allHTTPRoutes))

			// Debug: Check if parent references are being found
			if len(allHTTPRoutes) > 0 {
				route := allHTTPRoutes[0]
				ctx := routeInputs.WithCtx(krt.TestingDummyContext{})
				parentRefs := extractParentReferenceInfo(ctx, routeParents, route)
				t.Logf("Debug: Found %d parent references", len(parentRefs))
				for i, pr := range parentRefs {
					t.Logf("Debug: Parent ref %d: %+v", i, pr)
				}

				filteredParentRefs := filteredReferences(parentRefs)
				t.Logf("Debug: Found %d filtered parent references", len(filteredParentRefs))
			}

			// Debug: Check if gateways are being found
			allGateways := krt.Fetch(krt.TestingDummyContext{}, gateways)
			t.Logf("Debug: Found %d gateways", len(allGateways))

			// Debug: Check what parent keys are being created for gateways
			for _, gw := range allGateways {
				gvk := gw.Obj.GetObjectKind().GroupVersionKind()
				t.Logf("Debug: Gateway GVK: %+v", gvk)
				t.Logf("Debug: Gateway Name: %s, Namespace: %s", gw.Name, gw.Namespace)
			}

			// Debug: Check if services are being found
			allServices := krt.Fetch(krt.TestingDummyContext{}, services)
			t.Logf("Debug: Found %d services", len(allServices))

			// Debug: Check service index
			if len(allServices) > 0 {
				svc := allServices[0]
				key := svc.Namespace + "/" + svc.Name
				t.Logf("Debug: Service key: %s", key)
				indexedServices := krt.Fetch(krt.TestingDummyContext{}, services, krt.FilterIndex(servicesIndex, key))
				t.Logf("Debug: Found %d services by index for key %s", len(indexedServices), key)
			}

			// Debug: Check parent references
			if len(allHTTPRoutes) > 0 {
				route := allHTTPRoutes[0]
				t.Logf("Debug: HTTP route parent refs: %+v", route.Spec.ParentRefs)
				t.Logf("Debug: HTTP route backend refs: %+v", route.Spec.Rules[0].BackendRefs)

				// Debug: Check what parent key is being created
				if len(route.Spec.ParentRefs) > 0 {
					parentRef := route.Spec.ParentRefs[0]
					pk, err := toInternalParentReference(parentRef, route.Namespace)
					if err != nil {
						t.Logf("Debug: Error creating parent key: %v", err)
					} else {
						t.Logf("Debug: Created parent key: %+v", pk)
					}
				}
			}

			// Verify the results
			assert.NotNil(t, routeCollection)
			assert.NotNil(t, routeCollection.routes)
			assert.NotNil(t, routeCollection.routeIndex)

			// Fetch all routes to verify count
			assert.Equal(t, tt.expectedRoutes, len(allRoutes), "Expected %d routes, got %d", tt.expectedRoutes, len(allRoutes))

			// Verify route structure for the first test case
			if tt.name == "basic HTTP route with service backend" && len(allRoutes) > 0 {
				route := allRoutes[0]
				assert.NotNil(t, route.Route)
				assert.Equal(t, "default.test-route.0.0", route.Route.Key)
				assert.Equal(t, "default/test-route", route.Route.RouteName)
				assert.Equal(t, "", route.Route.ListenerKey)
				assert.Equal(t, "", route.Route.RuleName)
				assert.Len(t, route.Route.Matches, 1)
				assert.Len(t, route.Route.Backends, 1)
				assert.Len(t, route.Route.Hostnames, 1)
				assert.Equal(t, "example.com", route.Route.Hostnames[0])

				// Verify path match
				pathMatch := route.Route.Matches[0].Path
				assert.NotNil(t, pathMatch)
				assert.Equal(t, "/api", pathMatch.GetPathPrefix())

				// Verify backend
				backend := route.Route.Backends[0]
				assert.Equal(t, int32(8080), backend.Port)
				assert.Equal(t, int32(1), backend.Weight)
				assert.NotNil(t, backend.Kind)
				serviceBackend := backend.Kind.(*api.RouteBackend_Service)
				assert.Equal(t, "default/test-service.default.svc.cluster.local", serviceBackend.Service)
			}

			// Verify route index functionality
			if len(tt.gateways) > 0 {
				gw := tt.gateways[0]
				parentKey := parentKey{
					Kind:      gw.GetObjectKind().GroupVersionKind(),
					Name:      gw.Name,
					Namespace: gw.Namespace,
				}
				indexedRoutes := krt.Fetch(krt.TestingDummyContext{}, routeCollection.routes, krt.FilterIndex(routeCollection.routeIndex, parentKey))
				assert.GreaterOrEqual(t, len(indexedRoutes), 0, "Should be able to fetch routes by parent key")
			}
		})
	}
}

// TestConvertHTTPRouteToAGWRoute tests the convertHTTPRouteToAGWRoute function
func TestConvertHTTPRouteToAGWRoute(t *testing.T) {
	tests := []struct {
		name          string
		routeRule     gwv1.HTTPRouteRule
		httpRoute     *gwv1.HTTPRoute
		pos           int
		matchPos      int
		expectedError bool
		expectedRoute *api.Route
	}{
		{
			name: "basic route rule",
			routeRule: gwv1.HTTPRouteRule{
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
								Name:      "test-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"example.com"},
				},
			},
			pos:           0,
			matchPos:      0,
			expectedError: false,
			expectedRoute: &api.Route{
				Key:         "default.test-route.0.0",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/test-route",
				Hostnames:   []string{"example.com"},
			},
		},
		{
			name: "route rule with named rule",
			routeRule: gwv1.HTTPRouteRule{
				Name: (*gwv1.SectionName)(ptr.To("test-rule")),
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchExact),
							Value: ptr.To("/api/v1"),
						},
					},
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "test-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"example.com"},
				},
			},
			pos:      1,
			matchPos: 0,
			expectedRoute: &api.Route{
				Key:         "default.test-route.1.0",
				ListenerKey: "",
				RuleName:    "test-rule",
				RouteName:   "default/test-route",
				Hostnames:   []string{"example.com"},
			},
			expectedError: false,
		},
		{
			name: "route rule with regex path match",
			routeRule: gwv1.HTTPRouteRule{
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchRegularExpression),
							Value: ptr.To("/api/.*"),
						},
					},
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "test-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "regex-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"api.example.com"},
				},
			},
			pos:      0,
			matchPos: 0,
			expectedRoute: &api.Route{
				Key:         "default.regex-route.0.0",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/regex-route",
				Hostnames:   []string{"api.example.com"},
			},
			expectedError: false,
		},
		{
			name: "route rule with multiple matches",
			routeRule: gwv1.HTTPRouteRule{
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchPathPrefix),
							Value: ptr.To("/api"),
						},
					},
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchExact),
							Value: ptr.To("/health"),
						},
					},
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "health-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "multi-match-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"health.example.com"},
				},
			},
			pos:      0,
			matchPos: 1, // Testing second match
			expectedRoute: &api.Route{
				Key:         "default.multi-match-route.0.1",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/multi-match-route",
				Hostnames:   []string{"health.example.com"},
			},
			expectedError: false,
		},
		{
			name: "route rule with URL rewrite filter",
			routeRule: gwv1.HTTPRouteRule{
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchPathPrefix),
							Value: ptr.To("/api"),
						},
					},
				},
				Filters: []gwv1.HTTPRouteFilter{
					{
						Type: gwv1.HTTPRouteFilterURLRewrite,
						URLRewrite: &gwv1.HTTPURLRewriteFilter{
							Path: &gwv1.HTTPPathModifier{
								ReplaceFullPath: ptr.To("/new-path"),
							},
						},
					},
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "rewrite-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rewrite-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"rewrite.example.com"},
				},
			},
			pos:      0,
			matchPos: 0,
			expectedRoute: &api.Route{
				Key:         "default.rewrite-route.0.0",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/rewrite-route",
				Hostnames:   []string{"rewrite.example.com"},
			},
			expectedError: false,
		},
		{
			name: "route rule with timeouts",
			routeRule: gwv1.HTTPRouteRule{
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchPathPrefix),
							Value: ptr.To("/slow-api"),
						},
					},
				},
				Timeouts: &gwv1.HTTPRouteTimeouts{
					Request:        ptr.To(gwv1.Duration("30s")),
					BackendRequest: ptr.To(gwv1.Duration("10s")),
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "slow-service",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "timeout-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"timeout.example.com"},
				},
			},
			pos:      0,
			matchPos: 0,
			expectedRoute: &api.Route{
				Key:         "default.timeout-route.0.0",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/timeout-route",
				Hostnames:   []string{"timeout.example.com"},
			},
			expectedError: false,
		},
		{
			name: "route rule with multiple backends",
			routeRule: gwv1.HTTPRouteRule{
				Matches: []gwv1.HTTPRouteMatch{
					{
						Path: &gwv1.HTTPPathMatch{
							Type:  ptr.To(gwv1.PathMatchPathPrefix),
							Value: ptr.To("/load-balanced"),
						},
					},
				},
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "backend1",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8080)),
							},
							Weight: ptr.To(int32(60)),
						},
					},
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "backend2",
								Namespace: ptr.To(gwv1.Namespace("default")),
								Port:      ptr.To(gwv1.PortNumber(8081)),
							},
							Weight: ptr.To(int32(40)),
						},
					},
				},
			},
			httpRoute: &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "load-balanced-route",
					Namespace: "default",
				},
				Spec: gwv1.HTTPRouteSpec{
					Hostnames: []gwv1.Hostname{"lb.example.com"},
				},
			},
			pos:      0,
			matchPos: 0,
			expectedRoute: &api.Route{
				Key:         "default.load-balanced-route.0.0",
				ListenerKey: "",
				RuleName:    "",
				RouteName:   "default/load-balanced-route",
				Hostnames:   []string{"lb.example.com"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock context
			ctx := RouteContext{
				Krt:           krt.TestingDummyContext{},
				RouteParents:  RouteParents{},
				DomainSuffix:  "cluster.local",
				Services:      nil,
				ServicesIndex: nil,
				Namespaces:    nil,
			}

			// Call the function under test
			result, err := convertHTTPRouteToAGWRoute(ctx, tt.routeRule, tt.httpRoute, tt.pos, tt.matchPos)

			if tt.expectedError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedRoute.Key, result.Key)
				assert.Equal(t, "default/test-route", result.RouteName)
				assert.Equal(t, "", result.ListenerKey)
				assert.Equal(t, defaultString(tt.routeRule.Name, ""), result.RuleName)
				assert.Len(t, result.Hostnames, 1)
				assert.Equal(t, "example.com", result.Hostnames[0])

				// Verify matches
				if len(tt.routeRule.Matches) > 0 {
					assert.Len(t, result.Matches, len(tt.routeRule.Matches))
					for i, match := range tt.routeRule.Matches {
						if match.Path != nil {
							pathMatch := result.Matches[i].Path
							assert.NotNil(t, pathMatch)
							if match.Path.Type != nil && *match.Path.Type == gwv1.PathMatchExact {
								assert.Equal(t, *match.Path.Value, pathMatch.GetExact())
							} else {
								assert.Equal(t, *match.Path.Value, pathMatch.GetPathPrefix())
							}
						}
					}
				}

				// Verify backends
				if len(tt.routeRule.BackendRefs) > 0 {
					assert.Len(t, result.Backends, len(tt.routeRule.BackendRefs))
				}
			}
		})
	}
}

// TestCreateADPPathMatch tests the createADPPathMatch function
func TestCreateADPPathMatch(t *testing.T) {
	tests := []struct {
		name          string
		match         gwv1.HTTPRouteMatch
		expectedPath  string
		expectedType  string
		expectedError bool
	}{
		{
			name: "path prefix match",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/api"),
				},
			},
			expectedPath:  "/api",
			expectedType:  "path_prefix",
			expectedError: false,
		},
		{
			name: "exact path match",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchExact),
					Value: ptr.To("/api/v1"),
				},
			},
			expectedPath:  "/api/v1",
			expectedType:  "exact",
			expectedError: false,
		},
		{
			name: "regex path match",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchRegularExpression),
					Value: ptr.To("/api/.*"),
				},
			},
			expectedPath:  "/api/.*",
			expectedType:  "regex",
			expectedError: false,
		},
		{
			name: "path prefix with trailing slash",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/api/"),
				},
			},
			expectedPath:  "/api",
			expectedType:  "path_prefix",
			expectedError: false,
		},
		{
			name: "default path prefix",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Value: ptr.To("/api"),
				},
			},
			expectedPath:  "/api",
			expectedType:  "path_prefix",
			expectedError: false,
		},
		{
			name: "root path",
			match: gwv1.HTTPRouteMatch{
				Path: &gwv1.HTTPPathMatch{
					Type:  ptr.To(gwv1.PathMatchPathPrefix),
					Value: ptr.To("/"),
				},
			},
			expectedPath:  "/",
			expectedType:  "path_prefix",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createADPPathMatch(tt.match)

			if tt.expectedError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)

				switch tt.expectedType {
				case "path_prefix":
					assert.Equal(t, tt.expectedPath, result.GetPathPrefix())
				case "exact":
					assert.Equal(t, tt.expectedPath, result.GetExact())
				case "regex":
					assert.Equal(t, tt.expectedPath, result.GetRegex())
				}
			}
		})
	}
}

// TestCreateADPHeadersMatch tests the createADPHeadersMatch function
func TestCreateADPHeadersMatch(t *testing.T) {
	tests := []struct {
		name          string
		match         gwv1.HTTPRouteMatch
		expectedCount int
		expectedError bool
	}{
		{
			name: "exact header match",
			match: gwv1.HTTPRouteMatch{
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Name:  "Content-Type",
						Value: "application/json",
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "regex header match",
			match: gwv1.HTTPRouteMatch{
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Name:  "User-Agent",
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Value: ".*Chrome.*",
					},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "multiple header matches",
			match: gwv1.HTTPRouteMatch{
				Headers: []gwv1.HTTPHeaderMatch{
					{
						Name:  "Content-Type",
						Value: "application/json",
					},
					{
						Name:  "Authorization",
						Value: "Bearer token",
					},
				},
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "no headers",
			match:         gwv1.HTTPRouteMatch{},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createADPHeadersMatch(tt.match)

			if tt.expectedError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				if tt.expectedCount == 0 {
					assert.NotNil(t, err)
					assert.Nil(t, result)
				} else {
					assert.Nil(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, tt.expectedCount)

					for i, header := range tt.match.Headers {
						adpHeader := result[i]
						assert.Equal(t, string(header.Name), adpHeader.Name)

						if header.Type != nil && *header.Type == gwv1.HeaderMatchRegularExpression {
							assert.Equal(t, header.Value, adpHeader.GetRegex())
						} else {
							assert.Equal(t, header.Value, adpHeader.GetExact())
						}
					}
				}
			}
		})
	}
}
