package tagsfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// Parse reads the comma separated tags out of the file at path. A path that does
// not exist yields no tags and no error, since the file is usually written by an
// earlier job that may legitimately not have run. A pipe that can not go on
// without the tags asks for strict and gets an error instead.
func Parse(log *logrus.Entry, path string, strict bool) ([]string, error) {
	if _, err := os.Stat(path); err != nil {
		if strict && path != "" {
			return nil, fmt.Errorf("Tags file is set but does not exists: %s", path)
		}

		return nil, nil
	}

	log.Infof(
		"Tags file exists: %s",
		path,
	)

	content, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Can not read the tags file: %s -> %+v", path, err.Error())
	}

	tags := strings.Split(string(content), ",")

	if len(tags) == 0 {
		return nil, fmt.Errorf("Tags file does not contain any tags: %s", path)
	}

	newline := regexp.MustCompile(`\r?\n`)
	parsed := []string{}

	for _, tag := range tags {
		parsed = append(parsed, newline.ReplaceAllString(tag, ""))
	}

	return parsed, nil
}
