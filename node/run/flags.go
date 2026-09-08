package run

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

const (
	CategoryNodeCommand = "Command"
)

var Flags = []cli.Flag{

	// CategoryNodeCommand

	&cli.StringFlag{
		Category: CategoryNodeCommand,
		Name:     "node.run.script",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_RUN_SCRIPT"),
			cli.EnvVar("NODE_COMMAND_SCRIPT"),
		),
		Usage: fmt.Sprintf(
			"package.json script for the given command operation. %s",
			environment.HelpFormatTemplate,
		),
		Required:    false,
		Destination: &P.Run.Script,
	},

	&cli.StringFlag{
		Category: CategoryNodeCommand,
		Name:     "node.run.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_RUN_CWD"),
			cli.EnvVar("NODE_COMMAND_CWD"),
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
		UsageText:   "Arguments appended to the script.",
		Destination: &P.Run.Command,
	},
}
