package manifest

import (
	"fmt"
	"os"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
	"go.yaml.in/yaml/v4"
)

func DiscoverPublishedImageFiles(tl *TaskList) *Task {
	return tl.CreateTask("discover", "file").
		ShouldDisable(func(t *Task) bool {
			return len(P.ContainerManifest.Files) == 0
		}).
		Set(func(t *Task) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			matches := []string{}

			for _, v := range P.ContainerManifest.Files {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Warnf(
					"Can not match any files with the given pattern: %s",
					strings.Join(P.ContainerManifest.Files, ", "),
				)

				return nil
			}

			matches = utils.RemoveDuplicateStr(matches)

			t.Log.Debugf("Paths matched for given pattern: %s", strings.Join(matches, ", "))

			C.Matches = matches

			return nil
		})
}

func FetchPublishedImagesFromFiles(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "file").
		ShouldDisable(func(t *Task) bool {
			return len(C.Matches) == 0
		}).
		Set(func(t *Task) error {
			for _, f := range C.Matches {
				t.CreateSubtask(f).
					Set(func(t *Task) error {
						content, err := os.ReadFile(f)
						if err != nil {
							return err
						}

						parsed := &ContainerManifestMatrix{}
						if err := yaml.Unmarshal(content, parsed); err != nil {
							return fmt.Errorf("Can not unmarshal container manifest matrix: %w", err)
						}

						if parsed.Target == "" {
							return nil
						}

						t.Log.Debugf("Found published images: %v for %s in %s", parsed.Images, parsed.Target, f)

						t.Lock.Lock()
						C.ManifestedImages[parsed.Target] = append(C.ManifestedImages[parsed.Target], parsed.Images...)
						t.Lock.Unlock()

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}
			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func FetchUserPublishedImages(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "user").
		ShouldDisable(func(t *Task) bool {
			return len(P.ContainerManifest.Images) == 0
		}).
		Set(func(t *Task) error {
			if P.ContainerManifest.Target != "" && len(P.ContainerManifest.Images) > 0 {
				t.Lock.Lock()
				var err error
				if P.ContainerManifest.Target, err = InlineTemplate[any](P.ContainerManifest.Target, nil); err != nil {
					return err
				}

				C.ManifestedImages[P.ContainerManifest.Target] = append(C.ManifestedImages[P.ContainerManifest.Target], P.ContainerManifest.Images...)
				t.Lock.Unlock()

				t.Log.Debugf("Fetched direct image: %s -> %v", P.ContainerManifest.Target, P.ContainerManifest.Images)
			}

			for _, manifest := range P.ContainerManifest.Matrix {
				t.Lock.Lock()
				C.ManifestedImages[manifest.Target] = append(C.ManifestedImages[manifest.Target], manifest.Images...)
				t.Lock.Unlock()

				t.Log.Debugf("Fetched manifest from matrix: %s -> %v", manifest.Target, manifest.Images)
			}

			return nil
		})
}

func UpdateManifests(tl *TaskList) *Task {
	return tl.CreateTask("manifest").
		Set(func(t *Task) error {
			for target, images := range C.ManifestedImages {
				t.CreateSubtask(target).
					Set(func(t *Task) error {
						t.
							CreateCommand(
								"buildah",
								"manifest",
								"create",
								target,
							).
							AddSelfToTheTask()

						for _, image := range images {
							t.
								CreateCommand(
									"buildah",
									"manifest",
									"add",
									target,
									image,
								).
								AddSelfToTheTask()
						}

						t.
							CreateCommand(
								"buldah",
								"manifest",
								"push",
								"-p",
								target,
							).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobSequence()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
