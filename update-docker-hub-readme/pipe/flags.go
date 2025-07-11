package pipe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER_HUB = "DockerHub"
	CATEGORY_README     = "Readme"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker_hub.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_USERNAME"),
		),
		Usage:       "DockerHub username for updating the readme.",
		Required:    true,
		Destination: &P.DockerHub.Username,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker_hub.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_PASSWORD"),
		),
		Usage:       "DockerHub password for updating the readme.",
		Required:    true,
		Destination: &P.DockerHub.Password,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_HUB,
		Name:     "docker_hub.address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_HUB_ADDRESS"),
		),
		Usage:       "HTTP address for the DockerHub compatible service.",
		Value:       "https://hub.docker.com/v2/repositories",
		Destination: &P.DockerHub.Address,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "readme.repository",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_NAME"),
			cli.EnvVar("README_REPOSITORY"),
		),
		Usage:       "Repository for applying the readme on.",
		Required:    false,
		Value:       "",
		Destination: &P.Readme.Repository,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "readme.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("README_FILE"),
		),
		Usage:       "Readme file for the given repository.",
		Value:       "README.md",
		Destination: &P.Readme.File,
		Required:    false,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "readme.short_description",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("README_SHORT_DESCRIPTION"),
		),
		Usage:       "Short description to display on DockerHub.",
		Destination: &P.Readme.Description,
		Required:    false,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if len(P.Readme.Description) > 100 {
				return fmt.Errorf(
					"Readme short description can only be 100 characters long while you have: %d",
					len(P.Readme.Description),
				)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "readme.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("README_MATRIX"),
		),
		Usage:    "Matrix of multiple README files to update. json([]struct { repository: string, file: string, description?: string })",
		Required: false,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.Readme.Matrix); err != nil {
				return fmt.Errorf("Can not unmarshal Readme matrix: %w", err)
			}

			return nil
		},
	},
}
