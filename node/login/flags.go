package login

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "npm.login",
		Usage:       "npm registries to login to. Format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])",
		Required:    false,
		EnvVars:     []string{"NPM_LOGIN"},
		Value:       "",
		Destination: &Pipe.Npm.Login,
	},

	&cli.StringFlag{
		Name:        "npm.npmrc_file",
		Usage:       ".npmrc file to use.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC_FILE"},
		Value:       ".npmrc",
		Destination: &Pipe.Npm.NpmRc,
	},
}
