package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func ParseTagsFile(log *logrus.Entry, file string, strict bool) ([]string, error) {
	if _, err := os.Stat(file); err == nil {
		log.Infof(
			"Tags file exists: %s",
			file,
		)

		content, err := os.ReadFile(file)

		if strict && errors.Is(err, os.ErrNotExist) && file != "" {
			return nil, fmt.Errorf("Tags file is set but does not exists: %s", file)
		} else if err != nil {
			return nil, fmt.Errorf("Can not read the tags file: %s -> %+v", file, err.Error())
		}

		tags := strings.Split(string(content), ",")

		if len(tags) == 0 {
			return nil, fmt.Errorf("Tags file does not contain any tags: %s", file)
		}

		t := []string{}

		re := regexp.MustCompile(`\r?\n`)

		for _, v := range tags {
			t = append(t, re.ReplaceAllString(v, ""))
		}

		return t, nil
	}
	return nil, nil
}
