package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
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
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"fmt",
				"-check",
				"-diff",
				"-recursive",
			).
				Set(func(c *Command) error {
					if P.Lint.FormatCheckArgs != "" {
						c.AppendArgs(P.Lint.FormatCheckArgs)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func lintValidate(tl *TaskList) *Task {
	return tl.CreateTask("lint", "validate").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"validate",
			).
				Set(func(c *Command) error {
					if P.Lint.ValidateArgs != "" {
						c.AppendArgs(P.Lint.ValidateArgs)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
