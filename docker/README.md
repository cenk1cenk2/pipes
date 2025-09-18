# pipe-docker

Builds and publishes Docker images from CI/CD.

`pipe-docker [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>info</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code>[]</code> |

## Commands

### `pipe-docker login`

Login to the given Docker registries.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br />`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | <code>true</code> |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | <code>false</code> |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | <code>CI</code> |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | <code></code> |

### `pipe-docker build`

Build Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br />`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | <code>true</code> |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | <code>false</code> |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | <code>CI</code> |
| `$DOCKER_BUILDX_PLATFORMS` | Platform arguments for Docker BuildX. | `string` | `false` | <code>linux/amd64</code> |

**Docker Image**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME` | Image name for the will be built Docker image. | `string` | `true` | <code></code> |
| `$DOCKER_IMAGE_TAGS` | Image tag for the will be built Docker image. | `string[]` | `true` | <code>[]</code> |
| `$DOCKERFILE_CONTEXT` | Dockerfile context argument for build operation. | `string` | `false` | <code>.</code> |
| `$DOCKERFILE_NAME` | Dockerfile path for the build operation | `string` | `false` | <code>Dockerfile</code> |
| `$DOCKER_IMAGE_TAGS_AS_LATEST` | Regex pattern to tag the image as latest.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json(RegExp[])` | `false` | <code>[ "^tags/v?\\d+.\\d+.\\d+$" ]</code> |
| `$DOCKER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />      Template is interpolated with the given matches in the regular expression. | `string`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | <code>[<br />    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }<br />]</code> |
| `$DOCKER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br />      Template is interpolated with the given matches in the regular expression. | `string`<br/>`json([]struct { match: RegExp, template: Template[string](RegExpMatch) })` | `false` | <code>[]</code> |
| `$DOCKER_IMAGE_INSPECT` | Inspect after pushing the image. | `bool` | `false` | <code>true</code> |
| `$DOCKER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br />      You can use it as a template with environment variables as the context. | `string[]`<br/>`format(map[string]Template[string]())` | `false` | <code>[]</code> |
| `$DOCKER_IMAGE_PULL` | Pull before building the image. | `bool` | `false` | <code>true</code> |
| `$DOCKER_MANIFEST_OUTPUT_FILE` | Write all the images that are published in to a file for later use. | `string`<br/>`format(Template[string]([]string))` | `false` | <code>.published-docker-images_{{ $ | join "," | sha256sum }}</code> |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template[string]([]string))` | `false` | <code></code> |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | <code></code> |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `string` | `false` | <code></code> |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` | <code>false</code> |

### `pipe-docker manifest`

Update manifests of the Docker images.

#### Flags

**Docker**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_BUILDKIT`<br />`$DOCKER_USE_BUILDKIT` | Use Docker BuildKit for building images. | `bool` | `false` | <code>true</code> |
| `$DOCKER_USE_BUILDX` | Use Docker BuildX builder for multi-platform builds. | `bool` | `false` | <code>false</code> |
| `$DOCKER_BUILDX_INSTANCE` | Docker BuildX instance to be started or to use. | `string` | `false` | <code>CI</code> |

**Docker Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_MANIFEST_FILES` | Read published tags from a file. | `string[]`<br/>`format(glob)` | `false` | <code>[**/.published-docker-images*]</code> |
| `$DOCKER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template[string]())` | `false` | <code></code> |
| `$DOCKER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `string[]` | `false` | <code>[]</code> |
| `$DOCKER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `string`<br/>`json([]struct { target: string, images: []string })` | `false` | <code></code> |

**Docker Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_REGISTRY` | Docker registry url to login to. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_USERNAME` | Docker registry username for the given registry. | `string` | `false` | <code></code> |
| `$DOCKER_REGISTRY_PASSWORD` | Docker registry password for the given registry. | `string` | `false` | <code></code> |
