package session_persistence

import (
	"context"
	"strings"

	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	testdefaults "github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/tests/base"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is a suite for testing session persistence functionality.
type testingSuite struct {
	*base.BaseTestingSuite
}

// NewTestingSuite creates a new testing suite for session persistence functionality.
func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setup, testCases),
	}
}

func (s *testingSuite) TestCookieSessionPersistence() {
	s.assertSessionPersistence("cookie")
}

func (s *testingSuite) TestHeaderSessionPersistence() {
	s.assertSessionPersistence("header")
}

// assertSessionPersistence makes multiple requests and verifies they go to the same backend pod
func (s *testingSuite) assertSessionPersistence(persistenceType string) {
	// Make sure the curl and echo pods are running.
	s.TestInstallation.Assertions.EventuallyPodsRunning(s.Ctx, echoDeployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=echo",
	})
	s.TestInstallation.Assertions.EventuallyPodsRunning(s.Ctx, echoDeployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: testdefaults.CurlPodLabelSelector,
	})

	var (
		gatewayService metav1.ObjectMeta
		sessionHeader  string
	)

	switch persistenceType {
	case "cookie":
		gatewayService = cookieGateway.ObjectMeta
	case "header":
		gatewayService = headerGateway.ObjectMeta
		sessionHeader = "session-a"
	}

	firstCurlOpts := []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(gatewayService)),
		curl.WithHostHeader("echo.local"),
		curl.WithPort(8080),
		curl.Silent(),
		curl.WithArgs([]string{"-i"}),
	}

	firstResp, err := s.TestInstallation.ClusterContext.Cli.CurlFromPod(s.Ctx, testdefaults.CurlPodExecOpt, firstCurlOpts...)
	s.Assert().NoError(err, "first request should succeed")

	firstPodName := s.extractPodNameFromResponse(firstResp.StdOut)
	s.Assert().NotEmpty(firstPodName, "should be able to extract pod name from first response. Response was: %s", firstResp.StdOut)

	var subsequentCurlOpts []curl.Option
	if persistenceType == "cookie" {
		cookie := s.extractSessionCookieFromResponse(firstResp.StdOut)
		s.Assert().NotEmpty(cookie, "should have received a session cookie. Response was: %s", firstResp.StdOut)
		subsequentCurlOpts = []curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(gatewayService)),
			curl.WithHostHeader("echo.local"),
			curl.WithPort(8080),
			curl.Silent(),
			curl.WithHeader("Cookie", cookie),
		}
	} else {
		headerValue := s.extractSessionHeaderFromResponse(firstResp.StdOut)
		s.Assert().NotEmpty(headerValue, "should have received a session header. Response was: %s", firstResp.StdOut)
		subsequentCurlOpts = []curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(gatewayService)),
			curl.WithHostHeader("echo.local"),
			curl.WithPort(8080),
			curl.Silent(),
			curl.WithHeader(sessionHeader, headerValue),
		}
	}

	for i := 0; i < 10; i++ {
		resp, err := s.TestInstallation.ClusterContext.Cli.CurlFromPod(s.Ctx, testdefaults.CurlPodExecOpt, subsequentCurlOpts...)
		s.Assert().NoError(err, "request %d should succeed", i+2)

		podName := s.extractPodNameFromResponse(resp.StdOut)
		s.Assert().Equal(firstPodName, podName, "request %d should go to the same pod due to session persistence. Expected: %s, Got: %s, Response: %s", i+2, firstPodName, podName, resp.StdOut)
	}
}

// extractPodNameFromResponse extracts the pod name from the echo service response
func (s *testingSuite) extractPodNameFromResponse(response string) string {
	// The echo service returns something like "pod=echo-abc123-xyz"
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "pod=") {
			parts := strings.Split(line, "pod=")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// extractSessionCookieFromResponse extracts the session cookie from the curl response
func (s *testingSuite) extractSessionCookieFromResponse(response string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "set-cookie:") && strings.Contains(line, "Session-A") {
			// Extract cookie value from "Set-Cookie: Session-A=value; ..."
			parts := strings.Split(line, "set-cookie:")
			if len(parts) > 1 {
				cookiePart := strings.TrimSpace(parts[1])
				// Take only the cookie name=value part (before any semicolon)
				if idx := strings.Index(cookiePart, ";"); idx != -1 {
					cookiePart = cookiePart[:idx]
				}
				return strings.TrimSpace(cookiePart)
			}
		}
	}
	return ""
}

// extractSessionHeaderFromResponse extracts the session header from the curl response
func (s *testingSuite) extractSessionHeaderFromResponse(response string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "session-a:") {
			parts := strings.Split(line, "session-a:")
			if len(parts) > 1 {
				cookiePart := strings.TrimSpace(parts[1])
				if idx := strings.Index(cookiePart, ";"); idx != -1 {
					cookiePart = cookiePart[:idx]
				}
				return strings.TrimSpace(cookiePart)
			}
		}
	}
	return ""
}
