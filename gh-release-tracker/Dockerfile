FROM alpine:latest

COPY ./dist/gh-release-tracker /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
