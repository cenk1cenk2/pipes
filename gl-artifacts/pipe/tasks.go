package pipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type Ctx struct {
	StepsResponse       GLApiSuccessfulStepsResponse
	StepIds             []StepId
	Client              *http.Client
	DownloadedArtifacts []DownloadedArtifact
}

var Context Ctx

func VerifyVariables() utils.Task {
	metadata := utils.TaskMetadata{Context: "verify"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		reqUrl := fmt.Sprintf(
			"https://gitlab.kilic.dev/api/v4/projects/%s/pipelines/%s/jobs/?scope=success",
			Pipe.Gitlab.ParentProjectId,
			Pipe.Gitlab.ParentPipelineId,
		)

		utils.Log.Debugln(fmt.Sprintf("Pipeline steps API request URL: %s", reqUrl))

		req, err := http.NewRequest(
			http.MethodGet,
			reqUrl,
			nil,
		)

		if err != nil {
			utils.Log.Fatalln(err)
		}

		req.Header.Set("PRIVATE-TOKEN", Pipe.Gitlab.Token)

		Context.Client = &http.Client{}

		res, err := Context.Client.Do(req)

		if err != nil {
			return err
		}

		defer res.Body.Close()

		err = ParseGLApiResponseCode(reqUrl, res.StatusCode)

		if err != nil {
			return err
		}

		decoder := json.NewDecoder(res.Body)

		err = decoder.Decode(&Context.StepsResponse)

		if err != nil {
			return err
		}

		return nil
	}}
}

func DiscoverArtifacts() utils.Task {
	metadata := utils.TaskMetadata{Context: "discover"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		Context.StepIds = []StepId{}

		var wg sync.WaitGroup
		wg.Add(len(Pipe.Gitlab.DownloadArtifacts.Value()))

		for _, step := range Pipe.Gitlab.DownloadArtifacts.Value() {
			go func(step string) {
				defer wg.Done()

				found := false

				for _, v := range Context.StepsResponse {
					if v.Name == step {
						utils.Log.Debugln(
							fmt.Sprintf("Adding step artifacts: %s with id %d", step, v.ID),
						)

						Context.StepIds = append(Context.StepIds, StepId{id: v.ID, name: step})

						found = true
					}
				}

				if !found {
					utils.Log.Errorln(
						fmt.Sprintf(
							"Job with name is not found so artifacts are not downloaded: %s ",
							step,
						),
					)

					var available []string = []string{}

					for _, v := range Context.StepsResponse {
						available = append(available, v.Name)

					}

					utils.Log.Fatalln(
						fmt.Sprintf(
							"Available steps are: %s",
							strings.Join(available, ", "),
						),
					)
				}

			}(step)
		}

		wg.Wait()

		return nil
	}}
}

func DownloadArtifacts() utils.Task {
	metadata := utils.TaskMetadata{Context: "download"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		url := fmt.Sprintf(
			"https://gitlab.kilic.dev/api/v4/projects/%s/jobs/%s/artifacts/",
			Pipe.Gitlab.ParentProjectId,
			"%d",
		)

		var wg sync.WaitGroup
		wg.Add(len(Context.StepIds))

		for _, stepId := range Context.StepIds {
			go func(stepId StepId) {
				defer wg.Done()

				utils.Log.Debugln(
					fmt.Sprintf(
						"Will download artifact with parent job: %s -> %d",
						stepId.name,
						stepId.id,
					),
				)

				path, err := DownloadArtifact(fmt.Sprintf(url, stepId.id))

				if err != nil {
					utils.Log.Fatalln(
						fmt.Sprintf(
							"Can not download artifacts from stage: %s -> %d with error: %s",
							stepId.name,
							stepId.id,
							err,
						),
					)
				}

				Context.DownloadedArtifacts = append(
					Context.DownloadedArtifacts,
					DownloadedArtifact{name: stepId.name, path: path},
				)
			}(stepId)
		}

		wg.Wait()

		return nil
	}}
}

func UnarchiveArtifacts() utils.Task {
	metadata := utils.TaskMetadata{
		Context:        "unarchive",
		StdOutLogLevel: logrus.DebugLevel,
	}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		t.Commands = []utils.Command{}

		var wg sync.WaitGroup
		wg.Add(len(Context.DownloadedArtifacts))

		for _, artifact := range Context.DownloadedArtifacts {
			go func(artifact DownloadedArtifact) {
				defer wg.Done()

				utils.Log.Debugln(
					fmt.Sprintf(
						"Will Decompres artifact: %s -> %s",
						artifact.name,
						artifact.path,
					),
				)

				cmd := exec.Command("unzip", "-o")

				cmd.Args = append(cmd.Args, artifact.path, "-d", "./")

				t.Commands = append(t.Commands, cmd)
			}(artifact)
		}

		wg.Wait()

		return nil
	}}
}
