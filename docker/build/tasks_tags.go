package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gitlab.kilic.dev/devops/pipes/docker/manifest"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerTagsParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				tl.JobParallel(
					DockerTagsUser(tl).Job(),
					DockerTagsFile(tl).Job(),
					DockerTagsLatest(tl).Job(),
				),
				job,
				DockerTagsWriteManifestFile(tl).Job(),
			)
		}).
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.Tags = utils.RemoveDuplicateStr(
				utils.DeleteEmptyStringsFromSlice(t.Pipe.Ctx.Tags),
			)

			t.Log.Infof(
				"Image tags: %s", strings.Join(t.Pipe.Ctx.Tags, ", "),
			)

			return nil
		})
}

func DockerTagsWriteManifestFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "file").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerManifest.OutputFile == "" || t.Pipe.DockerManifest.Target == ""
		}).
		Set(func(t *Task[Pipe]) error {
			target, err := InlineTemplate(t.Pipe.DockerManifest.Target, t.Pipe.Ctx.Tags)
			if err != nil {
				return err
			}

			image, err := ProcessDockerTag(t, target)
			if err != nil {
				return err
			}

			tags, err := json.Marshal(&manifest.DockerManifestMatrixJson{
				Target: image,
				Images: t.Pipe.Ctx.Tags,
			})

			if err != nil {
				return err
			}

			filename, err := InlineTemplate(t.Pipe.DockerManifest.OutputFile, string(tags))

			t.Log.Debugf("Filename for outputting the tags to: %s", filename)

			if err != nil {
				return err
			}

			if err := os.WriteFile(filename, tags, 0600); err != nil {
				return err
			}

			t.Log.Infof("Wrote image manifest to file for later use: %s", filename)

			return nil
		})
}

func DockerTagsUser(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags", "user").
		Set(func(t *Task[Pipe]) error {
			// add all the specified tags
			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(t.Pipe.DockerImage.Tags)) {
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
				if t.Pipe.DockerImage.TagsFileStrict {
					t.Log.Warnf("Tags file is set but it does not exists: %s", t.Pipe.DockerImage.TagsFile)

					t.SendExit(0)

					return nil
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
			return t.Pipe.DockerImage.TagAsLatest == nil
		}).
		Set(func(t *Task[Pipe]) error {
		out:
			for _, expression := range t.Pipe.DockerImage.TagAsLatest {
				for _, reference := range t.Pipe.Ctx.References {
					re, err := regexp.Compile(expression)

					if err != nil {
						return fmt.Errorf("Can not process regular expression for latest tag: %w", err)
					}

					t.Log.Debugf("Trying to match condition for given reference: %s with %v", reference, re.String())

					if re.MatchString(reference) {
						if err := AddDockerTag(t, setup.DOCKER_LATEST_TAG); err != nil {
							return err
						}

						t.Log.Infof(
							"Will tag image as latest since tag regex matches: %s",
							expression,
						)

						break out
					}
				}
			}

			return nil
		})
}
