package release

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func release(tl *TaskList) *Task {
	return tl.CreateTask("release").
		Set(func(t *Task) error {
			if P.Workspace {
				C.Exe = MultiSemanticReleaseExe
			} else {
				C.Exe = SemanticReleaseExe
			}

			t.CreateCommand(
				C.Exe,
			).
				Set(func(c *Command) error {
					// --ignore-private-packages belongs to the original multi-semantic-release, not @qiwi/multi-semantic-release.

					if P.SemanticRelease.DryRun {
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
