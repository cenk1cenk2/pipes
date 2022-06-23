package pipe

import (
	"net/http"
)

type Ctx struct {
	StepsResponse       GLApiSuccessfulStepsResponse
	Steps               []Step
	Client              *http.Client
	DownloadedArtifacts []DownloadedArtifact
	JobNames            []string
}
