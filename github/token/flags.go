package token

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

//revive:disable:line-length-limit

const (
	CategoryToken = "Token"
)

var Flags = append(client.NewFlags(client.Options{Destination: &P.App}), []cli.Flag{
	&cli.StringSliceFlag{
		Category: CategoryToken,
		Name:     "token.repositories",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_TOKEN_REPOSITORIES"),
		),
		Usage:       "Repository names, without the owner, to narrow the token down to. Left empty, the token covers every repository of the installation.",
		Required:    false,
		Value:       []string{},
		Destination: &P.Token.Repositories,
	},

	flags.JSONFlag(&P.Token.Permissions, &cli.StringFlag{
		Category: CategoryToken,
		Name:     "token.permissions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_TOKEN_PERMISSIONS"),
		),
		Usage:    `Permissions to narrow the token down to, e.g. {"statuses":"write"}. Left empty, the token carries every permission of the installation. format(json(map[string]string))`,
		Required: false,
	}),

	&cli.StringFlag{
		Category: CategoryToken,
		Name:     "token.variable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_TOKEN_VARIABLE"),
		),
		Usage:       "Variable name the token is written under in the dotenv file.",
		Required:    false,
		Value:       "GH_TOKEN",
		Destination: &P.Token.Variable,
	},

	&cli.StringFlag{
		Category: CategoryToken,
		Name:     "token.git-credentials-variable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_TOKEN_GIT_CREDENTIALS_VARIABLE"),
		),
		Usage:       "Variable name the token is written under as an x-access-token git credential in the dotenv file, for semantic-release to push with. Left empty, no git credential is written.",
		Required:    false,
		Value:       "GIT_CREDENTIALS",
		Destination: &P.Token.GitCredentialsVariable,
	},

	&cli.StringFlag{
		Category: CategoryToken,
		Name:     "token.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GITHUB_TOKEN_FILE"),
		),
		Usage:       "Dotenv file the token is written into, for the job to expose as a dotenv report. Other variables in an existing file are kept.",
		Required:    false,
		Value:       "github.env",
		Destination: &P.Token.File,
	},
}...)
