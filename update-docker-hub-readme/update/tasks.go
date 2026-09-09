package update

import (
	"context"
	"fmt"
	"net/http"
	"os"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
)

func login(tl *TaskList) *Task {
	return tl.CreateTask("login").
		Set(func(_ context.Context, t *Task) error {
			token, err := C.Hub.Login(
				context.Background(),
				P.DockerHub.Username,
				P.DockerHub.Password,
			)

			if err != nil {
				return err
			}

			t.Log.Debug("Authentication token obtained.")

			C.Token = token

			return nil
		})
}

func discover(tl *TaskList) *Task {
	return tl.CreateTask("discover").
		Set(func(_ context.Context, t *Task) error {
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

// verifyReadme decides whether the readme actually landed. The service answers 200
// on a readme it did not take, so the response body and not the status proves it.
func verifyReadme(res hub.Result, repository string, readme ParsedReadme, content string) error {
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

func update(tl *TaskList) *Task {
	return tl.CreateTask("update").
		Set(func(_ context.Context, t *Task) error {
			for repository, readme := range C.Readme {
				t.CreateSubtask(repository).
					Set(func(_ context.Context, t *Task) error {
						t.Log.Debug(fmt.Sprintf("Running against repository: %s/%s",
							P.DockerHub.Address,
							repository),
						)

						t.Log.Debug(fmt.Sprintf("Trying to read file: %s", readme.File))

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

						t.Log.Debug(fmt.Sprintf("Status Code: %d", res.StatusCode))

						if err := verifyReadme(res, repository, readme, string(content)); err != nil {
							return err
						}

						t.Log.Info(fmt.Sprintf("Successfully pushed readme file to: %s > %s/%s",
							readme.File,
							P.DockerHub.Address,
							repository),
						)

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}
