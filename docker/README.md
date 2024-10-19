# pipe-docker

Builds and publishes Docker images from CI/CD.

`pipe-docker [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

## Commands

### `login`

Login to the given Docker registries.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USE_BUILDKIT`<br/>`$DOCKER_BUILDKIT` | Use Docker BuildKit for building images. | `Bool` | `false` | true |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `String` | `false` | CI |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `String` | `false` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `String` | `false` |  |

### `build`

Build Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USE_BUILDKIT`<br/>`$DOCKER_BUILDKIT` | Use Docker BuildKit for building images. | `Bool` | `false` | true |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `String` | `false` | CI |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for Docker BuildX. | `String` | `false` | linux/amd64 |

**Docker Image**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME` | Image name for the will be built Docker image. | `String` | `true` |  |
| `$DOCKER_IMAGE_TAGS` | Image tag for the will be built Docker image. | `StringSlice` | `true` |  |
| `$DOCKERFILE_CONTEXT` | Dockerfile context argument for build operation. | `String` | `false` | . |
| `$DOCKERFILE_NAME` | Dockerfile path for the build operation | `String` | `false` | Dockerfile |
| `$DOCKER_IMAGE_TAG_AS_LATEST` | Regex pattern to tag the image as latest.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json(RegExp[])` | `false` | [ "^tags/v?\\d+.\\d+.\\d+$" ] |
| `$DOCKER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />      Template is interpolated with the given matches in the regular expression. | `String`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | [<br />  { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }<br />] |
| `$DOCKER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br />      Template is interpolated with the given matches in the regular expression. | `String`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | [] |
| `$DOCKER_IMAGE_INSPECT` | Inspect after pushing the image. | `Bool` | `false` | true |
| `$DOCKER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br />      You can use it as a template with environment variables as the context. | `StringSlice`<br/>`format(map[string]Template[string]())` | `false` |  |
| `$DOCKER_IMAGE_PULL` | Pull before building the image. | `Bool` | `false` | true |
| `$DOCKER_MANIFEST_OUTPUT_FILE` | Write all the images that are published in to a file for later use. | `String`<br/>`format(Template[string]([]string))` | `false` | .published-docker-images_{{ $ | join "," | sha256sum }} |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `String`<br/>`format(Template[string]([]string))` | `false` |  |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `String` | `false` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `String` | `false` |  |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `String` | `false` |  |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `Bool` | `false` | false |

### `manifest`

Update manifests of the Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USE_BUILDKIT`<br/>`$DOCKER_BUILDKIT` | Use Docker BuildKit for building images. | `Bool` | `false` | true |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `Bool` | `false` | false |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `String` | `false` | CI |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_FILES` | Read published tags from a file. | `StringSlice`<br/>`format(glob)` | `false` | "**/.published-docker-images*" |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `String`<br/>`format(Template[string]())` | `false` |  |
| `$DOCKER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `StringSlice` | `false` |  |
| `$DOCKER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `String`<br/>`json([]struct { target: string, images: []string })` | `false` |  |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `String` | `false` |  |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `String` | `false` |  |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `String` | `false` |  |
