package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "docker_hub.username",
		Usage:       "Docker Hub username for updating the readme.",
		EnvVars:     []string{"DOCKER_USERNAME", "PLUGIN_DOCKER_USERNAME"},
		Required:    true,
		Destination: &Pipe.DockerHub.Username,
	},
	&cli.StringFlag{
		Name:        "docker_hub.password",
		Usage:       "Docker Hub password for updating the readme.",
		EnvVars:     []string{"DOCKER_PASSWORD", "PLUGIN_DOCKER_PASSWORD"},
		Required:    true,
		Destination: &Pipe.DockerHub.Password,
	},
	&cli.StringFlag{
		Name:        "docker_hub.address",
		Usage:       "HTTP address for the docker hub. There is only one!",
		EnvVars:     []string{"DOCKER_HUB_ADDRESS", "PLUGIN_DOCKER_HUB_ADDRESS"},
		Value:       "https://hub.docker.com/v2/repositories",
		Destination: &Pipe.DockerHub.Address,
	},
	&cli.StringFlag{
		Name:        "readme.repository",
		Usage:       "Repository for applying the readme on.",
		EnvVars:     []string{"README_REPOSITORY", "PLUGIN_README_REPOSITORY"},
		Required:    true,
		Destination: &Pipe.Readme.Repository,
	},
	&cli.StringFlag{
		Name:        "readme.file",
		Usage:       "Readme file for the given repossitory.",
		EnvVars:     []string{"README_FILE", "PLUGIN_README_FILE"},
		Value:       "README.md",
		Destination: &Pipe.Readme.File,
		Required:    false,
	},
	&cli.StringFlag{
		Name:        "readme.short_description",
		Usage:       "Pass in description to send it in the request.",
		EnvVars:     []string{"README_DESCRIPTION", "PLUGIN_README_DESCRIPTION"},
		Destination: &Pipe.Readme.Description,
		Required:    false,
	},
}
