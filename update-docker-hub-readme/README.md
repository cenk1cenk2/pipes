# pipes-update-docker-hub-readme

Updates the readme file on DockerHub or any compatible API.

`pipes-update-docker-hub-readme [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $$DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $$LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --docker_hub.username, $$DOCKER_USERNAME, $$PLUGIN_DOCKER_USERNAME | Docker Hub username for updating the readme. |  String  | true |  |
| --docker_hub.password, $$DOCKER_PASSWORD, $$PLUGIN_DOCKER_PASSWORD | Docker Hub password for updating the readme. |  String  | true |  |
| --docker_hub.address, $$DOCKER_HUB_ADDRESS, $$PLUGIN_DOCKER_HUB_ADDRESS | HTTP address for the docker hub. There is only one! |  String  | false | &#34;https://hub.docker.com/v2/repositories&#34; |
| --readme.repository, $$DOCKER_IMAGE_NAME, $$IMAGE_NAME, $$README_REPOSITORY, $$PLUGIN_README_REPOSITORY | Repository for applying the readme on. |  String  | true |  |
| --readme.file, $$README_FILE, $$PLUGIN_README_FILE | Readme file for the given repossitory. |  String  | false | &#34;README.md&#34; |
| --readme.short_description, $$README_DESCRIPTION, $$PLUGIN_README_DESCRIPTION | Pass in description to send it in the request. |  String  | false |  |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
