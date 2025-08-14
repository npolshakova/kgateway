package http_listener_policy

import (
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	// Use well-known nginx resources instead of custom ones
	nginxPod    = testdefaults.NginxPod
	nginxSvc    = testdefaults.NginxSvc
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
)
