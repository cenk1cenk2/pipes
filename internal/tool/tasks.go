package tool

import (
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

// Setup resolves the working directory and probes the tool version into ctx. The
// extra tasks run after those, in the order they are given, which is where a pipe
// puts whatever else has to be in place before its own commands.
func Setup(
	p *plumber.Plumber,
	spec Spec,
	cfg *Config,
	ctx *Ctx,
	extra ...func(*plumber.TaskList) *plumber.Task,
) *plumber.TaskList {
	tl := &plumber.TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return cli.Validated(p, cfg)
		}).
		Set(func(tl *plumber.TaskList) plumber.Job {
			jobs := []plumber.Job{
				initialize(tl, cfg, ctx).Job(),
				version(tl, spec, ctx).Job(),
			}

			for _, task := range extra {
				jobs = append(jobs, task(tl).Job())
			}

			return plumber.JobSequence(jobs...)
		})
}

func initialize(tl *plumber.TaskList, cfg *Config, ctx *Ctx) *plumber.Task {
	return tl.CreateTask("init").
		Set(func(t *plumber.Task) error {
			ctx.Cwd = cfg.Cwd

			t.Log.Debugf("Working directory: %s", ctx.Cwd)

			return nil
		})
}

func version(tl *plumber.TaskList, spec Spec, ctx *Ctx) *plumber.Task {
	return tl.CreateTask("version").
		Set(func(t *plumber.Task) error {
			t.CreateCommand(spec.Name, spec.VersionArgs...).
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *plumber.Command) error {
					ctx.Version = spec.ParseVersion(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Infof("%s version: %s", spec.Name, ctx.Version)

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
