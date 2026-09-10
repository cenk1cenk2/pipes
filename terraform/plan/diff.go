package plan

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

const (
	sensitiveValue  = terraform.Marker("(sensitive value)")
	knownAfterApply = terraform.Marker("(known after apply)")
)

// One half of a change next to the masks Terraform ships with it. A mask mirrors the
// shape of its value with a true where a leaf is sensitive or unknown, or where a
// whole subtree is, and leaves out what is neither.
type side struct {
	value     any
	present   bool
	sensitive any
	unknown   any
}

func (s side) child(key string) side {
	value, present := field(s.value, key)
	sensitive, _ := field(s.sensitive, key)
	unknown, _ := field(s.unknown, key)

	// an attribute that is not known yet is left out of the value and only kept in
	// its mask.
	return side{value: value, present: present || flagged(unknown), sensitive: sensitive, unknown: unknown}
}

func (s side) element(index int) side {
	value, present := element(s.value, index)
	sensitive, _ := element(s.sensitive, index)
	unknown, _ := element(s.unknown, index)

	return side{value: value, present: present, sensitive: sensitive, unknown: unknown}
}

func (s side) isSensitive() bool {
	return flagged(s.sensitive)
}

func (s side) isUnknown() bool {
	return flagged(s.unknown)
}

// masked is the value with everything the report must not show replaced, subtree by
// subtree, by the marker that says why.
func (s side) masked() any {
	switch {
	case s.isUnknown():
		return knownAfterApply
	case s.isSensitive():
		return sensitiveValue
	}

	switch value := s.value.(type) {
	case map[string]any:
		masked := make(map[string]any, len(value))
		for key := range value {
			masked[key] = s.child(key).masked()
		}
		if unknown, ok := s.unknown.(map[string]any); ok {
			for key, flag := range unknown {
				if flagged(flag) {
					masked[key] = knownAfterApply
				}
			}
		}

		return masked
	case []any:
		masked := make([]any, len(value))
		for index := range value {
			masked[index] = s.element(index).masked()
		}

		return masked
	}

	return s.value
}

func field(value any, key string) (any, bool) {
	object, ok := value.(map[string]any)
	if !ok {
		return nil, false
	}

	child, ok := object[key]

	return child, ok
}

func element(value any, index int) (any, bool) {
	list, ok := value.([]any)
	if !ok || index >= len(list) {
		return nil, false
	}

	return list[index], true
}

func flagged(mask any) bool {
	flag, ok := mask.(bool)

	return ok && flag
}

type changeDiff struct {
	// the attribute paths Terraform names as the cause of a replacement.
	replace map[string]bool
}

func resourceChanges(change *tfjson.Change) []terraform.Change {
	before := side{
		value:     change.Before,
		present:   change.Before != nil,
		sensitive: change.BeforeSensitive,
	}
	after := side{
		value:     change.After,
		present:   change.After != nil || flagged(change.AfterUnknown),
		sensitive: change.AfterSensitive,
		unknown:   change.AfterUnknown,
	}

	diff := changeDiff{replace: map[string]bool{}}
	for _, path := range change.ReplacePaths {
		if steps, ok := path.([]any); ok {
			diff.replace[joinPath(steps)] = true
		}
	}

	root, changed := diff.compare("", nil, before, after)
	if !changed {
		return nil
	}

	return root.Children
}

func (d changeDiff) compare(name string, path []any, before side, after side) (terraform.Change, bool) {
	switch {
	case !before.present && !after.present:
		return terraform.Change{}, false
	case !before.present:
		return d.note(path, terraform.Expand(terraform.ChangeCreate, name, after.masked())), true
	case !after.present:
		return d.note(path, terraform.Expand(terraform.ChangeDelete, name, before.masked())), true
	}

	change := terraform.Change{Name: name, Action: terraform.ChangeUpdate}

	// a side that is masked as a whole reads as one value, however deep it goes.
	opaque := before.isSensitive() || after.isSensitive() || after.isUnknown()

	switch beforeValue := before.value.(type) {
	case map[string]any:
		afterValue, ok := after.value.(map[string]any)
		if !ok || opaque {
			break
		}

		keys := slices.Concat(slices.Collect(maps.Keys(beforeValue)), slices.Collect(maps.Keys(afterValue)))
		if unknown, ok := after.unknown.(map[string]any); ok {
			keys = slices.Concat(keys, slices.Collect(maps.Keys(unknown)))
		}
		slices.Sort(keys)

		for _, key := range slices.Compact(keys) {
			child, changed := d.compare(key, append(path, key), before.child(key), after.child(key))
			if changed {
				change.Children = append(change.Children, child)
			}
		}

		if len(change.Children) == 0 {
			return terraform.Change{}, false
		}

		return d.note(path, change), true
	case []any:
		afterValue, ok := after.value.([]any)
		if !ok || opaque || (!terraform.Expandable(beforeValue) && !terraform.Expandable(afterValue)) {
			break
		}

		for index := range max(len(beforeValue), len(afterValue)) {
			child, changed := d.compare(
				fmt.Sprintf("[%d]", index),
				append(path, index),
				before.element(index),
				after.element(index),
			)
			if changed {
				change.Children = append(change.Children, child)
			}
		}

		if len(change.Children) == 0 {
			return terraform.Change{}, false
		}

		return d.note(path, change), true
	}

	if !after.isUnknown() && before.isSensitive() == after.isSensitive() && reflect.DeepEqual(before.value, after.value) {
		return terraform.Change{}, false
	}

	change.Before = terraform.FormatValue(before.masked())
	change.After = terraform.FormatValue(after.masked())

	return d.note(path, change), true
}

func (d changeDiff) note(path []any, change terraform.Change) terraform.Change {
	if d.replace[joinPath(path)] {
		change.Note = "forces replacement"
	}

	return change
}

// Terraform writes a path as a list of attribute names and element indexes, and the
// indexes come out of the JSON as floats.
func joinPath(path []any) string {
	parts := make([]string, 0, len(path))
	for _, step := range path {
		parts = append(parts, fmt.Sprintf("%v", step))
	}

	return strings.Join(parts, ".")
}
