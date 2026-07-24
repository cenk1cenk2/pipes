# syntax=docker/dockerfile-upstream:master-labs
FROM docker.io/alpine/helm:latest

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
