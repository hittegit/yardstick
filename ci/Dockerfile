# Minimal CI image for yardstick
# Base: debian:bookworm-slim
# Includes: Go 1.22, git, curl, ca-certificates

FROM debian:bookworm-slim

ARG GO_VERSION=1.22.0
ENV PATH=/usr/local/go/bin:${PATH}

RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        bash; \
    rm -rf /var/lib/apt/lists/*; \
    curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go.tgz; \
    tar -C /usr/local -xzf /tmp/go.tgz; \
    rm -f /tmp/go.tgz; \
    go version

WORKDIR /work

# The workflow will mount the source and run commands here
