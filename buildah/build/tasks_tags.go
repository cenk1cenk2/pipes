package build

import (
	"fmt"
	"os"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
	"go.yaml.in/yaml/v4"
)

// The collector reads the parsed flags, so it is only built from inside a task
// list, never at package level.
func ContainerImageTags() *versions.Collector {
	return &versions.Collector{
		Name: "tags",

		FromUser: P.Image.Tags,

		File:       P.Image.TagsFile,
		FileStrict: P.Image.TagsFileStrict,
		FileDir:    P.File.Context,

		LatestWhen:  P.Image.TagAsLatest,
		LatestValue: P.Image.LatestTag,
		References:  P.Git.References(),

		Templates: P.Image.TagsTemplate,
		Sanitize:  P.Image.TagsSanitize,

		Format: func(tag string) string {
			if login.P.Uri == "" {
				return fmt.Sprintf("%s:%s", P.Image.Name, tag)
			}

			return fmt.Sprintf("%s/%s:%s", login.P.Uri, P.Image.Name, tag)
		},
	}
}

// The manifest write hangs off the parent rather than the sequence around it, so
// it only ever sees a tag list every source has already been collected into.
func ContainerImageTagsParent(tl *TaskList) *Task {
	collector := ContainerImageTags()

	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				JobParallel(
					collector.UserTask(tl, &C.Tags).Job(),
					collector.FileTask(tl, &C.Tags).Job(),
				),
				collector.LatestTask(tl, &C.Tags).Job(),
				job,
				ContainerManifestFileWrite(tl, collector).Job(),
			)
		}).
		Set(func(t *Task) error {
			C.Tags = slices.Compact(C.Tags)

			t.Log.Infof("Image tags: %s", strings.Join(C.Tags, ", "))

			return nil
		})
}

func ContainerManifestFileWrite(tl *TaskList, collector *versions.Collector) *Task {
	return tl.CreateTask("tags", "manifest").
		ShouldDisable(func(t *Task) bool {
			return P.Manifest.File == "" || P.Manifest.Target == ""
		}).
		Set(func(t *Task) error {
			target, err := InlineTemplate(P.Manifest.Target, C.Tags)
			if err != nil {
				return err
			}

			image, err := collector.Process(t.Log, target)
			if err != nil {
				return err
			}

			tags, err := yaml.Marshal(&manifest.ManifestMatrix{
				Target: image,
				Images: C.Tags,
			})

			if err != nil {
				return err
			}

			filename, err := InlineTemplate(P.Manifest.File, C.Tags)

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
