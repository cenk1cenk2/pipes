package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/nochso/gomd/eol"
	"gitlab.kilic.dev/devops/pipes/node/pipe"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	u "gitlab.kilic.dev/libraries/go-utils/utils"
)

type Ctx struct {
	EnvironmentVariables []string
	SelectedEnvironment  string
	FallbackEnvironment  string
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify-build"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			Context.EnvironmentVariables = []string{}

			if Pipe.Git.Tag != "" {
				var envConditions map[string]string
				err := json.Unmarshal([]byte(Pipe.NodeBuild.EnvironmentConditions), &envConditions)

				if err != nil {
					return err
				}

				for name, re := range envConditions {
					m, err := regexp.Match(re, []byte(Pipe.Git.Tag))

					if err != nil {
						return err
					}

					if m {
						Context.SelectedEnvironment = name
					}
				}
			} else if Pipe.Git.Branch != "" {
				Context.SelectedEnvironment = Pipe.Git.Branch
			} else {
				t.Log.Fatalln("Can not set selected environment. Either tag name or branch name environment variable should be present.")
			}

			t.Log.Debugf("Selected environment set: %s", Context.SelectedEnvironment)

			if Pipe.NodeBuild.EnvironmentFallback != "" {
				Context.FallbackEnvironment = Pipe.NodeBuild.EnvironmentFallback
			} else if Pipe.Git.Branch != "" {
				Context.FallbackEnvironment = Pipe.Git.Branch
			} else {
				t.Log.Fatalln("Can not set fallback environment. Either manual fallback parameter should be set or brannch name environment variable should be present.")
			}

			t.Log.Debugf("Fallback environment set: %s", Context.FallbackEnvironment)

			return nil
		},
	}
}

func InjectEnvironmentVariables() utils.Task {
	return utils.Task{Metadata: utils.TaskMetadata{
		Context: "variables",
		Skip: len(
			u.DeleteEmptyStringsFromSlice(Pipe.NodeBuild.EnvironmentFiles.Value()),
		) == 0,
	}, Task: func(t *utils.Task) error {
		var wg sync.WaitGroup
		wg.Add(len(Pipe.NodeBuild.EnvironmentFiles.Value()))

		errs := []error{}
		EOL := eol.OSDefault().String()

		for i, file := range Pipe.NodeBuild.EnvironmentFiles.Value() {
			go func(i int, file string) {
				defer wg.Done()

				t.Log.Infof("Injecting environment variables from: %s", file)

				cmd := exec.Command("ta-gitlab-env")

				cmd.Args = append(
					cmd.Args,
					"--yml-file",
					file,
					"--prefix",
					Context.SelectedEnvironment,
					"--fallback",
					Context.FallbackEnvironment,
				)

				output, err := cmd.CombinedOutput()

				if err != nil {
					errs = append(errs, errors.New(string(output)))
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

				variables = u.DeleteEmptyStringsFromSlice(variables)

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

				Context.EnvironmentVariables = append(Context.EnvironmentVariables, variables...)
			}(i, file)
		}

		wg.Wait()

		if len(errs) > 0 {
			for _, v := range errs {
				t.Log.Errorln(v)
			}

			t.Log.Fatalln("Errors encountered while injecting environment variables.")
		}

		if len(Context.EnvironmentVariables) > 0 {
			t.Log.Debugf(
				"Injected Environment Variables:%s%s",
				EOL,
				strings.Join(Context.EnvironmentVariables, EOL),
			)
		} else {
			t.Log.Warningf("No variables are injected from any of the environment files: %s", strings.Join(Pipe.NodeBuild.EnvironmentFiles.Value(), ", "))
		}

		return nil
	}}
}

func BuildNodeApplication() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "build"},
		Task: func(t *utils.Task) error {
			args := []string{}

			cmd := exec.Command(pipe.Context.PackageManager.Exe)

			args = append(args, pipe.Context.PackageManager.Commands.Run...)
			args = append(args, Pipe.NodeBuild.Script)
			args = append(args, pipe.Context.PackageManager.Commands.RunDelimitter...)
			args = append(args, strings.Split(Pipe.NodeBuild.ScriptArgs, " ")...)

			cmd.Args = append(cmd.Args, args...)

			cmd.Dir = Pipe.NodeBuild.Cwd

			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Env = append(cmd.Env, Context.EnvironmentVariables...)

			t.Command = cmd

			return nil
		},
	}
}
