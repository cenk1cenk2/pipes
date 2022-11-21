# pipe-docker

Builds and publishes Docker images from CI/CD.

`pipe-docker [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### Docker

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USE_BUILDKIT` | Use Docker build kit for building images. | `Bool` | `false` | false |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for Docker BuildX. | `String` | `false` | linux/amd64 |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `String` | `false` | CI |

### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

### Image

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$IMAGE_NAME`<br/>`$DOCKER_IMAGE_NAME` | Image name for will be built Docker image. | `String` | `true` |  |
| `$IMAGE_TAGS`<br/>`$DOCKER_IMAGE_TAGS` | Image tag for will be built Docker image. | `StringSlice` | `true` |  |
| `$DOCKERFILE_CONTEXT` | Dockerfile context argument for build operation. | `String` | `false` | . |
| `$DOCKERFILE_NAME` | Dockerfile path for the build operation | `String` | `false` | Dockerfile |
| `$IMAGE_TAG_AS_LATEST`<br/>`$DOCKER_IMAGE_TAG_AS_LATEST` | Regex pattern to tag the image as latest.<br />Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json(RegExp[])` | `false` | [ "^tags/v?\\d+.\\d+.\\d+$" ] |
| `$IMAGE_SANITIZE_TAGS`<br/>`$DOCKER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />Template is interpolated with the given matches in the regular expression. | `String`<br/>`json([]struct{ match: RegExp, template: Template(map[string]string) })` | `false` | [<br />  { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }<br />] |
| `$IMAGE_INSPECT`<br/>`$DOCKER_IMAGE_INSPECT` | Inspect after pushing the image. | `Bool` | `false` | false |
| `$BUILD_ARGS`<br/>`$DOCKER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image. | `StringSlice` | `false` |  |
| `$IMAGE_PULL`<br/>`$DOCKER_IMAGE_PULL` | Pull before building the image. | `Bool` | `false` | false |

### Registry

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url for logging in. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `String` | `false` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `String` | `false` |  |

### Tags File

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `String` | `false` |  |
| `$TAGS_FILE_STRICT` | Strict mode does not tolorate the missing tags file. | `Bool` | `false` | false |
