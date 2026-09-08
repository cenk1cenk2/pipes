package node

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

// Options is what the package manager flags are built onto.
type Options struct {
	Destination *Config
}

func NewFlags(opts Options) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    CategoryPackageManager,
			Name:        "node.package-manager",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("NODE_PACKAGE_MANAGER")),
			Usage:       `Preferred package manager for nodejs. format(enum("npm", "yarn", "pnpm"))`,
			Required:    false,
			Value:       DefaultPackageManager,
			Destination: &opts.Destination.PackageManager,
		},
	}
}

// LoginOptions is what the npm login flags are built onto.
type LoginOptions struct {
	Destination *Login
}

func NewLoginFlags(opts LoginOptions) []cli.Flag {
	return []cli.Flag{
		flags.JSONFlag(&opts.Destination.Entries, &cli.StringFlag{
			Category: "Login",
			Name:     "npm.login",
			Sources:  cli.NewValueSourceChain(cli.EnvVar("NPM_LOGIN")),
			Usage:    "NPM registries to login. format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))",
			Required: false,
			Value:    "",
		}),

		&cli.StringSliceFlag{
			Category:    "Login",
			Name:        "npm.npmrc-file",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("NPM_NPMRC_FILE")),
			Usage:       ".npmrc file to use.",
			Required:    false,
			Value:       []string{".npmrc"},
			Destination: &opts.Destination.NpmRcFiles,
		},

		&cli.StringFlag{
			Category:    "Login",
			Name:        "npm.npmrc",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("NPM_NPMRC")),
			Usage:       "Direct contents of .npmrc file.",
			Required:    false,
			Value:       "",
			Destination: &opts.Destination.NpmRc,
		},
	}
}
