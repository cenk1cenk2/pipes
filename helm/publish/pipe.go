package publish

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

type (
	Chart struct {
		Target            string
		Versions          []string
		VersionFile       string
		VersionFileStrict bool
		VersionsSanitize  []versions.Match
		VersionsTemplate  []versions.Match
		Destination       string `validate:"dirpath"`
		AppVersion        string
	}

	Pipe struct {
		Chart
	}

	Ctx struct {
		Versions []string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				versionsTask(tl).Job(),
				packageTask(tl).Job(),
				publish(tl).Job(),
			)
		})
}
