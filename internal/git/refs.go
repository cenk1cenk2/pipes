package git

import (
	"fmt"
	"regexp"
	"slices"
)

const (
	REFERENCE_HEADS = "heads"
	REFERENCE_TAGS  = "tags"
)

// Refs is the source control position the pipe was triggered from.
type Refs struct {
	Branch string
	Tag    string
}

// References renders the refs as the "<kind>/<name>" strings that user supplied
// match patterns are written against. The tag comes first so a tagged pipeline,
// which also carries the branch it was tagged on, matches the tag rule.
func (r Refs) References() []string {
	references := []string{}

	if r.Tag != "" {
		references = append(references, fmt.Sprintf("%s/%s", REFERENCE_TAGS, r.Tag))
	}

	if r.Branch != "" {
		references = append(references, fmt.Sprintf("%s/%s", REFERENCE_HEADS, r.Branch))
	}

	return references
}

// MatchAny returns the index of the first pattern matching any of the references,
// or -1 when none of them do. Patterns are ranked ahead of references so the
// earliest rule the user wrote wins, whichever ref it happened to match.
func MatchAny(patterns []string, references []string) (int, error) {
	for index, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return -1, fmt.Errorf("Can not process regular expression: %s -> %w", pattern, err)
		}

		if slices.ContainsFunc(references, re.MatchString) {
			return index, nil
		}
	}

	return -1, nil
}
