package node

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PACKAGE_MANAGER = "Package Manager"

	DEFAULT_PACKAGE_MANAGER = "pnpm"
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
		ShouldRunBefore(func(_ *TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				setupPackageManager(tl, cfg, ctx).Job(),
				packageManagerVersion(tl, ctx).Job(),
			)
		})
}

func setupPackageManager(tl *TaskList, cfg *Config, ctx *Ctx) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			ctx.PackageManager = PackageManager{
				Exe:      cfg.PackageManager,
				Commands: PackageManagers[cfg.PackageManager],
			}

			t.Log.Infof("Using package manager: %s", cfg.PackageManager)

			return nil
		})
}

func packageManagerVersion(tl *TaskList, ctx *Ctx) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch node.js version.")

						return nil
					}

					t.Log.Infof("node.js version: %s", stream[0])

					return nil
				}).
				AddSelfToTheTask()

			t.CreateCommand(
				ctx.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					c.AppendArgs(ctx.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch package manager version.")

						return nil
					}

					t.Log.Infof("%s version: v%s", ctx.PackageManager.Exe, stream[0])

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
