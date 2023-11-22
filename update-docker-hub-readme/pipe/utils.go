package pipe

import (
	"fmt"
	"net/http"

	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func AddAuthenticationHeadersToRequest(t *Task[Pipe], req *http.Request) *http.Request {
	req.Header.Add("User-Agent", t.Plumber.Cli.Name)
	req.Header.Add("Content-Type", JSON_REQUEST)
	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", t.Pipe.Ctx.Token))

	return req
}
