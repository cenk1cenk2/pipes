# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

ARG BUILDOS
ARG BUILDARCH

RUN \
  unset NPM_CONFIG_USER

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
