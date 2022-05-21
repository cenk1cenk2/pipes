FROM alpine:latest

COPY ./dist/template /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
