package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "gh.token",
		Usage:       "Github token for the API requests.",
		Required:    false,
		EnvVars:     []string{"GH_TOKEN", "GITHUB_TOKEN"},
		Value:       "",
		Destination: &Pipe.Github.Token,
	},
	&cli.StringFlag{
		Name:        "gh.repository",
		Usage:       "Target repository to fetch the latest tag.",
		Required:    true,
		EnvVars:     []string{"GH_REPOSITORY", "GH_REPOSITORY"},
		Value:       "",
		Destination: &Pipe.Github.Repository,
	},
	&cli.StringFlag{
		Name:        "docker_image.tags_file",
		Usage:       "Read tags from a file.",
		Required:    true,
		EnvVars:     []string{"TAGS_FILE"},
		Value:       "",
		Destination: &Pipe.DockerImage.TagsFile,
	},
}
