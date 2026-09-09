# pipe-update-docker-hub-readme

Updates the readme file on DockerHub or any compatible API.

`pipe-update-docker-hub-readme [FLAGS]`

## Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

**DockerHub**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$DOCKER_HUB_USERNAME`**<br/>**`$DOCKER_USERNAME`**\* | DockerHub username for updating the README. | `string` |  |
| **`$DOCKER_HUB_PASSWORD`**<br/>**`$DOCKER_PASSWORD`**\* | DockerHub password for updating the README. | `string` |  |
| `$DOCKER_HUB_ADDRESS` | HTTP address for the DockerHub-compatible service. | `string` | `"https://hub.docker.com/v2/repositories"` |

\* required

**Readme**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$DOCKER_HUB_README_REPOSITORY`<br/>`$DOCKER_IMAGE_NAME`<br/>`$CONTAINER_IMAGE_NAME`<br/>`$README_REPOSITORY` | Repository to apply the README to. | `string` |  |
| `$DOCKER_HUB_README_FILE`<br/>`$README_FILE` | README file for the given repository. | `string` | `"README.md"` |
| `$DOCKER_HUB_README_DESCRIPTION`<br/>`$README_SHORT_DESCRIPTION` | Short description to display on DockerHub. | `string` |  |
| `$DOCKER_HUB_README_MATRIX`<br/>`$README_MATRIX` | Matrix of multiple README files to update. | `string`<br/>`format(json([]struct{ repository: string, file: string, description?: string }))` |  |
