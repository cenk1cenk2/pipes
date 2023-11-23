package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func DiscoverPublishedImageFiles(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("discover", "file").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.DockerManifest.Files) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			matches := []string{}

			for _, v := range t.Pipe.DockerManifest.Files {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Warnf(
					"Can not match any files with the given pattern: %s",
					strings.Join(t.Pipe.DockerManifest.Files, ", "),
				)

				return nil
			}

			matches = utils.RemoveDuplicateStr(matches)

			t.Log.Debugf("Paths matched for given pattern: %s", strings.Join(matches, ", "))

			t.Pipe.Ctx.Matches = matches

			return nil
		})
}

func FetchPublishedImagesFromFiles(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fetch", "file").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.Ctx.Matches) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, f := range t.Pipe.Ctx.Matches {
				func(f string) {
					t.CreateSubtask(f).
						Set(func(t *Task[Pipe]) error {
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
							t.Pipe.Ctx.ManifestedImages[parsed.Target] = append(t.Pipe.Ctx.ManifestedImages[parsed.Target], parsed.Images...)
							t.Lock.Unlock()

							return nil
						}).
						AddSelfToTheParentAsParallel()
				}(f)
			}
			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func FetchUserPublishedImages(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fetch", "user").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.DockerManifest.Images) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			if t.Pipe.DockerManifest.Target != "" && len(t.Pipe.DockerManifest.Images) > 0 {
				t.Lock.Lock()
				var err error
				if t.Pipe.DockerManifest.Target, err = InlineTemplate[any](t.Pipe.DockerManifest.Target, nil); err != nil {
					return err
				}

				t.Pipe.Ctx.ManifestedImages[t.Pipe.DockerManifest.Target] = append(t.Pipe.Ctx.ManifestedImages[t.Pipe.DockerManifest.Target], t.Pipe.DockerManifest.Images...)
				t.Lock.Unlock()

				t.Log.Debugf("Fetched direct image: %s -> %v", t.Pipe.DockerManifest.Target, t.Pipe.DockerManifest.Images)
			}

			for _, manifest := range t.Pipe.DockerManifest.Matrix {
				t.Lock.Lock()
				t.Pipe.Ctx.ManifestedImages[manifest.Target] = append(t.Pipe.Ctx.ManifestedImages[manifest.Target], manifest.Images...)
				t.Lock.Unlock()

				t.Log.Debugf("Fetched manifest from matrix: %s -> %v", manifest.Target, manifest.Images)
			}

			return nil
		})
}

func UpdateManifests(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("manifest").
		Set(func(t *Task[Pipe]) error {
			for target, images := range t.Pipe.Ctx.ManifestedImages {
				func(target string, images []string) {
					t.CreateSubtask(target).
						Set(func(t *Task[Pipe]) error {
							t.
								CreateCommand(
									setup.DOCKER_EXE,
									"manifest",
									"create",
									target,
								).
								Set(func(c *Command[Pipe]) error {
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
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobSequence()
						}).
						AddSelfToTheParentAsParallel()
				}(target, images)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
