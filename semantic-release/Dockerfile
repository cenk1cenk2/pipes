# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

ARG BUILDOS
ARG BUILDARCH
ARG PNPM_VERSION=10.33.2

ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}:${PNPM_HOME}/bin"

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate && \
  pnpm add -g semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker semantic-release-replace-plugin && \
  pnpm store prune && \
  # # https://github.com/npm/cli/issues/4605
  # npm i -g npm@8.4
  npm -g config set workspaces-update false

RUN \
  git init --bare --initial-branch=main /tmp/semantic-release-smoke.git && \
  mkdir -p /tmp/semantic-release-smoke && \
  cd /tmp/semantic-release-smoke && \
  git init --initial-branch=main && \
  git config user.email smoke@example.test && \
  git config user.name Smoke && \
  printf "# smoke\n" > README.md && \
  git add README.md && \
  git commit -m "feat: initial" && \
  git remote add origin /tmp/semantic-release-smoke.git && \
  git push -u origin main && \
  printf '%s\n' '{"branches":["main"],"plugins":["@semantic-release/changelog"]}' > .releaserc.json && \
  semantic-release --dry-run --no-ci --branches main > /tmp/semantic-release-smoke.log 2>&1 || (cat /tmp/semantic-release-smoke.log && exit 1) && \
  cat /tmp/semantic-release-smoke.log && \
  grep -q 'Loaded plugin "verifyConditions" from "@semantic-release/changelog"' /tmp/semantic-release-smoke.log && \
  ! grep -q "Cannot find module '@semantic-release/changelog'" /tmp/semantic-release-smoke.log && \
  rm -rf /tmp/semantic-release-smoke /tmp/semantic-release-smoke.git /tmp/semantic-release-smoke.log

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  node --version && \
  npm --version && \
  yarn --version && \
  pnpm --version && \
  pipe --help
