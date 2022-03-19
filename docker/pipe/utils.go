package pipe

import (
	"errors"
	"fmt"
)

func AddDockerTag(tag string) (string, error) {
	if tag == "" {
		return "", errors.New("Can not add empty tag to list.")
	}

	var t string

	if Pipe.DockerRegistry.Registry != "" {
		t = fmt.Sprintf("%s/%s:%s", Pipe.DockerRegistry.Registry, Pipe.DockerImage.Name, tag)
	} else {
		t = fmt.Sprintf("%s:%s", Pipe.DockerImage.Name, tag)
	}

	Context.Tags = append(Context.Tags, t)

	return t, nil
}
