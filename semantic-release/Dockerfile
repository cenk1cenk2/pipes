FROM node:16-alpine

RUN \
  apk add --no-cache --no-progress git git-lfs ssh && \
  yarn global add semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec semantic-release-docker @google/semantic-release-replace-plugin

COPY ./dist/semantic-release /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
