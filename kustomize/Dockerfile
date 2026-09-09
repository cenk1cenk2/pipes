# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/alpine/helm:4 AS helm
FROM registry.k8s.io/kustomize/kustomize:v5.8.1 AS kustomize
FROM docker.io/alpine:3

ARG BUILDOS
ARG BUILDARCH

# kustomize shells out to git to accumulate remote bases
RUN apk add --no-cache --no-progress git

COPY --from=helm /usr/bin/helm /usr/bin/helm
COPY --from=kustomize /app/kustomize /usr/bin/kustomize

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN helm version && kustomize version && pipe --help
