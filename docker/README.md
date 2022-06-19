<!-- clidocs -->

# NAME

gitlab-pipes-docker - Builds and publishes Docker images from CI/CD.

# SYNOPSIS

gitlab-pipes-docker

```
[--ci]
[--debug]
[--docker.buildx_platforms]=[value]
[--docker.use_buildx]
[--docker_file.context]=[value]
[--docker_file.name]=[value]
[--docker_image.build_args]=[value]
[--docker_image.inspect]
[--docker_image.name]=[value]
[--docker_image.pull]
[--docker_image.tag_as_latest_for_branches_regex]=[value]
[--docker_image.tag_as_latest_for_tags_regex]=[value]
[--docker_image.tags]=[value]
[--docker_image.tags_file]=[value]
[--docker_registry.password]=[value]
[--docker_registry.registry]=[value]
[--docker_registry.username]=[value]
[--git.branch]=[value]
[--git.tag]=[value]
[--help|-h]
[--log-level]=[value]
[--version|-v]
```

# DESCRIPTION

Builds and publishes Docker images from CI/CD.

**Usage**:

```
gitlab-pipes-docker [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--ci**: Sets whether this is running inside a CI/CD environment.

**--debug**: Enable debugging for the application.

**--docker.buildx_platforms**="": Platform arguments for docker buildx. (default: linux/amd64)

**--docker.use_buildx**: Use docker buildx builder.

**--docker_file.context**="": Context for Dockerfile. (default: .)

**--docker_file.name**="": Dockerfile name to build from. (default: Dockerfile)

**--docker_image.build_args**="": Pass in extra build arguments for image. (default: [])

**--docker_image.inspect**: Inspect after pushing the image.

**--docker_image.name**="": Image name for the to be built Docker image.

**--docker_image.pull**: Pull while building the image.

**--docker_image.tag_as_latest_for_branches_regex**="": Regex pattern to tag the image as latest. format: json(string[]) (default: [])

**--docker_image.tag_as_latest_for_tags_regex**="": Regex pattern to tag the image as latest. format: json(string[]) (default: ["^v\\d*\\.\\d*\\.\\d*$"])

**--docker_image.tags**="": Image tag for the to be built Docker image. (default: [])

**--docker_image.tags_file**="": Read tags from a file.

**--docker_registry.password**="": Docker registry password.

**--docker_registry.registry**="": Docker registry to login to.

**--docker_registry.username**="": Docker registry username.

**--git.branch**="": Source control management branch.

**--git.tag**="": Source control management tag.

**--help, -h**: show help

**--log-level**="": Define the log level for the application. (default: info)

**--version, -v**: print the version


# COMMANDS

## help, h

Shows a list of commands or help for one command

<!-- clidocsstop -->
