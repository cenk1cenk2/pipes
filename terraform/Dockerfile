# syntax=docker/dockerfile-upstream:master-labs
FROM hashicorp/terraform:latest

COPY --chmod=777 ./dist/pipe /usr/bin/pipe

RUN pipe --help
