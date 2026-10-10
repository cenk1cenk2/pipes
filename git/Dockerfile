# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/alpine:3

ARG BUILDOS
ARG BUILDARCH

# the sync clones partially and sparsely, which takes the git alpine ships rather than
# a pinned older one.
RUN apk add --no-cache --no-progress git

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
