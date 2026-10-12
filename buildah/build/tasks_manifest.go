package build

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	json "encoding/json/v2"
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	buildahmanifest "gitlab.kilic.dev/devops/pipes/buildah/manifest"
)

// The images are handed over only once they are pushed, since GitLab uploads a
// report even from a job that fails.
func manifest(tl *TaskList) *Task {
	return tl.CreateTask("manifest").
		ShouldDisable(func(t *Task) bool {
			return !P.Image.Push || P.Manifest.Target == ""
		}).
		Set(func(ctx context.Context, t *Task) error {
			target, err := InlineTemplate(P.Manifest.Target, C.Tags)
			if err != nil {
				return err
			}

			image, err := tagsCollector().Process(ctx, t, target)
			if err != nil {
				return err
			}

			value, err := json.Marshal([]buildahmanifest.ManifestMatrix{{Target: image, Images: C.Tags}})
			if err != nil {
				return err
			}

			key := manifestKey()

			if err := os.WriteFile(P.Manifest.EnvFile, fmt.Appendf(nil, "%s=%s\n", key, value), 0o600); err != nil {
				return err
			}

			t.Log.Info(fmt.Sprintf("Published images written to dotenv file: %s > %s", P.Manifest.EnvFile, key))

			return nil
		})
}

// Every build job of a pipeline hands its images over under a key of its own, since
// GitLab picks an undefined winner among the jobs that pass the same one. The hash
// of the tags tells apart the jobs whose slugs GitLab cut to the same length.
func manifestKey() string {
	sum := sha256.Sum256([]byte(strings.Join(C.Tags, ",")))
	hash := strings.ToUpper(hex.EncodeToString(sum[:])[:8])

	if P.Manifest.Job == "" {
		return P.Manifest.PublishedPrefix + hash
	}

	job := regexp.MustCompile(`[^A-Z0-9]`).ReplaceAllString(strings.ToUpper(P.Manifest.Job), "_")

	return fmt.Sprintf("%s%s_%s", P.Manifest.PublishedPrefix, job, hash)
}
