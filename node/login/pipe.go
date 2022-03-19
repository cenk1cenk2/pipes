package login

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type (
	Npm struct {
		Login string `validate:"json"`
		NpmRc string
	}

	Plugin struct {
		Npm Npm
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{VerifyVariables(), GenerateNpmRc(), VerifyNpmLogin()},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
