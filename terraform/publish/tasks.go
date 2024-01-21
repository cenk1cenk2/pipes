package publish

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gitlab.kilic.dev/devops/pipes/common/parser"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func TerraformTagsFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tags").
		Set(func(t *Task[Pipe]) error {
			tags, err := parser.ParseTagsFile(t.Log, t.Pipe.Module.TagsFile, true)

			if err != nil {
				return err
			}

			t.Pipe.Ctx.Tags = tags

			return nil
		})
}

func TerraformPackage(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("package").
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateSubtask(tag).
						Set(func(t *Task[Pipe]) error {
							output := fmt.Sprintf("%s/%s-%s-%s.tar.gz", TF_MODULE_OUTPUT_DIR, t.Pipe.Module.Name, t.Pipe.Module.System, tag)

							t.CreateCommand(
								"tar",
								"-vczf",
								output,
								"--exclude=./.git",
								".",
							).
								SetDir(t.Pipe.Module.Cwd).
								SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
								ShouldRunBefore(func(c *Command[Pipe]) error {
									c.Log.Infof("Creating package for tag: %s", tag)

									return nil
								}).
								ShouldRunAfter(func(c *Command[Pipe]) error {
									t.Lock.Lock()
									t.Pipe.Ctx.Packages = append(t.Pipe.Ctx.Packages, PublishablePackage{
										Tag:    tag,
										Output: output,
									})
									t.Lock.Unlock()

									return nil
								}).
								AddSelfToTheTask()

							return nil
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobSequence()
						}).
						AddSelfToTheParentAsParallel()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func TerraformPublish(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("publish").
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobParallel(
				TerraformPublishGitlab(tl).Job(),
			)
		})
}

func TerraformPublishGitlab(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask(TF_REGISTRY_GITLAB, tl.Pipe.Module.Name, tl.Pipe.Module.System).
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Registry.Name != TF_REGISTRY_GITLAB
		}).
		Set(func(t *Task[Pipe]) error {
			for _, p := range t.Pipe.Ctx.Packages {
				func(p PublishablePackage) {
					t.CreateSubtask(p.Tag).
						Set(func(t *Task[Pipe]) error {
							url := fmt.Sprintf(
								"%s/projects/%s/packages/terraform/modules/%s/%s/%s/file",
								t.Pipe.Registry.Gitlab.ApiUrl,
								t.Pipe.Registry.Gitlab.ProjectId,
								t.Pipe.Module.Name,
								t.Pipe.Module.System,
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

							if res.StatusCode == http.StatusOK {
								t.Log.Debugln(string(body))
							} else {
								t.Log.Warnln(string(body))
							}

							return nil
						}).
						AddSelfToTheParentAsParallel()
				}(p)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
