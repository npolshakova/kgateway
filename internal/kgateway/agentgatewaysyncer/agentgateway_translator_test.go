package agentgatewaysyncer

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

type translatorTestCase struct {
	inputFile     string
	outputFile    string
	assertReports AssertReports
	gwNN          types.NamespacedName
}

var _ = DescribeTable("Basic agentgateway Tests",
	func(in translatorTestCase, settingOpts ...SettingsOpts) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dir := fsutils.MustGetThisDir()

		inputFiles := []string{filepath.Join(dir, "testdata/inputs/", in.inputFile)}
		expectedProxyFile := filepath.Join(dir, "testdata/outputs/", in.outputFile)
		TestTranslation(GinkgoT(), ctx, inputFiles, expectedProxyFile, in.gwNN, in.assertReports, settingOpts...)
	},
	Entry(
		"http gateway with basic routing",
		translatorTestCase{
			inputFile:  "http-routing",
			outputFile: "http-routing-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "agentgateway-example",
			},
		}),
)
