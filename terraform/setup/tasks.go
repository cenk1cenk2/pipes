package setup

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
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

func DiscoverWorkspaces(tl *TaskList) *Task {
	return tl.CreateTask("workspaces").
		Set(func(t *Task) error {
			if len(P.Project.Workspaces) > 0 {
				fs := os.DirFS(P.Project.Cwd)

				t.Log.Debugf(
					"Trying to match patterns: %s",
					strings.Join(P.Project.Workspaces, ", "),
				)

				matches := []string{}

				for _, v := range P.Project.Workspaces {
					match, err := glob.Glob(fs, v)

					if err != nil {
						return err
					}

					matches = append(matches, match...)
				}

				if len(matches) == 0 {
					return fmt.Errorf(
						"Can not match any files with the given pattern: %s",
						strings.Join(P.Project.Workspaces, ", "),
					)
				}

				matches = utils.RemoveDuplicateStr(matches)

				t.Log.Infof("Paths matched for given pattern as workspace: %s", strings.Join(matches, ", "))

				P.Project.Workspaces = matches
			} else {
				P.Project.Workspaces = []string{P.Project.Cwd}

				t.Log.Debugln("Using project root as workspace since there is no defined.")
			}

			return nil
		})
}

func GenerateTerraformEnvVars(tl *TaskList) *Task {
	return tl.CreateTask("environment").
		Set(func(t *Task) error {
			C.EnvVars["TF_VAR_CI_JOB_ID"] = P.CiVariables.JobId
			C.EnvVars["TF_VAR_CI_COMMIT_SHA"] = P.CiVariables.CommitSha
			C.EnvVars["TF_VAR_CI_JOB_STAGE"] = P.CiVariables.JobStage
			C.EnvVars["TF_VAR_CI_PROJECT_ID"] = P.CiVariables.ProjectId
			C.EnvVars["TF_VAR_CI_PROJECT_NAME"] = P.CiVariables.ProjectName
			C.EnvVars["TF_VAR_CI_PROJECT_NAMESPACE"] = P.CiVariables.ProjectNamespace
			C.EnvVars["TF_VAR_CI_PROJECT_PATH"] = P.CiVariables.ProjectPath
			C.EnvVars["TF_VAR_CI_PROJECT_URL"] = P.CiVariables.ProjectUrl

			C.EnvVars["TF_IN_AUTOMATION"] = "true"

			C.EnvVars["TF_LOG"] = P.Config.LogLevel

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", C.EnvVars)

			return nil
		})
}
