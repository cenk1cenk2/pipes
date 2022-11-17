# pipe-update-docker-hub-readme

Updates the readme file on DockerHub or any compatible API.

`pipe-update-docker-hub-readme [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### DockerHub

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USERNAME`<br/>`$PLUGIN_DOCKER_USERNAME` | Docker Hub username for updating the readme. | `String` | `true` |  |
| `$DOCKER_PASSWORD`<br/>`$PLUGIN_DOCKER_PASSWORD` | Docker Hub password for updating the readme. | `String` | `true` |  |
| `$DOCKER_HUB_ADDRESS`<br/>`$PLUGIN_DOCKER_HUB_ADDRESS` | HTTP address for the docker hub. There is only one! | `String` | `false` | https://hub.docker.com/v2/repositories |

### Readme

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME`<br/>`$IMAGE_NAME`<br/>`$README_REPOSITORY`<br/>`$PLUGIN_README_REPOSITORY` | Repository for applying the readme on. | `String` | `true` |  |
| `$README_FILE`<br/>`$PLUGIN_README_FILE` | Readme file for the given repossitory. | `String` | `false` | README.md |
| `$README_DESCRIPTION`<br/>`$PLUGIN_README_DESCRIPTION` | Pass in description to send it in the request. | `String` | `false` |  |
