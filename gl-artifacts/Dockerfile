# syntax=docker/dockerfile-upstream:master-labs
FROM alpine:latest

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  apk add --no-cache --no-progress unzip && \
  # smoke test
  pipe --help
