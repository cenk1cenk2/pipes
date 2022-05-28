ARG NODE_VERSION

FROM node:${NODE_VERSION}

RUN yarn global add tailored-gitlab-env

COPY ./dist/pipe /usr/bin/pipe

RUN unset NPM_CONFIG_USER

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
