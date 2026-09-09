package manifest

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	. "github.com/cenk1cenk2/plumber/v7"
	"go.yaml.in/yaml/v4"
)

func discoverFile(tl *TaskList) *Task {
	return tl.CreateTask("discover", "file").
		ShouldDisable(func(t *Task) bool {
			return len(P.Manifest.Files) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			cwd, err := os.Getwd()

			if err != nil {
				return err
			}

			fs := os.DirFS(cwd)

			matches := []string{}

			for _, v := range P.Manifest.Files {
				match, err := glob.Glob(fs, v)

				if err != nil {
					return err
				}

				matches = append(matches, match...)
			}

			if len(matches) == 0 {
				t.Log.Warn(fmt.Sprintf("Can not match any files with the given pattern: %s",
					strings.Join(P.Manifest.Files, ", ")),
				)

				return nil
			}

			matches = slices.Compact(slices.Sorted(slices.Values(matches)))

			t.Log.Debug(fmt.Sprintf("Paths matched for given pattern: %s", strings.Join(matches, ", ")))

			C.Matches = matches

			return nil
		})
}

func fetchFile(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "file").
		ShouldDisable(func(t *Task) bool {
			return len(C.Matches) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, f := range C.Matches {
				t.CreateSubtask(f).
					Set(func(_ context.Context, t *Task) error {
						content, err := os.ReadFile(f)
						if err != nil {
							return err
						}

						parsed := &ManifestMatrix{}
						if err := yaml.Unmarshal(content, parsed); err != nil {
							return fmt.Errorf("Can not unmarshal container manifest matrix: %w", err)
						}

						if parsed.Target == "" {
							return nil
						}

						t.Log.Debug(fmt.Sprintf("Found published images: %v for %s in %s", parsed.Images, parsed.Target, f))

						t.Lock.Lock()
						C.ManifestedImages[parsed.Target] = append(C.ManifestedImages[parsed.Target], parsed.Images...)
						t.Lock.Unlock()

						return nil
					}).
					AddSelfToTheParentAsParallel()
			}
			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}

func fetchUser(tl *TaskList) *Task {
	return tl.CreateTask("fetch", "user").
		ShouldDisable(func(t *Task) bool {
			return len(P.Manifest.Images) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			if P.Manifest.Target != "" && len(P.Manifest.Images) > 0 {
				t.Lock.Lock()
				var err error
				if P.Manifest.Target, err = InlineTemplate[any](P.Manifest.Target, nil); err != nil {
					return err
				}

				C.ManifestedImages[P.Manifest.Target] = append(C.ManifestedImages[P.Manifest.Target], P.Manifest.Images...)
				t.Lock.Unlock()

				t.Log.Debug(fmt.Sprintf("Fetched direct image: %s -> %v", P.Manifest.Target, P.Manifest.Images))
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
