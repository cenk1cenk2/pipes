package tagsfile

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

const (
	CATEGORY_TAGS_FILE = "Tags File"
)

// NewFlags builds the tags file flags. A nil strict destination leaves the strict
// flag out, for the pipes that only ever read the file leniently.
func NewFlags(dst *string, value string, strict *bool, required bool) []cli.Flag {
	list := []cli.Flag{
		&cli.StringFlag{
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags-file",
			Sources:     flags.EnvVars("TAGS_FILE"),
			Usage:       "Read tags from a file.",
			Required:    required,
			Value:       value,
			Destination: dst,
		},
	}

	if strict == nil {
		return list
	}

	return append(list, &cli.BoolFlag{
		Category:    CATEGORY_TAGS_FILE,
		Name:        "tags-file.strict",
		Sources:     flags.EnvVars("TAGS_FILE_STRICT"),
		Usage:       "Fail on missing tags file.",
		Required:    false,
		Value:       false,
		Destination: strict,
	})
}
