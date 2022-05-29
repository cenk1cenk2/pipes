package run

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	NodeCommand struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		NodeCommand
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return P.New(p).ShouldRunBefore(func(tl *TaskList[Pipe], ctx *cli.Context) error {
		args := ctx.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("Arguments are needed to run a specific script.")
		} else {
			P.Pipe.NodeCommand.Script, P.Pipe.NodeCommand.ScriptArgs = args[0], strings.Join(args[1:], " ")
		}

		return nil
	}).SetTasks(
		P.JobSequence(
			RunNodeScript(&P).Job(),
		),
	)
}
