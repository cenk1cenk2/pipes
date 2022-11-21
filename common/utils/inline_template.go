package utils

import (
	"bytes"
	"fmt"
	"html/template"

	"gitlab.kilic.dev/libraries/plumber/v4"
)

func InlineTemplate[Ctx any](tmpl string, ctx Ctx) (string, error) {
	if tmpl == "" {
		return "", nil
	}

	// functions can be found here: https://go-task.github.io/slim-sprig/
	tmp, err := template.New("inline").Funcs(plumber.TemplateFuncMap()).Parse(tmpl)

	if err != nil {
		return "", fmt.Errorf("Can not create inline template: %w", err)
	}

	var w bytes.Buffer

	err = tmp.ExecuteTemplate(&w, "inline", ctx)

	if err != nil {
		return "", fmt.Errorf("Can not generate inline template: %w", err)
	}

	return w.String(), nil
}
