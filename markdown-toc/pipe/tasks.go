package pipe

import (
	"os"
	"strings"

	"gitlab.kilic.dev/libraries/go-utils/utils"

	glob "github.com/bmatcuk/doublestar/v4"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func FindMarkdownFiles(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("discover").
		Set(func(t *Task[Pipe]) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			t.Log.Debugf(
				"Trying to match patterns: %s",
				strings.Join(t.Pipe.Markdown.Patterns.Value(), ", "),
			)

			matches := []string{}

			for _, v := range t.Pipe.Markdown.Patterns.Value() {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Fatalf(
					"Can not match any files with the given pattern: %s",
					strings.Join(t.Pipe.Markdown.Patterns.Value(), ", "),
				)
			}

			matches = utils.RemoveDuplicateStr(matches)

			t.Log.Debugf("Paths matched for given pattern: %s", strings.Join(matches, ", "))

			t.Pipe.Ctx.Matches = matches

			return nil
		})
}

func RunMarkdownToc(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("markdown-toc").
		Set(func(t *Task[Pipe]) error {
			for _, match := range t.Pipe.Ctx.Matches {
				t.CreateCommand(
					MARKDOWN_TOC_COMMAND,
					t.Pipe.Markdown.Arguments,
					"-i",
					match,
				).
					AddSelfToTheTask()
			}

			return nil
		}).ShouldRunAfter(func(t *Task[Pipe]) error {
		return t.RunCommandJobAsJobParallel()
	})
}
