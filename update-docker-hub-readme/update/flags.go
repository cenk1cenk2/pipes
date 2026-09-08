package update

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER_HUB = "DockerHub"
	CATEGORY_README     = "Readme"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker-hub.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_USERNAME"),
			cli.EnvVar("DOCKER_USERNAME"),
		),
		Usage:       "DockerHub username for updating the readme.",
		Required:    true,
		Destination: &P.DockerHub.Username,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker-hub.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_PASSWORD"),
			cli.EnvVar("DOCKER_PASSWORD"),
		),
		Usage:       "DockerHub password for updating the readme.",
		Required:    true,
		Destination: &P.DockerHub.Password,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker-hub.address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_ADDRESS"),
		),
		Usage:       "HTTP address for the DockerHub compatible service.",
		Value:       "https://hub.docker.com/v2/repositories",
		Destination: &P.DockerHub.Address,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "docker-hub.readme.repository",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_README_REPOSITORY"),
			cli.EnvVar("DOCKER_IMAGE_NAME"),
			cli.EnvVar("CONTAINER_IMAGE_NAME"),
			cli.EnvVar("README_REPOSITORY"),
		),
		Usage:       "Repository for applying the readme on.",
		Required:    false,
		Value:       "",
		Destination: &P.Readme.Repository,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "docker-hub.readme.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_README_FILE"),
			cli.EnvVar("README_FILE"),
		),
		Usage:       "Readme file for the given repository.",
		Value:       "README.md",
		Destination: &P.Readme.File,
		Required:    false,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "docker-hub.readme.description",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_README_DESCRIPTION"),
			cli.EnvVar("README_SHORT_DESCRIPTION"),
		),
		Usage:       "Short description to display on DockerHub.",
		Destination: &P.Readme.Description,
		Required:    false,
	},

	flags.JSONFlag(&P.Readme.Matrix, &cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "docker-hub.readme.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_README_MATRIX"),
			cli.EnvVar("README_MATRIX"),
		),
		Usage:    "Matrix of multiple README files to update. format(json([]struct{ repository: string, file: string, description?: string }))",
		Required: false,
	}),
}
