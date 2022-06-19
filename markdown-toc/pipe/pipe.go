package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	Markdown struct {
		Patterns  cli.StringSlice
		Arguments string
	}

	Pipe struct {
		Ctx

		Markdown
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			FindMarkdownFiles(&TL).Job(),
			RunMarkdownToc(&TL).Job(),
		),
	)
}
