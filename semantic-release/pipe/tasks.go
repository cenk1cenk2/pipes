package pipe

import (
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
