# syntax=docker/dockerfile-upstream:master-labs
FROM node:20-alpine

ARG BUILDOS
ARG BUILDARCH

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  corepack enable pnpm && \
  npm -g config set workspaces-update false && \
  yarn global add release-it @release-it/bumper @release-it/conventional-changelog @release-it-plugins/workspaces @j-ulrich/release-it-regex-bumper && \
  yarn cache clean

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  node --version && \
  npm --version && \
  yarn --version && \
  pnpm --version && \
  pipe --help
