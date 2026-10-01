package terraform

import (
	json "encoding/json/v2"
	"fmt"
	"maps"
	"reflect"
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

	// What a string value that holds a JSON document is written out structurally with.
	JSONEncoded = "jsonencode"

	heredocOpen  = "<<-EOT"
	heredocClose = "EOT"

	// past this many pairs of lines the table the line diff of two heredocs is worked
	// out on grows too large, and the two are written out whole instead.
	lineDiffLimit = 1 << 20
)

// DiffLayout is how a resource diff is written out.
type DiffLayout int

const (
	// DiffPlan writes the diff the way the plan output of the tools does, an update
	// on one line marked with a tilde; the job log colors it with a lexer of its own.
	DiffPlan DiffLayout = iota
	// DiffUnified writes an update as its old line removed and its new line added,
	// the only shape a code host highlights in a diff block.
	DiffUnified
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

	if decoded, ok := DecodeJSON(value); ok {
		change = Expand(action, name, decoded)
		change.Note = JSONEncoded

		return change
	}

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

// Compare writes the change between two plain values, attribute by attribute where
// both are objects or lists, and line by line where both are values that span lines;
// it says nothing changed when the two are equal. A marker on the new side, a value
// that is not known yet or not shown, still shows the value it replaces in full.
func Compare(name string, before any, after any) (Change, bool) {
	change := Change{Name: name, Action: ChangeUpdate}

	if decoded, ok := DecodeJSON(before); ok {
		if other, ok := DecodeJSON(after); ok {
			before, after = decoded, other
			change.Note = JSONEncoded
		} else if _, ok := after.(Marker); ok {
			before = decoded
			change.Note = JSONEncoded
		}
	}

	if _, ok := after.(Marker); ok && Expandable(before) {
		switch before := before.(type) {
		case map[string]any:
			for _, key := range slices.Sorted(maps.Keys(before)) {
				child, _ := Compare(key, before[key], after)
				change.Children = append(change.Children, child)
			}
		case []any:
			for index, element := range before {
				child, _ := Compare(fmt.Sprintf("[%d]", index), element, after)
				change.Children = append(change.Children, child)
			}
		}

		return change, true
	}

	switch beforeValue := before.(type) {
	case map[string]any:
		afterValue, ok := after.(map[string]any)
		if !ok {
			break
		}

		keys := slices.Concat(slices.Collect(maps.Keys(beforeValue)), slices.Collect(maps.Keys(afterValue)))
		slices.Sort(keys)

		for _, key := range slices.Compact(keys) {
			if child, changed := compareChild(key, beforeValue, afterValue); changed {
				change.Children = append(change.Children, child)
			}
		}

		return change, len(change.Children) > 0
	case []any:
		afterValue, ok := after.([]any)
		if !ok || (!Expandable(beforeValue) && !Expandable(afterValue)) {
			break
		}

		for index := range max(len(beforeValue), len(afterValue)) {
			key := fmt.Sprintf("[%d]", index)
			beforeElements, afterElements := map[string]any{}, map[string]any{}
			if index < len(beforeValue) {
				beforeElements[key] = beforeValue[index]
			}
			if index < len(afterValue) {
				afterElements[key] = afterValue[index]
			}

			if child, changed := compareChild(key, beforeElements, afterElements); changed {
				change.Children = append(change.Children, child)
			}
		}

		return change, len(change.Children) > 0
	}

	if reflect.DeepEqual(before, after) {
		return Change{}, false
	}

	change.Before = FormatValue(before)
	change.After = FormatValue(after)

	return change, true
}

func compareChild(key string, before map[string]any, after map[string]any) (Change, bool) {
	beforeValue, inBefore := before[key]
	afterValue, inAfter := after[key]

	switch {
	case !inBefore:
		return Expand(ChangeCreate, key, afterValue), true
	case !inAfter:
		return Expand(ChangeDelete, key, beforeValue), true
	}

	return Compare(key, beforeValue, afterValue)
}

// Expandable says whether a value is written out one child per line, which a
// non-empty object always is and a list only is when an element is not a scalar or
// spans lines: a list of short scalars reads better on one line than as a run of
// indexes.
func Expandable(value any) bool {
	switch value := value.(type) {
	case map[string]any:
		return len(value) > 0
	case []any:
		return slices.ContainsFunc(value, func(element any) bool {
			switch element := element.(type) {
			case map[string]any, []any:
				return true
			case string:
				return strings.Contains(element, "\n")
			}

			return false
		})
	}

	return false
}

// DecodeJSON reads a string holding a JSON document into the value it encodes, as
// long as that value is written out one child per line; anything else stays the
// string it is, since decoding it would only change how it is quoted.
func DecodeJSON(value any) (any, bool) {
	text, ok := value.(string)
	if !ok {
		return nil, false
	}

	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return nil, false
	}

	var decoded any
	if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil || !Expandable(decoded) {
		return nil, false
	}

	return decoded, true
}

// FormatValue writes a value the way its plan JSON carries it, with markers left bare
// and objects sorted by key so the same value always reads the same. A string that
// spans lines is written as a heredoc instead, the way the tools show it themselves,
// since its escaped form runs a whole document onto one line.
func FormatValue(value any) string {
	if text, ok := value.(string); ok && strings.Contains(text, "\n") {
		return heredoc(text)
	}

	return formatInline(value)
}

func formatInline(value any) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case Marker:
		return string(value)
	case map[string]any:
		fields := make([]string, 0, len(value))
		for _, key := range slices.Sorted(maps.Keys(value)) {
			fields = append(fields, fmt.Sprintf("%q: %s", key, formatInline(value[key])))
		}

		return "{" + strings.Join(fields, ", ") + "}"
	case []any:
		elements := make([]string, 0, len(value))
		for _, element := range value {
			elements = append(elements, formatInline(element))
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

func (c Change) lines(depth int, layout DiffLayout, out []string) []string {
	indent := strings.Repeat("  ", depth)
	note := ""
	if c.Note != "" {
		note = " # " + c.Note
	}

	multiline := c.Action != ChangeDelete && c.Before != "" && c.After != "" &&
		(strings.Contains(c.Before, "\n") || strings.Contains(c.After, "\n"))

	action := c.Action
	if layout == DiffUnified && action == ChangeUpdate {
		switch {
		// what only changed inside reads as context around the lines that did.
		case len(c.Children) > 0, multiline:
			action = ""
		case c.Before != "" && c.After != "":
			out = append(out, sign(ChangeDelete)+indent+c.Name+": "+c.Before)

			return append(out, sign(ChangeCreate)+indent+c.Name+": "+c.After+note)
		case c.After == "" && c.Before != "":
			action = ChangeDelete
		default:
			action = ChangeCreate
		}
	}

	head := sign(action) + indent + c.Name

	if len(c.Children) > 0 {
		out = append(out, head+":"+note)
		for _, child := range c.Children {
			out = child.lines(depth+1, layout, out)
		}

		return out
	}

	if multiline {
		return c.multilineUpdate(head, note, indent, out)
	}

	value := c.value()
	if value == "" {
		return append(out, head+note)
	}

	// a heredoc carries its lines on, and every one of them keeps the action of the
	// change so that it reads as part of it.
	lines := strings.Split(value, "\n")
	out = append(out, head+": "+lines[0]+note)
	for _, line := range lines[1:] {
		out = append(out, sign(action)+indent+line)
	}

	return out
}

// multilineUpdate writes a change of a value that spans lines: between two heredocs it
// marks the lines that changed in place, so a reader finds the one line an edit to a
// long document touched; otherwise the old value is followed by the new one whole.
func (c Change) multilineUpdate(head string, note string, indent string, out []string) []string {
	before, after := strings.Split(c.Before, "\n"), strings.Split(c.After, "\n")

	if len(before) > 1 && len(after) > 1 {
		out = append(out, head+": "+heredocOpen+note)
		for _, line := range diffLines(before[1:len(before)-1], after[1:len(after)-1]) {
			out = append(out, sign(line.action)+indent+line.text)
		}

		return append(out, sign("")+indent+heredocClose)
	}

	out = append(out, head+":"+note)
	for _, line := range before {
		out = append(out, sign(ChangeDelete)+indent+"  "+line)
	}
	for _, line := range after {
		out = append(out, sign(ChangeCreate)+indent+"  "+line)
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

func sign(action string) string {
	if action == "" {
		return "  "
	}

	return action + " "
}

func heredoc(text string) string {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	for index, line := range lines {
		if line != "" {
			lines[index] = "  " + line
		}
	}

	return heredocOpen + "\n" + strings.Join(lines, "\n") + "\n" + heredocClose
}

type diffLine struct {
	action string
	text   string
}

// diffLines lines two versions of a text up along the longest run of lines they
// share, leaving those unmarked and marking every other line as removed or added.
func diffLines(before []string, after []string) []diffLine {
	lines := []diffLine{}
	index, other := 0, 0

	if len(before)*len(after) <= lineDiffLimit {
		// common[i][j] is how many lines before[i:] and after[j:] share in order.
		common := make([][]int, len(before)+1)
		for i := range common {
			common[i] = make([]int, len(after)+1)
		}
		for i := len(before) - 1; i >= 0; i-- {
			for j := len(after) - 1; j >= 0; j-- {
				if before[i] == after[j] {
					common[i][j] = common[i+1][j+1] + 1
				} else {
					common[i][j] = max(common[i+1][j], common[i][j+1])
				}
			}
		}

		for index < len(before) && other < len(after) {
			switch {
			case before[index] == after[other]:
				lines = append(lines, diffLine{text: before[index]})
				index++
				other++
			case common[index+1][other] >= common[index][other+1]:
				lines = append(lines, diffLine{action: ChangeDelete, text: before[index]})
				index++
			default:
				lines = append(lines, diffLine{action: ChangeCreate, text: after[other]})
				other++
			}
		}
	}

	for ; index < len(before); index++ {
		lines = append(lines, diffLine{action: ChangeDelete, text: before[index]})
	}
	for ; other < len(after); other++ {
		lines = append(lines, diffLine{action: ChangeCreate, text: after[other]})
	}

	return lines
}

// RenderChanges writes a resource diff out, one attribute per line.
func RenderChanges(changes []Change, layout DiffLayout) string {
	lines := []string{}
	for _, change := range changes {
		lines = change.lines(0, layout, lines)
	}

	return strings.Join(lines, "\n")
}
