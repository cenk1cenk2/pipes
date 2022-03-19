# gitlab-pipes-docker

# Description

Builds, tags and publishes a Docker image.

# Commands

```bash
NAME:
   docker-build - Builds and publishes Docker images from CI/CD.

USAGE:
   docker-build [global options] command [command options] [arguments...]

VERSION:
   latest

DESCRIPTION:
   Builds and publishes Docker images from CI/CD.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --utils.ci value                                       Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value                                    Set the log level debug for the application. [$DEBUG, $PLUGIN_DEBUG]
   --utils.log value                                      Define the log level for the application. (default: "info") [$LOG_LEVEL, $PLUGIN_LOG_LEVEL]
   --git.branch value                                     Source control management branch. [$CI_COMMIT_REF_NAME]
   --git.tag value                                        Source control management tag. [$CI_COMMIT_TAG]
   --docker_image.name value                              Image name for the to be built Docker image. [$IMAGE_NAME]
   --docker_image.tags value                              Image tag for the to be built Docker image. [$IMAGE_TAGS]
   --docker_file.context value                            Context for Dockerfile. (default: ".") [$DOCKERFILE_CONTEXT]
   --docker_file.name value                               Dockerfile name to build from. (default: "Dockerfile") [$DOCKERFILE_NAME]
   --docker_registry.registry value                       Docke registry to login to. [$DOCKER_REGISTRY]
   --docker_registry.username value                       Docker registry username. [$DOCKER_REGISTRY_USERNAME]
   --docker_registry.password value                       Docker registry password. [$DOCKER_REGISTRY_PASSWORD]
   --docker_image.tag_as_latest_for_tags_regex value      Regex pattern to tag the image as latest. format: json(string[]) [$TAG_AS_LATEST_FOR_TAGS_REGEX]
   --docker_image.tag_as_latest_for_branches_regex value  Regex pattern to tag the image as latest. format: json(string[]) [$TAG_AS_LATEST_FOR_BRANCHES_REGEX]
   --docker_image.pull                                    Pull while building the image. (default: true) [$IMAGE_PULL]
   --docker_image.inspect                                 Inspect after pushing the image. (default: true) [$IMAGE_INSPECT]
   --docker_image.build_args value                        Pass in extra build arguments for image. [$BUILD_ARGS]
   --help, -h                                             show help (default: false)
   --version, -v                                          print the version (default: false)
```
