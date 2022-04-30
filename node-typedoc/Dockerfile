FROM node:16-alpine

COPY ./dist/node-typedoc /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
