# syntax=docker/dockerfile-upstream:master-labs
ARG BASE_IMAGE

FROM ${BASE_IMAGE}

# renovate: datasource=github-releases depName=jdx/mise
ARG MISE_VERSION=v2026.9.15

ARG BUILDOS
ARG BUILDARCH

# the base decides the libc of every tool mise installs, so the published variants
# have to match the libc of the jobs that consume the cache.
RUN \
  if command -v apk > /dev/null 2>&1; then \
  apk add --no-cache --no-progress ca-certificates curl git; \
  else \
  apt-get update && \
  apt-get install -y --no-install-recommends ca-certificates curl git && \
  rm -rf /var/lib/apt/lists/*; \
  fi && \
  curl -fsSL https://mise.run | MISE_VERSION=${MISE_VERSION} MISE_INSTALL_PATH=/usr/local/bin/mise sh

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN mise --version && pipe --help
