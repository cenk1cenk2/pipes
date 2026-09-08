package environment

import (
	"fmt"
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
)

// Ctx is what the environment task list resolves; the values only land once the
// tasks have run, so a pipe holds on to the instance it handed to SetupTaskList.
type Ctx struct {
	References  []string
	Environment string
	EnvVars     map[string]string
}

// SetupTaskList selects an environment out of the source control references and
// reads the variables belonging to it into ctx.
func SetupTaskList(p *Plumber, cfg *Config, ctx *Ctx) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldDisable(func(_ *TaskList) bool {
			return !cfg.Enable
		}).
		ShouldRunBefore(func(_ *TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				parseReferences(tl, cfg, ctx).Job(),

				selectEnvironment(tl, cfg, ctx).Job(),
				fetchEnvironment(tl, ctx).Job(),
			)
		})
}

func parseReferences(tl *TaskList, cfg *Config, ctx *Ctx) *Task {
	return tl.CreateTask("init", "references").
		Set(func(t *Task) error {
			ctx.References = cfg.Git.References()

			if cfg.FailOnNoReference && len(ctx.References) == 0 {
				return fmt.Errorf("References for the given environment has not been found.")
			}

			t.Log.Debugf("References for environment selection: %v", ctx.References)

			return nil
		})
}

func selectEnvironment(tl *TaskList, cfg *Config, ctx *Ctx) *Task {
	return tl.CreateTask("environment", "select").
		Set(func(t *Task) error {
			t.Log.Debugf("Conditions for environment variable selection: %+v", cfg.Conditions)

			selected, err := Select(cfg.Conditions, ctx.References)

			if err != nil {
				return err
			}

			ctx.Environment = selected

			if ctx.Environment != "" {
				t.Log.Infof("Environment selected: %s", ctx.Environment)

				return nil
			}

			if cfg.Strict {
				return fmt.Errorf("Environment is not selected. Can not process further on strict mode.")
			}

			t.Log.Infof("No environment has been selected.")

			return nil
		})
}

func fetchEnvironment(tl *TaskList, ctx *Ctx) *Task {
	return tl.CreateTask("environment", "fetch").
		ShouldDisable(func(_ *Task) bool {
			return ctx.Environment == ""
		}).
		Set(func(t *Task) error {
			ctx.EnvVars = Fetch(os.Environ(), ctx.Environment)

			t.Log.Infof("Environment variables that matches the current environment: %s -> %+v", ctx.Environment, ctx.EnvVars)

			return nil
		})
}
