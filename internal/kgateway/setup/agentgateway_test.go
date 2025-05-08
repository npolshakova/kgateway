package setup_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/settings"
)

func TestAgentGateway(t *testing.T) {
	st, err := settings.BuildSettings()
	st.EnableAgentGateway = true

	if err != nil {
		t.Fatalf("can't get settings %v", err)
	}
	setupEnvTestAndRun(t, st, func(t *testing.T, ctx context.Context, kdbg *krt.DebugHandler, client istiokube.CLIClient, xdsPort int) {
		client.Kube().CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "gwtest"}}, metav1.CreateOptions{})

		err = client.ApplyYAMLContents("gwtest", `
apiVersion: v1
kind: Service
metadata:
  name: mcp
  namespace: gwtest
  labels:
    app: mcp
spec:
  clusterIP: "10.0.0.11"
  ports:
    - name: http
      port: 8080
      targetPort: 8080
      appProtocol: kgateway.dev/mcp
  selector:
    app: mcp
---
apiVersion: v1
kind: Service
metadata:
  name: a2a
  namespace: gwtest
  labels:
    app: a2a
spec:
  clusterIP: "10.0.0.12"
  ports:
    - name: http
      port: 8081
      targetPort: 8081
      appProtocol: kgateway.dev/a2a
  selector:
    app: a2a
---
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: agentgateway
spec:
  controllerName: kgateway.dev/kgateway
  parametersRef:
    group: gateway.kgateway.dev
    kind: GatewayParameters
    name: kgateway
    namespace: default
---
kind: GatewayParameters
apiVersion: gateway.kgateway.dev/v1alpha1
metadata:
  name: kgateway
spec:
  selfManaged: {}
---
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: http-gw
  namespace: gwtest
spec:
  gatewayClassName: agentgateway
  listeners:
  - protocol: kgateway.dev/mcp
    port: 8080
    name: mcp
    allowedRoutes:
      namespaces:
        from: All
  - protocol: kgateway.dev/a2a
    port: 8081
    name: a2a
    allowedRoutes:
      namespaces:
        from: All
`)

		if err != nil {
			t.Fatalf("failed to apply yamls: %v", err)
		}

		time.Sleep(time.Second / 2)

		dumper := newAgentGatewayXdsDumper(t, ctx, xdsPort, "http-gw", "gwtest")
		t.Cleanup(dumper.Close)
		t.Cleanup(func() {
			if t.Failed() {
				logKrtState(t, fmt.Sprintf("krt state for failed test: %s", t.Name()), kdbg)
			} else if os.Getenv("KGW_DUMP_KRT_ON_SUCCESS") == "true" {
				logKrtState(t, fmt.Sprintf("krt state for successful test: %s", t.Name()), kdbg)
			}
		})

		dump := dumper.DumpAgentGateway(t, ctx)
		if len(dump.McpTargets) != 1 {
			t.Fatalf("expected 1 mcp target config, got %d", len(dump.McpTargets))
		}
		if len(dump.A2ATargets) != 1 {
			t.Fatalf("expected 1 a2a target config, got %d", len(dump.A2ATargets))
		}
		if len(dump.Listeners) != 2 {
			t.Fatalf("expected 2 listener config, got %d", len(dump.Listeners))
		}
		t.Logf("%s finished", t.Name())
	})
}
