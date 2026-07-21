# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION
ARG PULUMI_VERSION=latest

FROM docker.io/pulumi/pulumi-nodejs:${PULUMI_VERSION} AS pulumi
FROM node:${NODE_VERSION}

ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}"

ARG BUILDOS
ARG BUILDARCH

COPY --from=pulumi /pulumi/bin/ /usr/bin/

RUN \
  unset NPM_CONFIG_USER && \
  (command -v corepack > /dev/null 2>&1 || npm install --global corepack@latest) && \
  corepack enable && corepack prepare pnpm@latest --activate

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
