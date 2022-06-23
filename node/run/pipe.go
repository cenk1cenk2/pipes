package run

import (
	"fmt"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	NodeCommand struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		NodeCommand
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).ShouldRunBefore(func(tl *TaskList[Pipe]) error {
		args := tl.CliContext.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("Arguments are needed to run a specific script.")
		} else {
			TL.Pipe.NodeCommand.Script, TL.Pipe.NodeCommand.ScriptArgs = args[0], strings.Join(args[1:], " ")
		}

		return nil
	}).SetTasks(
		TL.JobSequence(
			RunNodeScript(&TL).Job(),
		),
	)
}
