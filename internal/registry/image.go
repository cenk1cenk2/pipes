package registry

import (
	"fmt"
	"strings"
)

// Image is a fully qualified image name, split at the registry host it publishes to.
type Image struct {
	Host string
	Path string
}

// ParseImage splits a fully qualified image name at its registry host. Like the
// container tooling, the first component only counts as a host when it carries a
// dot or a port, or is localhost.
func ParseImage(name string) (*Image, error) {
	host, path, found := strings.Cut(name, "/")

	if !found || (!strings.ContainsAny(host, ".:") && host != "localhost") {
		return nil, fmt.Errorf("Image name is not fully qualified, it has to start with the registry host: %s", name)
	}

	return &Image{Host: host, Path: path}, nil
}

func (i *Image) String() string {
	return i.Host + "/" + i.Path
}
