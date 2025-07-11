package run

import (
	"fmt"

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
