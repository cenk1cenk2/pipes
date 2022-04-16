FROM node:16-alpine

RUN \
  apk add --no-cache --no-progress git && \
  yarn global add semantic-release @qiwi/multi-semantic-release

COPY dist/semantic-release /pipe

RUN chmod +x /pipe
