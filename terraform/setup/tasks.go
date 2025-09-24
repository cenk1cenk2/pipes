package setup

import (
	"regexp"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func Version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"terraform",
				"version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
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
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func GenerateTerraformEnvVars(tl *TaskList) *Task {
	return tl.CreateTask("environment").
		Set(func(t *Task) error {
			C.EnvVars["TF_IN_AUTOMATION"] = "true"

			C.EnvVars["TF_LOG"] = P.Config.LogLevel

			C.EnvVars["TF_VAR_CI_API_V4_URL"] = P.CiVariables.ApiUrl
			C.EnvVars["TF_VAR_CI_PROJECT_ID"] = P.CiVariables.ProjectId

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", C.EnvVars)

			return nil
		})
}
