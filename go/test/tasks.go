package test

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	"gitlab.kilic.dev/devops/pipes/internal/report/cobertura"
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
			C.Failure = t.RunCommandJobAsJobSequence(ctx)

			return nil
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
			// a report path that leaves the module would have every module overwrite the
			// report of the module before.
			for _, report := range []struct{ env, path string }{
				{"GO_TEST_COVERAGE", P.Coverage},
				{"GO_TEST_JUNIT", P.JUnit},
				{"GO_TEST_COBERTURA", P.Cobertura},
			} {
				if report.path != "" && !filepath.IsLocal(report.path) {
					return fmt.Errorf("%s has to be a path inside every module of the workspace, otherwise the modules overwrite each other's report: %s", report.env, report.path)
				}
			}

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
			C.Failure = t.RunCommandJobAsJobSequence(ctx)

			return nil
		})
}

// Converts the coverage profile every tested directory holds into its Cobertura report.
func convert(tl *TaskList) *Task {
	return tl.CreateTask("cobertura").
		ShouldDisable(func(t *Task) bool {
			if P.Coverage == "" || P.Cobertura == "" {
				t.Log.Info("Skipping the Cobertura report because no coverage profile or report is configured.")

				return true
			}

			return false
		}).
		Set(func(_ context.Context, t *Task) error {
			dirs := []string{setup.C.Cwd}
			if setup.C.Workspace {
				dirs = setup.C.Modules
			}

			src := cobertura.Source{
				Adapter: cobertura.GoProfileConverter{},
				Input:   P.Coverage,
				Output:  P.Cobertura,
				Dirs:    dirs,
			}

			return src.Convert(t)
		})
}

// Fails the pipe with the outcome of the test run, which the tasks before it held back
// so the reports of a failing run are still written.
func verdict(tl *TaskList) *Task {
	return tl.CreateTask("verdict").
		ShouldDisable(func(_ *Task) bool {
			return C.Failure == nil
		}).
		Set(func(_ context.Context, _ *Task) error {
			return C.Failure
		})
}
