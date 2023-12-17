package setup

import (
	"regexp"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.EnvVars = make(map[string]string)

			return nil
		})
}

func Version(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"terraform",
				"version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command[Pipe]) error {
					stream := c.GetCombinedStream()

					joined := strings.Join(stream, "\n")

					regex := regexp.MustCompile(`Terraform (v\d+\.\d+\.\d+)`)

					matches := regex.FindStringSubmatch(joined)

					c.Log.Infof("Terraform version: %s", matches[1])

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func GenerateTerraformEnvVars(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_JOB_ID"] = t.Pipe.CiVariables.JobId
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_COMMIT_SHA"] = t.Pipe.CiVariables.CommitSha
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_JOB_STAGE"] = t.Pipe.CiVariables.JobStage
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_PROJECT_ID"] = t.Pipe.CiVariables.ProjectId
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_PROJECT_NAME"] = t.Pipe.CiVariables.ProjectName
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_PROJECT_NAMESPACE"] = t.Pipe.CiVariables.ProjectNamespace
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_PROJECT_PATH"] = t.Pipe.CiVariables.ProjectPath
			t.Pipe.Ctx.EnvVars["TF_VAR_CI_PROJECT_URL"] = t.Pipe.CiVariables.ProjectUrl

			t.Pipe.Ctx.EnvVars["TF_IN_AUTOMATION"] = "true"

			return nil
		})
}
