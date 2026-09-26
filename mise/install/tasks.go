package install

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/mise/setup"
)

// A copy that no longer runs, as one built for another platform, reports no
// version and is replaced like any other stale one.
func binary(tl *TaskList) *Task {
	var output string

	return tl.CreateTask("binary").
		Set(func(_ context.Context, t *Task) error {
			output = ""

			if _, err := os.Stat(setup.C.Binary); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("Cannot probe for the mise binary: %s -> %w", setup.C.Binary, err)
				}

				return nil
			}

			t.CreateCommand(setup.C.Binary, "--version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				SetIgnoreError().
				CaptureStdout(&output).
				SetMaskOsEnvironment().
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				return err
			}

			if setup.ParseVersion(output) == setup.C.Version {
				t.Log.Info(fmt.Sprintf("mise binary is already at version %s: %s", setup.C.Version, setup.C.Binary))

				return nil
			}

			t.Log.Info(fmt.Sprintf("Copying mise binary: %s -> %s", setup.C.Executable, setup.C.Binary))

			return copyBinary(setup.C.Executable, setup.C.Binary)
		})
}

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				setup.C.Binary,
				"install",
			).
				SetLogLevel(LogLevelDefault, LogLevelInfo, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Installing tools: in %s", setup.C.Cwd))

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				SetMaskOsEnvironment().
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func prune(tl *TaskList) *Task {
	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return !P.Prune
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				setup.C.Binary,
				"prune",
				"--yes",
			).
				SetLogLevel(LogLevelDefault, LogLevelInfo, LogLevelDefault).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info(fmt.Sprintf("Pruning unused tool versions: in %s", setup.C.Cwd))

					return nil
				}).
				SetMaskOsEnvironment().
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func list(tl *TaskList) *Task {
	return tl.CreateTask("list").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				setup.C.Binary,
				"ls",
				"--current",
			).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDebug).
				SetDir(setup.C.Cwd).
				Set(func(_ context.Context, c *Command) error {
					t.Log.Info("Installed tool versions:")

					return nil
				}).
				SetMaskOsEnvironment().
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func copyBinary(source string, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("Cannot create the directory of the mise binary: %s -> %w", filepath.Dir(destination), err)
	}

	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Cannot open the mise binary: %s -> %w", source, err)
	}
	defer src.Close()

	dst, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("Cannot create the mise binary: %s -> %w", destination, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("Cannot copy the mise binary: %s -> %w", destination, err)
	}

	// a file that already existed keeps its mode through the open above.
	if err := dst.Chmod(0o755); err != nil {
		return fmt.Errorf("Cannot make the mise binary executable: %s -> %w", destination, err)
	}

	return dst.Close()
}
