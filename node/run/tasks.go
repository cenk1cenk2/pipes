package run

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/shell"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func run(tl *TaskList) *Task {
	return tl.CreateTask("run", C.Script).
		Set(func(ctx context.Context, t *Task) error {
			t.CreateCommand(
				setup.NodeCtx.PackageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.Run...)

					if C.Script != "" {
						tmpl, err := InlineTemplate[any](C.Script, nil)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.RunDelimiter...)

					if C.ScriptArgs != "" {
						tmpl, err := InlineTemplate[any](C.ScriptArgs, nil)

						if err != nil {
							return err
						}

						args, err := shell.Split(tmpl)
						if err != nil {
							return err
						}

						c.AppendArgs(args...)
					}

					// arguments after the double dash reach the pipe already separated, so each one is rendered on its own.
					args, err := InlineTemplates[any](C.Args, nil)
					if err != nil {
						return err
					}

					c.AppendArgs(args...)

					c.SetDir(P.Run.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
