# syntax=docker/dockerfile-upstream:master-labs
FROM docker:latest

COPY --from=docker/buildx-bin:latest /buildx /usr/libexec/docker/cli-plugins/docker-buildx

COPY --chmod=777 ./dist/pipe /usr/bin/pipe

RUN pipe --help
