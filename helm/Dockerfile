# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/alpine/helm:4

# renovate: datasource=github-releases depName=helm-unittest/helm-unittest
ARG HELM_UNITTEST_VERSION=v1.2.1

ARG BUILDOS
ARG BUILDARCH

RUN helm plugin install https://github.com/helm-unittest/helm-unittest --version ${HELM_UNITTEST_VERSION} --verify=false

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
