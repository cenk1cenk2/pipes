package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			packageManager := setup.NodeCtx.PackageManager

			t.CreateCommand(
				packageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Install.UseLockFile {
						c.AppendArgs(packageManager.Commands.InstallWithLock...)

						t.Log.Info("Using lockfile for installation.")
					} else {
						c.AppendArgs(packageManager.Commands.Install...)

						t.Log.Info("Installing dependencies without a lockfile.")
					}

					c.AppendArgs(strings.Split(P.Install.Args, " ")...)

					if P.Install.Cache {
						cacheDir := fmt.Sprintf(".%s", packageManager.Exe)
						t.Log.Info(fmt.Sprintf("Setting up cache: %s", cacheDir))

						c.AppendArgs(packageManager.Commands.Cache...)
						c.AppendArgs(cacheDir)
					}

					c.SetDir(P.Install.Cwd)

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

// pnpm keeps every version it ever fetched in its store, so the cached store only
// stays the size of the project when it drops what nothing references anymore.
func prune(tl *TaskList) *Task {
	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return !P.Install.Cache || setup.NodeCtx.PackageManager.Exe != "pnpm"
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"pnpm",
				"store",
				"prune",
				"--store-dir",
				".pnpm",
			).
				SetDir(P.Install.Cwd).
				AppendDirectEnvironment(os.Environ()...).
				AddSelfToTheTask()

			home := os.Getenv("COREPACK_HOME")
			if home == "" {
				return nil
			}

			var version string

			t.CreateCommand("pnpm", "--version").
				SetDir(P.Install.Cwd).
				AppendDirectEnvironment(os.Environ()...).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				SetIgnoreError().
				CaptureStdout(&version).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					version = strings.TrimSpace(version)
					if c.HasFailed() {
						return nil
					}

					cacheDir := filepath.Join(home, "v1", "pnpm")
					for _, dir := range []string{filepath.Clean(home), filepath.Join(home, "v1"), cacheDir} {
						info, err := os.Lstat(dir)
						if err != nil || !info.IsDir() {
							return nil
						}
					}

					entries, err := os.ReadDir(cacheDir)
					if err != nil {
						return nil
					}

					matched := false
					for _, entry := range entries {
						if entry.IsDir() && entry.Name() == version {
							matched = true
							break
						}
					}
					if !matched {
						return nil
					}

					var removed []string
					for _, entry := range entries {
						if !entry.IsDir() || entry.Name() == version {
							continue
						}

						if err := os.RemoveAll(filepath.Join(cacheDir, entry.Name())); err != nil {
							return err
						}

						removed = append(removed, entry.Name())
					}
					if len(removed) > 0 {
						t.Log.Info(fmt.Sprintf("Pruned corepack pnpm versions: %s", strings.Join(removed, ", ")))
					}

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
