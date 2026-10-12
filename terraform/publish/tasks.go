package publish

import (
	"context"
	"fmt"
	"os"
	"slices"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

// The collector reads the parsed flags, so it is only built from inside a task list.
func tags(tl *TaskList) *Task {
	collector := &versions.Collector{
		Name: "tags",

		FromUser:        P.Module.Tags,
		FromEnvironment: P.Module.TagsFromEnvironment,
	}

	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				collector.UserTask(tl, &C.Tags).Job(),
				collector.EnvironmentTask(tl, &C.Tags).Job(),
				job,
			)
		}).
		Set(func(_ context.Context, t *Task) error {
			C.Tags = slices.Compact(C.Tags)

			if len(C.Tags) > 0 {
				t.Log.Info(fmt.Sprintf("Module tags: %+v", C.Tags))
			} else {
				t.Log.Warn("No tags to publish the module with, doing nothing.")
			}

			return nil
		})
}

func packageTask(tl *TaskList) *Task {
	return tl.CreateTask("package", P.Module.Name, P.Module.System).
		Set(func(ctx context.Context, t *Task) error {
			for _, tag := range C.Tags {
				t.CreateSubtask(tag).
					Set(func(_ context.Context, t *Task) error {
						output := fmt.Sprintf("%s/%s-%s-%s.tar.gz", TFModuleOutputDir, P.Module.Name, P.Module.System, tag)

						t.CreateCommand(
							"tar",
							"-vczf",
							output,
							"--exclude=./.git",
							".",
						).
							SetDir(P.Module.Cwd).
							SetLogLevel(LogLevelDebug, LogLevelDefault, LogLevelDefault).
							ShouldRunBefore(func(_ context.Context, c *Command) error {
								c.Log.Info(fmt.Sprintf("Creating package for tag: %s", tag))

								return nil
							}).
							ShouldRunAfter(func(_ context.Context, c *Command) error {
								t.Lock.Lock()
								C.Packages = append(C.Packages, PublishablePackage{
									Tag:    tag,
									Output: output,
								})
								t.Lock.Unlock()

								return nil
							}).
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

func publish(tl *TaskList) *Task {
	return tl.CreateTask("publish").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				publishGitlab(tl).Job(),
			)
		})
}

func publishGitlab(tl *TaskList) *Task {
	return tl.CreateTask("publish", TFRegistryGitLab, P.Module.Name, P.Module.System).
		ShouldDisable(func(t *Task) bool {
			return P.Registry.Name != TFRegistryGitLab
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, p := range C.Packages {
				t.CreateSubtask(p.Tag).
					Set(func(_ context.Context, t *Task) error {
						file, err := os.Open(p.Output)
						if err != nil {
							return err
						}

						defer file.Close()

						if err := C.Registry.UploadModule(
							context.Background(),
							P.Module.Name,
							P.Module.System,
							p.Tag,
							file,
						); err != nil {
							return err
						}

						t.Log.Info(fmt.Sprintf("Package has been published: %s@%s", P.Module.Name, p.Tag))

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
