package pipe

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/libraries/go-utils/utils"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func DockerTagsParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags:parent").
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobParallel(
					DockerTagsUser(tl).Job(),
					DockerTagsFile(tl).Job(),
					DockerTagsLatestFromTag(tl).Job(),
					DockerTagsLatestFromBranch(tl).Job(),
				),
			)

			return nil
		}).ShouldRunAfter(func(t *Task[Pipe]) error {
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
	return tl.CreateTask("tags:user").
		Set(func(t *Task[Pipe]) error {
			// add all the specified tags
			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(t.Pipe.DockerImage.Tags.Value())) {
				if err := AddDockerTag(v); err != nil {
					return err
				}
			}

			return nil
		})
}

func DockerTagsFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags:file").
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

				content, err := ioutil.ReadFile(t.Pipe.DockerImage.TagsFile)
				if err != nil {
					return err
				}

				tags := strings.Split(string(content), ",")

				t.Pipe.Ctx.Tags = []string{}

				re := regexp.MustCompile(`\r?\n`)

				for _, v := range tags {
					func(v string) {
						t.CreateSubtask("").
							Set(func(t *Task[Pipe]) error {
								return AddDockerTag(re.ReplaceAllString(v, ""))
							}).
							AddSelfToParent(func(pt, st *Task[Pipe]) {
								pt.ExtendSubtask(func(j Job) Job {
									return tl.JobParallel(j, st.Job())
								})
							})
					}(v)
				}

				if err = t.RunSubtasks(); err != nil {
					return err
				}
			} else if errors.Is(err, os.ErrNotExist) && t.Pipe.DockerImage.TagsFile != "" {
				t.Log.Warnf("Tags file is set but it does not exists: %s", t.Pipe.DockerImage.TagsFile)

				if !t.Pipe.DockerImage.TagsFileIgnoreMissing {
					t.SendExit(0)
				}

				return nil
			} else {
				t.Log.Warnf("Can not read the tags file: %s", t.Pipe.DockerImage.TagsFile)
			}

			return nil
		})
}

func DockerTagsLatestFromTag(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags:latest:tag").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerImage.TagAsLatestForTagsRegex == "" || t.Pipe.Git.Tag == ""
		}).
		Set(func(t *Task[Pipe]) error {
			TagAsLatestForTagsRegex := []string{}

			if err := json.Unmarshal([]byte(t.Pipe.DockerImage.TagAsLatestForTagsRegex), &TagAsLatestForTagsRegex); err != nil {
				return err
			}

			for _, re := range TagAsLatestForTagsRegex {
				m, err := regexp.Match(re, []byte(t.Pipe.Git.Tag))

				if err != nil {
					return err
				}

				if m {
					if err := AddDockerTag(DOCKER_LATEST_TAG); err != nil {
						return err
					}

					t.Log.Infof(
						"Will tag image as latest since tag regex matches: %s",
						re,
					)

					break
				}
			}

			return nil
		})
}

func DockerTagsLatestFromBranch(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags:latest:branch").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerImage.TagAsLatestForBranchesRegex == "" || t.Pipe.Git.Branch == ""
		}).
		Set(func(t *Task[Pipe]) error {
			TagAsLatestForBranchesRegex := []string{}

			if err := json.Unmarshal([]byte(t.Pipe.DockerImage.TagAsLatestForBranchesRegex), &TagAsLatestForBranchesRegex); err != nil {
				return err
			}

			for _, re := range TagAsLatestForBranchesRegex {
				m, err := regexp.Match(re, []byte(t.Pipe.Git.Branch))

				if err != nil {
					return err
				}

				if m {
					if err := AddDockerTag(DOCKER_LATEST_TAG); err != nil {
						return err
					}

					t.Log.Infof(
						"Will tag image as latest since branch regex matches: %s",
						re,
					)

					break
				}
			}

			return nil
		})
}
