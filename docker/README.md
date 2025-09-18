# pipe-docker

Builds and publishes Docker images from CI/CD.

`pipe-docker [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | `info` |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | `[]` |

## Commands

### `pipe-docker login`

Login to the given Docker registries.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br/>`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | `true` |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | `false` |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | `CI` |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | `` |

### `pipe-docker build`

Build Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br/>`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | `true` |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | `false` |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | `CI` |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for Docker BuildX. | `string` | `false` | `linux/amd64` |

**Docker Image**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME` | Image name for the will be built Docker image. | `string` | `true` | `` |
| `$DOCKER_IMAGE_TAGS` | Image tag for the will be built Docker image. | `string[]` | `true` | `[]` |
| `$DOCKERFILE_CONTEXT` | Dockerfile context argument for build operation. | `string` | `false` | `.` |
| `$DOCKERFILE_NAME` | Dockerfile path for the build operation | `string` | `false` | `Dockerfile` |
| `$DOCKER_IMAGE_TAGS_AS_LATEST` | Regex pattern to tag the image as latest.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json(RegExp[])` | `false` | `[ "^tags/v?\\d+.\\d+.\\d+$" ]` |
| `$DOCKER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />      Template is interpolated with the given matches in the regular expression. | `string`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]` |
| `$DOCKER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br />      Template is interpolated with the given matches in the regular expression. | `string`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | `[]` |
| `$DOCKER_IMAGE_INSPECT` | Inspect after pushing the image. | `bool` | `false` | `true` |
| `$DOCKER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br />      You can use it as a template with environment variables as the context. | `string[]`<br/>`format(map[string]Template[string]())` | `false` | `[]` |
| `$DOCKER_IMAGE_PULL` | Pull before building the image. | `bool` | `false` | `true` |
| `$DOCKER_MANIFEST_OUTPUT_FILE` | Write all the images that are published in to a file for later use. | `string`<br/>`format(Template[string]([]string))` | `false` | `.published-docker-images_{{ $ | join "," | sha256sum }}` |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template[string]([]string))` | `false` | `` |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | `` |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | `` |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | `` |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `string` | `false` | `` |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` | `false` |

### `pipe-docker manifest`

Update manifests of the Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br/>`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | `true` |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | `false` |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | `CI` |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_FILES` | Read published tags from a file. | `string[]`<br/>`format(glob)` | `false` | `[**/.published-docker-images*]` |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template[string]())` | `false` | `` |
| `$DOCKER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `string[]` | `false` | `[]` |
| `$DOCKER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `string`<br/>`json([]struct { target: string, images: []string })` | `false` | `` |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | `` |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | `` |
