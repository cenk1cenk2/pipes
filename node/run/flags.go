package run

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_COMMAND = "Command"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_COMMAND

	&cli.StringFlag{
		Category: CATEGORY_NODE_COMMAND,
		Name:     "node.command_script",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_COMMAND_SCRIPT"),
		),
		Usage: fmt.Sprintf(
			"package.json script for given command operation. %s",
			environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE,
		),
		Required:    false,
		Destination: &P.NodeCommand.Script,
		Action: func(ctx context.Context, c *cli.Command, v string) error {
			if v == "" {
				args := c.Args().Slice()

				if len(args) < 1 {
					return fmt.Errorf("Arguments are needed to run a specific script.")
				}

				C.Script = args[0]
				C.ScriptArgs = strings.Join(args[1:], " ")
			} else {
				C.Script = strings.Split(v, " ")[0]
				C.ScriptArgs = strings.Join(strings.Split(v, " ")[1:], " ")
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_COMMAND,
		Name:     "node.command_cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_COMMAND_CWD"),
		),
		Usage:       "Working directory for the given command operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.NodeCommand.Cwd,
	},
}
