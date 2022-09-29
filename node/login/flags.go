package login

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_node_login = "Login"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_node_login,
		Name:        "npm.login",
		Usage:       "npm registries to login to. format(json({ username: string, password: string, registry?: string, useHttps?: boolean }[]))",
		Required:    false,
		EnvVars:     []string{"NPM_LOGIN"},
		Value:       "",
		Destination: &TL.Pipe.Npm.Login,
	},

	&cli.StringSliceFlag{
		Category:    category_node_login,
		Name:        "npm.npmrc_file",
		Usage:       ".npmrc file to use.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC_FILE"},
		Value:       cli.NewStringSlice(".npmrc"),
		Destination: &TL.Pipe.Npm.NpmRcFile,
	},

	&cli.StringFlag{
		Category:    category_node_login,
		Name:        "npm.npmrc",
		Usage:       "Pass direct contents of the NPMRC file.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC"},
		Value:       "",
		Destination: &TL.Pipe.Npm.NpmRc,
	},
}
