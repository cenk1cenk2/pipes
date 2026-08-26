package iac

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

// Everything a plan report needs that differs between the tools. Read is what
// actually reaches for the plan -- a command for one pipe, a file for the other --
// so the two tasks below stay the same shape whichever tool produced it.
type Source struct {
	Read          func(t *plumber.Task) (Report, error)
	Summary       func(Report) Summary
	SummaryOutput string
	Cwd           string
	MergeRequest  gitlab.MergeRequestReportConfig
	Notes         gitlab.NotesFactory
	// Names what this job reports on, so concurrent report jobs on one merge
	// request do not overwrite each other's note.
	Discriminators func() []string
	Metadata       Metadata
}

func SummaryTask(tl *plumber.TaskList, src Source) *plumber.Task {
	return tl.CreateTask("summary").
		ShouldDisable(func(t *plumber.Task) bool {
			if src.SummaryOutput == "" {
				t.Log.Debugln("Skipping plan summary because no summary output file is configured.")

				return true
			}

			return false
		}).
		Set(func(t *plumber.Task) error {
			report, err := src.Read(t)
			if err != nil {
				return err
			}

			output := src.SummaryOutput
			if !filepath.IsAbs(output) {
				output = filepath.Join(src.Cwd, output)
			}

			if err := os.WriteFile(output, []byte(RenderSummary(src.Summary(report))), 0o644); err != nil {
				return fmt.Errorf("write plan summary %s: %w", output, err)
			}

			t.Log.Infof("Wrote plan summary: %s", output)

			return nil
		})
}

func MergeRequestReportTask(tl *plumber.TaskList, src Source) *plumber.Task {
	return tl.CreateTask("merge-request-report").
		ShouldDisable(func(t *plumber.Task) bool {
			if !src.MergeRequest.Enabled {
				return true
			}

			if src.MergeRequest.MergeRequestIid == 0 {
				t.Log.Debugln("Skipping GitLab merge request report because this is not a merge request pipeline.")

				return true
			}

			return false
		}).
		Set(func(t *plumber.Task) error {
			report, err := src.Read(t)
			if err != nil {
				return err
			}

			body, err := RenderMergeRequestReport(report)
			if err != nil {
				return err
			}

			config := src.MergeRequest
			config.Identifier = gitlab.ResolveReportIdentifier(
				config.Identifier,
				src.Metadata.JobName,
				src.Discriminators()...,
			)
			config.LegacyIdentifiers = []string{src.Metadata.JobName}

			notes, err := src.Notes(config)
			if err != nil {
				return err
			}

			result, err := gitlab.UpsertMergeRequestReport(context.Background(), notes, config, body)
			if err != nil {
				return err
			}

			t.Log.Infof(
				"Merge request report note %s: %d (identifier: %s)",
				result.Action(),
				result.NoteId,
				result.Identifier,
			)

			return nil
		})
}
