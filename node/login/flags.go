package login

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_LOGIN = "Login"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_LOGIN

	&cli.StringFlag{
		Category: CATEGORY_NODE_LOGIN,
		Name:     "npm.login",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NPM_LOGIN"),
		),
		Usage:    "NPM registries to login. json([]struct { username: string, password: string, registry?: string, useHttps?: bool })",
		Required: false,
		Value:    "",
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.Npm.Login); err != nil {
				return fmt.Errorf("Can not unmarshal Npm registry login credentials: %w", err)
			}

			return nil
		},
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_NODE_LOGIN,
		Name:     "npm.npmrc_file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NPM_NPMRC_FILE"),
		),
		Usage:    ".npmrc file to use.",
		Required: false,
		Value:    []string{".npmrc"},
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_LOGIN,
		Name:     "npm.npmrc",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NPM_NPMRC"),
		),
		Usage:       "Direct contents of .npmrc file.",
		Required:    false,
		Value:       "",
		Destination: &P.Npm.NpmRc,
	},
}
