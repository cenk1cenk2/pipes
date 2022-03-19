FROM alpine:latest

COPY ./dist/gitlab-pipes-gl-artifacts /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  apk add --no-cache --no-progress unzip && \
  # smoke test
  pipe --help
