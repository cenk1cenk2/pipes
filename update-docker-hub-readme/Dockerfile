# syntax=docker/dockerfile-upstream:master-labs
FROM alpine:latest

ARG BUILDOS=linux
ARG BUILDARCH=amd64

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN pipe --help
