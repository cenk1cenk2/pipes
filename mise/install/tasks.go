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

// The copy lands beside the binary and replaces it by rename, since a binary
// that is still running cannot be opened for writing but can be replaced.
func copyBinary(source string, destination string) (err error) {
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("Cannot create the directory of the mise binary: %s -> %w", dir, err)
	}

	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Cannot open the mise binary: %s -> %w", source, err)
	}
	defer src.Close()

	dst, err := os.CreateTemp(dir, ".mise-*")
	if err != nil {
		return fmt.Errorf("Cannot create the mise binary: %s -> %w", destination, err)
	}
	defer func() {
		if err != nil {
			_ = dst.Close()
			_ = os.Remove(dst.Name())
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("Cannot copy the mise binary: %s -> %w", destination, err)
	}

	if err = dst.Chmod(0o755); err != nil {
		return fmt.Errorf("Cannot make the mise binary executable: %s -> %w", destination, err)
	}

	if err = dst.Close(); err != nil {
		return fmt.Errorf("Cannot write the mise binary: %s -> %w", destination, err)
	}

	if err = os.Rename(dst.Name(), destination); err != nil {
		return fmt.Errorf("Cannot replace the mise binary: %s -> %w", destination, err)
	}

	return nil
}
