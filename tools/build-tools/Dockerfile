# Based on Istio's build-tools image: https://github.com/istio/tools/tree/228d21452fd640bf7389d7d41fecaa715ce73249/docker/build-tools 
ARG GO_VERSION=1.25.6
FROM golang:${GO_VERSION}-bookworm AS go_context

FROM ubuntu:noble AS build_tools

ARG TARGETARCH
ARG VERSION=dev

LABEL "io.kgateway.repo"="https://github.com/kgateway-dev/kgateway"
LABEL "io.kgateway.version"="${VERSION}"

ENV DEBIAN_FRONTEND=noninteractive

# Pinned tool versions (keep aligned with `go.mod` and repo tooling where possible)
ENV RUST_VERSION=1.86.0
ENV KUBECTL_VERSION=1.35.0
ENV KIND_VERSION=v0.31.0
ENV HELM_VERSION=v3.19.2
ENV PROTOC_VERSION=33.2
ENV BUF_VERSION=v1.62.1
ENV YQ_VERSION=4.50.1
ENV GOLANGCI_LINT_VERSION=v2.8.0
ENV GOIMPORTS_VERSION=v0.40.0
ENV SU_EXEC_VERSION=0.3.1

# General env
ENV HOME=/home
ENV LANG=C.UTF-8
ENV CARGO_HOME=/home/.cargo
ENV RUSTUP_HOME=/home/.rustup

# Go env
ENV GO111MODULE=on
ENV GOTOOLCHAIN=local
ENV GOPROXY=https://proxy.golang.org
ENV GOSUMDB=sum.golang.org
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GOCACHE=/gocache
ENV GOBIN=/gobin
# Include /usr/sbin and /sbin so devcontainer features can find admin tools like groupadd.
ENV PATH=/usr/local/go/bin:/gobin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
      ca-certificates \
      curl \
      wget \
      git \
      openssh-client \
      rsync \
      jq \
      unzip \
      xz-utils \
      file \
      make \
      gcc \
      g++ \
      libc-dev \
      pkg-config \
      passwd \
      python3 \
      python3-pip \
      python3-setuptools \
      python3-yaml \
      sudo \
      gnupg \
      lsb-release \
    && rm -rf /var/lib/apt/lists/*

# Install Docker CLI + buildx plugin (Codespaces uses docker-outside-of-docker)
RUN --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    install -m 0755 -d /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc \
    && chmod a+r /etc/apt/keyrings/docker.asc \
    && echo "deb [arch=${TARGETARCH} signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
      docker-ce-cli \
      docker-buildx-plugin \
    && rm -rf /var/lib/apt/lists/*

# Go toolchain (copy from official golang image to match go.mod)
COPY --from=go_context /usr/local/go /usr/local/go

# Install protoc
RUN set -eux; \
    case "${TARGETARCH}" in \
      amd64) PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-x86_64.zip" ;; \
      arm64) PROTOC_ZIP="protoc-${PROTOC_VERSION}-linux-aarch_64.zip" ;; \
      *) echo "unsupported TARGETARCH ${TARGETARCH}"; exit 1 ;; \
    esac; \
    wget -nv -O "/tmp/${PROTOC_ZIP}" "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
    unzip -q "/tmp/${PROTOC_ZIP}" -d /tmp/protoc; \
    mv /tmp/protoc/bin/protoc /usr/local/bin/protoc; \
    chmod 555 /usr/local/bin/protoc; \
    rm -rf /tmp/protoc "/tmp/${PROTOC_ZIP}"

# kubectl
RUN set -eux; \
    curl -fsSL -o /usr/local/bin/kubectl "https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/${TARGETARCH}/kubectl"; \
    chmod 555 /usr/local/bin/kubectl

# kind
RUN set -eux; \
    curl -fsSL -o /usr/local/bin/kind "https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-linux-${TARGETARCH}"; \
    chmod 555 /usr/local/bin/kind

# helm
RUN set -eux; \
    curl -fsSL -o /tmp/helm.tgz "https://get.helm.sh/helm-${HELM_VERSION}-linux-${TARGETARCH}.tar.gz"; \
    mkdir -p /tmp/helm; \
    tar -xzf /tmp/helm.tgz -C /tmp/helm; \
    mv "/tmp/helm/linux-${TARGETARCH}/helm" /usr/local/bin/helm; \
    chmod 555 /usr/local/bin/helm; \
    rm -rf /tmp/helm /tmp/helm.tgz

# buf
RUN set -eux; \
    arch="$(uname -m)"; \
    curl -fsSL -o /usr/local/bin/buf "https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-Linux-${arch}"; \
    chmod 555 /usr/local/bin/buf

# yq
RUN set -eux; \
    curl -fsSL -o /usr/local/bin/yq "https://github.com/mikefarah/yq/releases/download/v${YQ_VERSION}/yq_linux_${TARGETARCH}"; \
    chmod 555 /usr/local/bin/yq

# Install common Go dev tools (pinned)
# Note: we set GOBIN explicitly here so we don't rely on GOPATH defaults (/root/go/bin),
# which may not match our global env (GOBIN=/gobin).
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -eux; \
    GOBIN=/usr/local/bin CGO_ENABLED=0 go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION}; \
    GOBIN=/usr/local/bin CGO_ENABLED=0 go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}; \
    chmod 555 /usr/local/bin/goimports /usr/local/bin/golangci-lint

# Rust toolchain (for internal/envoyinit)
RUN set -eux; \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain "${RUST_VERSION}" --profile minimal --component rustfmt --component clippy; \
    ln -s "${CARGO_HOME}/bin/"* /usr/local/bin/

# Install su-exec (like Istio: small, reliable privilege drop for entrypoint)
RUN set -eux; \
    curl -fsSL -o "/tmp/su-exec.tgz" "https://github.com/NobodyXu/su-exec/archive/refs/tags/v${SU_EXEC_VERSION}.tar.gz"; \
    tar -xzf /tmp/su-exec.tgz -C /tmp; \
    pushd "/tmp/su-exec-${SU_EXEC_VERSION}" >/dev/null; \
      LDFLAGS="-fvisibility=hidden -Wl,-O2 -Wl,--discard-all -Wl,--strip-all -Wl,--as-needed -Wl,--gc-sections" make; \
      cp -a su-exec /usr/local/bin/su-exec; \
      chmod u+sx /usr/local/bin/su-exec; \
    popd >/dev/null; \
    rm -rf /tmp/su-exec.tgz "/tmp/su-exec-${SU_EXEC_VERSION}"

# Note: devcontainer build contexts can vary (and may not include the full repo tree).
# To make this image robust, generate small helper scripts directly in the image
# rather than relying on COPY from the build context.
RUN cat > /usr/local/bin/docker-entrypoint <<'EOF' \
set -euo pipefail

# Copy credentials from mountpoints using su-exec
uid=$(id -u)
gid=$(id -g)

shopt -s dotglob

if [[ -d /config ]]; then
  # Make a copy of the host's config secrets. Do not copy docker sockets.
  su-exec 0:0 rsync -a --exclude=docker*.sock --exclude=/config/.config/gcloud/logs/* /config/ /config-copy/ || true

  # Set the ownership of the host's config secrets to that of the container
  su-exec 0:0 chown -R "${uid}":"${gid}" /config-copy || true

  # Permit only the UID:GID to read the copy of the host's config secrets
  chmod -R 700 /config-copy || true

  # If docker_for_mac plaintext-passwords.json exists, import it into config.json
  if [[ -f /config-copy/.docker/plaintext-passwords.json ]]; then
    auth_value=$(jq -r '.auths."https://index.docker.io/v1/".auth' /config-copy/.docker/plaintext-passwords.json)
    if [[ "${auth_value}" == "null" ]]; then
      echo "Missing docker credentials."
    fi
    encode_value=$(echo "${auth_value}" | base64 --decode | base64)
    jq --arg auth "${encode_value}" '.auths."https://index.docker.io/v1/".auth=$auth' /config-copy/.docker/config.json > /config-copy/.docker/config-tmp.json
    jq 'del(.credsStore)' /config-copy/.docker/config-tmp.json > /config-copy/.docker/config.json
  fi
fi

# Add user based upon passed UID. Skip if run as root.
if [[ "${uid}" -ne 0 ]]; then
  su-exec 0:0 useradd --uid "${uid}" --system user || true
fi

# Set ownership of /home to UID:GID
su-exec 0:0 chown "${uid}":"${gid}" /home || true

# Copy the config secrets without changing permissions nor ownership
if [[ -d /config-copy ]]; then
  cp -R /config-copy/* /home/ 2>/dev/null || true
fi

exec "$@"
EOF
RUN chmod +x /usr/local/bin/docker-entrypoint

RUN cat > /usr/local/bin/prow-entrypoint <<'EOF' \

set -euo pipefail

# Minimal entrypoint for CI environments that want a running dockerd inside the container.
# Codespaces generally uses docker-outside-of-docker; this file exists for parity with Istio.

if command -v dockerd >/dev/null 2>&1; then
  dockerd --host=unix:///var/run/docker.sock &
fi

exec "$@"
EOF
RUN chmod +x /usr/local/bin/prow-entrypoint

RUN cat > /home/.bashrc <<'EOF'
# Copyright Istio Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Override prompt to avoid the annoying "I have no name!" in the default prompt
PS1="build-tools:\w\\$ \[$(tput sgr0)\]"
export PS1

# Tab completion for make that works with monorepos.
# In Codespaces the repo is typically under /workspaces/<repo>.
root="${REPO_ROOT:-}"
if [[ -z "${root}" ]]; then
  if [[ -d /work ]]; then
    root="/work"
  elif [[ -d /workspaces ]]; then
    root="/workspaces"
  else
    root="/"
  fi
fi

make_options=$(find "${root}" -maxdepth 6 -iname "Makefile*" 2>/dev/null | xargs -r -I {} grep -hoE '^[a-zA-Z0-9_.-]+:([^=]|$)' {} | sed 's/[^a-zA-Z0-9_.-]*$//' | sort -u)
complete -W "$make_options" make
EOF

# Mountpoints for host mounts (Codespaces & local devcontainers)
RUN mkdir -p /go /gocache /gobin \
    && mkdir -p /config/.docker /config/.config/gcloud /config/.kube \
    && mkdir -p /config-copy \
    && mkdir -p /home/.cache /home/.cargo/registry /home/.cargo/git \
    && chmod -R 777 /go /gocache /gobin /config /config-copy /home/.cache /home/.cargo

WORKDIR /
ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
CMD ["bash"]


