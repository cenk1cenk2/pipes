package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Plugin struct {
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
