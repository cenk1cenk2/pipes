package run

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "node.command_script",
		Usage:       "package.json script for given command operation.",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_SCRIPT"},
		Destination: &Pipe.NodeCommand.Script,
	},

	&cli.StringFlag{
		Name:        "node.command_script_args",
		Usage:       "package.json script arguments for given command operation.",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_SCRIPT_ARGS"},
		Destination: &Pipe.NodeCommand.ScriptArgs,
	},

	&cli.StringFlag{
		Name:        "node.command_cwd",
		Usage:       "Working directory for the given command operation.",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_CWD"},
		Value:       ".",
		Destination: &Pipe.NodeCommand.Cwd,
	},
}
