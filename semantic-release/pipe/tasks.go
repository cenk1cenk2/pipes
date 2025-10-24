package pipe

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

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

					if len(P.SemanticRelease.IsolateWorkspaces) > 0 {
						c.AppendArgs(
							"--ignore-packages",
							strings.Join(P.SemanticRelease.IsolateWorkspaces, ","),
						)
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
