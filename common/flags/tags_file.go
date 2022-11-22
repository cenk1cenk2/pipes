package flags

import (
	"github.com/urfave/cli/v2"
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
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags-file",
			Usage:       "Read tags from a file.",
			Required:    setup.TagsFileRequired,
			EnvVars:     []string{"TAGS_FILE"},
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
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags-file.strict",
			Usage:       "Fail on missing tags file.",
			Required:    setup.TagsFileStrictRequired,
			EnvVars:     []string{"TAGS_FILE_STRICT"},
			Value:       setup.TagsFileStrictValue,
			Destination: setup.TagsFileStrictDestination,
		},
	}
}
