package run

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_COMMAND = "Command"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_COMMAND

	&cli.StringFlag{
		Category:    CATEGORY_NODE_COMMAND,
		Name:        "node.command_script",
		Usage:       "package.json script for given command operation. format(Template(struct{ Environment: string }))",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_SCRIPT"},
		Destination: &TL.Pipe.NodeCommand.Script,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_COMMAND,
		Name:        "node.command_script_args",
		Usage:       "package.json script arguments for given command operation. format(Template(struct{ Environment: string }))",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_SCRIPT_ARGS"},
		Destination: &TL.Pipe.NodeCommand.ScriptArgs,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_COMMAND,
		Name:        "node.command_cwd",
		Usage:       "Working directory for the given command operation.",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeCommand.Cwd,
	},
}
