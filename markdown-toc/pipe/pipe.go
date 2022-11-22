package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
