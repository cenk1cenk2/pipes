package publish

import (
	"fmt"
	"regexp"

	. "github.com/cenk1cenk2/plumber/v6"
)

func AddHelmChartVersion(t *Task, tag string) error {
	p, err := ProcessHelmChartVersion(t, tag)

	if err != nil {
		return err
	}

	return AppendHelmChartVersion(t, p)
}

func AppendHelmChartVersion(t *Task, tag string) error {
	t.Lock.Lock()
	C.Versions = append(C.Versions, tag)
	t.Lock.Unlock()

	return nil
}

func ProcessHelmChartVersion(t *Task, tag string) (string, error) {
	var err error
	if tag, err = ApplyHelmChartVersionTemplate(t, tag); err != nil {
		return "", err
	}

	if tag, err = SanitizeHelmChartVersion(t, tag); err != nil {
		return "", err
	}

	if tag == "" {
		return tag, fmt.Errorf("Can not add empty tag to list.")
	}

	return tag, nil
}

func SanitizeHelmChartVersion(t *Task, tag string) (string, error) {
	for _, s := range P.HelmChart.VersionsSanitize {
		re, err := regexp.Compile(s.Match)

		if err != nil {
			return "", fmt.Errorf("Can not compile sanitize regular expression: %w", err)
		}

		matches := re.FindStringSubmatch(tag)

		t.Log.Debugf("Trying to sanitize tag: %s with %v", tag, re.String())

		if matches == nil {
			continue
		}

		t.Log.Debugf("Sanitizing since condition matched for given tag: %s -> %s with %v", tag, re.String(), matches)

		return InlineTemplate(s.Template, matches)
	}

	return tag, nil
}

func ApplyHelmChartVersionTemplate(t *Task, tag string) (string, error) {
	for _, s := range P.HelmChart.VersionsTemplate {
		re, err := regexp.Compile(s.Match)

		if err != nil {
			return "", fmt.Errorf("Can not compile tag template regular expression: %w", err)
		}

		matches := re.FindStringSubmatch(tag)

		t.Log.Tracef("Trying to apply template to tag: %s with %v", tag, re.String())

		if matches == nil {
			continue
		}

		t.Log.Debugf("Applying template since condition matched for given tag: %s -> %s with %v", tag, re.String(), matches)

		return InlineTemplate(s.Template, matches)
	}

	return InlineTemplate[any](tag, nil)
}
