package setup

import (
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/git"
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
			C.References = P.Git.References()

			if P.Environment.FailOnNoReference && len(C.References) == 0 {
				return fmt.Errorf("References for the given environment has not been found.")
			}

			t.Log.Debugf("References for environment selection: %v", C.References)

			return nil
		})
}

func SelectEnvironment(tl *TaskList) *Task {
	return tl.CreateTask("environment", "select").
		Set(func(t *Task) error {
			t.Log.Debugf("Conditions for environment variable selection: %+v", P.Environment.Conditions)

			matches := make([]string, 0, len(P.Environment.Conditions))
			for _, c := range P.Environment.Conditions {
				matches = append(matches, c.Match)
			}

			matched, err := git.MatchAny(matches, C.References)
			if err != nil {
				return fmt.Errorf("Can not process regular expression for environment: %w", err)
			}

			if matched >= 0 {
				C.Environment = P.Environment.Conditions[matched].Environment

				t.Log.Infof("Environment selected: %s", C.Environment)
			}

			if P.Environment.Strict && C.Environment == "" {
				return fmt.Errorf("Environment is not selected. Can not process further on strict mode.")
			} else if C.Environment == "" {
				t.Log.Infof("No environment has been selected.")
			}

			return nil
		})
}

func FetchEnvironment(tl *TaskList) *Task {
	return tl.CreateTask("environment", "fetch").
		ShouldDisable(func(t *Task) bool {
			return C.Environment == ""
		}).
		ShouldRunBefore(func(t *Task) error {
			C.EnvVars = map[string]string{}

			return nil
		}).
		Set(func(t *Task) error {
			vars := os.Environ()

			prefix := strings.ToUpper(C.Environment)

			prefix = prefix + "_"

			for _, v := range vars {
				pair := strings.SplitN(v, "=", 2)

				key := pair[0]
				value := pair[1]

				trimmed, _ := strings.CutPrefix(key, prefix)

				C.EnvVars[trimmed] = value
			}

			C.EnvVars["ENVIRONMENT"] = C.Environment

			t.Log.Infof("Environment variables that matches the current environment: %s -> %+v", C.Environment, C.EnvVars)

			return nil
		})
}
