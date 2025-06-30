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
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/test/util/retry"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/agentgatewaysyncer"

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

func newAgentGatewayXdsDumper(t *testing.T, ctx context.Context, xdsPort int, gwname, gwnamespace string) xdsDumper {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", xdsPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(time.Second*10),
	)
	if err != nil {
		t.Fatalf("failed to connect to xds server: %v", err)
	}

	d := xdsDumper{
		conn: conn,
		dr: &discovery_v3.DiscoveryRequest{
			Node: &envoycore.Node{
				Id: "gateway.gwtest",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"role": structpb.NewStringValue(fmt.Sprintf("%s~%s", gwnamespace, gwname)),
					},
				},
			},
		},
	}

	ads := discovery_v3.NewAggregatedDiscoveryServiceClient(d.conn)
	ctx, cancel := context.WithTimeout(ctx, time.Second*30) // long timeout - just in case. we should never reach it.
	adsClient, err := ads.StreamAggregatedResources(ctx)
	if err != nil {
		t.Fatalf("failed to get ads client: %v", err)
	}
	d.adsClient = adsClient
	d.cancel = cancel

	return d
}

type agentGwDump struct {
	Resources []*api.Resource
	Addresses []*api.Address
}

func (x xdsDumper) DumpAgentGateway(t *testing.T, ctx context.Context) agentGwDump {
	// get resources
	resources := x.GetResources(t, ctx)
	addresses := x.GetAddress(t, ctx)

	return agentGwDump{
		Resources: resources,
		Addresses: addresses,
	}
}

func (x xdsDumper) GetResources(t *testing.T, ctx context.Context) []*api.Resource {
	dr := proto.Clone(x.dr).(*discovery_v3.DiscoveryRequest)
	dr.TypeUrl = agentgatewaysyncer.TargetTypeResourceUrl
	x.adsClient.Send(dr)
	var resources []*api.Resource
	// run this in parallel with a 5s timeout
	done := make(chan struct{})
	go func() {
		defer close(done)
		sent := 1
		for i := 0; i < sent; i++ {
			dresp, err := x.adsClient.Recv()
			if err != nil {
				t.Errorf("failed to get response from xds server: %v", err)
			}
			t.Logf("got response: %s len: %d", dresp.GetTypeUrl(), len(dresp.GetResources()))
			if dresp.GetTypeUrl() == agentgatewaysyncer.TargetTypeResourceUrl {
				for _, anyResource := range dresp.GetResources() {
					var resource api.Resource
					if err := anyResource.UnmarshalTo(&resource); err != nil {
						t.Errorf("failed to unmarshal resource: %v", err)
					}
					resources = append(resources, &resource)
				}
			}
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		// don't fatal yet as we want to dump the state while still connected
		t.Error("timed out waiting for resources for agent gateway xds dump")
		return nil
	}
	if len(resources) == 0 {
		t.Error("no resources found")
		return nil
	}
	t.Logf("xds: found %d resources", len(resources))
	return resources
}

func (x xdsDumper) GetAddress(t *testing.T, ctx context.Context) []*api.Address {
	dr := proto.Clone(x.dr).(*discovery_v3.DiscoveryRequest)
	dr.TypeUrl = agentgatewaysyncer.TargetTypeAddressUrl
	x.adsClient.Send(dr)
	var address []*api.Address
	// run this in parallel with a 5s timeout
	done := make(chan struct{})
	go func() {
		defer close(done)
		sent := 1
		for i := 0; i < sent; i++ {
			dresp, err := x.adsClient.Recv()
			if err != nil {
				t.Errorf("failed to get response from xds server: %v", err)
			}
			t.Logf("got address response: %s len: %d", dresp.GetTypeUrl(), len(dresp.GetResources()))
			if dresp.GetTypeUrl() == agentgatewaysyncer.TargetTypeAddressUrl {
				for _, anyResource := range dresp.GetResources() {
					var resource api.Address
					if err := anyResource.UnmarshalTo(&resource); err != nil {
						t.Errorf("failed to unmarshal resource: %v", err)
					}
					address = append(address, &resource)
				}
			}
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		// don't fatal yet as we want to dump the state while still connected
		t.Error("timed out waiting for address resources for agent gateway xds dump")
		return nil
	}
	if len(address) == 0 {
		t.Error("no address resources found")
		return nil
	}
	t.Logf("xds: found %d address resources", len(address))
	return address
}
