package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/markdown"
)

// Everything a plan report needs that differs between the tools. Read reaches for the
// plan -- a command for one pipe, a file for the other -- so the tasks below keep one
// shape. The tasks take the pipe's own instance, which is only filled once its flags
// have been parsed.
type Source struct {
	Read          func(ctx context.Context, t *Task) (Report, error)
	Summary       func(Report) Summary
	SummaryOutput string
	Cwd           string
	MergeRequest  gitlab.MergeRequestReportConfig
	Notes         gitlab.NotesFactory
	// names what this job reports on, so concurrent report jobs on one merge
	// request do not overwrite each other's note.
	Discriminators func() []string
	Metadata       Metadata
	Log            LogConfig

	// held between the tasks below, since reaching for the plan runs a command for one
	// of the pipes and they run in sequence anyway.
	cache *Report
}

func (s *Source) read(ctx context.Context, t *Task) (Report, error) {
	if s.cache != nil {
		return *s.cache, nil
	}

	report, err := s.Read(ctx, t)
	if err != nil {
		return Report{}, err
	}

	s.cache = &report

	return report, nil
}

// LogTask writes the report into the job log, which is what a pipeline is read with
// long before anyone opens the merge request the note goes on.
func LogTask(tl *TaskList, src *Source) *Task {
	return tl.CreateTask("report").
		ShouldDisable(func(t *Task) bool {
			if !src.Log.Enabled {
				t.Log.Debug("Skipping the plan report in the job log because it is disabled.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			report, err := src.read(ctx, t)
			if err != nil {
				return err
			}

			body, err := RenderReport(report)
			if err != nil {
				return err
			}

			return markdown.Log(t.Log, body)
		})
}

func SummaryTask(tl *TaskList, src *Source) *Task {
	return tl.CreateTask("summary").
		ShouldDisable(func(t *Task) bool {
			if src.SummaryOutput == "" {
				t.Log.Debug("Skipping plan summary because no summary output file is configured.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			report, err := src.read(ctx, t)
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

			t.Log.Info(fmt.Sprintf("Wrote plan summary: %s", output))

			return nil
		})
}

func MergeRequestReportTask(tl *TaskList, src *Source) *Task {
	return tl.CreateTask("merge-request-report").
		ShouldDisable(func(t *Task) bool {
			if !src.MergeRequest.Enabled {
				return true
			}

			if src.MergeRequest.MergeRequestIid == 0 {
				t.Log.Debug("Skipping GitLab merge request report because this is not a merge request pipeline.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			report, err := src.read(ctx, t)
			if err != nil {
				return err
			}

			body, err := RenderReport(report)
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

			t.Log.Info(fmt.Sprintf("Merge request report note %s: %d (identifier: %s)",
				result.Action(),
				result.NoteId,
				result.Identifier),
			)

			return nil
		})
}
