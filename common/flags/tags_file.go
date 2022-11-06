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
			Name:        "tags.file",
			Usage:       "Read tags from a file.",
			Required:    setup.TagsFileRequired,
			EnvVars:     []string{"TAGS_FILE"},
			Value:       setup.TagsFileValue,
			Destination: setup.TagsFileDestination,
		},
	}
}

type TagsFileIgnoreMissingSetup struct {
	TagsFileIgnoreMissingDestination *bool
	TagsFileIgnoreMissingRequired    bool
	TagsFileIgnoreMissingValue       bool
}

func NewTagsFileIgnoreMissingFlags(setup TagsFileIgnoreMissingSetup) []cli.Flag {
	return []cli.Flag{
		// CATEGORY_TAGS_FILE
		&cli.BoolFlag{
			Category:    CATEGORY_TAGS_FILE,
			Name:        "tags.ignore-missing",
			Usage:       "Ignore the missing tags file and contunie operation as expected in that case.",
			Required:    setup.TagsFileIgnoreMissingRequired,
			EnvVars:     []string{"TAGS_FILE_IGNORE_MISSING"},
			Value:       setup.TagsFileIgnoreMissingValue,
			Destination: setup.TagsFileIgnoreMissingDestination,
		},
	}
}
