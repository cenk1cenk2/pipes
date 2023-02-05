package setup

import (
	"github.com/docker/docker/client"
)

type Ctx struct {
	Client *client.Client
}
