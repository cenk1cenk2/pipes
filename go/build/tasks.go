package build

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func GoBuild(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			if len(P.BuildTargets) == 0 {
				P.BuildTargets = append(P.BuildTargets, GoBuildTarget{Os: runtime.GOOS, Arch: runtime.GOARCH})
			}

			for _, target := range P.BuildTargets {
				if target.Os == "" {
					target.Os = runtime.GOOS
				}
				if target.Arch == "" {
					target.Os = runtime.GOARCH
				}

				t.CreateSubtask(fmt.Sprintf("%s/%s", target.Os, target.Arch)).
					Set(func(t *Task) error {
						t.CreateCommand(
							"go",
							"build",
							"-mod=vendor",
						).
							SetDir(P.Cwd).
							Set(func(c *Command) error {
								if !P.EnableCGO {
									c.AppendEnvironment(map[string]string{
										"CGO_ENABLED": "0",
									})
								}

								if P.LdFlags != "" {
									c.AppendArgs(fmt.Sprintf("-ldflags=%s", P.LdFlags))
								}

								c.AppendArgs(strings.Split(P.Args, " ")...)

								output, err := InlineTemplate(P.BinaryTemplate, map[string]string{
									"os":   target.Os,
									"arch": target.Arch,
									"name": P.BinaryName,
								})
								if err != nil {
									return fmt.Errorf("Cannot template binary name from template: %s -> %w", P.BinaryTemplate, err)
								}

								c.AppendArgs(
									"-o",
									filepath.Join(
										P.Output,
										output,
									),
								)

								return nil
							}).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobParallel()
					}).
					AddSelfToTheParentAsParallel()

			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
