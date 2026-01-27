## build-tools (Codespaces/devcontainer image)

This directory contains the Docker image definition for the `build-tools` devcontainer used by this repo.

It is **inspired by Istio's `build-tools` image** (from `istio/common-files`) and is intended to be
published to GitHub Container Registry (GHCR) so GitHub Codespaces can pull it quickly and reliably.

### What’s included (high level)

- Go (version matches `go.mod`)
- Rust toolchain (for `internal/envoyinit/`)
- Common build tooling: `git`, `make`, `gcc`, `jq`, `yq`, `kubectl`, `kind`, `helm`, `protoc`, `buf`
- Docker CLI (for `docker-outside-of-docker` feature)

### Building locally

```bash
docker build -t kgateway-build-tools:dev -f tools/build-tools/Dockerfile .
```


