package pipe

import (
	"errors"
	"fmt"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func AddDockerTag(_ *Task[Pipe], tag string) error {
	if tag == "" {
		return errors.New("Can not add empty tag to list.")
	}

	if TL.Pipe.DockerRegistry.Registry != "" {
		tag = fmt.Sprintf("%s/%s:%s", TL.Pipe.DockerRegistry.Registry, TL.Pipe.DockerImage.Name, tag)
	} else {
		tag = fmt.Sprintf("%s:%s", TL.Pipe.DockerImage.Name, tag)
	}

	TL.Lock.Lock()
	TL.Pipe.Ctx.Tags = append(TL.Pipe.Ctx.Tags, tag)
	TL.Lock.Unlock()

	return nil
}
