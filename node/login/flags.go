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
		Destination: &P.Pipe.Npm.Login,
	},

	&cli.StringSliceFlag{
		Name:        "npm.npmrc_file",
		Usage:       ".npmrc file to use.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC_FILE"},
		Value:       cli.NewStringSlice(".npmrc"),
		Destination: &P.Pipe.Npm.NpmRcFile,
	},

	&cli.StringFlag{
		Name:        "npm.npmrc",
		Usage:       "Pass direct contents of the NPMRC file.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC"},
		Value:       "",
		Destination: &P.Pipe.Npm.NpmRc,
	},
}
