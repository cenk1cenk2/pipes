package publish

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/parser"
)

func TerraformTagsFile(tl *TaskList) *Task {
	return tl.CreateTask("tags").
		Set(func(t *Task) error {
			tags, err := parser.ParseTagsFile(t.Log, path.Join(P.Module.Cwd, P.Module.TagsFile), false)

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
						url := fmt.Sprintf(
							"%s/projects/%s/packages/terraform/modules/%s/%s/%s/file",
							P.Registry.Gitlab.ApiUrl,
							P.Registry.Gitlab.ProjectId,
							P.Module.Name,
							P.Module.System,
							p.Tag,
						)

						file, err := os.Open(p.Output)
						if err != nil {
							return err
						}

						defer file.Close()

						req, err := http.NewRequest(http.MethodPut, url, file)

						if err != nil {
							return err
						}

						req.Header.Set("Content-Type", "application/tar+gzip")
						req.Header.Set("JOB-TOKEN", P.Registry.Gitlab.Token)

						client := &http.Client{}

						res, err := client.Do(req)

						if err != nil {
							return err
						}

						defer res.Body.Close()

						body, err := io.ReadAll(res.Body)

						if err != nil {
							return err
						}

						if res.StatusCode == http.StatusCreated {
							t.Log.Infof("Package has been published: %s@%s", P.Module.Name, p.Tag)

							t.Log.Debugln(string(body))
						} else {
							t.Log.Warnln(string(body))
						}

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
