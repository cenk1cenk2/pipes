package pipe

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DownloadArtifact(t *Task[Pipe], url string) (string, error) {
	path := fmt.Sprintf("/tmp/%s/", uuid.New().String())
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("Can not create temporary directory: %s with error %w", path, err)
	}

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(path, url)

	if t.Pipe.Gitlab.Token != "" {
		req.HTTPRequest.Header.Set("PRIVATE-TOKEN", t.Pipe.Gitlab.Token)
	} else {
		req.HTTPRequest.Header.Set("JOB-TOKEN", t.Pipe.Gitlab.JobToken)
	}

	// start download
	res := client.Do(req)

	if res.Filename == "" {
		res.Filename = fmt.Sprintf("%s.zip", uuid.New().String())
	}

	t.Log.Infof("Downloading file: %s > %s", url, res.Filename)

	code, err := strconv.Atoi(strings.Split(res.HTTPResponse.Status, " ")[0])

	if err != nil {
		return "", err
	}

	err = ParseGLApiResponseCode(t, url, code)

	if err != nil {
		return "", err
	}

	// start UI loop
	timer := time.NewTicker(500 * time.Millisecond)
	defer timer.Stop()

Loop:
	for {
		select {
		case <-timer.C:
			t.Log.Infof("%s > %s: transferred %s / %s (%.2f%%) [%s]",
				url,
				res.Filename,
				humanize.Bytes(uint64(res.BytesComplete())),
				humanize.Bytes(uint64(res.Size())),
				100*res.Progress(),
				humanize.Bytes(uint64(res.BytesPerSecond())),
			)

		case <-res.Done:
			t.Log.Infof(
				"Download completed: %s to %s in %s > %s with %s/s",
				url,
				res.Filename,
				res.Duration(),
				humanize.Bytes(uint64(res.Size())),
				humanize.Bytes(uint64(res.BytesPerSecond())),
			)

			break Loop
		}
	}

	return res.Filename, nil
}

func ParseGLApiResponseCode(t *Task[Pipe], url string, code int) error {
	switch code {
	case http.StatusUnauthorized:
		return fmt.Errorf(
			"Given token does not have access to the given project through API.",
		)
	case http.StatusNotFound:
		return fmt.Errorf("Given API path is not found.")
	default:
		t.Log.Debugf("Status code: %d from %s", code, url)
	}

	return nil
}
