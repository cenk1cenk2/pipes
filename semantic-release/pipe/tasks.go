package pipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func IsolateWorkspaces(tl *TaskList) *Task {
	return tl.CreateTask("isolate").
		ShouldDisable(func(t *Task) bool {
			return len(P.SemanticRelease.IsolateWorkspaces) == 0
		}).
		Set(func(t *Task) error {
			t.Log.Infof("Isolating workspaces: %s", strings.Join(P.SemanticRelease.IsolateWorkspaces, ", "))

			t.Log.Debugf("Reading package.json file.")

			f, err := os.ReadFile("package.json")
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("package.json file does not exist, skipping isolation.")
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(f, &parsed); err != nil {
				return fmt.Errorf("failed to parse package.json: %v", err)
			}

			parsed["workspaces"].(map[string]interface{})["packages"] = P.SemanticRelease.IsolateWorkspaces

			t.Log.Debugf("Writing updated package.json file.")
			updated, err := json.MarshalIndent(parsed, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal updated package.json: %v", err)
			}

			if err := os.WriteFile("package.json", updated, 0644); err != nil {
				return fmt.Errorf("failed to write updated package.json: %v", err)
			}

			t.Log.Infof("Updated package.json file workspaces: %s", strings.Join(P.SemanticRelease.IsolateWorkspaces, ", "))

			return nil
		})
}

func RunSemanticRelease(tl *TaskList) *Task {
	return tl.CreateTask("release").
		Set(func(t *Task) error {
			if P.Workspace {
				C.Exe = MULTI_SEMANTIC_RELEASE_EXE
			} else {
				C.Exe = SEMANTIC_RELEASE_EXE
			}

			t.CreateCommand(
				C.Exe,
			).
				Set(func(c *Command) error {
					// this should be added for original multi-semantic-release and not the @qiwi/multi-semantic-release
					// if P.SemanticRelease.Workspace {
					// 	c.AppendArgs("--ignore-private-packages")
					// }

					if P.SemanticRelease.IsDryRun {
						c.AppendEnvironment(map[string]string{
							// detected by the following rules, have to disable them to trick
							// https://github.com/semantic-release/env-ci/tree/master/services
							"CI":        "false",
							"GITLAB_CI": "false",
						})
						c.AppendArgs("--dry-run", "--no-ci", "--branches", P.CI.CommitReference)
					}

					if t.Plumber.Environment.Debug {
						c.AppendArgs("--debug")
					}

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
