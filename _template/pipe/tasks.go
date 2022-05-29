package pipe

import (
	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
}

func DefaultTask(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("default").
		Set(func(t *Task[Pipe], c floc.Control) error {
			t.CreateCommand("echo").
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs("hello")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunCommandJobAsJobSequence()
		})
}
