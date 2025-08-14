package http_listener_policy

import (
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/tests/base"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
	testdefaults "github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
)

var (
	gatewayManifest                 = filepath.Join(fsutils.MustGetThisDir(), "testdata", "gateway.yaml")
	httpRouteManifest               = filepath.Join(fsutils.MustGetThisDir(), "testdata", "httproute.yaml")
	allFieldsManifest               = filepath.Join(fsutils.MustGetThisDir(), "testdata", "http-listener-policy-all-fields.yaml")
	serverHeaderManifest            = filepath.Join(fsutils.MustGetThisDir(), "testdata", "http-listener-policy-server-header.yaml")
	preserveHttp1HeaderCaseManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "preserve-http1-header-case.yaml")

	// Gateway proxy resources created dynamically per test
	proxyObjectMeta = metav1.ObjectMeta{
		Name:      "gw",
		Namespace: "default",
	}
	proxyService    = &corev1.Service{ObjectMeta: proxyObjectMeta}
	proxyDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gw",
			Namespace: "default",
		},
	}
	echoService = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "raw-header-echo",
			Namespace: "default",
		},
	}
	echoDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "raw-header-echo",
			Namespace: "default",
		},
	}

	setup = base.TestCase{
		Manifests: []string{testdefaults.CurlPodManifest, testdefaults.NginxPodManifest},
		Resources: []client.Object{testdefaults.CurlPod, testdefaults.NginxPod},
	}

	// test cases
	testCases = map[string]base.TestCase{
		"TestHttpListenerPolicyAllFields": base.TestCase{
			Manifests: []string{gatewayManifest, httpRouteManifest, allFieldsManifest},
			Resources: []client.Object{proxyService, proxyDeployment},
		},
		"TestHttpListenerPolicyServerHeader": base.TestCase{
			Manifests: []string{gatewayManifest, httpRouteManifest, serverHeaderManifest},
			Resources: []client.Object{proxyService, proxyDeployment},
		},
		"TestPreserveHttp1HeaderCase": base.TestCase{
			Manifests: []string{gatewayManifest, preserveHttp1HeaderCaseManifest},
			Resources: []client.Object{proxyService, proxyDeployment, echoService, echoDeployment},
		},
	}
)
