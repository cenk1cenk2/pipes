package pipe

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/libraries/go-utils/v2/utils"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v6"
	toc "github.com/ekalinin/github-markdown-toc.go"
)

func FindMarkdownFiles(tl *TaskList) *Task {
	return tl.CreateTask("discover").
		Set(func(t *Task) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			t.Log.Debugf(
				"Trying to match patterns: %s",
				strings.Join(P.Markdown.Patterns, ", "),
			)

			matches := []string{}

			for _, v := range P.Markdown.Patterns {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				return fmt.Errorf(
					"Can not match any files with the given pattern: %s",
					strings.Join(P.Markdown.Patterns, ", "),
				)
			}

			matches = utils.RemoveDuplicateStr(matches)

			t.Log.Debugf("Paths matched for given pattern: %s", strings.Join(matches, ", "))

			C.Matches = matches

			return nil
		})
}

func RunMarkdownToc(tl *TaskList) *Task {
	return tl.CreateTask("markdown-toc").
		Set(func(t *Task) error {
			const start = "<!-- toc -->"
			const end = "<!-- tocstop -->"
			expr := fmt.Sprintf(`(?s)%s(.*)%s`, start, end)

			t.Log.Debugf("Using expression: %s", expr)

			for _, match := range C.Matches {
				t.CreateSubtask(match).
					Set(func(t *Task) error {
						parser := toc.NewGHDoc(match, false, P.Markdown.StartDepth, P.Markdown.EndDepth, false, "", P.Markdown.Indentation, false)

						p := parser.GetToc()

						var b bytes.Buffer

						if err := p.Print(&b); err != nil {
							return fmt.Errorf("failed to generate table of contents: %w", err)
						}

						s := b.String()

						marker := regexp.MustCompile(`(?m)^(\s+)?\*`)

						s = marker.ReplaceAllString(s, fmt.Sprintf(`$1%s`, P.Markdown.ListIdentifier))

						t.Log.Debugf("Trying to read file: %s", match)

						content, err := os.ReadFile(match)

						if err != nil {
							return err
						}

						readme := string(content)

						r := regexp.MustCompile(expr)

						replace := strings.Join([]string{start, "", strings.TrimSpace(s), "", end}, "\n")

						result := r.ReplaceAllString(readme, replace)

						if err := os.WriteFile(match, []byte(result), 0600); err != nil {
							return err
						}

						t.Log.Infof("Processed file: %s", match)

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).ShouldRunAfter(func(t *Task) error {
		return t.RunSubtasks()
	})
}
