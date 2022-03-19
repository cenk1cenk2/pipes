package pipe

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/urfave/cli/v2"
)

type (
	Git struct {
		Branch string
		Tag    string
	}

	DockerImage struct {
		Name                        string
		Tags                        cli.StringSlice
		TagAsLatestForTagsRegex     string
		TagAsLatestForBranchesRegex string
		Pull                        bool
		Inspect                     bool
		BuildArgs                   cli.StringSlice
	}

	DockerFile struct {
		Context string
		Name    string
	}

	DockerRegistry struct {
		Registry string
		Username string
		Password string
	}

	Plugin struct {
		Git            Git
		DockerImage    DockerImage
		DockerFile     DockerFile
		DockerRegistry DockerRegistry
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{
			VerifyVariables(),
			DockerVersion(),
			DockerLogin(),
			DockerBuild(),
			DockerPush(),
			DockerInspect(),
		},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
