package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/parser"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

func Setup(tl *TaskList) *Task {
	return tl.CreateTask("init").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				job,
				ParseReferences(tl).Job(),
			)
		})
}

func ParseReferences(tl *TaskList) *Task {
	return tl.CreateTask("init", "references").
		Set(func(t *Task) error {
			C.References = parser.ParseGitReferences(P.Git.Tag, P.Git.Branch)

			t.Log.Debugf("References for environment selection: %v", C.References)

			return nil
		})
}

func DockerInspect(tl *TaskList) *Task {
	return tl.CreateTask("inspect").
		ShouldDisable(func(t *Task) bool {
			return !P.DockerImage.Inspect
		}).
		Set(func(t *Task) error {
			for _, tag := range C.Tags {
				t.CreateCommand(
					setup.DOCKER_EXE,
					"manifest",
					"inspect",
					tag,
				).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
					Set(func(c *Command) error {
						c.Log.Infof(
							"Inspecting Docker image: %s",
							tag,
						)

						return nil
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
