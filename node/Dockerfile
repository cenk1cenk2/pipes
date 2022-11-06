ARG NODE_VERSION

FROM node:${NODE_VERSION}

COPY ./dist/pipe /usr/bin/pipe

RUN \
  unset NPM_CONFIG_USER && \
  npm config set unsafe-perm true

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
