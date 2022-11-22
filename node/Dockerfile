# syntax=docker/dockerfile-upstream:master-labs
ARG NODE_VERSION

FROM node:${NODE_VERSION}

RUN \
  unset NPM_CONFIG_USER && \
  npm config set unsafe-perm true

COPY --chmod=777 ./dist/pipe /usr/bin/pipe

RUN pipe --help
