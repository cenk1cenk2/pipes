package utils

import (
	"bytes"
	"html/template"
	"strings"
)

func InlineTemplate[Ctx any](tmpl string, ctx Ctx) (string, error) {
	tmp, err := template.New("inline").Funcs(template.FuncMap{"join": strings.Join, "to_upper_case": strings.ToUpper, "to_lower_case": strings.ToLower}).Parse(tmpl)

	if err != nil {
		return "", err
	}

	var w bytes.Buffer

	err = tmp.ExecuteTemplate(&w, "inline", ctx)

	if err != nil {
		return "", err
	}

	return w.String(), nil
}
