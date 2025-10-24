package lint

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func HelmLint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"lint",
			).
				Set(func(c *Command) error {
					if P.Kubernetes.Version != "" {
						c.AppendArgs("--kube-version", P.Kubernetes.Version)
					}

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func HelmTemplate(tl *TaskList) *Task {
	return tl.CreateTask("template").
		ShouldDisable(func(t *Task) bool {
			return !P.ShouldTemplate
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"template",
				".",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(setup.P.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
