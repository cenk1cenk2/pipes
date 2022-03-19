package pipe

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
)

func DownloadArtifact(url string) (string, error) {
	path := fmt.Sprintf("/tmp/%s/", uuid.New().String())
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("Can not create temporary directory: %s with error %s", path, err)
	}

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(path, url)

	if Pipe.Gitlab.JobToken != "" {
		req.HTTPRequest.Header.Set("JOB-TOKEN", Pipe.Gitlab.JobToken)
	} else {
		req.HTTPRequest.Header.Set("PRIVATE-TOKEN", Pipe.Gitlab.Token)
	}

	// start download
	res := client.Do(req)

	utils.Log.Infoln(fmt.Sprintf("Downloading file: %s -> %s", url, res.Filename))

	code, err := strconv.Atoi(strings.Split(res.HTTPResponse.Status, " ")[0])

	if err != nil {
		return "", err
	}

	err = ParseGLApiResponseCode(url, code)

	if err != nil {
		return "", err
	}

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			utils.Log.Infoln(fmt.Sprintf("%s -> %s: transferred %s / %s (%.2f%%) [%s]",
				url,
				res.Filename,
				humanize.Bytes(uint64(res.BytesComplete())),
				humanize.Bytes(uint64(res.Size())),
				100*res.Progress(),
				humanize.Bytes(uint64(res.BytesPerSecond())),
			),
			)

		case <-res.Done:
			utils.Log.Infoln(
				fmt.Sprintf(
					"Download completed: %s to %s in %s",
					url,
					res.Filename,
					res.Duration(),
				),
			)

			break Loop
		}
	}

	return res.Filename, nil
}

func ParseGLApiResponseCode(url string, code int) error {
	switch code {
	case http.StatusUnauthorized:
		return fmt.Errorf(
			"Given token does not have access to the given project through API.",
		)
	case http.StatusNotFound:
		return fmt.Errorf("Given API path is not found.")
	default:
		utils.Log.Debugln(fmt.Sprintf("Status code: %d from %s", code, url))
	}

	return nil
}
