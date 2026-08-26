package run

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_COMMAND = "Command"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_COMMAND

	&cli.StringFlag{
		Category: CATEGORY_NODE_COMMAND,
		Name:     "node.run.script",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_COMMAND_SCRIPT"),
			cli.EnvVar("NODE_RUN_SCRIPT"),
		),
		Usage: fmt.Sprintf(
			"package.json script for given command operation. %s",
			environment.HELP_FORMAT_TEMPLATE,
		),
		Required:    false,
		Destination: &P.Run.Script,
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_COMMAND,
		Name:     "node.run.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_COMMAND_CWD"),
			cli.EnvVar("NODE_RUN_CWD"),
		),
		Usage:       "Working directory for the given command operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Run.Cwd,
	},
}

var Arguments = []cli.Argument{
	&cli.StringArgs{
		Name:        "arg",
		Min:         0,
		Max:         -1,
		UsageText:   "Tool to run.",
		Destination: &P.Run.Command,
	},
}
