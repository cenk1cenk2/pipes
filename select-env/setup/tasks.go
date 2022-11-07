package setup

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/devops/pipes/common/parser"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				job,
				ParseReferences(tl).Job(),
			)
		})
}

func ParseReferences(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init", "references").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.References = parser.ParseGitReferences(t.Pipe.Git.Tag, t.Pipe.Git.Branch)

			if t.Pipe.Environment.FailOnNoReference && len(t.Pipe.Ctx.References) == 0 {
				return fmt.Errorf("References for the given environment has not been found.")
			}

			t.Log.Debugf("References for environment selection: %v", t.Pipe.Ctx.References)

			return nil
		})
}

func SelectEnvironment(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "select").
		Set(func(t *Task[Pipe]) error {
			t.Log.Debugf("Conditions for environment variable selection: %+v", t.Pipe.Environment.Conditions)

		out:
			for _, c := range t.Pipe.Environment.Conditions {
				for _, reference := range t.Pipe.Ctx.References {
					re, err := regexp.Compile(c.Match)

					if err != nil {
						return fmt.Errorf("Can not process regular expression for environment: %s -> %w", c.Environment, err)
					}

					if re.MatchString(reference) {
						t.Pipe.Ctx.Environment = c.Environment

						t.Log.Infof("Environment selected: %s", c.Environment)

						break out
					}
				}
			}

			if t.Pipe.Environment.Strict && t.Pipe.Ctx.Environment == "" {
				return fmt.Errorf("Environment is not selected. Can not process further on strict mode.")
			} else if t.Pipe.Ctx.Environment == "" {
				t.Log.Infof("No environment has been selected.")
			}

			return nil
		})
}

func FetchEnvironment(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "fetch").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Ctx.Environment == ""
		}).
		ShouldRunBefore(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.EnvVars = map[string]string{}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			vars := os.Environ()

			prefix := strings.ToUpper(t.Pipe.Ctx.Environment)

			prefix = prefix + "_"

			for _, v := range vars {
				pair := strings.SplitN(v, "=", 2)

				key := pair[0]
				value := pair[1]

				if strings.HasPrefix(key, prefix) {
					trimmed := strings.TrimPrefix(key, prefix)

					t.Pipe.Ctx.EnvVars[trimmed] = value
				}
			}

			t.Pipe.Ctx.EnvVars["ENVIRONMENT"] = t.Pipe.Ctx.Environment

			t.Log.Debugf("Environment variables that matches the current environment: %s -> %+v", t.Pipe.Ctx.Environment, t.Pipe.Ctx.EnvVars)

			return nil
		})
}
