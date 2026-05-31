# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

ARG BUILDOS
ARG BUILDARCH

ENV PNPM_HOME="/root/.local/share/pnpm"
ENV SEMANTIC_RELEASE_HOME="/opt/semantic-release"
ENV PATH="${PATH}:${PNPM_HOME}/bin:${SEMANTIC_RELEASE_HOME}/node_modules/.bin"

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  corepack enable && corepack prepare pnpm@latest --activate && \
  mkdir -p ${SEMANTIC_RELEASE_HOME} && \
  cd ${SEMANTIC_RELEASE_HOME} && \
  pnpm add --prod semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker semantic-release-replace-plugin && \
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
