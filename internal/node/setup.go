package node

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
)

//revive:disable:line-length-limit

const (
	CategoryPackageManager = "Package Manager"

	DefaultPackageManager = "pnpm"
)

// Config is the package manager a pipe was configured with.
type Config struct {
	PackageManager string `validate:"oneof=npm yarn pnpm"`
}

// Ctx is the resolved package manager the tasks build their commands out of; it
// only carries anything once SetupTaskList has run.
type Ctx struct {
	PackageManager
}

// SetupTaskList resolves the configured package manager into ctx and reports
// the versions the rest of the pipe is going to run against.
func SetupTaskList(p *Plumber, cfg *Config, ctx *Ctx) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl, cfg, ctx).Job(),
				version(tl, ctx).Job(),
			)
		})
}

func initialize(tl *TaskList, cfg *Config, ctx *Ctx) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			ctx.PackageManager = PackageManager{
				Exe:      cfg.PackageManager,
				Commands: PackageManagers[cfg.PackageManager],
			}

			t.Log.Info(fmt.Sprintf("Using package manager: %s", cfg.PackageManager))

			return nil
		})
}

func version(tl *TaskList, ctx *Ctx) *Task {
	return tl.CreateTask("version").
		Set(func(_ context.Context, t *Task) error {
			var nodeVersion string

			t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				CaptureOutput(&nodeVersion).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					if nodeVersion == "" {
						t.Log.Debug("Can not fetch node.js version.")

						return nil
					}

					t.Log.Info(fmt.Sprintf("node.js version: %s", nodeVersion))

					return nil
				}).
				AddSelfToTheTask()

			var packageManagerVersion string

			t.CreateCommand(
				ctx.PackageManager.Exe,
			).
				Set(func(_ context.Context, c *Command) error {
					c.AppendArgs(ctx.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				CaptureOutput(&packageManagerVersion).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					if packageManagerVersion == "" {
						t.Log.Debug("Can not fetch package manager version.")

						return nil
					}

					t.Log.Info(fmt.Sprintf("%s version: v%s", ctx.PackageManager.Exe, packageManagerVersion))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobParallel(ctx)
		})
}
