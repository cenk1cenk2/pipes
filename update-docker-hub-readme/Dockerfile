FROM alpine:latest

COPY ./dist/update-docker-hub-readme /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
