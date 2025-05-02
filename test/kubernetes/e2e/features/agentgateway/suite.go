package agentgateway

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	"github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

var (
	manifests = []string{
		agentgatewayManifest,
		a2aAgentManifest,
		mcpManifest,
		gatewayManifest,
		defaults.CurlPodManifest,
	}

	proxyObjMeta = metav1.ObjectMeta{
		Name:      "agent-gateway",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: proxyObjMeta}
	proxyService    = &corev1.Service{ObjectMeta: proxyObjMeta}
)

type testingSuite struct {
	suite.Suite
	ctx              context.Context
	testInstallation *e2e.TestInstallation
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) TestAgentGateway() {
	s.T().Cleanup(func() {
		for _, manifest := range manifests {
			err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, manifest)
			s.Require().NoError(err)
		}
		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
	})

	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Require().NoError(err)
	}

	// assert the expected resources are created and running before attempting to send traffic
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)
	// check curl pod is running
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, defaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=curl",
	})
	// Check agentgateway, a2a-agent and mcp-tool are running
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, "default", metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=a2a-agent",
	})
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, "default", metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=mcp-tool",
	})
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, proxyObjMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=agentgateway",
	})

	// Check MCP SSE endpoint is reachable through the Agent Gateway
	// curl -v http://localhost:8080/sse
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(a2aPort),
			curl.WithPath("/sse"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		})

	// Check A2A Agent endpoint is reachable through the Agent Gateway
	/*
		curl -X POST http://localhost:9090/default-a2a-agent \
		  -H "Content-Type: application/json" \
		  -v \
		  -d '{
		    "jsonrpc": "2.0",
		    "id": "1",
		    "method": "tasks/send",
		    "params": {
		      "id": "1",
		      "message": {
		        "role": "user",
		        "parts": [
		          {
		            "type": "text",
		            "text": "hello gateway!"
		          }
		        ]
		      }
		    }
		  }'
	*/
	data := `{"jsonrpc":"2.0","id":"1","method":"tasks/send","params":{"id":"1","message":{"role":"user","parts":[{"type":"text","text":"hello gateway!"}]}}}`
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(a2aPort),
			curl.WithPath("/default-a2a-agent"),
			curl.WithContentType("application/json"),
			curl.WithMethod(http.MethodPost),
			curl.WithBody(data),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
		})
}
