# syntax=docker/dockerfile-upstream:master-labs
FROM alpine:latest

COPY --chmod=777 ./dist/pipe /usr/bin/pipe

RUN pipe --help
