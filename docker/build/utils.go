package build

import (
	"fmt"
	"regexp"

	"gitlab.kilic.dev/devops/pipes/common/utils"
	"gitlab.kilic.dev/devops/pipes/docker/login"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func AddDockerTag(t *Task[Pipe], tag string) error {
	if tag == "" {
		return fmt.Errorf("Can not add empty tag to list.")
	}

	var err error
	if tag, err = SanitizeDockerTag(t, tag); err != nil {
		return err
	}

	if tag, err = ApplyTagTemplate(t, tag); err != nil {
		return err
	}

	if login.TL.Pipe.DockerRegistry.Registry != "" {
		tag = fmt.Sprintf("%s/%s:%s", login.TL.Pipe.DockerRegistry.Registry, t.Pipe.DockerImage.Name, tag)
	} else {
		tag = fmt.Sprintf("%s:%s", t.Pipe.DockerImage.Name, tag)
	}

	t.Lock.Lock()
	t.Pipe.Ctx.Tags = append(t.Pipe.Ctx.Tags, tag)
	t.Lock.Unlock()

	return nil
}

func SanitizeDockerTag(t *Task[Pipe], tag string) (string, error) {
	for _, s := range t.Pipe.DockerImage.TagsSanitize {
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

		return utils.InlineTemplate(s.Template, matches)
	}

	return tag, nil
}

func ApplyTagTemplate(t *Task[Pipe], tag string) (string, error) {
	for _, s := range t.Pipe.DockerImage.TagsTemplate {
		re, err := regexp.Compile(s.Match)

		if err != nil {
			return "", fmt.Errorf("Can not compile tag template regular expression: %w", err)
		}

		matches := re.FindStringSubmatch(tag)

		t.Log.Debugf("Trying to apply template to tag: %s with %v", tag, re.String())

		if matches == nil {
			continue
		}

		t.Log.Debugf("Applying template since condition matched for given tag: %s -> %s with %v", tag, re.String(), matches)

		return utils.InlineTemplate(s.Template, matches)
	}

	return tag, nil
}
