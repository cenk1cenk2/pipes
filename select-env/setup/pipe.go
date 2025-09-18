package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

type (
	Environment struct {
		Enable            bool
		Conditions        []EnvironmentConditionJson
		FailOnNoReference bool
		Strict            bool
	}

	Git flags.GitFlags

	Pipe struct {
		Environment
		Git
	}

	Ctx struct {
		References  []string
		Environment string
		EnvVars     map[string]string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldDisable(func(tl *TaskList) bool {
			return !P.Environment.Enable
		}).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				Setup(tl).Job(),

				SelectEnvironment(tl).Job(),
				FetchEnvironment(tl).Job(),
			)
		})
}
