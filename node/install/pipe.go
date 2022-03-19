package install

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type (
	NodeInstall struct {
		Cwd         string `validate:"dir"`
		UseLockFile bool
	}

	Plugin struct {
		NodeInstall NodeInstall
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			InstallNodeDependencies(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
