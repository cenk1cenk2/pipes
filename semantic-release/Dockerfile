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
  (command -v corepack > /dev/null 2>&1 || npm install --global corepack@latest) && \
  corepack enable && corepack prepare pnpm@latest --activate && \
  mkdir -p ${SEMANTIC_RELEASE_HOME}

# the plugin set is a committed manifest instead of a bare `pnpm add` so the versions
# baked into the image are reproducible and renovate can raise the bumps.
COPY ./package.json ./pnpm-lock.yaml ${SEMANTIC_RELEASE_HOME}/

RUN \
  cd ${SEMANTIC_RELEASE_HOME} && \
  pnpm install --prod --frozen-lockfile && \
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
  multi-semantic-release --version && \
  pipe --help
