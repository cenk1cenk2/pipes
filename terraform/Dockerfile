# syntax=docker/dockerfile-upstream:master-labs
ARG TF_VERSION=latest

FROM hashicorp/terraform:${TF_VERSION}

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help

ENTRYPOINT [ ]
