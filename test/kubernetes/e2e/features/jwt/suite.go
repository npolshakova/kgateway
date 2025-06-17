package jwt

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	testdefaults "github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is a suite of tests for jwt functionality
type testingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of kgateway
	testInstallation *e2e.TestInstallation

	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string

	// Track core objects for cleanup
	coreObjects []client.Object
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

// SetupSuite runs before all tests in the suite
func (s *testingSuite) SetupSuite() {
	// Initialize test manifest mappings
	s.manifests = map[string][]string{
		"TestJwtAuthentication": {jwtManifest},
		"TestJWTAuthorization":  {jwtRbacManifest},
	}

	// Apply core infrastructure
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, setupManifest)
	s.Require().NoError(err)

	// Apply curl pod for testing
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.Require().NoError(err)

	// Track core objects
	s.coreObjects = []client.Object{
		testdefaults.CurlPod,              // curl
		httpbinDeployment,                 // httpbin
		gatewayService, gatewayDeployment, // gateway service
	}

	// Wait for core infrastructure to be ready
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, s.coreObjects...)
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: testdefaults.CurlPodLabelSelector,
	})
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, httpbinDeployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin",
	})
	s.testInstallation.Assertions.EventuallyPodsRunning(
		s.ctx,
		gatewayDeployment.ObjectMeta.GetNamespace(),
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", gatewayObjectMeta.GetName()),
		},
	)
	s.testInstallation.Assertions.EventuallyHTTPRouteCondition(s.ctx, "httpbin", "httpbin", gwv1.RouteConditionAccepted, metav1.ConditionTrue)
}

// TearDownSuite cleans up any remaining resources
func (s *testingSuite) TearDownSuite() {
	// Clean up core infrastructure
	err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, setupManifest)
	s.Require().NoError(err)

	// Clean up curl pod
	err = s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, testdefaults.CurlPodManifest)
	s.Require().NoError(err)

	s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, s.coreObjects...)
	s.testInstallation.Assertions.EventuallyPodsNotExist(s.ctx, gatewayObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", gatewayObjectMeta.GetName()),
	})
	s.testInstallation.Assertions.EventuallyPodsNotExist(s.ctx, httpbinObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin",
	})
}

// BeforeTest runs before each test
func (s *testingSuite) BeforeTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Require().NoError(err)
	}
}

// AfterTest runs after each test
func (s *testingSuite) AfterTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, manifest)
		s.Require().NoError(err)
	}
}

// TestJwtAuthentication tests the jwt is valid
func (s *testingSuite) TestJwtAuthentication() {
	// Send request to route with no JWT config applied, should get 200 OK
	s.T().Log("send request to route with no JWT config applied, should get 200 OK")
	statusReqCurlOpts := []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(gatewayService.ObjectMeta)),
		curl.WithHostHeader("httpbin"),
		curl.WithPort(8080),
		curl.WithPath("/status/200"),
	}
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		statusReqCurlOpts,
		expectStatus200Success)

	// The /get route does have a JWT config applied, should get 401 Unauthorized
	s.T().Log("The /get route does have a JWT config applied, should fail when no JWT is provided")
	getReqCurlOpts := []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(gatewayService.ObjectMeta)),
		curl.WithHostHeader("httpbin"),
		curl.WithPort(8080),
		curl.WithPath("/get"),
	}
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		getReqCurlOpts,
		expectedJwtMissingFailedResponse)

	// correct JWT is used should result in 200 OK
	s.T().Log("The /get route does have a JWT config applied, should fail when incorrect JWT is provided")
	getReqBadJwtCurlOpts := append(getReqCurlOpts, curl.WithHeader("Authorization", "Bearer "+badJwtToken))
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		getReqBadJwtCurlOpts,
		expectedJwtVerificationFailedResponse,
	)

	// correct JWT is used should result in 200 OK
	s.T().Log("The /get route does have a JWT config applied, should succeed when correct JWT is provided")
	getReqJwtCurlOpts := append(getReqCurlOpts, curl.WithHeader("Authorization", "Bearer "+dev1JwtToken))
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		getReqJwtCurlOpts,
		expectStatus200Success,
	)
}

// TestJwtAuthentication tests the jwt claims have permissions
func (s *testingSuite) TestJwtAuthorization() {
	getReqCurlOpts := []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(gatewayService.ObjectMeta)),
		curl.WithHostHeader("httpbin"),
		curl.WithPort(8080),
		curl.WithPath("/get"),
	}

	// correct JWT, but incorrect claims should be denied
	s.T().Log("The /get route has a JWT applies at the route level, should fail when correct JWT is provided but incorrect claims")
	getReqDev1JwtCurlOpts := append(getReqCurlOpts, curl.WithHeader("Authorization", "Bearer "+dev1JwtToken))
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		getReqDev1JwtCurlOpts,
		expectRbacDeniedWithJwt,
	)
	// correct JWT is used should result in 200 OK
	s.T().Log("The /get route has a JWT applies at the route level, should succeed when correct JWT is provided with correct claims")
	getReqDev2JwtCurlOpts := append(getReqCurlOpts, curl.WithHeader("Authorization", "Bearer "+dev2JwtToken))
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		getReqDev2JwtCurlOpts,
		expectStatus200Success,
	)
}
