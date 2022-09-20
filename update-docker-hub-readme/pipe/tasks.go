package pipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			if len(t.Pipe.Readme.Description) > 100 {
				return fmt.Errorf(
					"Readme short description can only be 100 characters long while you have: %d",
					len(t.Pipe.Readme.Description),
				)
			}

			return nil
		})
}

func LoginToDockerHubRegistry(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		Set(func(t *Task[Pipe]) error {
			login, err := json.Marshal(types.AuthConfig{
				Username: t.Pipe.DockerHub.Username,
				Password: t.Pipe.DockerHub.Password,
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

			t.Pipe.Ctx.Token = b.Token

			return nil
		})
}

func ReadReadmeFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("read").
		Set(func(t *Task[Pipe]) error {
			t.Log.Debugf("Trying to read file: %s", t.Pipe.Readme.File)

			content, err := os.ReadFile(t.Pipe.Readme.File)

			if err != nil {
				return err
			}

			t.Pipe.Ctx.Readme = string(content)

			t.Log.Debugf("File read: %s", t.Pipe.Readme.File)

			return nil
		})
}

func UpdateDockerReadme(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("update").
		Set(func(t *Task[Pipe]) error {
			t.Log.Debugf(
				"Running against repository: %s/%s",
				t.Pipe.DockerHub.Address,
				t.Pipe.Readme.Repository,
			)

			update := DockerHubUpdateReadmeRequest{
				Readme: t.Pipe.Ctx.Readme,
			}

			if t.Pipe.Readme.Description != "" {
				update.Description = t.Pipe.Readme.Description
			}

			body, err := json.Marshal(
				update,
			)

			if err != nil {
				return err
			}

			req, err := http.NewRequest(http.MethodPatch,
				fmt.Sprintf("%s/%s/", t.Pipe.DockerHub.Address, t.Pipe.Readme.Repository),
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
				if b.FullDescription != t.Pipe.Ctx.Readme {
					return fmt.Errorf("Uploaded README does not match with current repository README file.")
				}

				if t.Pipe.Readme.Description != "" && b.Description != t.Pipe.Readme.Description {
					return fmt.Errorf("Uploaded README does not match with current repository README file.")
				}

				t.Log.Infof(
					"Successfully pushed readme file to: %s > %s/%s",
					t.Pipe.Readme.File,
					t.Pipe.DockerHub.Address,
					t.Pipe.Readme.Repository,
				)
			case http.StatusNotFound:
				return fmt.Errorf(
					"Repository does not exists: %s/%s",
					t.Pipe.DockerHub.Address,
					t.Pipe.Readme.Repository,
				)
			default:
				if !b.CanEdit {
					return fmt.Errorf(
						"Given user credentials do not have permission to edit repository: %s/%s",
						t.Pipe.DockerHub.Address,
						t.Pipe.Readme.Repository,
					)
				}

				return fmt.Errorf(
					"Pushing readme failed with code: %d",
					res.StatusCode,
				)
			}

			return nil
		})
}
