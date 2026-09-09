package tool

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Pipe struct {
		Tool    string `validate:"required"`
		Args    string
		Command []string
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if len(P.Command) > 0 {
				if P.Tool == "" {
					P.Tool = P.Command[0]
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(P.Command[1:], " "))
				} else {
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(P.Command, " "))
				}
			}

			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				tool(tl).Job(),
			)
		})
}
