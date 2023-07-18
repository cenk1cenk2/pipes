# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

ARG BUILDOS=linux
ARG BUILDARCH=amd64

RUN \
  unset NPM_CONFIG_USER && \
  yarn global add pnpm

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
