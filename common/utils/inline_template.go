package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

func InlineTemplate[Ctx any](tmpl string, ctx Ctx) (string, error) {
	if tmpl == "" {
		return "", nil
	}

	tmp, err := template.New("inline").Funcs(template.FuncMap{"join": strings.Join, "to_upper_case": strings.ToUpper, "to_lower_case": strings.ToLower}).Parse(tmpl)

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
