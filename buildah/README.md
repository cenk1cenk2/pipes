# pipe-buildah

Builds and publishes container images from CI/CD with buildah.io

`pipe-buildah [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | `info` |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | `[]` |

## Commands

### `pipe-buildah login`

Login to the given container registries.

#### Flags

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | `` |

### `pipe-buildah build`

Build container images.

#### Flags

**Container Image**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_IMAGE_PLATFORMS` | Container image platforms to be built. | `string[]` | `false` | `[]` |
| `$CONTAINER_IMAGE_NAME` | Image name for the container image to be built. | `string` | `true` | `` |
| `$CONTAINER_IMAGE_TAGS` | Image tags for the container image to be built. | `string[]` | `true` | `[]` |
| `$CONTAINER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | `[]` |
| `$CONTAINER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]` |
| `$CONTAINER_IMAGE_TAGS_AS_LATEST` | Regex pattern to tag the image as latest.<br />    Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(yaml([]RegExp))` | `false` | `[ "^tags/v?\\d+.\\d+.\\d+$" ]` |
| `$CONTAINER_IMAGE_PULL` | Pull before building the image. | `bool` | `false` | `true` |
| `$CONTAINER_IMAGE_PUSH` | Push the image after building. | `bool` | `false` | `true` |
| `$CONTAINER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br />    You can use it as a template with environment variables as the context. | `string`<br/>`format(yaml(map[string]Template()))` | `false` | `` |
| `$CONTAINER_IMAGE_INSPECT` | Inspect after pushing the image. | `bool` | `false` | `true` |
| `$CONTAINER_IMAGE_LATEST_TAG` | Latest tag for the container image where it is marked as latest. | `string` | `false` | `latest` |
| `$CONTAINER_IMAGE_CACHE` | Specify the cache for the container image. | `string` | `false` | `` |
| `$CONTAINER_IMAGE_FORMAT` | Specify the format for Container Image. | `string` | `false` | `oci` |
| `$CONTAINER_IMAGE_STORAGE_DRIVER`<br/>`$BUILDAH_STORAGE_DRIVER` | Specify the storage driver for Buildah. | `string` | `false` | `vfs` |

**Container Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template([]string))` | `false` | `` |
| `$CONTAINER_MANIFEST_FILE` | Write all the images that are published in to a file for later use. | `string`<br/>`format(Template([]string))` | `false` | `.published-container-images_{{ $ | join "," | sha256sum }}` |

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | `` |

**Containerfile**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINERFILE_CONTEXT` | Containerfile context argument for build operation. | `string` | `false` | `.` |
| `$CONTAINERFILE_NAME` | Containerfile path for the build operation | `string` | `false` | `Dockerfile` |

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

### `pipe-buildah manifest`

Update manifests of the container images.

#### Flags

**Container Manifest**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_MANIFEST_FILES` | Read published tags from a file. | `string[]`<br/>`format(glob)` | `false` | `[**/.published-container-images*]` |
| `$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template())` | `false` | `` |
| `$CONTAINER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `string[]` | `false` | `[]` |
| `$CONTAINER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `string`<br/>`format(yaml([]struct { target: string, images: []string }))` | `false` | `` |

**Container Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CONTAINER_REGISTRY_URI` | Container registry url to login to. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` | `false` | `` |
| `$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` | `false` | `` |
