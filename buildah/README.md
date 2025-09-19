# pipe-buildah

Builds and publishes container images from CI/CD with buildah.io

`pipe-buildah [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>info</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code>[]</code> |

## Commands

### `pipe-buildah login`

Login to the given container registries.

#### Flags

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | <code>docker.io</code> |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | <code></code> |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | <code></code> |

### `pipe-buildah build`

Build container images.

#### Flags

**Container Image**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_IMAGE_PLATFORMS` | Container image platforms to be built. | `string[]` | `false` | <code>[]</code> |
| `$CONTAINER_IMAGE_NAME` | Image name for the container image to be built. | `string` | `true` | <code></code> |
| `$CONTAINER_IMAGE_TAGS` | Image tags for the container image to be built. | `string[]` | `true` | <code>[]</code> |
| `$CONTAINER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>[]</code> |
| `$CONTAINER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>[<br />    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }<br />]</code> |
| `$CONTAINER_IMAGE_TAGS_AS_LATEST` | Regex pattern to tag the image as latest.<br />    Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(yaml([]RegExp))` | `false` | <code>[ "^tags/v?\\d+.\\d+.\\d+$" ]</code> |
| `$CONTAINER_IMAGE_PULL` | Pull before building the image. | `bool` | `false` | <code>true</code> |
| `$CONTAINER_IMAGE_PUSH` | Push the image after building. | `bool` | `false` | <code>true</code> |
| `$CONTAINER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br />    You can use it as a template with environment variables as the context. | `string`<br/>`format(yaml(map[string]Template()))` | `false` | <code></code> |
| `$CONTAINER_IMAGE_LATEST_TAG` | Latest tag for the container image where it is marked as latest. | `string` | `false` | <code>latest</code> |
| `$CONTAINER_IMAGE_CACHE` | Specify the cache for the container image. | `string` | `false` | <code></code> |
| `$CONTAINER_IMAGE_FORMAT` | Specify the format for Container Image. | `string` | `false` | <code>oci</code> |
| `$CONTAINER_IMAGE_STORAGE_DRIVER`<br />`$BUILDAH_STORAGE_DRIVER` | Specify the storage driver for Buildah. | `string` | `false` | <code>vfs</code> |

**Container Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template([]string))` | `false` | <code></code> |
| `$CONTAINER_MANIFEST_FILE` | Write all the images that are published in to a file for later use. | `string`<br/>`format(Template([]string))` | `false` | <code>.published-container-images_{{ $ | join "," | sha256sum }}</code> |

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | <code>docker.io</code> |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | <code></code> |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | <code></code> |

**Containerfile**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_FILE_CONTEXT` | Containerfile context argument for build operation. | `string` | `false` | <code>.</code> |
| `$CONTAINER_FILE_NAME` | Containerfile path for the build operation | `string` | `false` | <code>Dockerfile</code> |

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

### `pipe-buildah manifest`

Update manifests of the container images.

#### Flags

**Container Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_MANIFEST_FILES` | Read published tags from a file. | `string[]`<br/>`format(glob)` | `false` | <code>[**/.published-container-images*]</code> |
| `$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template())` | `false` | <code></code> |
| `$CONTAINER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `string[]` | `false` | <code>[]</code> |
| `$CONTAINER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `string`<br/>`format(yaml([]struct { target: string, images: []string }))` | `false` | <code></code> |

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | <code>docker.io</code> |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | <code></code> |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | <code></code> |
