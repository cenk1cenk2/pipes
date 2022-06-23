package build

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v3"
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

	Pipe struct {
		Ctx

		Git
		NodeBuild
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			SelectEnvironment(&TL).Job(),
			InjectEnvironmentVariables(&TL).Job(),
			BuildNodeApplication(&TL).Job(),
		),
	)
}
