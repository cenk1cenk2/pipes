ARG NODE_VERSION

FROM node:${NODE_VERSION}

RUN yarn global add tailored-gitlab-env

COPY ./dist/gitlab-pipes-node /usr/bin/pipe

RUN unset NPM_CONFIG_USER

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
