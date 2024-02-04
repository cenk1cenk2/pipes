package setup

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

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

func DiscoverWorkspaces(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("workspaces").
		Set(func(t *Task[Pipe]) error {
			if len(t.Pipe.Project.Workspaces) > 0 {
				fs := os.DirFS(t.Pipe.Project.Cwd)

				t.Log.Debugf(
					"Trying to match patterns: %s",
					strings.Join(t.Pipe.Project.Workspaces, ", "),
				)

				matches := []string{}

				for _, v := range t.Pipe.Project.Workspaces {
					match, err := glob.Glob(fs, v)

					if err != nil {
						return err
					}

					matches = append(matches, match...)
				}

				if len(matches) == 0 {
					return fmt.Errorf(
						"Can not match any files with the given pattern: %s",
						strings.Join(t.Pipe.Project.Workspaces, ", "),
					)
				}

				matches = utils.RemoveDuplicateStr(matches)

				t.Log.Infof("Paths matched for given pattern as workspace: %s", strings.Join(matches, ", "))

				t.Pipe.Project.Workspaces = matches
			} else {
				t.Pipe.Project.Workspaces = []string{t.Pipe.Project.Cwd}

				t.Log.Debugln("Using project root as workspace since there is no defined.")
			}

			return nil
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

			t.Pipe.Ctx.EnvVars["TF_LOG"] = t.Pipe.Config.LogLevel

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", t.Pipe.Ctx.EnvVars)

			return nil
		})
}
