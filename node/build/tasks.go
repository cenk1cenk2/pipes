package build

import (
	"context"
	"os"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func build(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(ctx context.Context, t *Task) error {
			t.CreateCommand(
				setup.NodeCtx.PackageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.Run...)

					if P.Build.Script != "" {
						tmpl, err := InlineTemplate[any](P.Build.Script, nil)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.RunDelimiter...)

					if P.Build.ScriptArgs != "" {
						tmpl, err := InlineTemplate[any](P.Build.ScriptArgs, nil)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.Build.Cwd)

					c.AppendDirectEnvironment(os.Environ()...)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
