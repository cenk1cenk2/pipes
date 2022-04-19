package login

import (
	"github.com/urfave/cli/v2"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Npm struct {
		Login     string `validate:"json"`
		NpmRcFile cli.StringSlice
		NpmRc     string
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
