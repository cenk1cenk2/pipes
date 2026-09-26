# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/debian:trixie-slim

# renovate: datasource=github-releases depName=jdx/mise
ARG MISE_VERSION=v2026.9.14

ARG BUILDOS
ARG BUILDARCH

# most mise backends ship glibc builds, which is why this pipe is not on alpine.
RUN \
  apt-get update && \
  apt-get install -y --no-install-recommends ca-certificates curl git && \
  rm -rf /var/lib/apt/lists/* && \
  curl -fsSL https://mise.run | MISE_VERSION=${MISE_VERSION} MISE_INSTALL_PATH=/usr/local/bin/mise sh

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN mise --version && pipe --help
