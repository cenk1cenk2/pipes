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
	return tl.CreateTask("setup", "references").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.References = parser.ParseGitReferences(t.Pipe.Git.Tag, t.Pipe.Git.Branch)

			if t.Pipe.Environment.FailOnNoReference && t.Pipe.Ctx.References == nil {
				return fmt.Errorf("References for the given environment has not been found.")
			}

			t.Log.Debugf("References for environment selection: %v", t.Pipe.Ctx.References)

			return nil
		})
}

func SelectEnvironment(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "select").
		Set(func(t *Task[Pipe]) error {
		out:
			for _, c := range t.Pipe.Conditions {
				for _, reference := range t.Pipe.Ctx.References {
					re, err := regexp.Compile(c.Condition)

					if err != nil {
						return fmt.Errorf("Can not process regular expression for environment: %s -> %w", c.Environment, err)
					}

					t.Log.Debugf("Trying to match condition for given reference: %s with %v", reference, re.String())

					if re.MatchString(reference) {
						t.Pipe.Ctx.Environment = c.Environment

						t.Log.Infof("Environment selected: %s", c.Environment)

						break out
					}
				}
			}

			return nil
		})
}

func FetchEnvironment(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "fetch").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Ctx.Environment == ""
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
