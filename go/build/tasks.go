package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func build(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(ctx context.Context, t *Task) error {
			if len(P.BuildTargets) == 0 {
				P.BuildTargets = append(P.BuildTargets, GoBuildTarget{Os: runtime.GOOS, Arch: runtime.GOARCH})
			}

			flags := P.LinkerFlags

			for k, v := range P.BuildVariables {
				flags = fmt.Sprintf("%s -X %s=%s", flags, k, v)
			}

			linker, err := InlineTemplate[any](strings.TrimSpace(flags), nil)
			if err != nil {
				return fmt.Errorf("Cannot template linker flags: %s", flags)
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
						target.Arch = runtime.GOARCH
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
						Set(func(_ context.Context, t *Task) error {
							t.CreateCommand(
								"go",
								"build",
								"-mod=vendor",
								"-v",
							).
								SetDir(module).
								Set(func(_ context.Context, c *Command) error {
									t.Log.Info(fmt.Sprintf("Building: %s in %s for %s/%s", P.BinaryName, module, target.Os, target.Arch))

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
						ShouldRunAfter(func(ctx context.Context, t *Task) error {
							return t.RunCommandJobAsJobParallel(ctx)
						}).
						AddSelfToTheParentAsParallel()
				}
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}

// packages resolves what the workspace has to build. Workspace mode only holds
// when the workspace file sits at the working directory itself, so a child
// pipeline building one module of a bigger workspace stays on the single module.
// The modules need no emptiness guard here, since the setup errors out when it
// resolves none of them. Each module is asked from inside its own directory,
// since the go tool drops an underscored directory out of a package pattern.
func packages(tl *TaskList) *Task {
	return tl.CreateTask("packages").
		ShouldDisable(func(_ *Task) bool {
			return !setup.C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			C.Packages = nil

			for _, module := range setup.C.Modules {
				t.CreateCommand(
					"go",
					"list",
					"-f",
					`{{if eq .Name "main"}}{{.Dir}}{{end}}`,
					"./...",
				).
					SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
					SetDir(module).
					EnableStreamRecording().
					ShouldRunAfter(func(_ context.Context, c *Command) error {
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
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				return err
			}

			if len(C.Packages) == 0 {
				return fmt.Errorf("Can not resolve any commands to build out of the go workspace.")
			}

			t.Log.Info(fmt.Sprintf("Commands of the workspace: %s", strings.Join(C.Packages, ", ")))

			return nil
		})
}
