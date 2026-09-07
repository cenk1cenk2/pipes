package build

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoBuild(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			if len(P.BuildTargets) == 0 {
				P.BuildTargets = append(P.BuildTargets, GoBuildTarget{Os: runtime.GOOS, Arch: runtime.GOARCH})
			}

			linker, err := linkerFlags()
			if err != nil {
				return err
			}

			packages := C.Packages

			if len(packages) == 0 {
				packages = []string{setup.C.Cwd}
			}

			for _, module := range packages {
				for _, target := range P.BuildTargets {
					if target.Os == "" {
						target.Os = runtime.GOOS
					}
					if target.Arch == "" {
						target.Os = runtime.GOARCH
					}

					output, err := InlineTemplate(P.BinaryTemplate, map[string]string{
						"os":   target.Os,
						"arch": target.Arch,
						"name": P.BinaryName,
					})
					if err != nil {
						return fmt.Errorf("Cannot template binary name from template: %s -> %w", P.BinaryTemplate, err)
					}

					name := fmt.Sprintf("%s@%s/%s", output, target.Os, target.Arch)

					if len(packages) > 1 {
						name = fmt.Sprintf("%s:%s", filepath.Base(module), name)
					}

					t.CreateSubtask(name).
						Set(func(t *Task) error {
							t.CreateCommand(
								"go",
								"build",
								"-mod=vendor",
								"-v",
							).
								SetDir(module).
								Set(func(c *Command) error {
									t.Log.Infof("Building: %s in %s for %s/%s", P.BinaryName, module, target.Os, target.Arch)

									if !P.EnableCGO {
										c.AppendEnvironment(map[string]string{
											"CGO_ENABLED": "0",
										})
									}

									c.AppendArgs(fmt.Sprintf("-ldflags=%s", linker))

									if len(P.BuildTags) > 0 {
										c.AppendArgs("-tags", strings.Join(P.BuildTags, ","))
									}

									c.AppendArgs(strings.Split(P.Args, " ")...)

									c.AppendArgs(
										"-o",
										filepath.Join(
											P.Output,
											output,
										),
									)

									return nil
								}).
								AppendEnvironment(setup.C.Env).
								AddSelfToTheTask()

							return nil
						}).
						ShouldRunAfter(func(t *Task) error {
							return t.RunCommandJobAsJobParallel()
						}).
						AddSelfToTheParentAsParallel()
				}
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

// GoBuildPackages resolves what the workspace actually has to build. The whole
// workspace is built only when the pipeline asked for it through the flag:
// setup.C.Workspace is on whenever the toolchain merely finds itself inside a
// workspace, which is true of every child pipeline that builds a single module
// of this repository, and those have to keep building only their own.
//
// Each module is asked from inside its own directory, both because a workspace
// carries library modules that have no command to build and because the go tool
// drops a directory whose name starts with an underscore out of a package
// pattern, which is what the scaffold module lives under.
func GoBuildPackages(tl *TaskList) *Task {
	return tl.CreateTask("packages").
		ShouldDisable(func(_ *Task) bool {
			return !setup.P.Workspace || len(setup.C.Modules) == 0
		}).
		Set(func(t *Task) error {
			C.Packages = nil

			for _, module := range setup.C.Modules {
				t.CreateCommand(
					"go",
					"list",
					"-f",
					`{{if eq .Name "main"}}{{.Dir}}{{end}}`,
					"./...",
				).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
					SetDir(module).
					EnableStreamRecording().
					ShouldRunAfter(func(c *Command) error {
						for _, dir := range c.GetStdoutStream() {
							if dir := strings.TrimSpace(dir); dir != "" {
								C.Packages = append(C.Packages, dir)
							}
						}

						return nil
					}).
					AppendEnvironment(setup.C.Env).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			if len(C.Packages) == 0 {
				return fmt.Errorf("Can not resolve any commands to build out of the go workspace.")
			}

			t.Log.Infof("Commands of the workspace: %s", strings.Join(C.Packages, ", "))

			return nil
		})
}

func linkerFlags() (string, error) {
	flags := P.LinkerFlags

	for k, v := range P.BuildVariables {
		flags = fmt.Sprintf("%s -X %s=%s", flags, k, v)
	}

	linker, err := InlineTemplate[any](strings.TrimSpace(flags), nil)
	if err != nil {
		return "", fmt.Errorf("Cannot template linker flags: %s", flags)
	}

	return linker, nil
}
