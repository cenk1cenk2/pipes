# syntax=docker/dockerfile-upstream:master-labs
ARG BUILDAH_VERSION=latest

FROM quay.io/buildah/stable:${BUILDAH_VERSION}

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
