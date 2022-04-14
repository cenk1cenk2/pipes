ARG NODE_VERSION

FROM node:${NODE_VERSION}

COPY ./dist/gitlab-pipes-node /usr/bin/pipe

RUN unset NPM_CONFIG_USER

RUN chmod +x /usr/bin/pipe && \
  yarn global add tailored-gitlab-env && \
  # smoke test
  pipe --help
