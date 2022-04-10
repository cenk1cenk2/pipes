package pipe

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/sirupsen/logrus"
)

type Ctx struct {
	Tags []string
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify"},
		Task: func(t *utils.Task) error {
			Context.Tags = []string{}
			TagAsLatestForTagsRegex := []string{}
			TagAsLatestForBranchesRegex := []string{}

			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(Pipe.DockerImage.Tags.Value())) {
				tag, err := AddDockerTag(v)

				if err != nil {
					return err
				}

				t.Log.Debugln(fmt.Sprintf("Docker image tag: %s",
					tag,
				))
			}

			if Pipe.DockerImage.TagAsLatestForTagsRegex != "" {

				err := json.Unmarshal(
					[]byte(Pipe.DockerImage.TagAsLatestForTagsRegex),
					&TagAsLatestForTagsRegex,
				)

				if err != nil {
					return err
				}

				if Pipe.Git.Tag != "" {
					for _, re := range TagAsLatestForTagsRegex {
						m, err := regexp.Match(re, []byte(Pipe.Git.Tag))

						if err != nil {
							return err
						}

						if m {
							tag, err := AddDockerTag(DOCKER_LATEST_TAG)

							if err != nil {
								return err
							}

							t.Log.Infoln(
								fmt.Sprintf(
									"Will tag image as latest since tag regex matches: %s",
									tag,
								),
							)

							break
						}
					}
				}
			}

			if Pipe.DockerImage.TagAsLatestForBranchesRegex != "" {
				err := json.Unmarshal(
					[]byte(Pipe.DockerImage.TagAsLatestForBranchesRegex),
					&TagAsLatestForBranchesRegex,
				)

				if err != nil {
					return err
				}

				if Pipe.Git.Branch != "" {
					for _, re := range TagAsLatestForBranchesRegex {
						m, err := regexp.Match(re, []byte(Pipe.Git.Branch))

						if err != nil {
							return err
						}

						if m {
							tag, err := AddDockerTag(DOCKER_LATEST_TAG)

							if err != nil {
								return err
							}

							t.Log.Infoln(
								fmt.Sprintf(
									"Will tag image as latest since branch regex matches: %s",
									tag,
								),
							)

							break
						}
					}
				}
			}

			Context.Tags = utils.RemoveDuplicateStr(
				utils.DeleteEmptyStringsFromSlice(Context.Tags),
			)

			return nil
		},
	}
}

func DockerVersion() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "version"},
		Task: func(t *utils.Task) error {
			cmd := exec.Command(DOCKER_EXE, "--version")

			t.Command = cmd

			return nil
		},
	}

}

func DockerLogin() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{
			Context:        "login",
			StdOutLogLevel: logrus.DebugLevel,
			Skip: Pipe.DockerRegistry.Username == "" ||
				Pipe.DockerRegistry.Password == "",
		},
		Task: func(t *utils.Task) error {
			t.Log.Infoln(
				fmt.Sprintf("Logging in to Docker registry: %s", Pipe.DockerRegistry.Registry),
			)

			cmd := exec.Command(DOCKER_EXE, "login")

			cmd.Args = append(cmd.Args, Pipe.DockerRegistry.Registry)
			cmd.Args = append(
				cmd.Args,
				"--username",
				Pipe.DockerRegistry.Username,
				"--password-stdin",
			)

			cmd.Stdin = strings.NewReader(Pipe.DockerRegistry.Password)

			t.Command = cmd

			return nil
		},
	}
}

func DockerBuild() utils.Task {
	return utils.Task{Metadata: utils.TaskMetadata{Context: "build"},
		Task: func(t *utils.Task) error {
			t.Log.Infoln(
				fmt.Sprintf(
					"Building Docker image: %s in %s",
					Pipe.DockerFile.Name,
					Pipe.DockerFile.Context,
				),
			)

			cmd := exec.Command(DOCKER_EXE, "build")

			cmd.Args = append(cmd.Args, Pipe.DockerImage.BuildArgs.Value()...)

			if Pipe.DockerImage.Pull {
				cmd.Args = append(cmd.Args, "--pull")
			}

			for _, tag := range Context.Tags {
				cmd.Args = append(cmd.Args, "-t", tag)
			}

			cmd.Dir = Pipe.DockerFile.Context

			cmd.Args = append(
				cmd.Args,
				"--file",
				Pipe.DockerFile.Name,
				".",
			)

			t.Log.Debugln(fmt.Sprintf("CWD set as: %s", cmd.Dir))

			t.Command = cmd

			return nil
		}}
}

func DockerPush() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "push"},
		Task: func(t *utils.Task) error {
			t.Commands = []utils.Command{}

			for _, tag := range Context.Tags {
				t.Log.Infoln(
					fmt.Sprintf(
						"Pushing Docker image: %s",
						tag,
					),
				)

				cmd := exec.Command(DOCKER_EXE, "push")

				cmd.Args = append(cmd.Args, tag)

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		},
	}
}

func DockerInspect() utils.Task {
	return utils.Task{Metadata: utils.TaskMetadata{
		Context:        "inspect",
		Skip:           !Pipe.DockerImage.Inspect,
		StdOutLogLevel: logrus.DebugLevel,
	},
		Task: func(t *utils.Task) error {
			t.Commands = []utils.Command{}

			for _, tag := range Context.Tags {
				t.Log.Infoln(
					fmt.Sprintf(
						"Inspecting Docker image: %s",
						tag,
					),
				)

				cmd := exec.Command(DOCKER_EXE, "manifest", "inspect")

				cmd.Args = append(cmd.Args, tag)

				t.Commands = append(t.Commands, cmd)
			}

			return nil
		}}

}
