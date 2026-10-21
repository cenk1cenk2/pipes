# syntax=docker/dockerfile-upstream:master-labs
ARG PYTHON_VERSION

FROM ghcr.io/astral-sh/uv:0.13.4 AS uv
FROM docker.io/python:${PYTHON_VERSION}

ARG BUILDOS
ARG BUILDARCH

# the build directory of a job is a mount of its own, so uv copies out of its cache
# instead of hardlinking, and it runs on the interpreter of the image.
ENV UV_LINK_MODE=copy
ENV UV_PYTHON_DOWNLOADS=never

# uv clones git dependencies with the git of the system.
RUN \
  if command -v apk > /dev/null 2>&1; then \
  apk add --no-cache --no-progress git; \
  fi

COPY --from=uv /uv /uvx /usr/bin/

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN uv --version && pipe --help
