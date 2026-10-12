package build

import (
	"context"
	"fmt"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

// The collector reads the parsed flags, so it is only built from inside a task list.
func tagsCollector() *versions.Collector {
	return &versions.Collector{
		Name: "tags",

		FromUser:        P.Image.Tags,
		FromEnvironment: P.Image.TagsFromEnvironment,

		LatestWhen:  P.Image.TagAsLatest,
		LatestValue: P.Image.LatestTag,
		References:  P.Git.References(),

		Templates: P.Image.TagsTemplate,
		Sanitize:  P.Image.TagsSanitize,

		Format: func(tag string) string {
			if login.P.Uri == "" {
				return fmt.Sprintf("%s:%s", P.Image.Name, tag)
			}

			return fmt.Sprintf("%s/%s:%s", login.P.Uri, P.Image.Name, tag)
		},
	}
}

func tags(tl *TaskList) *Task {
	collector := tagsCollector()

	return tl.CreateTask("tags").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				collector.UserTask(tl, &C.Tags).Job(),
				collector.EnvironmentTask(tl, &C.Tags).Job(),
				collector.LatestTask(tl, &C.Tags).Job(),
				job,
			)
		}).
		Set(func(_ context.Context, t *Task) error {
			C.Tags = slices.Compact(C.Tags)

			if len(C.Tags) == 0 {
				return fmt.Errorf("Can not build the image without any tags.")
			}

			t.Log.Info(fmt.Sprintf("Image tags: %s", strings.Join(C.Tags, ", ")))

			return nil
		})
}
