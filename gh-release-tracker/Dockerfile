FROM alpine:latest

COPY ./dist/pipe /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
