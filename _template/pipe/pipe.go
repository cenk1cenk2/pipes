package pipe

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
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
