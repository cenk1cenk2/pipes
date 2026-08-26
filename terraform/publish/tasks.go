package publish

import (
	"context"
	"fmt"
	"os"
	"path"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
)

func TerraformTagsFile(tl *TaskList) *Task {
	return tl.CreateTask("tags").
		Set(func(t *Task) error {
			tags, err := tagsfile.Parse(t.Log, path.Join(P.Module.Cwd, P.Module.TagsFile), false)

			if err != nil {
				return err
			}

			C.Tags = tags

			if len(C.Tags) > 0 {
				t.Log.Infof("Tags file has been parsed: %+v", C.Tags)
			} else {
				t.Log.Warnln("Tags file does not contain any tags, doing nothing.")
			}

			return nil
		})
}

func TerraformPackage(tl *TaskList) *Task {
	return tl.CreateTask("package", P.Module.Name, P.Module.System).
		Set(func(t *Task) error {
			for _, tag := range C.Tags {
				t.CreateSubtask(tag).
					Set(func(t *Task) error {
						output := fmt.Sprintf("%s/%s-%s-%s.tar.gz", TF_MODULE_OUTPUT_DIR, P.Module.Name, P.Module.System, tag)

						t.CreateCommand(
							"tar",
							"-vczf",
							output,
							"--exclude=./.git",
							".",
						).
							SetDir(P.Module.Cwd).
							SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							ShouldRunBefore(func(c *Command) error {
								c.Log.Infof("Creating package for tag: %s", tag)

								return nil
							}).
							ShouldRunAfter(func(c *Command) error {
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

func TerraformPublish(tl *TaskList) *Task {
	return tl.CreateTask("publish").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				TerraformPublishGitlab(tl).Job(),
			)
		})
}

func TerraformPublishGitlab(tl *TaskList) *Task {
	return tl.CreateTask("publish", TF_REGISTRY_GITLAB, P.Module.Name, P.Module.System).
		ShouldDisable(func(t *Task) bool {
			return P.Registry.Name != TF_REGISTRY_GITLAB
		}).
		Set(func(t *Task) error {
			for _, p := range C.Packages {
				t.CreateSubtask(p.Tag).
					Set(func(t *Task) error {
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

						t.Log.Infof("Package has been published: %s@%s", P.Module.Name, p.Tag)

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
