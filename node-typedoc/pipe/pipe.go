package pipe

import (
	"github.com/urfave/cli/v2"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	TypeDoc struct {
		Patterns  cli.StringSlice
		Arguments string
	}

	Plugin struct {
		TypeDoc
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			FindPackages(),
			RunTypeDoc(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
