# syntax=docker/dockerfile-upstream:master-labs
ARG GO_VERSION

FROM golang:${GO_VERSION}

ARG BUILDOS
ARG BUILDARCH

RUN apk add --no-cache --no-progress git

# for legacy reasons
RUN \
  apk add --no-cache --no-progress git make && \
  wget -O- -nv https://taskfile.dev/install.sh | sh /dev/stdin -d -b /usr/bin && \
  wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | BINDIR=/usr/bin sh -s

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help
