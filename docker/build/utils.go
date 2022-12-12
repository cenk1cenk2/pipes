package build

import (
	"fmt"
	"regexp"

	"gitlab.kilic.dev/devops/pipes/common/utils"
	"gitlab.kilic.dev/devops/pipes/docker/login"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func AddDockerTag(t *Task[Pipe], tag string) error {
	p, err := ProcessDockerTag(t, tag)

	if err != nil {
		return err
	}

	t.Lock.Lock()
	t.Pipe.Ctx.Tags = append(t.Pipe.Ctx.Tags, p)
	t.Lock.Unlock()

	return nil
}

func ProcessDockerTag(t *Task[Pipe], tag string) (string, error) {
	if tag == "" {
		return "", fmt.Errorf("Can not add empty tag to list.")
	}

	var err error
	if tag, err = SanitizeDockerTag(t, tag); err != nil {
		return "", err
	}

	if tag, err = ApplyTagTemplate(t, tag); err != nil {
		return "", err
	}

	if login.TL.Pipe.DockerRegistry.Registry != "" {
		tag = fmt.Sprintf("%s/%s:%s", login.TL.Pipe.DockerRegistry.Registry, t.Pipe.DockerImage.Name, tag)
	} else {
		tag = fmt.Sprintf("%s:%s", t.Pipe.DockerImage.Name, tag)
	}

	return tag, nil
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

		t.Log.Tracef("Trying to apply template to tag: %s with %v", tag, re.String())

		if matches == nil {
			continue
		}

		t.Log.Debugf("Applying template since condition matched for given tag: %s -> %s with %v", tag, re.String(), matches)

		return utils.InlineTemplate(s.Template, matches)
	}

	return tag, nil
}

func ApplyBuildArgsTemplate(t *Task[Pipe]) ([]string, error) {
	args := []string{}

	if len(t.Pipe.DockerImage.BuildArgs) == 0 {
		return args, nil
	}

	vars := ParseEnvironmentVariablesToMap()

	for _, v := range t.Pipe.DockerImage.BuildArgs {
		result, err := utils.InlineTemplate(v, vars)

		if err != nil {
			return nil, err
		}

		args = append(args, result)
	}

	return args, nil
}
