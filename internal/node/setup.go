package node

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
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

// Ctx is the resolved package manager the tasks build their commands out of. A
// pipe holds on to the same instance it handed to SetupTaskList, since it only
// carries anything once that list has run.
type Ctx struct {
	PackageManager
}

func NewFlags(cfg *Config) []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    CATEGORY_PACKAGE_MANAGER,
			Name:        "node.package-manager",
			Sources:     cli.EnvVars("NODE_PACKAGE_MANAGER"),
			Usage:       `Preferred Package manager for nodejs. enum("npm", "yarn", "pnpm")`,
			Required:    false,
			Value:       DEFAULT_PACKAGE_MANAGER,
			Destination: &cfg.PackageManager,
		},
	}
}

// SetupTaskList resolves the configured package manager into ctx and reports
// the versions the rest of the pipe is going to run against.
func SetupTaskList(p *plumber.Plumber, cfg *Config, ctx *Ctx) *plumber.TaskList {
	tl := &plumber.TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *plumber.TaskList) plumber.Job {
			return plumber.JobSequence(
				setupPackageManager(tl, cfg, ctx).Job(),
				packageManagerVersion(tl, ctx).Job(),
			)
		})
}

func setupPackageManager(tl *plumber.TaskList, cfg *Config, ctx *Ctx) *plumber.Task {
	return tl.CreateTask("init").
		Set(func(t *plumber.Task) error {
			ctx.PackageManager = PackageManager{
				Exe:      cfg.PackageManager,
				Commands: PackageManagers[cfg.PackageManager],
			}

			t.Log.Infof("Using package manager: %s", cfg.PackageManager)

			return nil
		})
}

func packageManagerVersion(tl *plumber.TaskList, ctx *Ctx) *plumber.Task {
	return tl.CreateTask("version").
		Set(func(t *plumber.Task) error {
			t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *plumber.Command) error {
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
				Set(func(c *plumber.Command) error {
					c.AppendArgs(ctx.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *plumber.Command) error {
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
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
