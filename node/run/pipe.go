package run

import (
	"fmt"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Environment struct {
		Enable bool
	}

	NodeCommand struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		Environment
		NodeCommand
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetName("node", "run").
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			if TL.Pipe.NodeCommand.Script == "" {
				args := tl.CliContext.Args().Slice()
				if len(args) < 1 {
					return fmt.Errorf("Arguments are needed to run a specific script.")
				}

				TL.Pipe.NodeCommand.Script, TL.Pipe.NodeCommand.ScriptArgs = args[0], strings.Join(args[1:], " ")
			}

			return nil
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				RunNodeScript(tl).Job(),
			)
		})
}
