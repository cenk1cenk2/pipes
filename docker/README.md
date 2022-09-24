# pipes-docker

Builds and publishes Docker images from CI/CD.

`pipes-docker [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DEBUG` | Enable debugging for the application. | `Bool` | `false` | false |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |
| `$CI_COMMIT_REF_NAME` | Source control management branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG` | Source control management tag. | `String` | `false` |  |
| `$IMAGE_NAME` | Image name for the to be built Docker image. | `String` | `true` |  |
| `$IMAGE_TAGS` | Image tag for the to be built Docker image. | `StringSlice` | `true` |  |
| `$DOCKERFILE_CONTEXT` | Context for Dockerfile. | `String` | `false` | &#34;.&#34; |
| `$DOCKERFILE_NAME` | Dockerfile name to build from. | `String` | `false` | &#34;Dockerfile&#34; |
| `$DOCKER_REGISTRY` | Docker registry to login to. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username. | `String` | `false` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password. | `String` | `false` |  |
| `$DOCKER_USE_BUILDX` | Use docker buildx builder. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for docker buildx. | `String` | `false` | &#34;linux/amd64&#34; |
| `$TAG_AS_LATEST_FOR_TAGS_REGEX` | Regex pattern to tag the image as latest.  | `String`<br/>format(json(string[])) | `false` | &#34;[\&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*$\&#34;]&#34; |
| `$TAG_AS_LATEST_FOR_BRANCHES_REGEX` | Regex pattern to tag the image as latest.  | `String`<br/>format(json(string[])) | `false` | &#34;[]&#34; |
| `$IMAGE_PULL` | Pull while building the image. | `Bool` | `false` | true |
| `$TAGS_FILE` | Read tags from a file. | `String` | `false` |  |
| `$TAGS_FILE_IGNORE_MISSING` | Dont finish the task if tags file is set and missing. | `Bool` | `false` | false |
| `$IMAGE_INSPECT` | Inspect after pushing the image. | `Bool` | `false` | true |
| `$BUILD_ARGS` | Pass in extra build arguments for image. | `StringSlice` | `false` |  |

## Commands

### `help` , `h`

`Shows a list of commands or help for one command`
