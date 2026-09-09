# syntax=docker/dockerfile-upstream:master-labs
FROM registry.k8s.io/kustomize/kustomize:v5.8.1 AS kustomize
FROM docker.io/alpine/helm:4

ARG BUILDOS
ARG BUILDARCH

COPY --from=kustomize /app/kustomize /usr/bin/kustomize

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN kustomize version && pipe --help
