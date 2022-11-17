package run

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_COMMAND = "Command"
)

var Flags = TL.Plumber.AppendFlags(flags.NewSelectEnvEnableFlag(
	flags.SelectEnvEnableFlagSetup{
		SelectEnvEnableDestination: &TL.Pipe.Environment.Enable,
		SelectEnvEnableRequired:    false,
		SelectEnvEnableValue:       false,
	},
), []cli.Flag{

	// CATEGORY_NODE_COMMAND

	&cli.StringFlag{
		Category: CATEGORY_NODE_COMMAND,
		Name:     "node.command_script",
		Usage: fmt.Sprintf(
			"package.json script for given command operation. format(%s)",
			environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE,
		),
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_SCRIPT"},
		Destination: &TL.Pipe.NodeCommand.Script,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_COMMAND,
		Name:        "node.command_script_args",
		Usage:       fmt.Sprintf("package.json script arguments for given command operation. format(%s)", environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE),
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
})
