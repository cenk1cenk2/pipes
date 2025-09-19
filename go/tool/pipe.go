package tool

import (
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Tool string `validate:"required"`
		Args string
	}

	Ctx struct {
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			args := p.Cli.Args().Tail()
			if len(args) > 0 {
				if P.Tool == "" {

					if len(args) < 1 {
						return fmt.Errorf("Arguments are needed to run a specific script.")
					}

					P.Tool = args[0]
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(args[1:], " "))
				} else {
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(args, " "))
				}
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoTool(tl).Job(),
			)
		})
}
