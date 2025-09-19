# syntax=docker/dockerfile-upstream:master-labs
ARG GO_VERSION

FROM golang:${GO_VERSION}

RUN apk add --no-cache --no-progress git

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
