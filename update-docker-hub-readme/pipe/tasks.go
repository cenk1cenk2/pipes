package pipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/docker/docker/api/types/registry"
)

func LoginToDockerHubRegistry(tl *TaskList) *Task {
	return tl.CreateTask("login").
		Set(func(t *Task) error {
			login, err := json.Marshal(registry.AuthConfig{
				Username: P.DockerHub.Username,
				Password: P.DockerHub.Password,
			})

			if err != nil {
				return err
			}

			res, err := http.Post(
				"https://hub.docker.com/v2/users/login/",
				JSON_REQUEST,
				bytes.NewReader(login),
			)

			if err != nil {
				return err
			}

			t.Log.Debugln("Authentication token obtained.")

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)

			if err != nil {
				return err
			}

			b := DockerHubLoginResponse{}
			if err = json.Unmarshal(body, &b); err != nil {
				return err
			}

			C.Token = b.Token

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

//ignore:funlen
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

						update := DockerHubUpdateReadmeRequest{
							Readme: string(content),
						}

						if readme.Description != "" {
							update.Description = readme.Description
						}

						body, err := json.Marshal(
							update,
						)

						if err != nil {
							return err
						}

						req, err := http.NewRequest(http.MethodPatch,
							fmt.Sprintf("%s/%s/", P.DockerHub.Address, repository),
							bytes.NewReader(body),
						)

						req = AddAuthenticationHeadersToRequest(t, req)

						if err != nil {
							return err
						}

						res, err := http.DefaultClient.Do(req)

						if err != nil {
							return err
						}

						t.Log.Debugf("Status Code: %d", res.StatusCode)

						defer res.Body.Close()

						body, err = io.ReadAll(res.Body)

						if err != nil {
							return err
						}

						t.Log.Debugf("Response body: %s", string(body))

						b := DockerHubUpdateReadmeResponse{}
						err = json.Unmarshal(body, &b)

						if err != nil {
							return fmt.Errorf("Response unexpected: %w > %s", err, string(body))
						}

						switch res.StatusCode {
						case http.StatusOK:
							if b.FullDescription != string(content) {
								return fmt.Errorf("Uploaded README does not match with current repository README file.")
							}

							if readme.Description != "" && b.Description != readme.Description {
								return fmt.Errorf("Uploaded README does not match with current repository README file.")
							}

							t.Log.Infof(
								"Successfully pushed readme file to: %s > %s/%s",
								readme.File,
								P.DockerHub.Address,
								repository,
							)
						case http.StatusNotFound:
							return fmt.Errorf(
								"Repository does not exists: %s/%s",
								P.DockerHub.Address,
								repository,
							)
						default:
							if !b.CanEdit {
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
