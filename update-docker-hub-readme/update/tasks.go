package update

import (
	"context"
	"fmt"
	"net/http"
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
)

func LoginToDockerHubRegistry(tl *TaskList) *Task {
	return tl.CreateTask("login").
		Set(func(t *Task) error {
			token, err := C.Hub.Login(
				context.Background(),
				P.DockerHub.Username,
				P.DockerHub.Password,
			)

			if err != nil {
				return err
			}

			t.Log.Debugln("Authentication token obtained.")

			C.Token = token

			return nil
		})
}

func DiscoverJobs(tl *TaskList) *Task {
	return tl.CreateTask("discover").
		Set(func(t *Task) error {
			if P.Readme.Repository != "" {
				C.Readme[P.Readme.Repository] = ParsedReadme{
					File:        P.Readme.File,
					Description: P.Readme.Description,
				}
			}

			if len(P.Readme.Matrix) > 0 {
				for _, readme := range P.Readme.Matrix {
					C.Readme[readme.Repository] = ParsedReadme{
						File:        readme.File,
						Description: readme.Description,
					}
				}
			}

			return nil
		})
}

// VerifyReadme decides whether the readme actually landed on the repository.
// The service answers 200 on a readme it did not take, so what came back rather
// than the status alone is what proves it.
func VerifyReadme(res hub.Result, repository string, readme ParsedReadme, content string) error {
	switch res.StatusCode {
	case http.StatusOK:
		if res.FullDescription != content {
			return fmt.Errorf("Uploaded README does not match with current repository README file.")
		}

		if readme.Description != "" && res.Description != readme.Description {
			return fmt.Errorf("Uploaded README does not match with current repository README file.")
		}

		return nil
	case http.StatusNotFound:
		return fmt.Errorf(
			"Repository does not exists: %s/%s",
			P.DockerHub.Address,
			repository,
		)
	default:
		if !res.CanEdit {
			return fmt.Errorf(
				"Given user credentials do not have permission to edit repository: %s/%s",
				P.DockerHub.Address,
				repository,
			)
		}

		return fmt.Errorf(
			"Pushing readme failed with code: %d",
			res.StatusCode,
		)
	}
}

func UpdateDockerReadme(tl *TaskList) *Task {
	return tl.CreateTask("update").
		Set(func(t *Task) error {
			for repository, readme := range C.Readme {
				t.CreateSubtask(repository).
					Set(func(t *Task) error {
						t.Log.Debugf(
							"Running against repository: %s/%s",
							P.DockerHub.Address,
							repository,
						)

						t.Log.Debugf("Trying to read file: %s", readme.File)

						content, err := os.ReadFile(readme.File)

						if err != nil {
							return err
						}

						res, err := C.Hub.UpdateReadme(
							context.Background(),
							C.Token,
							repository,
							hub.Readme{
								Description: readme.Description,
								Full:        string(content),
							},
						)

						if err != nil {
							return err
						}

						t.Log.Debugf("Status Code: %d", res.StatusCode)

						if err := VerifyReadme(res, repository, readme, string(content)); err != nil {
							return err
						}

						t.Log.Infof(
							"Successfully pushed readme file to: %s > %s/%s",
							readme.File,
							P.DockerHub.Address,
							repository,
						)

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
