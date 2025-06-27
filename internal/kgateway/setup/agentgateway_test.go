package setup_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/agentgateway/agentgateway/go/api"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/test/util/retry"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/settings"
)

func TestAgentGatewayScenarioDump(t *testing.T) {
	st, err := settings.BuildSettings()
	st.EnableAgentGateway = true

	if err != nil {
		t.Fatalf("can't get settings %v", err)
	}

	// Use the runScenario approach to test agent gateway scenarios
	runAgentGatewayScenario(t, "testdata/agentgateway", st)
}

func runAgentGatewayScenario(t *testing.T, scenarioDir string, globalSettings *settings.Settings) {
	setupEnvTestAndRun(t, globalSettings, func(t *testing.T, ctx context.Context, kdbg *krt.DebugHandler, client istiokube.CLIClient, xdsPort int) {
		// list all yamls in test data
		files, err := os.ReadDir(scenarioDir)
		if err != nil {
			t.Fatalf("failed to read dir: %v", err)
		}
		for _, f := range files {
			// run tests with the yaml files (but not -out.yaml files)
			parentT := t
			if strings.HasSuffix(f.Name(), ".yaml") && !strings.HasSuffix(f.Name(), "-out.yaml") {
				if os.Getenv("TEST_PREFIX") != "" && !strings.HasPrefix(f.Name(), os.Getenv("TEST_PREFIX")) {
					continue
				}
				fullpath := filepath.Join(scenarioDir, f.Name())
				t.Run(strings.TrimSuffix(f.Name(), ".yaml"), func(t *testing.T) {
					writer.set(t)
					t.Cleanup(func() {
						writer.set(parentT)
					})
					testAgentGatewayScenario(t, ctx, kdbg, client, xdsPort, fullpath)
				})
			}
		}
	})
}

func testAgentGatewayScenario(
	t *testing.T,
	ctx context.Context,
	kdbg *krt.DebugHandler,
	client istiokube.CLIClient,
	xdsPort int,
	f string,
) {
	fext := filepath.Ext(f)
	fpre := strings.TrimSuffix(f, fext)
	t.Logf("running agent gateway scenario for test file: %s", f)

	// read the out file
	fout := fpre + "-out" + fext
	write := false
	_, err := os.ReadFile(fout)
	// if not exist
	if os.IsNotExist(err) {
		write = true
		err = nil
	}
	if os.Getenv("REFRESH_GOLDEN") == "true" {
		write = true
	}
	if err != nil {
		t.Fatalf("failed to read file %s: %v", fout, err)
	}

	const gwname = "http-gw-for-test"
	testgwname := "http-" + filepath.Base(fpre)
	testyamlbytes, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	// change the gw name, so we could potentially run multiple tests in parallel (though currently
	// it has other issues, so we don't run them in parallel)
	testyaml := strings.ReplaceAll(string(testyamlbytes), gwname, testgwname)

	yamlfile := filepath.Join(t.TempDir(), "test.yaml")
	os.WriteFile(yamlfile, []byte(testyaml), 0o644)

	err = client.ApplyYAMLFiles("", yamlfile)

	t.Cleanup(func() {
		// always delete yamls, even if there was an error applying them; to prevent test pollution.
		err := client.DeleteYAMLFiles("", yamlfile)
		if err != nil {
			t.Fatalf("failed to delete yaml: %v", err)
		}
		t.Log("deleted yamls", t.Name())
	})

	if err != nil {
		t.Fatalf("failed to apply yaml: %v", err)
	}
	t.Log("applied yamls", t.Name())

	// wait at least a second before the first check
	// to give the CP time to process
	time.Sleep(time.Second)

	t.Cleanup(func() {
		if t.Failed() {
			logKrtState(t, fmt.Sprintf("krt state for failed test: %s", t.Name()), kdbg)
		} else if os.Getenv("KGW_DUMP_KRT_ON_SUCCESS") == "true" {
			logKrtState(t, fmt.Sprintf("krt state for successful test: %s", t.Name()), kdbg)
		}
	})

	// Use retry to wait for the agent gateway to be ready
	retry.UntilSuccessOrFail(t, func() error {
		dumper := newAgentGatewayXdsDumper(t, ctx, xdsPort, testgwname, "gwtest")
		defer dumper.Close()
		dump := dumper.DumpAgentGateway(t, ctx)
		if len(dump.Resources) == 0 {
			return fmt.Errorf("timed out waiting for agent gateway resources")
		}

		if write {
			t.Logf("writing out file")
			// Use proto dump instead of manual YAML writing
			dumpProtoToJSON(t, dump, fpre)
			return fmt.Errorf("wrote out file - nothing to test")
		}

		// Output the config dump
		t.Logf("Agent Gateway Config Dump for %s:", testgwname)
		t.Logf("Total resources: %d", len(dump.Resources))

		// Count different types of resources
		var bindCount, listenerCount, routeCount, worklodCount, serviceCount int
		for _, resource := range dump.Resources {
			switch resource.GetKind().(type) {
			case *api.Resource_Bind:
				bindCount++
				t.Logf("Bind resource: %+v", resource.GetBind())
			case *api.Resource_Listener:
				listenerCount++
				t.Logf("Listener resource: %+v", resource.GetListener())
			case *api.Resource_Route:
				routeCount++
				t.Logf("Route resource: %+v", resource.GetRoute())
			}
		}
		t.Logf("Resource counts - Binds: %d, Listeners: %d, Routes: %d", bindCount, listenerCount, routeCount)

		for _, resource := range dump.Addresses {
			switch resource.Type.(type) {
			case *api.Address_Workload:
				worklodCount++
				t.Logf("workload resource: %+v", resource.GetWorkload())
			case *api.Address_Service:
				serviceCount++
				t.Logf("service resource: %+v", resource.GetService())
			}
		}
		t.Logf("Address counts - Workload: %d, Service: %d", worklodCount, serviceCount)

		return nil
	}, retry.Converge(2), retry.BackoffDelay(2*time.Second), retry.Timeout(10*time.Second))

	t.Logf("%s finished", t.Name())
}

// dumpProtoToJSON dumps the agentgateway resources to JSON format
func dumpProtoToJSON(t *testing.T, dump agentGwDump, fpre string) {
	jsonFile := fpre + "-out.json"

	// Create a structured dump map
	dumpMap := map[string]interface{}{
		"resources": dump.Resources,
		"addresses": dump.Addresses,
	}

	// Marshal to JSON using regular JSON marshaling
	jsonData, err := json.MarshalIndent(dumpMap, "", "  ")
	if err != nil {
		t.Logf("failed to marshal to JSON: %v", err)
		return
	}

	err = os.WriteFile(jsonFile, jsonData, 0o644)
	if err != nil {
		t.Logf("failed to write JSON file: %v", err)
		return
	}

	t.Logf("wrote JSON dump to: %s", jsonFile)
}
