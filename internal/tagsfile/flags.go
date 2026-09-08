package tagsfile

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_TAGS_FILE = "Tags File"
)

// Options is what the tags file flags are built onto. A nil Strict destination
// leaves the strict flag out, for the pipes that only ever read the file
// leniently.
type Options struct {
	Destination *string
	Value       string
	Strict      *bool
	Required    bool
}

func NewFlags(opts Options) []cli.Flag {
	list := []cli.Flag{
		&cli.StringFlag{
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags-file",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("TAGS_FILE")),
			Usage:       "Read tags from a comma separated file.",
			Required:    opts.Required,
			Value:       opts.Value,
			Destination: opts.Destination,
		},
	}

	if opts.Strict == nil {
		return list
	}

	return append(list, &cli.BoolFlag{
		Category:    CATEGORY_TAGS_FILE,
		Name:        "tags-file.strict",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("TAGS_FILE_STRICT")),
		Usage:       "Fail on missing tags file.",
		Required:    false,
		Value:       false,
		Destination: opts.Strict,
	})
}
