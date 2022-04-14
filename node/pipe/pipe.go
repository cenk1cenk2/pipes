package pipe

import (
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Node struct {
		PackageManager string `validate:"oneof=npm yarn"`
	}

	Plugin struct {
		Node Node
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{VerifyVariables()},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
