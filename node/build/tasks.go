package build

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nochso/gomd/eol"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/libraries/go-utils/utils"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	EnvironmentVariables []string
	SelectedEnvironment  string
	FallbackEnvironment  string
}

func SelectEnvironment(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(
				utils.DeleteEmptyStringsFromSlice(t.Pipe.NodeBuild.EnvironmentFiles.Value()),
			) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.EnvironmentVariables = []string{}

			if t.Pipe.Git.Tag != "" {
				var envConditions map[string]string
				err := json.Unmarshal(
					[]byte(t.Pipe.NodeBuild.EnvironmentConditions),
					&envConditions,
				)

				if err != nil {
					return err
				}

				for name, re := range envConditions {
					m, err := regexp.Match(re, []byte(t.Pipe.Git.Tag))

					if err != nil {
						return err
					}

					if m {
						t.Pipe.Ctx.SelectedEnvironment = name
					}
				}
			} else if t.Pipe.Git.Branch != "" {
				t.Pipe.Ctx.SelectedEnvironment = t.Pipe.Git.Branch
			} else {
				return fmt.Errorf("Can not set selected environment. Either tag name or branch name environment variable should be present.")
			}

			t.Log.Debugf("Selected environment set: %s", t.Pipe.Ctx.SelectedEnvironment)

			if t.Pipe.NodeBuild.EnvironmentFallback != "" {
				t.Pipe.Ctx.FallbackEnvironment = t.Pipe.NodeBuild.EnvironmentFallback
			} else if t.Pipe.Git.Branch != "" {
				t.Pipe.Ctx.FallbackEnvironment = t.Pipe.Git.Branch
			} else {
				t.Log.Fatalln("Can not set fallback environment. Either manual fallback parameter should be set or brannch name environment variable should be present.")
			}

			t.Log.Debugf("Fallback environment set: %s", t.Pipe.Ctx.FallbackEnvironment)

			return nil
		})
}

func InjectEnvironmentVariables(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("variables").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(
				utils.DeleteEmptyStringsFromSlice(t.Pipe.NodeBuild.EnvironmentFiles.Value()),
			) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			EOL := eol.OSDefault().String()

			for _, file := range t.Pipe.NodeBuild.EnvironmentFiles.Value() {
				t.CreateSubtask("inject").Set(func(st *Task[Pipe]) error {
					st.Log.Infof("Injecting environment variables from: %s", file)

					st.CreateCommand("ta-gitlab-env").Set(func(c *Command[Pipe]) error {
						c.AppendArgs(
							"--yml-file",
							file,
							"--prefix",
							st.Pipe.Ctx.SelectedEnvironment,
							"--fallback",
							st.Pipe.Ctx.FallbackEnvironment,
						)

						output, err := c.Command.CombinedOutput()

						if err != nil {
							return err
						}

						variables := strings.Split(string(output), EOL)

						for i, v := range variables {
							v := strings.TrimSpace(v)

							if v == "" {
								continue
							}

							m := regexp.MustCompile(`^export ([^=]*)="([^"]*)"$`)

							matches := m.FindStringSubmatch(v)

							if len(matches) != 3 {
								st.Log.Fatalf(
									"Can not fetch the environment variable from: %s", v,
								)
							}

							variables[i] = fmt.Sprintf("%s=%s", matches[1], matches[2])

							st.Log.Debugf("Matched from environment variable: %s -> %s",
								v,
								variables[i])
						}

						variables = utils.DeleteEmptyStringsFromSlice(variables)

						if len(variables) > 0 {
							st.Log.Debugf(
								"Injected Variables from environment file: %s%s%s",
								file,
								EOL,
								strings.Join(variables, EOL),
							)
						} else {
							st.Log.Warningf("No variables are injected from environment file: %s", file)
						}

						st.Lock.Lock()
						st.Pipe.Ctx.EnvironmentVariables = append(
							st.Pipe.Ctx.EnvironmentVariables,
							variables...)
						st.Lock.Unlock()

						return nil
					})

					return nil
				}).ToParent(t, func(pt *Task[Pipe], st *Task[Pipe]) {
					pt.ExtendSubtask(func(j Job) Job {
						return tl.JobParallel(j, st.Job())
					})
				})
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func BuildNodeApplication(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(setup.P.Pipe.Ctx.PackageManager.Exe).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs(setup.P.Pipe.Ctx.PackageManager.Commands.Run...).
						AppendArgs(t.Pipe.NodeBuild.Script).
						AppendArgs(setup.P.Pipe.Ctx.PackageManager.Commands.RunDelimitter...).
						AppendArgs(strings.Split(t.Pipe.NodeBuild.ScriptArgs, " ")...)

					c.SetDir(t.Pipe.NodeBuild.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendDirectEnvironment(t.Pipe.Ctx.EnvironmentVariables...)

					return nil
				}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
