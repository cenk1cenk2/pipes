package run

import (
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	NodeCommand struct {
		Script string
		Cwd    string `validate:"dir"`
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
			if err := p.Validate(P); err != nil {
				return err
			}

			if P.NodeCommand.Script == "" {
				args := p.Cli.Args().Slice()

				if len(args) < 1 {
					return fmt.Errorf("Arguments are needed to run a specific script.")
				}

				C.Script = args[0]
				C.ScriptArgs = strings.Join(args[1:], " ")
			} else {
				C.Script = strings.Split(P.NodeCommand.Script, " ")[0]
				C.ScriptArgs = strings.Join(strings.Split(P.NodeCommand.Script, " ")[1:], " ")
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				RunNodeScript(tl).Job(),
			)
		})
}
