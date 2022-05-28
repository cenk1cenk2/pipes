package build

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nochso/gomd/eol"
	"github.com/workanator/go-floc/v3"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/libraries/go-utils/utils"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	EnvironmentVariables []string
	SelectedEnvironment  string
	FallbackEnvironment  string
}

func SelectEnvironment(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "environment").ShouldDisable(func(t *Task[Pipe, Ctx]) bool {
		return len(utils.DeleteEmptyStringsFromSlice(t.Pipe.NodeBuild.EnvironmentFiles.Value())) == 0
	}).Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		t.Lock.Lock()
		t.Context.EnvironmentVariables = []string{}
		t.Lock.Unlock()

		if t.Pipe.Git.Tag != "" {
			var envConditions map[string]string
			err := json.Unmarshal([]byte(t.Pipe.NodeBuild.EnvironmentConditions), &envConditions)

			if err != nil {
				return err
			}

			for name, re := range envConditions {
				m, err := regexp.Match(re, []byte(t.Pipe.Git.Tag))

				if err != nil {
					return err
				}

				if m {
					t.Lock.Lock()
					t.Context.SelectedEnvironment = name
					t.Lock.Unlock()
				}
			}
		} else if t.Pipe.Git.Branch != "" {
			t.Lock.Lock()
			t.Context.SelectedEnvironment = t.Pipe.Git.Branch
			t.Lock.Unlock()
		} else {
			return fmt.Errorf("Can not set selected environment. Either tag name or branch name environment variable should be present.")
		}

		t.Log.Debugf("Selected environment set: %s", t.Context.SelectedEnvironment)

		if t.Pipe.NodeBuild.EnvironmentFallback != "" {
			t.Lock.Lock()
			t.Context.FallbackEnvironment = t.Pipe.NodeBuild.EnvironmentFallback
			t.Lock.Unlock()
		} else if t.Pipe.Git.Branch != "" {
			t.Lock.Lock()
			t.Context.FallbackEnvironment = t.Pipe.Git.Branch
			t.Lock.Unlock()
		} else {
			t.Log.Fatalln("Can not set fallback environment. Either manual fallback parameter should be set or brannch name environment variable should be present.")
		}

		t.Log.Debugf("Fallback environment set: %s", t.Context.FallbackEnvironment)

		return nil
	})
}

func InjectEnvironmentVariables(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "variables").ShouldDisable(func(t *Task[Pipe, Ctx]) bool {
		return len(utils.DeleteEmptyStringsFromSlice(t.Pipe.NodeBuild.EnvironmentFiles.Value())) == 0
	}).Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		EOL := eol.OSDefault().String()

		for _, file := range t.Pipe.NodeBuild.EnvironmentFiles.Value() {
			st := Task[Pipe, Ctx]{}

			st.New(tl, "inject").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
				t.Log.Infof("Injecting environment variables from: %s", file)

				cmd := Command[Pipe, Ctx]{}

				cmd.New(t, "ta-gitlab-env").Set(func(c *Command[Pipe, Ctx]) error {
					c.AppendArgs(
						"--yml-file",
						file,
						"--prefix",
						t.Context.SelectedEnvironment,
						"--fallback",
						t.Context.FallbackEnvironment,
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
							t.Log.Fatalf(
								"Can not fetch the environment variable from: %s", v,
							)
						}

						variables[i] = fmt.Sprintf("%s=%s", matches[1], matches[2])

						t.Log.Debugf("Matched from environment variable: %s -> %s",
							v,
							variables[i])
					}

					variables = utils.DeleteEmptyStringsFromSlice(variables)

					if len(variables) > 0 {
						t.Log.Debugf(
							"Injected Variables from environment file: %s%s%s",
							file,
							EOL,
							strings.Join(variables, EOL),
						)
					} else {
						t.Log.Warningf("No variables are injected from environment file: %s", file)
					}

					t.Lock.Lock()
					t.Context.EnvironmentVariables = append(t.Context.EnvironmentVariables, variables...)
					t.Lock.Unlock()

					return nil
				})

				return nil
			})

			t.ExtendSubtask(func(j floc.Job) floc.Job {
				return t.TaskList.JobParallel(j, st.Job())
			})
		}

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunSubtasks()
	})
}

func BuildNodeApplication(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "build").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		cmd := Command[Pipe, Ctx]{}

		cmd.New(t, setup.P.Context.PackageManager.Exe).Set(func(c *Command[Pipe, Ctx]) error {
			c.AppendArgs(setup.P.Context.PackageManager.Commands.Run...).
				AppendArgs(t.Pipe.NodeBuild.Script).
				AppendArgs(setup.P.Context.PackageManager.Commands.RunDelimitter...).
				AppendArgs(strings.Split(t.Pipe.NodeBuild.ScriptArgs, " ")...)

			c.SetDir(t.Pipe.NodeBuild.Cwd)

			c.AppendDirectEnvironment(os.Environ()...).
				AppendDirectEnvironment(t.Context.EnvironmentVariables...)

			return nil
		})

		t.AddCommands(cmd)

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.TaskList.RunJobs(t.GetCommandJobAsJobSequence())
	})
}
