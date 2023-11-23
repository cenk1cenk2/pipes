package login

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
		Usage:    "NPM registries to login. json([]struct { username: string, password: string, registry?: string, useHttps?: bool })",
		Required: false,
		EnvVars:  []string{"NPM_LOGIN"},
		Value:    "",
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_NODE_LOGIN,
		Name:     "npm.npmrc_file",
		Usage:    ".npmrc file to use.",
		Required: false,
		EnvVars:  []string{"NPM_NPMRC_FILE"},
		Value:    cli.NewStringSlice(".npmrc"),
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_LOGIN,
		Name:        "npm.npmrc",
		Usage:       "Direct contents of .npmrc file.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC"},
		Value:       "",
		Destination: &TL.Pipe.Npm.NpmRc,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	if v := tl.CliContext.String("npm.login"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.Npm.Login); err != nil {
			return fmt.Errorf("Can not unmarshal Npm registry login credentials: %w", err)
		}
	}

	tl.Pipe.Npm.NpmRcFile = tl.CliContext.StringSlice("npm.npmrc_file")

	return nil
}
