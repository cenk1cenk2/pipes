package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Markdown struct {
		Patterns       []string
		StartDepth     int
		EndDepth       int
		Indentation    int
		ListIdentifier string
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
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				FindMarkdownFiles(tl).Job(),
				RunMarkdownToc(tl).Job(),
			)
		})
}
