# pipes-update-docker-hub-readme

Updates the readme file on DockerHub or any compatible API.

`pipes-update-docker-hub-readme [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |

### DockerHub

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_USERNAME`<br/>`$PLUGIN_DOCKER_USERNAME` | Docker Hub username for updating the readme. | `String` | `true` |  |
| `$DOCKER_PASSWORD`<br/>`$PLUGIN_DOCKER_PASSWORD` | Docker Hub password for updating the readme. | `String` | `true` |  |
| `$DOCKER_HUB_ADDRESS`<br/>`$PLUGIN_DOCKER_HUB_ADDRESS` | HTTP address for the docker hub. There is only one! | `String` | `false` | &#34;https://hub.docker.com/v2/repositories&#34; |

### Readme

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DOCKER_IMAGE_NAME`<br/>`$IMAGE_NAME`<br/>`$README_REPOSITORY`<br/>`$PLUGIN_README_REPOSITORY` | Repository for applying the readme on. | `String` | `true` |  |
| `$README_FILE`<br/>`$PLUGIN_README_FILE` | Readme file for the given repossitory. | `String` | `false` | &#34;README.md&#34; |
| `$README_DESCRIPTION`<br/>`$PLUGIN_README_DESCRIPTION` | Pass in description to send it in the request. | `String` | `false` |  |
