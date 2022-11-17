package run

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
		Name:        "node.command_cwd",
		Usage:       "Working directory for the given command operation.",
		Required:    false,
		EnvVars:     []string{"NODE_COMMAND_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeCommand.Cwd,
	},
})

func ProcessFlags(tl *TaskList[Pipe]) error {
	if tl.Pipe.NodeCommand.Script == "" {
		args := tl.CliContext.Args().Slice()

		if len(args) < 1 {
			return fmt.Errorf("Arguments are needed to run a specific script.")
		}

		tl.Pipe.Ctx.Script = args[0]
		tl.Pipe.Ctx.ScriptArgs = strings.Join(args[1:], " ")
	} else {
		tl.Pipe.Ctx.Script = strings.Split(tl.Pipe.NodeCommand.Script, " ")[0]
		tl.Pipe.Ctx.ScriptArgs = strings.Join(strings.Split(tl.Pipe.NodeCommand.Script, " ")[1:], " ")
	}

	return nil
}
