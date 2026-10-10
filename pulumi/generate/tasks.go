package generate

import (
	"context"
	json "encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

// The state is thrown away on every run, so a backend that keeps it anywhere but the
// working directory would collect a stack per pipeline, have its existing one
// overwritten, or take the Pulumi home along when its state is removed.
func backend(tl *TaskList) *Task {
	return tl.CreateTask("backend").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand("pulumi", "whoami", "--output", "json").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				AppendEnvironment(setup.C.Env).
				SetDir(setup.C.Cwd).
				CaptureStdout(&C.Whoami).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					var whoami struct {
						URL string `json:"url"`
					}
					if err := json.Unmarshal([]byte(C.Whoami), &whoami); err != nil {
						return fmt.Errorf("parse Pulumi backend: %w", err)
					}

					state, err := localState(whoami.URL)
					if err != nil {
						return fmt.Errorf("Stack %s has to use the file:// backend of the working directory: %w", stack.P.Stack, err)
					}

					C.State = state

					c.Log.Info(fmt.Sprintf("Pulumi backend: %s", whoami.URL))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func localState(url string) (string, error) {
	path, ok := strings.CutPrefix(url, "file://")
	if !ok {
		return "", fmt.Errorf("resolves to: %s", url)
	}

	cwd, err := filepath.Abs(setup.C.Cwd)
	if err != nil {
		return "", err
	}

	// pulumi resolves a relative backend against the directory it runs in, which is the
	// working directory; a home directory backend does not resolve to it.
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}

	if filepath.Clean(path) != cwd {
		return "", fmt.Errorf("resolves to: %s", url)
	}

	return cwd, nil
}

func build(tl *TaskList) *Task {
	return tl.CreateTask("build").
		ShouldDisable(func(_ *Task) bool {
			return P.Build == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand("sh", "-c", P.Build).
				AppendEnvironment(setup.C.Env).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func restate(tl *TaskList) *Task {
	return tl.CreateTask("restate").
		Set(func(_ context.Context, t *Task) error {
			paths := []string{filepath.Join(C.State, ".pulumi")}
			for _, path := range P.Paths {
				paths = append(paths, filepath.Join(setup.C.Cwd, path))
			}

			for _, path := range paths {
				t.Log.Debug(fmt.Sprintf("Removing: %s", path))

				if err := os.RemoveAll(path); err != nil {
					return err
				}
			}

			t.CreateCommand("pulumi", "stack", "init", "-s", stack.P.Stack).
				AppendEnvironment(setup.C.Env).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func generate(tl *TaskList) *Task {
	return tl.CreateTask("generate").
		Set(func(_ context.Context, t *Task) error {
			var command *Command
			if P.Command != "" {
				command = t.CreateCommand("sh", "-c", P.Command)
			} else {
				command = t.CreateCommand("pulumi", "up", "--diff", "--yes", "-f", "--stack", stack.P.Stack)
			}

			command.
				AppendEnvironment(setup.C.Env).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
