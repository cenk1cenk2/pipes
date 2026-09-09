package install

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

// Only one of the vendor tasks fires, since the gates below are the two halves
// of the workspace condition the setup resolved.
func vendor(tl *TaskList) *Task {
	return tl.CreateTask("vendor").
		SetJobWrapper(func(_ Job, t *Task) Job {
			return JobParallel(
				vendorModule(tl).Job(),
				vendorWorkspace(tl).Job(),
			)
		})
}

func vendorModule(tl *TaskList) *Task {
	return tl.CreateTask("vendor", "module").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"vendor",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Vendoring: in %s", setup.C.Cwd))

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func vendorWorkspace(tl *TaskList) *Task {
	return tl.CreateTask("vendor", "workspace").
		ShouldDisable(func(_ *Task) bool {
			return !setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"work",
				"vendor",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Vendoring workspace: in %s", setup.C.Cwd))

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func verify(tl *TaskList) *Task {
	return tl.CreateTask("verify").
		ShouldDisable(func(t *Task) bool {
			return !P.Verify
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"verify",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Verifying modules: in %s", setup.C.Cwd))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
