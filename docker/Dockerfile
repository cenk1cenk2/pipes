FROM docker:latest

COPY ./dist/gitlab-pipes-docker /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help
