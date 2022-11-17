package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
	return TL.New(p).
		SetName("markdown-toc").
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				FindMarkdownFiles(tl).Job(),
				RunMarkdownToc(tl).Job(),
			)
		})
}
