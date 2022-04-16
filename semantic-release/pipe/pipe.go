package pipe

import (
	"github.com/urfave/cli/v2"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Packages struct {
		Apk  cli.StringSlice
		Node cli.StringSlice
	}

	SemanticRelease struct {
		IsDryRun bool
		UseMulti bool
	}

	Plugin struct {
		Packages
		SemanticRelease
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			InstallPackages(),
			RunSemanticRelease(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
