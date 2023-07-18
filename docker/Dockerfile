# syntax=docker/dockerfile-upstream:master-labs
FROM docker:latest

ARG BUILDOS=linux
ARG BUILDARCH=amd64

RUN apk add --no-cache --no-progress git

ENV DOCKER_API_VERSION=1.41

COPY --from=docker/buildx-bin:latest /buildx /usr/libexec/docker/cli-plugins/docker-buildx

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
