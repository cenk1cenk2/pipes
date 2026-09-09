# pipe-buildah

Builds and publishes container images from CI with buildah.io

`pipe-buildah [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-buildah login`](#pipe-buildah-login)
- [`pipe-buildah build`](#pipe-buildah-build)
- [`pipe-buildah manifest`](#pipe-buildah-manifest)

### `pipe-buildah login`

Login to the given container registries.

`pipe-buildah login [FLAGS]`

#### Flags

**Buildah**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_CWD`<br/>`$CONTAINER_CWD` | Working directory for buildah commands. | `string` | `"."` |

**Container Registry**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_LOGIN_REGISTRY_URI`<br/>`$CONTAINER_REGISTRY_URI` | Container registry URL to login to. | `string` | `"docker.io"` |
| `$BUILDAH_LOGIN_REGISTRY_USERNAME`<br/>`$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` |  |
| `$BUILDAH_LOGIN_REGISTRY_PASSWORD`<br/>`$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` |  |

### `pipe-buildah build`

Build container images.

`pipe-buildah build [FLAGS]`

#### Flags

<details>
<summary>Buildah</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_CWD`<br/>`$CONTAINER_CWD` | Working directory for buildah commands. | `string` | `"."` |

</details>

**Container Image**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_BUILD_IMAGE_PLATFORMS`<br/>`$CONTAINER_IMAGE_PLATFORMS` | Container image platforms to be built. | `string[]` |  |
| **`$BUILDAH_BUILD_IMAGE_NAME`**<br/>**`$CONTAINER_IMAGE_NAME`**\* | Image name for the container image to be built. | `string` |  |
| **`$BUILDAH_BUILD_IMAGE_TAGS`**<br/>**`$CONTAINER_IMAGE_TAGS`**\* | Image tags for the container image to be built. | `string[]` |  |
| `$BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE`<br/>`$CONTAINER_IMAGE_TAGS_TEMPLATE` | Modifies every tag that matches a certain condition.<br/>Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `"[]"` |
| `$BUILDAH_BUILD_IMAGE_TAGS_SANITIZE`<br/>`$CONTAINER_IMAGE_SANITIZE_TAGS` | Sanitizes the given regex pattern out of tag name.<br/>Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `"[\n    { \"match\": \"([^/]*)/(.*)\", \"template\": \"{{ index $ 1 \| upper }}_{{ index $ 2 }}\" }\n]"` |
| `$BUILDAH_BUILD_IMAGE_TAG_AS_LATEST`<br/>`$CONTAINER_IMAGE_TAGS_AS_LATEST` | Regex pattern to tag the image as latest.<br/>Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(yaml([]RegExp))` | `"[ \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\" ]"` |
| `$BUILDAH_BUILD_IMAGE_PULL`<br/>`$CONTAINER_IMAGE_PULL` | Pull before building the image. | `bool` | `true` |
| `$BUILDAH_BUILD_IMAGE_PUSH`<br/>`$CONTAINER_IMAGE_PUSH` | Push the image after building. | `bool` | `true` |
| `$BUILDAH_BUILD_IMAGE_BUILD_ARGS`<br/>`$CONTAINER_IMAGE_BUILD_ARGS` | Pass in extra build arguments for image.<br/>You can use it as a template with environment variables as the context. | `string`<br/>`format(yaml(map[string]Template()))` |  |
| `$BUILDAH_BUILD_IMAGE_LATEST_TAG`<br/>`$CONTAINER_IMAGE_LATEST_TAG` | Latest tag for the container image where it is marked as latest. | `string` | `"latest"` |
| `$BUILDAH_BUILD_IMAGE_CACHE`<br/>`$CONTAINER_IMAGE_CACHE` | Specify the cache for the container image. | `string` |  |
| `$BUILDAH_BUILD_IMAGE_FORMAT`<br/>`$CONTAINER_IMAGE_FORMAT` | Specify the format for Container Image. | `string`<br/>`format(enum("oci", "docker"))` | `"oci"` |
| `$BUILDAH_BUILD_IMAGE_STORAGE_DRIVER`<br/>`$CONTAINER_IMAGE_STORAGE_DRIVER`<br/>`$BUILDAH_STORAGE_DRIVER` | Specify the storage driver for Buildah. | `string`<br/>`format(enum("overlay", "overlay2", "vfs"))` | `"vfs"` |

\* required

**Container Manifest**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_BUILD_MANIFEST_TARGET`<br/>`$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template([]string))` |  |
| `$BUILDAH_BUILD_MANIFEST_FILE`<br/>`$CONTAINER_MANIFEST_FILE` | Write all the published images into a file for later use. | `string`<br/>`format(Template([]string))` | `".published-container-images_{{ $ \| join \",\" \| sha256sum }}"` |

<details>
<summary>Container Registry</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_LOGIN_REGISTRY_URI`<br/>`$CONTAINER_REGISTRY_URI` | Container registry URL to login to. | `string` | `"docker.io"` |
| `$BUILDAH_LOGIN_REGISTRY_USERNAME`<br/>`$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` |  |
| `$BUILDAH_LOGIN_REGISTRY_PASSWORD`<br/>`$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` |  |

</details>

**Containerfile**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_BUILD_FILE_CONTEXT`<br/>`$CONTAINER_FILE_CONTEXT` | Containerfile context argument for build operation. | `string` | `"."` |
| `$BUILDAH_BUILD_FILE_NAME`<br/>`$CONTAINER_FILE_NAME` | Containerfile path for the build operation. | `string` | `"Dockerfile"` |

**GIT**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` |  |

**Tags File**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$TAGS_FILE` | Read tags from a comma separated file. | `string` |  |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` |

### `pipe-buildah manifest`

Update manifests of the container images.

`pipe-buildah manifest [FLAGS]`

#### Flags

<details>
<summary>Buildah</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_CWD`<br/>`$CONTAINER_CWD` | Working directory for buildah commands. | `string` | `"."` |

</details>

**Container Manifest**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_MANIFEST_FILES`<br/>`$CONTAINER_MANIFEST_FILES` | Read published tags from a file. | `string[]`<br/>`format(glob)` | `"**/.published-container-images*"` |
| `$BUILDAH_MANIFEST_TARGET`<br/>`$CONTAINER_MANIFEST_TARGET` | Target image names for patching the manifest. | `string`<br/>`format(Template())` |  |
| `$BUILDAH_MANIFEST_IMAGES`<br/>`$CONTAINER_MANIFEST_IMAGES` | Image names for patching the manifest with the given target. | `string[]` |  |
| `$BUILDAH_MANIFEST_MATRIX`<br/>`$CONTAINER_MANIFEST_MATRIX` | Matrix of all the images that should be manifested. | `string`<br/>`format(yaml([]struct{ target: string, images: []string }))` |  |

<details>
<summary>Container Registry</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$BUILDAH_LOGIN_REGISTRY_URI`<br/>`$CONTAINER_REGISTRY_URI` | Container registry URL to login to. | `string` | `"docker.io"` |
| `$BUILDAH_LOGIN_REGISTRY_USERNAME`<br/>`$CONTAINER_REGISTRY_USERNAME` | Container registry username for the given registry. | `string` |  |
| `$BUILDAH_LOGIN_REGISTRY_PASSWORD`<br/>`$CONTAINER_REGISTRY_PASSWORD` | Container registry password for the given registry. | `string` |  |

</details>
