# syntax=docker/dockerfile-upstream:master-labs
FROM node:18-alpine

ARG BUILDOS
ARG BUILDARCH

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  yarn global add semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker semantic-release-replace-plugin && \
  # # https://github.com/npm/cli/issues/4605
  # npm i -g npm@8.4
  npm -g config set workspaces-update false && \
  yarn global add pnpm && \
  yarn cache clean


COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  node --version && \
  npm --version && \
  yarn --version && \
  pnpm --version && \
  pipe --help
