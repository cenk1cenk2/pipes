package pipe

import (
	"github.com/urfave/cli/v2"

	login "gitlab.kilic.dev/devops/pipes/node/login"
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
	if err := login.Pipe.Exec(); err != nil {
		return err
	}

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
