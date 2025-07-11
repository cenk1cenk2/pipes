package pipe

import (
	"fmt"
	"net/http"

	. "github.com/cenk1cenk2/plumber/v6"
)

func AddAuthenticationHeadersToRequest(t *Task, req *http.Request) *http.Request {
	req.Header.Add("User-Agent", t.Plumber.Cli.Name)
	req.Header.Add("Content-Type", JSON_REQUEST)
	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", C.Token))

	return req
}
