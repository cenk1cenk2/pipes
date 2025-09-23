package setup

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func Setup(tl *TaskList) *Task {
	return tl.CreateTask("setup").
		Set(func(t *Task) error {
			t.CreateCommand(
				"pulumi",
				"version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetCombinedStream()

					joined := strings.Join(stream, "\n")

					c.Log.Infof("Pulumi version: %s", joined)

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
