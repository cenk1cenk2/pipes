package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
)

func DiscoverPublishedImageFiles(tl *TaskList) *Task {
	return tl.CreateTask("discover", "file").
		ShouldDisable(func(t *Task) bool {
			return len(P.DockerManifest.Files) == 0
		}).
		Set(func(t *Task) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			matches := []string{}

			for _, v := range P.DockerManifest.Files {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Warnf(
					"Can not match any files with the given pattern: %s",
					strings.Join(P.DockerManifest.Files, ", "),
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

						parsed := &DockerManifestMatrixJson{}
						if err := json.Unmarshal(content, parsed); err != nil {
							return fmt.Errorf("Can not unmarshal Docker manifest matrix: %w", err)
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
			return len(P.DockerManifest.Images) == 0
		}).
		Set(func(t *Task) error {
			if P.DockerManifest.Target != "" && len(P.DockerManifest.Images) > 0 {
				t.Lock.Lock()
				var err error
				if P.DockerManifest.Target, err = InlineTemplate[any](P.DockerManifest.Target, nil); err != nil {
					return err
				}

				C.ManifestedImages[P.DockerManifest.Target] = append(C.ManifestedImages[P.DockerManifest.Target], P.DockerManifest.Images...)
				t.Lock.Unlock()

				t.Log.Debugf("Fetched direct image: %s -> %v", P.DockerManifest.Target, P.DockerManifest.Images)
			}

			for _, manifest := range P.DockerManifest.Matrix {
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
								setup.DOCKER_EXE,
								"manifest",
								"create",
								target,
							).
							Set(func(c *Command) error {
								for _, image := range images {
									c.AppendArgs(fmt.Sprintf("-a %s", image))
								}

								return nil
							}).
							AddSelfToTheTask()

						t.
							CreateCommand(
								setup.DOCKER_EXE,
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
