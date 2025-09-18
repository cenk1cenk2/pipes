package build

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/common/parser"
	"go.yaml.in/yaml/v4"
)

func ContainerImageTagsParent(tl *TaskList) *Task {
	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				JobParallel(
					ContainerImageTagsFromUser(tl).Job(),
					ContainerImageTagsFromFile(tl).Job(),
				),
				ContainerImageTagsFromLatest(tl).Job(),
				job,
				ContainerManifestFileWrite(tl).Job(),
			)
		}).
		Set(func(t *Task) error {
			C.Tags = slices.Compact(C.Tags)

			t.Log.Infof(
				"Image tags: %s", strings.Join(C.Tags, ", "),
			)

			return nil
		})
}

func ContainerImageTagsFromUser(tl *TaskList) *Task {
	return tl.CreateTask("tags", "user").
		Set(func(t *Task) error {
			// add all the specified tags
			for _, v := range slices.Compact(P.ContainerImage.Tags) {
				if err := AddContainerImageTag(t, v); err != nil {
					return err
				}
			}

			return nil
		})
}

func ContainerImageTagsFromFile(tl *TaskList) *Task {
	return tl.CreateTask("tags", "file").
		ShouldDisable(func(t *Task) bool {
			return P.ContainerImage.TagsFile == ""
		}).
		Set(func(t *Task) error {
			// add tags through tags file
			tags, err := parser.ParseTagsFile(t.Log, path.Join(P.ContainerFile.Context, P.ContainerImage.TagsFile), P.ContainerImage.TagsFileStrict)

			if err != nil {
				return err
			}

			for _, v := range tags {
				t.CreateSubtask(v).
					Set(func(t *Task) error {
						return AddContainerImageTag(t, v)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func ContainerImageTagsFromLatest(tl *TaskList) *Task {
	return tl.CreateTask("tags", "latest").
		ShouldDisable(func(t *Task) bool {
			return P.ContainerImage.TagAsLatest == nil
		}).
		Set(func(t *Task) error {
		out:
			for _, expression := range P.ContainerImage.TagAsLatest {
				for _, reference := range C.References {
					re, err := regexp.Compile(expression)

					if err != nil {
						return fmt.Errorf("Can not process regular expression for latest tag: %w", err)
					}

					t.Log.Debugf("Trying to match condition for given reference: %s with %v", reference, re.String())

					if re.MatchString(reference) {
						if err := AddContainerImageTag(t, P.ContainerImage.LatestTag); err != nil {
							return err
						}

						t.Log.Infof(
							"Will tag image as latest since tag regex matches: %s -> %s",
							P.ContainerImage.LatestTag,
							expression,
						)

						break out
					}
				}
			}

			return nil
		})
}

func ContainerManifestFileWrite(tl *TaskList) *Task {
	return tl.CreateTask("tags", "manifest").
		ShouldDisable(func(t *Task) bool {
			return P.ContainerManifest.File == "" || P.ContainerManifest.Target == ""
		}).
		Set(func(t *Task) error {
			target, err := InlineTemplate(P.ContainerManifest.Target, C.Tags)
			if err != nil {
				return err
			}

			image, err := ProcessContainerImageTag(t, target)
			if err != nil {
				return err
			}

			tags, err := yaml.Marshal(&manifest.ContainerManifestMatrix{
				Target: image,
				Images: C.Tags,
			})

			if err != nil {
				return err
			}

			filename, err := InlineTemplate(P.ContainerManifest.File, C.Tags)

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
