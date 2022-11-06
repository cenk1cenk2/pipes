package parser

import (
	"fmt"

	"gitlab.kilic.dev/devops/pipes/common/constants"
)

func ParseGitReferences(tag string, branch string) []string {
	references := []string{}

	if tag != "" {
		references = append(references, fmt.Sprintf("%s/%s", constants.GIT_REFERENCE_TAGS, tag))
	}

	if branch != "" {
		references = append(references, fmt.Sprintf("%s/%s", constants.GIT_REFERENCE_BRANCH, branch))
	}

	return references
}
