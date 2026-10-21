package plugins

import (
	"context"
	json "encoding/json/v2"
	"fmt"
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
)

type (
	Plugin struct {
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)

// pulumi downloads the plugin versions a program needs into the plugin cache and never
// removes the old ones, so every provider bump would leave one more directory in the
// cached home forever; only the newest version of each plugin is kept.
func Prune(tl *TaskList) *Task {
	var list string

	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return home() == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			list = ""

			t.CreateCommand("pulumi", "plugin", "ls", "--output", "json").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				CaptureStdout(&list).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			// the cache is only housekeeping, so it never fails a run that otherwise worked.
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				t.Log.Warn(fmt.Sprintf("Skipping the plugin cache prune, since the plugins can not be listed: %s", err))

				return nil
			}

			outdated, err := Outdated([]byte(list))
			if err != nil {
				t.Log.Warn(fmt.Sprintf("Skipping the plugin cache prune, since the plugins can not be read: %s", err))

				return nil
			}

			return t.RunCommandJob(ctx, func(t *Task) Job {
				jobs := []Job{}

				for _, plugin := range outdated {
					jobs = append(
						jobs,
						t.CreateCommand("pulumi", "plugin", "rm", plugin.Kind, plugin.Name, plugin.Version, "--yes").
							SetDir(setup.C.Cwd).
							AppendEnvironment(setup.C.Env).
							SetIgnoreError().
							Job(),
					)
				}

				return JobSequence(jobs...)
			})
		})
}

func home() string {
	if home := setup.C.Env["PULUMI_HOME"]; home != "" {
		return home
	}

	return os.Getenv("PULUMI_HOME")
}

// everything but the newest version of each kind and name is outdated, and a version
// that is no semver can not be compared, so it stays where it is.
func Outdated(data []byte) ([]Plugin, error) {
	var plugins []Plugin
	if err := json.Unmarshal(data, &plugins); err != nil {
		return nil, fmt.Errorf("parse Pulumi plugins: %w", err)
	}

	type (
		key struct {
			kind string
			name string
		}
		parsed struct {
			plugin  Plugin
			version *semver.Version
		}
	)

	groups := map[key][]parsed{}
	order := []key{}

	for _, plugin := range plugins {
		version, err := semver.NewVersion(plugin.Version)
		if err != nil {
			continue
		}

		k := key{kind: plugin.Kind, name: plugin.Name}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}

		groups[k] = append(groups[k], parsed{plugin: plugin, version: version})
	}

	outdated := []Plugin{}

	for _, k := range order {
		group := groups[k]

		sort.SliceStable(group, func(i, j int) bool {
			return group[i].version.GreaterThan(group[j].version)
		})

		for _, old := range group[1:] {
			// the same version listed twice is not older than itself.
			if old.version.Equal(group[0].version) {
				continue
			}

			outdated = append(outdated, old.plugin)
		}
	}

	return outdated, nil
}
