# syntax=docker/dockerfile-upstream:master-labs
FROM quay.io/buildah/stable:latest

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
