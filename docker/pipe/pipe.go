package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v3"
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

	Pipe struct {
		Ctx

		Git
		Docker
		DockerImage
		DockerFile
		DockerRegistry
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			Setup(&TL).Job(),
			DockerVersion(&TL).Job(),
			DockerLogin(&TL).Job(),
			DockerSetupBuildX(&TL).Job(),
			DockerBuildX(&TL).Job(),
			DockerBuild(&TL).Job(),
			DockerPush(&TL).Job(),
			DockerInspect(&TL).Job(),
		),
	)
}
