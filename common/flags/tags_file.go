package flags

import (
	"github.com/urfave/cli/v3"
)

type TagsFileFlagsSetup struct {
	TagsFileDestination *string
	TagsFileRequired    bool
	TagsFileValue       string
}

func NewTagsFileFlags(setup TagsFileFlagsSetup) []cli.Flag {
	return []cli.Flag{
		// CATEGORY_TAGS_FILE
		&cli.StringFlag{
			Category: CATEGORY_TAGS_FILE,
			Name:     "tags-file",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TAGS_FILE"),
			),
			Usage:       "Read tags from a file.",
			Required:    setup.TagsFileRequired,
			Value:       setup.TagsFileValue,
			Destination: setup.TagsFileDestination,
		},
	}
}

type TagsFileStrictFlagsSetup struct {
	TagsFileStrictDestination *bool
	TagsFileStrictRequired    bool
	TagsFileStrictValue       bool
}

func NewTagsFileStrictFlags(setup TagsFileStrictFlagsSetup) []cli.Flag {
	return []cli.Flag{
		// CATEGORY_TAGS_FILE
		&cli.BoolFlag{
			Category: CATEGORY_TAGS_FILE,
			Name:     "tags-file.strict",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("TAGS_FILE_STRICT"),
			),
			Usage:       "Fail on missing tags file.",
			Required:    setup.TagsFileStrictRequired,
			Value:       setup.TagsFileStrictValue,
			Destination: setup.TagsFileStrictDestination,
		},
	}
}
