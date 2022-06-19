FROM node:16-alpine

RUN yarn global add markdown-toc

COPY ./dist/pipe /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
