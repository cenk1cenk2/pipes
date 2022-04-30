package pipe

import (
	"os"
	"os/exec"
	"strings"

	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	u "gitlab.kilic.dev/libraries/go-utils/utils"

	glob "github.com/bmatcuk/doublestar/v4"
)

type Ctx struct {
	Matches []string
}

var Context Ctx

func FindPackages() utils.Task {
	metadata := utils.TaskMetadata{Context: "discover"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {
		log := utils.Log.WithField("context", t.Metadata.Context)

		cwd, err := os.Getwd()

		if err != nil {
			return err
		}

		fs := os.DirFS(cwd)

		log.Debugf(
			"Trying to match patterns: %s",
			strings.Join(Pipe.TypeDoc.Patterns.Value(), ", "),
		)

		matches := []string{}

		for _, v := range Pipe.TypeDoc.Patterns.Value() {
			match, err := glob.Glob(fs, v)

			if err != nil {
				return err
			}

			matches = append(matches, match...)
		}

		if len(matches) == 0 {
			log.Fatalf(
				"Can not match any files with the given pattern: %s",
				strings.Join(Pipe.TypeDoc.Patterns.Value(), ", "),
			)
		}

		matches = u.RemoveDuplicateStr(matches)

		log.Debugf("Paths matched for given pattern: %s", strings.Join(matches, ", "))

		Context.Matches = matches

		return nil
	}}
}

func RunTypeDoc() utils.Task {
	metadata := utils.TaskMetadata{Context: "generate"}

	return utils.Task{Metadata: metadata, Task: func(t *utils.Task) error {

		for _, match := range Context.Matches {
			cmd := exec.Command("yarn", "exec", TYPEDOC_COMMAND, Pipe.TypeDoc.Arguments)

			cmd.Dir = match

			utils.AddTask(utils.Task{Metadata: metadata, Command: cmd})
		}

		return nil
	}}
}
