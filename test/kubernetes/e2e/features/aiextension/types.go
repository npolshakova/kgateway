package aiextension

import (
	"path/filepath"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

var (
	// common setup manifest (Gateway and Curl pod)
	commonManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "common.yaml")

	// backends with Token Auth
	backendManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "backend-token.yaml")

	// backends with passthrough
	backendPassthroughManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "backend-passthrough.yaml")

	// routes to LLM backends
	routesBasicManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "routes-basic.yaml")

	// routes to LLM backends with extension ref
	routesWithExtensionManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "routes-with-extension-ref.yaml")

	// routes options for streaming
	routeOptionStreamingManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "streaming.yaml")
)
