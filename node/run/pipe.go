package run

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
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

	// Deps is the package manager the script is run through and the environment
	// it is templated against.
	Deps struct {
		Node        *node.Ctx
		Environment *environment.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if P.NodeCommand.Script == "" {
				C.Script = P.Command[0]
				C.ScriptArgs = strings.Join(P.Command[1:], " ")
			} else {
				C.Script, C.ScriptArgs, _ = strings.Cut(P.NodeCommand.Script, " ")
			}

			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				RunNodeScript(tl, deps).Job(),
			)
		})
}

// Step carries the arguments as well, since the script this pipe runs is the
// argument list when no script flag was given.
func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags:     Flags,
		Arguments: Arguments,
		New:       func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
