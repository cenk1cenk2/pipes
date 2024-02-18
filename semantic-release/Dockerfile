# syntax=docker/dockerfile-upstream:master-labs
FROM node:20-alpine

ARG BUILDOS
ARG BUILDARCH

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  corepack enable pnpm && pnpm setup && \
  pnpm add -g semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker semantic-release-replace-plugin && \
  pnpm store prune && \
  # # https://github.com/npm/cli/issues/4605
  # npm i -g npm@8.4
  npm -g config set workspaces-update false

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  node --version && \
  npm --version && \
  yarn --version && \
  pnpm --version && \
  pipe --help
