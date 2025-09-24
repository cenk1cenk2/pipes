package run

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	NodeCommand struct {
		Script  string
		Cwd     string `validate:"dir"`
		Command []string
	}

	Pipe struct {
		NodeCommand
	}

	Ctx struct {
		Script     string
		ScriptArgs string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if P.NodeCommand.Script == "" {
				C.Script = P.Command[0]
				C.ScriptArgs = strings.Join(P.Command[1:], " ")
			} else {
				C.Script = strings.Split(P.NodeCommand.Script, " ")[0]
				C.ScriptArgs = strings.Join(strings.Split(P.NodeCommand.Script, " ")[1:], " ")
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				RunNodeScript(tl).Job(),
			)
		})
}
