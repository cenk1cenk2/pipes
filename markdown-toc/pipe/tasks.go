package pipe

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"

	glob "github.com/bmatcuk/doublestar/v4"
)

type Ctx struct {
	Matches []string
}

var Context Ctx

func FindMarkdownFiles() utils.Task {
	metadata := utils.TaskMetadata{Context: "discover"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		log := utils.Log.WithField("context", t.Metadata.Context)

		cwd, err := os.Getwd()

		if err != nil {
			return err
		}

		fs := os.DirFS(cwd)

		log.Debugln(fmt.Sprintf("Trying to match patterns: %s", Pipe.Markdown.Patterns))

		matches := []string{}

		for _, v := range Pipe.Markdown.Patterns.Value() {
			match, err := glob.Glob(fs, v)

			if err != nil {
				return err
			}

			matches = append(matches, match...)
		}

		if len(matches) == 0 {
			log.Fatalln(
				fmt.Sprintf(
					"Can not match any files with the given pattern: %s",
					Pipe.Markdown.Patterns,
				),
			)
		}

		log.Debugln(fmt.Sprintf("Paths matched for given pattern: %s", strings.Join(matches, ", ")))

		Context.Matches = matches

		return nil
	}}
}

func RunMarkdownToc() utils.Task {
	metadata := utils.TaskMetadata{Context: "generate"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {

		for _, match := range Context.Matches {
			cmd := exec.Command(MARKDOWN_TOC_COMMAND, Pipe.Markdown.Arguments, "-i")

			cmd.Args = append(cmd.Args, match)

			utils.AddTask(utils.Task{Metadata: metadata, Command: cmd})
		}

		return nil
	}}
}
