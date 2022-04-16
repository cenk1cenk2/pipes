FROM node:16-alpine

RUN \
  apk add --no-cache --no-progress git && \
  yarn global add semantic-release @qiwi/multi-semantic-release @semantic-release/github @semantic-release/npm @semantic-release/git @semantic-release/changelog @semantic-release/gitlab @semantic-release/exec

COPY dist/semantic-release /pipe

RUN chmod +x /pipe
