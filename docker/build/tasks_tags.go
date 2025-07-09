package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/parser"
	"gitlab.kilic.dev/devops/pipes/docker/manifest"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
)

func DockerTagsParent(tl *TaskList) *Task {
	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				JobParallel(
					DockerTagsUser(tl).Job(),
					DockerTagsFile(tl).Job(),
				),
				DockerTagsLatest(tl).Job(),
				job,
				DockerTagsWriteManifestFile(tl).Job(),
			)
		}).
		Set(func(t *Task) error {
			C.Tags = utils.RemoveDuplicateStr(
				utils.DeleteEmptyStringsFromSlice(C.Tags),
			)

			t.Log.Infof(
				"Image tags: %s", strings.Join(C.Tags, ", "),
			)

			return nil
		})
}

func DockerTagsWriteManifestFile(tl *TaskList) *Task {
	return tl.CreateTask("tags", "manifest").
		ShouldDisable(func(t *Task) bool {
			return P.DockerManifest.OutputFile == "" || P.DockerManifest.Target == ""
		}).
		Set(func(t *Task) error {
			target, err := InlineTemplate(P.DockerManifest.Target, C.Tags)
			if err != nil {
				return err
			}

			image, err := ProcessDockerTag(t, target)
			if err != nil {
				return err
			}

			tags, err := json.Marshal(&manifest.DockerManifestMatrixJson{
				Target: image,
				Images: C.Tags,
			})

			if err != nil {
				return err
			}

			filename, err := InlineTemplate(P.DockerManifest.OutputFile, C.Tags)

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

func DockerTagsUser(tl *TaskList) *Task {
	return tl.CreateTask("tags", "user").
		Set(func(t *Task) error {
			// add all the specified tags
			for _, v := range utils.RemoveDuplicateStr(utils.DeleteEmptyStringsFromSlice(P.DockerImage.Tags)) {
				if err := AddDockerTag(t, v); err != nil {
					return err
				}
			}

			return nil
		})
}

func DockerTagsFile(tl *TaskList) *Task {
	return tl.CreateTask("tags", "file").
		ShouldDisable(func(t *Task) bool {
			return P.DockerImage.TagsFile == ""
		}).
		Set(func(t *Task) error {
			// add tags through tags file
			tags, err := parser.ParseTagsFile(t.Log, path.Join(P.DockerFile.Context, P.DockerImage.TagsFile), P.DockerImage.TagsFileStrict)

			if err != nil {
				return err
			}

			for _, v := range tags {
				t.CreateSubtask(v).
					Set(func(t *Task) error {
						return AddDockerTag(t, v)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func DockerTagsLatest(tl *TaskList) *Task {
	return tl.CreateTask("tags", "latest").
		ShouldDisable(func(t *Task) bool {
			return P.DockerImage.TagAsLatest == nil
		}).
		Set(func(t *Task) error {
		out:
			for _, expression := range P.DockerImage.TagAsLatest {
				for _, reference := range C.References {
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
