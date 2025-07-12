# pipe-update-docker-hub-readme

Updates the readme file on DockerHub or any compatible API.

`pipe-update-docker-hub-readme [FLAGS]`

## Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | [] |

**DockerHub**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USERNAME` | DockerHub username for updating the readme. | `string` | `true` |  |
| `$DOCKER_PASSWORD` | DockerHub password for updating the readme. | `string` | `true` |  |
| `$DOCKER_HUB_ADDRESS` | HTTP address for the DockerHub compatible service. | `string` | `false` | https://hub.docker.com/v2/repositories |

**Readme**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME`<br/>`$README_REPOSITORY` | Repository for applying the readme on. | `string` | `false` |  |
| `$README_FILE` | Readme file for the given repository. | `string` | `false` | README.md |
| `$README_SHORT_DESCRIPTION` | Short description to display on DockerHub. | `string` | `false` |  |
| `$README_MATRIX` | Matrix of multiple README files to update. | `string`<br/>`json([]struct { repository: string, file: string, description?: string })` | `false` |  |
