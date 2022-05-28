FROM node:16-alpine

RUN \
  apk add --no-cache --no-progress git git-lfs openssh && \
  yarn global add semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker @google/semantic-release-replace-plugin && \
  # https://github.com/npm/cli/issues/4605
  npm i -g npm@8.4


COPY ./dist/pipe /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  node --version && \
  npm --version && \
  yarn --version && \
  pipe --help
