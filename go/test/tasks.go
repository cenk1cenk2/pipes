package test

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

// Only one of the test tasks fires, since the gates below are the two halves of
// the workspace condition the setup resolved.
func test(tl *TaskList) *Task {
	return tl.CreateTask("test").
		SetJobWrapper(func(_ Job, t *Task) Job {
			return JobParallel(
				testModule(tl).Job(),
				testWorkspace(tl).Job(),
			)
		})
}

func testModule(tl *TaskList) *Task {
	return tl.CreateTask("test", "module").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			name, args := command()

			t.CreateCommand(name, args...).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDebug).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// Every module is tested from inside its own directory, since the go tool
// resolves a "<module>/..." pattern to nothing for an underscored directory.
func testWorkspace(tl *TaskList) *Task {
	return tl.CreateTask("test", "workspace").
		ShouldDisable(func(_ *Task) bool {
			return !setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			name, args := command()

			for _, module := range setup.C.Modules {
				t.CreateCommand(name, args...).
					SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDebug).
					SetDir(module).
					AppendEnvironment(setup.C.Env).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
