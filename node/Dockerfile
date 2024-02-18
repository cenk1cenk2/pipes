# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

ARG BUILDOS
ARG BUILDARCH

ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}"

RUN \
  unset NPM_CONFIG_USER && \
  corepack enable && corepack prepare pnpm --activate

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
