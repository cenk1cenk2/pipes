package pipe

import (
	"errors"
	"fmt"
)

func AddDockerTag(tag string) error {
	if tag == "" {
		return errors.New("Can not add empty tag to list.")
	}

	var t string

	if TL.Pipe.DockerRegistry.Registry != "" {
		t = fmt.Sprintf("%s/%s:%s", TL.Pipe.DockerRegistry.Registry, TL.Pipe.DockerImage.Name, tag)
	} else {
		t = fmt.Sprintf("%s:%s", TL.Pipe.DockerImage.Name, tag)
	}

	TL.Lock.Lock()
	TL.Pipe.Ctx.Tags = append(TL.Pipe.Ctx.Tags, t)
	TL.Lock.Unlock()

	return nil
}
