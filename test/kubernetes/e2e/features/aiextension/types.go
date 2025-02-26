package aiextension

import (
	"path/filepath"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

var (
	// common setup manifest (Gateway and Curl pod)
	commonManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "common.yaml")

	// Upstreams with Token Auth
	upstreamManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "upstream-token.yaml")

	// Upstreams with passthrough
	upstreamPassthroughManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "upstream-passthrough.yaml")

	// routes to LLM backends
	routesBasicManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "routes-basic.yaml")

	// routes options for streaming
	routeOptionStreamingManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "streaming.yaml")

	// enable ratelimiting on the routes
	ratelimitManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "ratelimit.yaml")

	// prompt guard with webhook on the routes
	promptGuardWebhookManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "prompt-guard-webhook.yaml")

	// prompt guard with webhook on the routes
	promptGuardWebhookStreamingManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "prompt-guard-webhook-streaming.yaml")

	// prompt guard on the routes
	promptGuardManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "prompt-guard.yaml")

	// prompt guard (streaming response) on the routes
	promptGuardStreamingManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "prompt-guard-streaming.yaml")
)
