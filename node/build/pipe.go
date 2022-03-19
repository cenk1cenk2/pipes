package build

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/urfave/cli/v2"
)

type (
	Git struct {
		Branch string
		Tag    string
	}

	NodeBuild struct {
		Script                string
		ScriptArgs            string
		Cwd                   string `validate:"dir"`
		EnvironmentFiles      cli.StringSlice
		EnvironmentFallback   string
		EnvironmentConditions string `validate:"json"`
	}

	Plugin struct {
		Git       Git
		NodeBuild NodeBuild
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			InjectEnvironmentVariables(),
			BuildNodeApplication(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
