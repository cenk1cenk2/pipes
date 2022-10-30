package pipe

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

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

	if t.Pipe.DockerRegistry.Registry != "" {
		tag = fmt.Sprintf("%s/%s:%s", t.Pipe.DockerRegistry.Registry, t.Pipe.DockerImage.Name, tag)
	} else {
		tag = fmt.Sprintf("%s:%s", t.Pipe.DockerImage.Name, tag)
	}

	t.Lock.Lock()
	t.Pipe.Ctx.Tags = append(t.Pipe.Ctx.Tags, tag)
	t.Lock.Unlock()

	return nil
}

func SanitizeDockerTag(t *Task[Pipe], tag string) (string, error) {
	for expression, tmpl := range t.Pipe.Ctx.SanitizedRegularExpression {
		re, err := regexp.Compile(expression)

		if err != nil {
			return "", fmt.Errorf("Can not compile sanitize regular expression: %w", err)
		}

		matches := re.FindStringSubmatch(tag)

		t.Log.Debugf("Trying to sanitize tag: %s with %v", tag, re.String())

		if matches == nil {
			continue
		}

		t.Log.Debugf("Sanitizing since condition matched for given tag: %s -> %s with %v", tag, re.String(), matches)

		tmp, err := template.New("inline").Funcs(template.FuncMap{"join": strings.Join, "to_upper_case": strings.ToUpper, "to_lower_case": strings.ToLower}).Parse(tmpl)

		if err != nil {
			return "", err
		}

		var w bytes.Buffer

		err = tmp.ExecuteTemplate(&w, "inline", matches)

		if err != nil {
			return "", err
		}

		return w.String(), nil
	}

	return tag, nil
}
