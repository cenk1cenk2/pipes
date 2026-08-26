package tagsfile

import (
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const (
	CATEGORY_TAGS_FILE = "Tags File"
)

// NewFlags builds the tags file flags. A nil strict destination leaves the strict
// flag out, for the pipes that only ever read the file leniently.
func NewFlags(dst *string, value string, strict *bool, required bool) []ucli.Flag {
	flags := []ucli.Flag{
		&ucli.StringFlag{
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags-file",
			Sources:     cli.EnvVars("TAGS_FILE"),
			Usage:       "Read tags from a file.",
			Required:    required,
			Value:       value,
			Destination: dst,
		},
	}

	if strict == nil {
		return flags
	}

	return append(flags, &ucli.BoolFlag{
		Category:    CATEGORY_TAGS_FILE,
		Name:        "tags-file.strict",
		Sources:     cli.EnvVars("TAGS_FILE_STRICT"),
		Usage:       "Fail on missing tags file.",
		Required:    false,
		Value:       false,
		Destination: strict,
	})
}
