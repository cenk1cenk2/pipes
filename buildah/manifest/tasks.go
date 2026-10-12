package manifest

import (
	"context"
	json "encoding/json/v2"
	"fmt"
	"os"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

// Every build job of the pipeline hands the images it pushed over in a dotenv
// variable of its own, which GitLab exports into this job.
func fetchEnv(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "env").
		Set(func(_ context.Context, t *Task) error {
			for _, env := range os.Environ() {
				key, value, _ := strings.Cut(env, "=")

				if !strings.HasPrefix(key, P.Manifest.PublishedPrefix) {
					continue
				}

				parsed := []ManifestMatrix{}
				if err := json.Unmarshal([]byte(value), &parsed); err != nil {
					return fmt.Errorf("Can not unmarshal container manifest matrix: %s -> %w", key, err)
				}

				for _, entry := range parsed {
					if entry.Target == "" {
						continue
					}

					t.Lock.Lock()
					C.ManifestedImages[entry.Target] = append(C.ManifestedImages[entry.Target], entry.Images...)
					t.Lock.Unlock()

					t.Log.Debug(fmt.Sprintf("Found published images: %v for %s in %s", entry.Images, entry.Target, key))
				}
			}

			return nil
		})
}

func fetchUser(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "user").
		ShouldDisable(func(t *Task) bool {
			return len(P.Manifest.Images) == 0 && len(P.Manifest.Matrix) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			if P.Manifest.Target != "" && len(P.Manifest.Images) > 0 {
				target, err := InlineTemplate[any](P.Manifest.Target, nil)
				if err != nil {
					return err
				}

				t.Lock.Lock()
				C.ManifestedImages[target] = append(C.ManifestedImages[target], P.Manifest.Images...)
				t.Lock.Unlock()

				t.Log.Debug(fmt.Sprintf("Fetched direct image: %s -> %v", target, P.Manifest.Images))
			}

			for _, entry := range P.Manifest.Matrix {
				t.Lock.Lock()
				C.ManifestedImages[entry.Target] = append(C.ManifestedImages[entry.Target], entry.Images...)
				t.Lock.Unlock()

				t.Log.Debug(fmt.Sprintf("Fetched manifest from matrix: %s -> %v", entry.Target, entry.Images))
			}

			return nil
		})
}

func manifest(tl *TaskList) *Task {
	return tl.CreateTask("manifest").
		Set(func(ctx context.Context, t *Task) error {
			count := 0

			// the fetchers run in parallel and may each carry the same images, so they are
			// only deduplicated once all of them are in.
			for target, images := range C.ManifestedImages {
				C.ManifestedImages[target] = slices.Compact(slices.Sorted(slices.Values(images)))
				count += len(C.ManifestedImages[target])
			}

			if count == 0 {
				return fmt.Errorf("Can not find any published images to manifest.")
			}

			for target, images := range C.ManifestedImages {
				t.CreateSubtask(target).
					Set(func(_ context.Context, t *Task) error {
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
								"buildah",
								"manifest",
								"push",
								"--rm",
								target,
							).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunCommandJobAsJobSequence(ctx)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}
