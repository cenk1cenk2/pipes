package run

import (
	"strings"

	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	"github.com/urfave/cli/v2"
)

type (
	NodeCommand struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Plugin struct {
		NodeCommand NodeCommand
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec(c *cli.Context) error {
	args := c.Args().Slice()
	Pipe.NodeCommand.Script, Pipe.NodeCommand.ScriptArgs = args[0], strings.Join(args[1:], " ")

	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			RunNodeScript(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
