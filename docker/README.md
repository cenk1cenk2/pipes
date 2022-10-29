# pipes-docker

Builds and publishes Docker images from CI/CD.

`pipes-docker [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |

### Docker

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USE_BUILDKIT` | Use Docker build kit for building images. | `Bool` | `false` | false |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for Docker BuildX. | `String` | `false` |  |

### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

### Image

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$IMAGE_NAME` | Image name for will be built Docker image. | `String` | `true` |  |
| `$IMAGE_TAGS` | Image tag for will be built Docker image. | `StringSlice` | `true` |  |
| `$DOCKERFILE_CONTEXT` | Dockerfile context argument for build operation. | `String` | `false` |  |
| `$DOCKERFILE_NAME` | Dockerfile path for the build operation | `String` | `false` |  |
| `$IMAGE_TAG_AS_LATEST` | Regex pattern to tag the image as latest. Use either &#34;heads/&#34; for narrowing the search to branches or &#34;tags/&#34; for narrowing the search to tags.  | `String`<br/>format(json(RegExp[])) | `false` |  |
| `$IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name. Template is interpolated with the given matches in the regular expression.  | `String`<br/>format(json(map[RegExp]Template[[]string])) | `false` |  |
| `$IMAGE_PULL` | Pull before building the image. | `Bool` | `false` | false |
| `$TAGS_FILE` | Read tags from a file. | `String` | `false` |  |
| `$TAGS_FILE_IGNORE_MISSING` | Ignore the missing tags file and contunie operation as expected in that case. | `Bool` | `false` | false |
| `$IMAGE_INSPECT` | Inspect after pushing the image. | `Bool` | `false` | false |
| `$BUILD_ARGS` | Pass in extra build arguments for image. | `StringSlice` | `false` |  |

### Registry

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url for logging in. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `String` | `true` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `String` | `true` |  |
