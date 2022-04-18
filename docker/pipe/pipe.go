package pipe

import (
	"github.com/urfave/cli/v2"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type (
	Git struct {
		Branch string
		Tag    string
	}

	Docker struct {
		UseBuildx       bool
		BuildxPlatforms string
	}

	DockerImage struct {
		Name                        string
		Tags                        cli.StringSlice
		TagAsLatestForTagsRegex     string
		TagAsLatestForBranchesRegex string
		TagsFile                    string
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
		Git
		Docker
		DockerImage
		DockerFile
		DockerRegistry
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
