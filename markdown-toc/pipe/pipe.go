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
		Markdown
	}

	Ctx struct {
		Matches []string
	}
)

var TL = TaskList{}
var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				FindMarkdownFiles(tl).Job(),
				RunMarkdownToc(tl).Job(),
			)
		})
}
