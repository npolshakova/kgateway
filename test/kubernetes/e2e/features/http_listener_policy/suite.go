package http_listener_policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	"github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	testdefaults "github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/tests/base"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "HttpListenerPolicy" feature
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setup, testCases),
	}
}

func (s *testingSuite) TestHttpListenerPolicyAllFields() {
	// Test that the HTTPListenerPolicy with all additional fields is applied correctly
	// The test verifies that the gateway is working and all policy fields are applied
	fmt.Println("TestHttpListenerPolicyAllFields")

	// Assert that the HTTPRoute is accepted
	s.TestInstallation.Assertions.EventuallyHTTPRouteCondition(s.Ctx, "example-route", "default", gwv1.RouteConditionAccepted, metav1.ConditionTrue)

	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("Welcome to nginx!"),
		})

	// Check the health check path is working
	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPath("/health_check"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.BeEmpty(),
		})
}

func (s *testingSuite) TestHttpListenerPolicyServerHeader() {
	// Test that the HTTPListenerPolicy with serverHeaderTransformation field is applied correctly
	// The test verifies that the server header is transformed as expected
	// With PassThrough, the server header should be the backend server's header (nginx/1.28.0)
	// instead of Envoy's default (envoy)

	// Assert that the HTTPRoute is accepted
	s.TestInstallation.Assertions.EventuallyHTTPRouteCondition(s.Ctx, "example-route", "default", gwv1.RouteConditionAccepted, metav1.ConditionTrue)

	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("Welcome to nginx!"),
			Headers: map[string]any{
				"server": "nginx/1.28.0", // Should be the backend server header, not "envoy"
			},
		})
}

func (s *testingSuite) TestPreserveHttp1HeaderCase() {
	// The test verifies that the HTTP1 headers are preserved as expected in the request and response
	// The HTTPListenerPolicy ensures that the header is preserved in the request,
	// and the BackendConfigPolicy ensures that the header is preserved in the response.

	// Assert that the HTTPRoute is accepted
	s.TestInstallation.Assertions.EventuallyHTTPRouteCondition(s.Ctx, "echo-route", "default", gwv1.RouteConditionAccepted, metav1.ConditionTrue)

	s.TestInstallation.Assertions.AssertEventualCurlResponse(
		s.Ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithHeader("X-CaSeD-HeAdEr", "test"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("X-CaSeD-HeAdEr"),
			Headers: map[string]any{
				"ReSpOnSe-miXed-CaSe-hEaDeR": "Foo",
			},
		},
	)
}
