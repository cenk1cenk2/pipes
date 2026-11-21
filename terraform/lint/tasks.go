package lint

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/shell"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func lint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				lintFmt(t.TL).Job(),
				lintValidate(t.TL).Job(),
			)
		})
}

func lintFmt(tl *TaskList) *Task {
	return tl.CreateTask("lint", "fmt").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"terraform",
				"fmt",
				"-check",
				"-diff",
				"-recursive",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Lint.FormatCheckArgs != "" {
						args, err := shell.Split(P.Lint.FormatCheckArgs)
						if err != nil {
							return err
						}

						c.AppendArgs(args...)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func lintValidate(tl *TaskList) *Task {
	return tl.CreateTask("lint", "validate").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"terraform",
				"validate",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Lint.ValidateArgs != "" {
						args, err := shell.Split(P.Lint.ValidateArgs)
						if err != nil {
							return err
						}

						c.AppendArgs(args...)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
