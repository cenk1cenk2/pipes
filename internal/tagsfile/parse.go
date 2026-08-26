package tagsfile

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var newline = regexp.MustCompile(`\r?\n`)

// Parse reads the comma separated tags out of the file at path. A path that does
// not exist yields no tags and no error, since the file is usually written by an
// earlier job that may legitimately not have run.
func Parse(log *logrus.Entry, path string, strict bool) ([]string, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, nil
	}

	log.Infof(
		"Tags file exists: %s",
		path,
	)

	content, err := os.ReadFile(path)

	if strict && errors.Is(err, os.ErrNotExist) && path != "" {
		return nil, fmt.Errorf("Tags file is set but does not exists: %s", path)
	} else if err != nil {
		return nil, fmt.Errorf("Can not read the tags file: %s -> %+v", path, err.Error())
	}

	tags := strings.Split(string(content), ",")

	if len(tags) == 0 {
		return nil, fmt.Errorf("Tags file does not contain any tags: %s", path)
	}

	parsed := []string{}

	for _, tag := range tags {
		parsed = append(parsed, newline.ReplaceAllString(tag, ""))
	}

	return parsed, nil
}
