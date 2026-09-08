package setup

import (
	"regexp"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debugf("Working directory: %s", C.Cwd)

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			pattern := regexp.MustCompile(`Terraform (v\d+\.\d+\.\d+)`)

			t.CreateCommand("terraform", "version").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
					output := strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					// the banner is only ever logged, and terraform has already proven it
					// runs by answering at all, so an unrecognised one is reported whole
					// rather than treated as an error.
					if matches := pattern.FindStringSubmatch(output); len(matches) > 1 {
						C.Version = matches[1]
					} else {
						C.Version = output
					}

					c.Log.Infof("terraform version: %s", C.Version)

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
			C.Env["TF_IN_AUTOMATION"] = "true"

			C.Env["TF_LOG"] = P.LogLevel

			C.Env["TF_VAR_CI_API_V4_URL"] = P.CiVariables.ApiUrl
			C.Env["TF_VAR_CI_PROJECT_ID"] = P.CiVariables.ProjectId

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", C.Env)

			return nil
		})
}
