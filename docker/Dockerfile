# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/library/docker:latest

ARG BUILDOS
ARG BUILDARCH

RUN apk add --no-cache --no-progress git

ENV DOCKER_API_VERSION=1.41

COPY --from=docker.io/docker/buildx-bin:latest /buildx /usr/libexec/docker/cli-plugins/docker-buildx

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
