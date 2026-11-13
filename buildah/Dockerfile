# syntax=docker/dockerfile-upstream:master-labs
# Buildah 1.44+ requires a context overlay even with vfs in Kubernetes runners.
# TODO: @cenk1cenk2 once the fix for https://github.com/podman-container-tools/buildah/issues/6910 is released and verified on the Kubernetes runner, update this and the bootstrap image pin and remove the Renovate restriction.
ARG BUILDAH_VERSION=v1.43.5

FROM quay.io/buildah/stable:${BUILDAH_VERSION}

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
