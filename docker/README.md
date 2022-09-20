# pipes-docker

Builds and publishes Docker images from CI/CD.

`pipes-docker [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $$DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $$LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --git.branch, $$CI_COMMIT_REF_NAME | Source control management branch. |  String  | false |  |
| --git.tag, $$CI_COMMIT_TAG | Source control management tag. |  String  | false |  |
| --docker_image.name, $$IMAGE_NAME | Image name for the to be built Docker image. |  String  | true |  |
| --docker_image.tags, $$IMAGE_TAGS | Image tag for the to be built Docker image. |  StringSlice  | true |  |
| --docker_file.context, $$DOCKERFILE_CONTEXT | Context for Dockerfile. |  String  | false | &#34;.&#34; |
| --docker_file.name, $$DOCKERFILE_NAME | Dockerfile name to build from. |  String  | false | &#34;Dockerfile&#34; |
| --docker_registry.registry, $$DOCKER_REGISTRY | Docker registry to login to. |  String  | false |  |
| --docker_registry.username, $$DOCKER_REGISTRY_USERNAME | Docker registry username. |  String  | false |  |
| --docker_registry.password, $$DOCKER_REGISTRY_PASSWORD | Docker registry password. |  String  | false |  |
| --docker.use_buildx, $$DOCKER_USE_BUILDX | Use docker buildx builder. |  Bool  | false | false |
| --docker.buildx_platforms, $$DOCKER_BUILDX_PLATFORMS | Platform arguments for docker buildx. |  String  | false | &#34;linux/amd64&#34; |
| --docker_image.tag_as_latest_for_tags_regex, $$TAG_AS_LATEST_FOR_TAGS_REGEX | Regex pattern to tag the image as latest. (format: json(string[])) |  String  | false | &#34;[\&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*$\&#34;]&#34; |
| --docker_image.tag_as_latest_for_branches_regex, $$TAG_AS_LATEST_FOR_BRANCHES_REGEX | Regex pattern to tag the image as latest. (format: json(string[])) |  String  | false | &#34;[]&#34; |
| --docker_image.pull, $$IMAGE_PULL | Pull while building the image. |  Bool  | false | true |
| --docker_image.tags_file, $$TAGS_FILE | Read tags from a file. |  String  | false |  |
| --docker_image.tags_file_ignore_missing, $$TAGS_FILE_IGNORE_MISSING | Dont finish the task if tags file is set and missing. |  Bool  | false | false |
| --docker_image.inspect, $$IMAGE_INSPECT | Inspect after pushing the image. |  Bool  | false | true |
| --docker_image.build_args, $$BUILD_ARGS | Pass in extra build arguments for image. |  StringSlice  | false |  |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
