# Based on Istio's build-tools image: https://github.com/istio/tools/tree/228d21452fd640bf7389d7d41fecaa715ce73249/docker/build-tools 
set -euo pipefail

# Minimal entrypoint for CI environments that want a running dockerd inside the container.
# Codespaces generally uses docker-outside-of-docker; this file exists for parity with Istio.

if command -v dockerd >/dev/null 2>&1; then
  dockerd --host=unix:///var/run/docker.sock &
fi

exec "$@"


