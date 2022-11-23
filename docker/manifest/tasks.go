package manifest

import (
	"fmt"
	"os"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	"gitlab.kilic.dev/libraries/go-utils/v2/utils"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
			return len(t.Pipe.DockerManifest.Files) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, f := range t.Pipe.DockerManifest.Files {
				func(f string) {
					t.CreateSubtask(f).
						Set(func(t *Task[Pipe]) error {
							content, err := os.ReadFile(f)
							if err != nil {
								return err
							}

							tags := strings.Split(string(content), ",")
							t.Log.Debugf("Found published images: %v in %s", tags, f)

							t.Lock.Lock()
							t.Pipe.Ctx.PublishedImages = append(t.Pipe.Ctx.PublishedImages, tags...)
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
			t.Lock.Lock()
			t.Pipe.Ctx.PublishedImages = append(t.Pipe.Ctx.PublishedImages, t.Pipe.DockerManifest.Images...)
			t.Lock.Unlock()

			return nil
		})
}

func UpdateManifests(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("manifest").
		Set(func(t *Task[Pipe]) error {
			for _, target := range t.Pipe.DockerManifest.Targets {
				func(target string) {
					t.CreateSubtask(target).
						Set(func(t *Task[Pipe]) error {
							var matches = []string{}

							parsed := strings.Split(target, ":")
							targetName, targetTag := parsed[0], parsed[1]

							if targetName == "" || targetTag == "" {
								return fmt.Errorf("Error while parsing target tag: %s", target)
							}

							for _, image := range t.Pipe.Ctx.PublishedImages {
								parsed := strings.Split(image, ":")
								name, tag := parsed[0], parsed[1]

								if name == "" || tag == "" {
									return fmt.Errorf("Error while parsing image tag: %s", image)
								}

								if targetName == name {
									t.Log.Debugf("Image matches with target: %s -> %s", image, target)

									matches = append(matches, image)
								}
							}

							if len(matches) == 0 {
								t.Log.Warnf("No matches for given target tag, doing nothing: %s", target)
							}

							t.
								CreateCommand(
									setup.DOCKER_EXE,
									"manifest",
									"create",
									target,
								).
								Set(func(c *Command[Pipe]) error {
									for _, image := range matches {
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
						AddSelfToTheParentAsParallel()
				}(target)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
