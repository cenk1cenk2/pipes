package build

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/markdown"
	"gitlab.kilic.dev/devops/pipes/internal/report/manifest"
	"gitlab.kilic.dev/devops/pipes/internal/report/patch"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
	"sigs.k8s.io/kustomize/api/konfig"
)

//go:embed assets/report.md.gotmpl
var reportTemplate string

// GitLab refuses a note past a million characters, and the rest of the report has to
// fit next to the diff.
var reportDiffLimit = 900_000

type (
	reportView struct {
		Base           string
		ReportMetadata terraform.Metadata
		Overlays       []overlayView
		Files          patch.Section
		Truncated      bool
	}

	overlayView struct {
		Name   string
		Counts manifest.Counts
		Err    error
	}
)

// Renders every overlay again on the base commit and reports what the change does to
// the rendered manifests, in the job log and as a note on the merge request, which
// every later push updates in place.
func report(tl *TaskList) *Task {
	return tl.CreateTask("report").
		ShouldDisable(func(t *Task) bool {
			if P.DiffBase == "" {
				t.Log.Debug("Skipping the rendered manifest report because there is no base commit to diff against, which only a merge request pipeline has.")

				return true
			}

			if !P.ReportLog.Enabled && P.MergeRequestReport.MergeRequestIid == 0 {
				t.Log.Debug("Skipping the rendered manifest report because both the job log and the merge request note are disabled.")

				return true
			}

			return len(C.Results) == 0
		}).
		Set(func(ctx context.Context, t *Task) error {
			repository := git.Repository{Dir: setup.C.Cwd}

			var top string
			if err := t.CreateCommand("git", "rev-parse", "--show-toplevel").
				SetDir(setup.C.Cwd).
				SetLogLevel(LogLevelDebug, LogLevelWarn, LogLevelDebug).
				CaptureOutput(&top).
				Run(ctx); err != nil {
				return fmt.Errorf("Can not find the root of the repository: %w", err)
			}
			top = strings.TrimSpace(top)

			dir, err := os.MkdirTemp("", "kustomize-base-")
			if err != nil {
				return fmt.Errorf("Can not create a directory for the base: %w", err)
			}
			defer os.RemoveAll(dir)

			worktree := filepath.Join(dir, "base")
			if _, err := repository.Detach(ctx, t, P.DiffBase, worktree); err != nil {
				return err
			}
			defer func() {
				if err := repository.Remove(context.WithoutCancel(ctx), t, worktree); err != nil {
					t.Log.Warn(fmt.Sprintf("Can not remove the worktree of the base: %s", err))
				}
			}()

			groups := []manifest.Group{}
			for _, result := range slices.SortedFunc(slices.Values(C.Results), func(a, b OverlayResult) int {
				return strings.Compare(a.Overlay, b.Overlay)
			}) {
				relative, base, err := baseOverlay(top, worktree, result.Overlay)
				if err != nil {
					return err
				}

				// named from the root of the repository, which reads the same in every job
				// whatever directory it builds from.
				group := manifest.Group{Name: overlayName(relative), Head: result.Resources}

				if base != "" {
					rendered := renderOverlay(base, P)
					group.Base, group.BaseErr = rendered.Resources, rendered.Err
				}

				t.Log.Debug(fmt.Sprintf("Rendered the base of the overlay: %s", result.Overlay))

				groups = append(groups, group)
			}

			diff, err := manifest.Diff(ctx, t, groups)
			if err != nil {
				return err
			}

			view := reportView{Base: P.DiffBase, ReportMetadata: P.ReportMetadata}
			if len(view.Base) > 8 {
				view.Base = view.Base[:8]
			}

			for _, group := range groups {
				view.Overlays = append(view.Overlays, overlayView{Name: group.Name, Counts: diff.Groups[group.Name], Err: group.BaseErr})
			}

			view.Files, _, view.Truncated = patch.Sections(patch.Split(diff.Patch), reportDiffLimit)

			body, err := InlineTemplate(reportTemplate, view, map[string]any{"link": link})
			if err != nil {
				return err
			}

			if P.ReportLog.Enabled {
				t.Log.Info("Generating the rendered manifest report.")

				if err := markdown.Log(t.Log, body); err != nil {
					return err
				}
			}

			if P.MergeRequestReport.MergeRequestIid == 0 {
				return nil
			}

			config := P.MergeRequestReport
			config.Identifier = gitlab.ResolveReportIdentifier(config.Identifier, P.ReportMetadata.JobName, setup.C.Cwd)

			notes, err := newNotes(config)
			if err != nil {
				return err
			}

			upserted, err := gitlab.UpsertMergeRequestReport(ctx, notes, config, body)
			if err != nil {
				return err
			}

			t.Log.Info(fmt.Sprintf("Merge request report note %s: %d (identifier: %s)", upserted.Action(), upserted.NoteId, upserted.Identifier))

			return nil
		})
}

// The overlay relative to the root of the repository, and the same place in the
// worktree of the base, or nothing when the base has no kustomization there, so
// everything the overlay renders is added.
func baseOverlay(top string, worktree string, overlay string) (string, string, error) {
	absolute, err := filepath.Abs(overlay)
	if err != nil {
		return "", "", err
	}

	// git names the root with its links resolved, which a runner building under a
	// linked directory would otherwise never match.
	if absolute, err = filepath.EvalSymlinks(absolute); err != nil {
		return "", "", err
	}

	relative, err := filepath.Rel(top, absolute)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", "", fmt.Errorf("The overlay is outside of the repository: %s", overlay)
	}

	base := filepath.Join(worktree, relative)
	for _, name := range konfig.RecognizedKustomizationFileNames() {
		if _, err := os.Stat(filepath.Join(base, name)); err == nil {
			return relative, base, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", "", err
		}
	}

	return relative, "", nil
}

func link(text string, url string) string {
	if url == "" {
		return text
	}

	return fmt.Sprintf("[%s](%s)", text, url)
}
