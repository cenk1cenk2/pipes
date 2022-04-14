# update-docker-hub-readme

## Description

Updates DockerHub short and long description for a given repository.

## Usage

```bash
GLOBAL OPTIONS:
   --utils.ci value                  Indicates this is running inside a CI/CD environment to act accordingly. [$CI]
   --utils.debug value               Set the log level debug for the application. [$DEBUG]
   --utils.log value                 Define the log level for the application. (default: "info") [$LOG_LEVEL]
   --docker_hub.username value       Docker Hub username for updating the readme. [$DOCKER_USERNAME, $PLUGIN_DOCKER_USERNAME]
   --docker_hub.password value       Docker Hub password for updating the readme. [$DOCKER_PASSWORD, $PLUGIN_DOCKER_PASSWORD]
   --docker_hub.address value        HTTP address for the docker hub. There is only one! (default: "https://hub.docker.com/v2/repositories") [$DOCKER_HUB_ADDRESS, $PLUGIN_DOCKER_HUB_ADDRESS]
   --readme.repository value         Repository for applying the readme on. [$README_REPOSITORY, $PLUGIN_README_REPOSITORY]
   --readme.file value               Readme file for the given repossitory. (default: "README.md") [$README_FILE, $PLUGIN_README_FILE]
   --readme.short_description value  Pass in description to send it in the request. [$README_DESCRIPTION, $PLUGIN_README_DESCRIPTION]
   --help, -h                        show help
   --version, -v                     print the version
```
