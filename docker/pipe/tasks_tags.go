package pipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/libraries/go-utils/utils"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerTagsParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "parent").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			// setup sanitizing the tags first
			if err := json.Unmarshal([]byte(t.Pipe.DockerImage.TagsSanitize), &t.Pipe.Ctx.SanitizedRegularExpression); err != nil {
				return err
			}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobParallel(
					DockerTagsUser(tl).Job(),
					DockerTagsFile(tl).Job(),
					DockerTagsLatest(tl).Job(),
				),
			)

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Pipe.Ctx.Tags = utils.RemoveDuplicateStr(
				utils.DeleteEmptyStringsFromSlice(t.Pipe.Ctx.Tags),
			)

			t.Log.Infof(
				"Image tags: %s", strings.Join(t.Pipe.Ctx.Tags, ", "),
			)

			return nil
		})
}

func DockerTagsUser(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "user").
		Set(func(t *Task[Pipe]) error {
			// add all the specified tags
			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(t.Pipe.DockerImage.Tags.Value())) {
				if err := AddDockerTag(t, v); err != nil {
					return err
				}
			}

			return nil
		})
}

func DockerTagsFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "file").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerImage.TagsFile == ""
		}).
		Set(func(t *Task[Pipe]) error {
			// add tags through tags file
			if _, err := os.Stat(t.Pipe.DockerImage.TagsFile); err == nil {
				t.Log.Infof(
					"Tags file exists: %s",
					t.Pipe.DockerImage.TagsFile,
				)

				content, err := os.ReadFile(t.Pipe.DockerImage.TagsFile)
				if err != nil {
					return err
				}

				tags := strings.Split(string(content), ",")

				t.Pipe.Ctx.Tags = []string{}

				re := regexp.MustCompile(`\r?\n`)

				for _, v := range tags {
					func(v string) {
						t.CreateSubtask(v).
							Set(func(t *Task[Pipe]) error {
								return AddDockerTag(t, re.ReplaceAllString(v, ""))
							}).
							AddSelfToTheParentAsParallel()
					}(v)
				}
			} else if errors.Is(err, os.ErrNotExist) && t.Pipe.DockerImage.TagsFile != "" {
				if !t.Pipe.DockerImage.TagsFileIgnoreMissing {
					t.Log.Warnf("Tags file is set but it does not exists: %s", t.Pipe.DockerImage.TagsFile)

					t.SendExit(0)
				}

				return nil
			} else {
				return fmt.Errorf("Can not read the tags file: %s", t.Pipe.DockerImage.TagsFile)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func DockerTagsLatest(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "latest").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerImage.TagAsLatest == ""
		}).
		Set(func(t *Task[Pipe]) error {
			tagAsLatestRes := []string{}

			if err := json.Unmarshal([]byte(t.Pipe.DockerImage.TagAsLatest), &tagAsLatestRes); err != nil {
				return err
			}

			references := []string{}

			if t.Pipe.Git.Tag != "" {
				references = append(references, fmt.Sprintf("%s/%s", GIT_REFERENCE_TAGS, t.Pipe.Git.Tag))
			}

			if t.Pipe.Git.Branch != "" {
				references = append(references, fmt.Sprintf("%s/%s", GIT_REFERENCE_BRANCH, t.Pipe.Git.Branch))
			}

			if len(references) == 0 {
				return nil
			}

			t.Log.Debugf("References for latest search: %v", references)

		out:
			for _, re := range tagAsLatestRes {
				for _, reference := range references {
					m, err := regexp.Match(re, []byte(reference))

					if err != nil {
						return err
					}

					if m {
						if err := AddDockerTag(t, DOCKER_LATEST_TAG); err != nil {
							return err
						}

						t.Log.Infof(
							"Will tag image as latest since tag regex matches: %s",
							re,
						)

						break out
					}
				}
			}

			return nil
		})
}
