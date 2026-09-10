package terraform

import (
	json "encoding/json/v2"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const (
	ChangeCreate = "+"
	ChangeDelete = "-"
	ChangeUpdate = "~"

	// What stands in a resource's block where the plan proposes no attribute of its
	// own, so that the block a reader opens is never empty.
	NoAttributeChanges = "No attribute changes."
)

type (
	// Change is one attribute of a resource diff, as far as the tool's plan carries
	// it: a leaf holds the rendered values, a container holds its children. Before
	// stays empty when the plan has no old value to show.
	Change struct {
		Name     string
		Action   string
		Before   string
		After    string
		Note     string
		Children []Change
	}

	// Marker stands in for a value the report must not show, so it renders bare
	// where a string value would be quoted.
	Marker string
)

// Expand turns a value into the subtree of changes that writes it out in full, with
// every line carrying the action; it is how a created or deleted resource, and every
// value the plan shows only one side of, is written.
func Expand(action string, name string, value any) Change {
	change := Change{Name: name, Action: action}

	switch value := value.(type) {
	case map[string]any:
		if len(value) == 0 {
			break
		}

		for _, key := range slices.Sorted(maps.Keys(value)) {
			change.Children = append(change.Children, Expand(action, key, value[key]))
		}

		return change
	case []any:
		if !Expandable(value) {
			break
		}

		for index, element := range value {
			change.Children = append(change.Children, Expand(action, fmt.Sprintf("[%d]", index), element))
		}

		return change
	}

	change.After = FormatValue(value)
	if action == ChangeDelete {
		change.Before, change.After = change.After, ""
	}

	return change
}

// Expandable says whether a value is written out one child per line, which a
// non-empty object always is and a list only is when its elements are not scalars:
// a list of scalars reads better on one line than as a run of indexes.
func Expandable(value any) bool {
	switch value := value.(type) {
	case map[string]any:
		return len(value) > 0
	case []any:
		return slices.ContainsFunc(value, func(element any) bool {
			switch element.(type) {
			case map[string]any, []any:
				return true
			}

			return false
		})
	}

	return false
}

// FormatValue writes a value the way its plan JSON carries it, with markers left bare
// and objects sorted by key so the same value always reads the same.
func FormatValue(value any) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case Marker:
		return string(value)
	case map[string]any:
		fields := make([]string, 0, len(value))
		for _, key := range slices.Sorted(maps.Keys(value)) {
			fields = append(fields, fmt.Sprintf("%q: %s", key, FormatValue(value[key])))
		}

		return "{" + strings.Join(fields, ", ") + "}"
	case []any:
		elements := make([]string, 0, len(value))
		for _, element := range value {
			elements = append(elements, FormatValue(element))
		}

		return "[" + strings.Join(elements, ", ") + "]"
	}

	// the remaining values are the JSON scalars, and the encoder cannot fail on those.
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(encoded)
}

// Glyph is the action sign a resource row leads with, in the vocabulary of the plan
// output of the tools themselves.
func Glyph(action string) string {
	switch action {
	case "create", "create-replacement", "import":
		return "+"
	case "update", "update-replacement":
		return "~"
	case "delete", "delete-replaced", "discard", "remove-pending-replace":
		return "-"
	case "replace":
		return "-/+"
	case "read", "read-replacement":
		return "<="
	case "move":
		return ">"
	case "forget":
		return "."
	}

	return "?"
}

func (c Change) lines(depth int, out []string) []string {
	prefix := "  "
	if c.Action != "" {
		prefix = c.Action + " "
	}

	line := prefix + strings.Repeat("  ", depth) + c.Name
	switch value := c.value(); {
	case len(c.Children) > 0:
		line += ":"
	case value != "":
		line += ": " + value
	}
	if c.Note != "" {
		line += " # " + c.Note
	}

	out = append(out, line)
	for _, child := range c.Children {
		out = child.lines(depth+1, out)
	}

	return out
}

func (c Change) value() string {
	before, after := c.Before, c.After

	switch {
	case c.Action == ChangeDelete:
		return before
	case before != "" && after != "":
		return before + " -> " + after
	case before != "":
		return before
	}

	return after
}

// RenderChanges writes a resource diff out, one attribute per line.
func RenderChanges(changes []Change) string {
	lines := []string{}
	for _, change := range changes {
		lines = change.lines(0, lines)
	}

	return strings.Join(lines, "\n")
}
