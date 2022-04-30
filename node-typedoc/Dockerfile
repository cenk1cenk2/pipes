FROM node:16-alpine

RUN yarn global add typedoc typedoc-plugin-markdown

COPY ./dist/node-typedoc /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
