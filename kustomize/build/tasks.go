package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

func RenderOverlays(tl *TaskList) *Task {
	return tl.CreateTask("build").
		ShouldDisable(func(t *Task) bool {
			return len(setup.C.Overlays) == 0
		}).
		Set(func(t *Task) error {
			if P.Output != "" {
				if err := os.MkdirAll(P.Output, 0o755); err != nil {
					return fmt.Errorf("Can not create Kustomize output directory: %s: %w", P.Output, err)
				}
			}

			for _, overlay := range setup.C.Overlays {
				t.CreateSubtask(overlay).
					Set(func(t *Task) error {
						result := renderOverlay(overlay, P)

						t.Lock.Lock()
						C.Results = append(C.Results, result)
						t.Lock.Unlock()

						if result.Err != nil {
							if P.FailFast {
								return fmt.Errorf("Can not build Kustomize overlay: %s: %w", overlay, result.Err)
							}

							t.Log.Warnf("Failed to build Kustomize overlay: %s: %v", overlay, result.Err)

							return nil
						}

						t.Log.Infof("Built Kustomize overlay: %s (%d resources)", overlay, result.DocCount)

						if P.Output != "" {
							if err := writeOverlay(overlay, result); err != nil {
								return err
							}
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

func writeOverlay(overlay string, result OverlayResult) error {
	name := strings.ReplaceAll(filepath.Clean(overlay), string(os.PathSeparator), "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = "root"
	}

	output := filepath.Join(P.Output, fmt.Sprintf("%s.yaml", name))
	if err := os.WriteFile(output, result.Yaml, 0o644); err != nil {
		return fmt.Errorf("Can not write rendered Kustomize overlay: %s: %w", output, err)
	}

	return nil
}

func MergeRequestReport(tl *TaskList) *Task {
	return tl.CreateTask("merge-request-report").
		ShouldDisable(func(t *Task) bool {
			if !P.MergeRequestReport.Enabled {
				return true
			}

			if P.MergeRequestReport.MergeRequestId == 0 {
				t.Log.Debugln("Skipping GitLab merge request report because this is not a merge request pipeline.")

				return true
			}

			return false
		}).
		Set(func(t *Task) error {
			report := newMergeRequestReport(C.Results)

			for _, item := range []mergeRequestReportMetadata{
				{Name: "Kustomize working directory", Value: setup.P.Cwd},
				{Name: "Overlays built", Value: fmt.Sprintf("%d", report.Summary.Total)},
				{Name: "GitLab project id", Value: P.MergeRequestReport.ProjectId},
				{Name: "GitLab merge request", Value: fmt.Sprintf("!%d", P.MergeRequestReport.MergeRequestId)},
				{Name: "Report identifier", Value: P.MergeRequestReport.Identifier},
			} {
				if item.Value != "" {
					report.Metadata = append(report.Metadata, item)
				}
			}

			body, err := renderMergeRequestReport(report)
			if err != nil {
				return err
			}

			result, err := gitlab.UpsertMergeRequestReport(
				context.Background(),
				P.MergeRequestReport,
				body,
			)
			if err != nil {
				return err
			}

			t.Log.Infof("Upserted GitLab merge request report note: %d", result.NoteId)

			return nil
		})
}
