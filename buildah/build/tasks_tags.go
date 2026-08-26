package build

import (
	"fmt"
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
	"go.yaml.in/yaml/v4"
)

// The collector reads the parsed flags, so it is only built from inside a task
// list, never at package level.
func ContainerImageTags(deps Deps) *versions.Collector {
	return &versions.Collector{
		Name:  "tags",
		Label: "Image tags",

		FromUser: P.ContainerImage.Tags,

		File:       P.ContainerImage.TagsFile,
		FileStrict: P.ContainerImage.TagsFileStrict,
		FileDir:    P.ContainerFile.Context,

		LatestWhen:  P.ContainerImage.TagAsLatest,
		LatestValue: P.ContainerImage.LatestTag,
		References:  P.Git.References(),

		Templates: P.ContainerImage.TagsTemplate,
		Sanitize:  P.ContainerImage.TagsSanitize,

		Format: func(tag string) string {
			if deps.Registry.Uri == "" {
				return fmt.Sprintf("%s:%s", P.ContainerImage.Name, tag)
			}

			return fmt.Sprintf("%s/%s:%s", deps.Registry.Uri, P.ContainerImage.Name, tag)
		},
	}
}

func ContainerManifestFileWrite(tl *TaskList, collector *versions.Collector) *Task {
	return tl.CreateTask("tags", "manifest").
		ShouldDisable(func(t *Task) bool {
			return P.ContainerManifest.File == "" || P.ContainerManifest.Target == ""
		}).
		Set(func(t *Task) error {
			target, err := InlineTemplate(P.ContainerManifest.Target, C.Tags)
			if err != nil {
				return err
			}

			image, err := collector.Process(t.Log, target)
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
