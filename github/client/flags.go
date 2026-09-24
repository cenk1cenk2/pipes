package client

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryGitHubApp = "GitHub App"
)

// Options is what the GitHub App flags are built onto.
type Options struct {
	Destination *Config
}

// NewFlags builds the GitHub App flags onto cfg. None of them is required here,
// since posting a status can run on a ready token alone; a command that always
// mints enforces them through Config.Enabled.
func NewFlags(opts Options) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category: CategoryGitHubApp,
			Name:     "github-app.id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITHUB_APP_ID"),
			),
			Usage:       "GitHub App id the token is minted for.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.Id,
		},

		&cli.StringFlag{
			Category: CategoryGitHubApp,
			Name:     "github-app.installation-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITHUB_APP_INSTALLATION_ID"),
			),
			Usage:       "Installation id of the GitHub App on the owner of the repository.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.InstallationId,
		},

		&cli.StringFlag{
			Category: CategoryGitHubApp,
			Name:     "github-app.private-key",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITHUB_APP_PRIVATE_KEY"),
			),
			Usage:       "Path to the private key of the GitHub App, as a GitLab file variable exports it, or the PEM contents themselves.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.PrivateKey,
		},

		&cli.StringFlag{
			Category: CategoryGitHubApp,
			Name:     "github-app.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITHUB_API_URL"),
			),
			Usage:       "GitHub API URL.",
			Required:    false,
			Value:       "https://api.github.com",
			Destination: &opts.Destination.ApiUrl,
		},
	}
}
