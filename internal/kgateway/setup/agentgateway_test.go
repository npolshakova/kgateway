package setup_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/settings"
)

func TestAgentGatewaySelfManaged(t *testing.T) {
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

		// Count different types of resources
		var bindCount, listenerCount, routeCount int
		for _, resource := range dump.Resources {
			switch resource.GetKind().(type) {
			case *api.Resource_Bind:
				bindCount++
			case *api.Resource_Listener:
				listenerCount++
			case *api.Resource_Route:
				routeCount++
			}
		}

		// We expect 2 binds (one for each port), 2 listeners, and at least 2 routes (mcp and a2a)
		if bindCount != 2 {
			t.Fatalf("expected 2 bind resources, got %d", bindCount)
		}
		if listenerCount != 2 {
			t.Fatalf("expected 2 listener resources, got %d", listenerCount)
		}
		if routeCount < 2 {
			t.Fatalf("expected at least 2 route resources, got %d", routeCount)
		}

		t.Logf("%s finished", t.Name())
	})
}

func TestAgentGatewayAllowedRoutes(t *testing.T) {
	st, err := settings.BuildSettings()
	st.EnableAgentGateway = true

	if err != nil {
		t.Fatalf("can't get settings %v", err)
	}
	setupEnvTestAndRun(t, st, func(t *testing.T, ctx context.Context, kdbg *krt.DebugHandler, client istiokube.CLIClient, xdsPort int) {
		client.Kube().CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "gwtest"}}, metav1.CreateOptions{})

		err = client.ApplyYAMLContents("", `
apiVersion: v1
kind: Namespace
metadata:
  name: othernamespace
---
apiVersion: v1
kind: Service
metadata:
  name: mcp-other
  namespace: othernamespace
  labels:
    app: mcp-other
spec:
  clusterIP: "10.0.0.11"
  ports:
    - name: http
      port: 8080
      targetPort: 8080
      appProtocol: kgateway.dev/mcp
  selector:
    app: mcp-other
---
apiVersion: v1
kind: Service
metadata:
  name: a2a-other
  namespace: othernamespace
  labels:
    app: a2a-other
spec:
  clusterIP: "10.0.0.12"
  ports:
    - name: http
      port: 8081
      targetPort: 8081
      appProtocol: kgateway.dev/a2a
  selector:
    app: a2a-other
---
apiVersion: v1
kind: Service
metadata:
  name: mcp-allowed
  namespace: gwtest
  labels:
    app: mcp-allowed
spec:
  clusterIP: "10.0.0.13"
  ports:
    - name: http
      port: 8080
      targetPort: 8080
      appProtocol: kgateway.dev/mcp
  selector:
    app: mcp-allowed
---
apiVersion: v1
kind: Service
metadata:
  name: a2a-allowed
  namespace: gwtest
  labels:
    app: a2a-allowed
spec:
  clusterIP: "10.0.0.14"
  ports:
    - name: http
      port: 8081
      targetPort: 8081
      appProtocol: kgateway.dev/a2a
  selector:
    app: a2a-allowed
---
kind: GatewayParameters
apiVersion: gateway.kgateway.dev/v1alpha1
metadata:
  name: kgateway
  namespace: gwtest
spec:
  kube:
    agentGateway:
      enabled: true
---
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: agentgateway
  namespace: gwtest
spec:
  controllerName: kgateway.dev/kgateway
  parametersRef:
    group: gateway.kgateway.dev
    kind: GatewayParameters
    name: kgateway
    namespace: gwtest
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
        from: Same
  - protocol: kgateway.dev/a2a
    port: 8081
    name: a2a
    allowedRoutes:
      namespaces:
        from: Same
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

		// Count different types of resources
		var bindCount, listenerCount, routeCount int
		for _, resource := range dump.Resources {
			switch resource.GetKind().(type) {
			case *api.Resource_Bind:
				bindCount++
			case *api.Resource_Listener:
				listenerCount++
			case *api.Resource_Route:
				routeCount++
			}
		}

		// We expect 2 binds (one for each port), 2 listeners, and at least 2 routes (mcp and a2a)
		if bindCount != 2 {
			t.Fatalf("expected 2 bind resources, got %d", bindCount)
		}
		if listenerCount != 2 {
			t.Fatalf("expected 2 listener resources, got %d", listenerCount)
		}
		if routeCount < 2 {
			t.Fatalf("expected at least 2 route resources, got %d", routeCount)
		}

		// Check that routes are properly configured for the allowed services
		var mcpRoutes, a2aRoutes []*api.Route
		for _, resource := range dump.Resources {
			if route := resource.GetRoute(); route != nil {
				if strings.Contains(route.RouteName, "mcp") {
					mcpRoutes = append(mcpRoutes, route)
				} else if strings.Contains(route.RouteName, "a2a") {
					a2aRoutes = append(a2aRoutes, route)
				}
			}
		}

		// Verify that we have routes for both protocols
		if len(mcpRoutes) == 0 {
			t.Fatalf("expected at least 1 MCP route, got 0")
		}
		if len(a2aRoutes) == 0 {
			t.Fatalf("expected at least 1 A2A route, got 0")
		}

		t.Logf("%s finished", t.Name())
	})
}
