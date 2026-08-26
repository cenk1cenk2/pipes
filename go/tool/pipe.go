package tool

import (
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	itool "gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Pipe struct {
		Tool    string `validate:"required"`
		Args    string
		Command []string
	}

	// Deps is the resolved go tool: the directory the tool runs in and the
	// environment the cache setup has written into.
	Deps struct {
		Tool *itool.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if len(P.Command) > 0 {
				if P.Tool == "" {
					P.Tool = P.Command[0]
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(P.Command[1:], " "))
				} else {
					P.Args = fmt.Sprintf("%s %s", P.Args, strings.Join(P.Command, " "))
				}
			}

			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GoTool(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags:     Flags,
		Arguments: Arguments,
		New:       func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
