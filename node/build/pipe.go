package build

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v2"
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
		Git
		NodeBuild
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return P.New(p).SetTasks(
		P.JobSequence(
			SelectEnvironment(&P).Job(),
			InjectEnvironmentVariables(&P).Job(),
			BuildNodeApplication(&P).Job(),
		),
	)
}
