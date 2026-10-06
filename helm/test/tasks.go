package test

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func test(tl *TaskList) *Task {
	return tl.CreateTask("test").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"helm",
				"unittest",
			).
				Set(func(_ context.Context, c *Command) error {
					for _, file := range P.Files {
						c.AppendArgs("--file", file)
					}

					for _, values := range P.Values {
						c.AppendArgs("--values", values)
					}

					if P.Strict {
						c.AppendArgs("--strict")
					}

					if P.Output.File != "" {
						c.AppendArgs("--output-file", P.Output.File, "--output-type", P.Output.Type)
					}

					c.AppendArgs(".")

					return nil
				}).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
