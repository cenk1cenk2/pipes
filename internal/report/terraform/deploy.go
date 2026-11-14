package terraform

import (
	"context"
	"fmt"
	"strconv"

	. "github.com/cenk1cenk2/plumber/v7"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

// DeployReportTask leaves the plan of the default branch on the merge requests merged
// into its commit, with the job that deploys it, since their own report was planned
// against the state before every other merge and drift in the meantime.
func DeployReportTask(tl *TaskList, src *Source) *Task {
	return tl.CreateTask("deploy-report").
		ShouldDisable(func(t *Task) bool {
			if !src.deploying() {
				t.Log.Debug("Skipping the deploy report, since this is not a pipeline of the default branch.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			report, err := src.read(ctx, t)
			if err != nil {
				return err
			}

			mergeRequests, err := src.mergeRequests()
			if err != nil {
				return err
			}

			deploy, err := src.pipelineJob(ctx, t, mergeRequests, src.MergeRequest.DeployJob)
			if err != nil {
				return err
			}

			body, err := RenderReport(report, DiffUnified)
			if err != nil {
				return err
			}

			lead := fmt.Sprintf(
				"Merged into `%s`. %s planned the changes below%s, which %s applies once it is run.",
				src.MergeRequest.DefaultBranch,
				src.job("plan"),
				src.target(),
				src.other(deploy, "the deploy job"),
			)

			return src.announce(ctx, t, mergeRequests, lead+"\n\n"+body, ":deploy-planned")
		})
}

// Apply runs the apply of the plan the default branch saved, and then tells the merge
// requests merged into the commit whether it went through. A deploy is never failed by
// its report.
func (s *Source) Apply(ctx context.Context, t *Task, apply func(context.Context) error) error {
	applied := apply(ctx)

	if !s.deploying() {
		return applied
	}

	err := func() error {
		mergeRequests, err := s.mergeRequests()
		if err != nil {
			return err
		}

		plan, err := s.pipelineJob(ctx, t, mergeRequests, s.MergeRequest.PlanJob)
		if err != nil {
			return err
		}

		body := fmt.Sprintf(
			"Deployed what was merged into `%s`. %s applied the plan of %s%s.\n",
			s.MergeRequest.DefaultBranch,
			s.job("deploy"),
			s.other(plan, "the plan job"),
			s.target(),
		)
		if applied != nil {
			body = fmt.Sprintf(
				"Deploying what was merged into `%s` failed. %s failed while applying the plan of %s%s, so it may be applied only in part.\n",
				s.MergeRequest.DefaultBranch,
				s.job("deploy"),
				s.other(plan, "the plan job"),
				s.target(),
			)
		}

		// a retried deploy replaces the outcome of the one before it.
		return s.announce(ctx, t, mergeRequests, body, ":deploy")
	}()
	if err != nil {
		t.Log.Warn(fmt.Sprintf("Can not report the deploy on the merged merge requests: %s", err))
	}

	return applied
}

func (s *Source) deploying() bool {
	config := s.MergeRequest

	return config.Enabled &&
		config.MergeRequestIid == 0 &&
		config.Branch != "" &&
		config.Branch == config.DefaultBranch &&
		config.ProjectId != "" &&
		s.Metadata.CommitSha != ""
}

func (s *Source) mergeRequests() (gitlab.MergeRequestsAdapter, error) {
	return s.MergeRequests(gitlab.MergeRequestConfig{Token: s.MergeRequest.Token, ApiUrl: s.MergeRequest.ApiUrl})
}

// The other job of the pipeline the note links to, found by its name, since the plan
// and the deploy do not know of each other otherwise. Nil leaves the note pointing at
// the pipeline instead.
func (s *Source) pipelineJob(ctx context.Context, t *Task, mergeRequests gitlab.MergeRequestsAdapter, name string) (*clientgitlab.Job, error) {
	if name == "" || s.Metadata.PipelineId == "" {
		return nil, nil
	}

	pipeline, err := strconv.ParseInt(s.Metadata.PipelineId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Can not parse the pipeline id: %s -> %w", s.Metadata.PipelineId, err)
	}

	for page := int64(1); page != 0; {
		jobs, response, err := mergeRequests.ListPipelineJobs(
			s.MergeRequest.ProjectId,
			pipeline,
			&clientgitlab.ListJobsOptions{ListOptions: clientgitlab.ListOptions{Page: page, PerPage: 100}},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("list GitLab jobs of the pipeline %d: %w", pipeline, err)
		}

		for _, job := range jobs {
			if job.Name == name {
				return job, nil
			}
		}

		page = 0
		if response != nil {
			page = response.NextPage
		}
	}

	t.Log.Warn(fmt.Sprintf("Linking the pipeline instead, since it has no job named: %s", name))

	return nil, nil
}

func (s *Source) announce(ctx context.Context, t *Task, mergeRequests gitlab.MergeRequestsAdapter, body string, suffix string) error {
	config := s.MergeRequest

	merged, err := gitlab.MergedMergeRequests(ctx, mergeRequests, config.ProjectId, s.Metadata.CommitSha, config.DefaultBranch)
	if err != nil {
		return err
	}

	if len(merged) == 0 {
		t.Log.Debug("Skipping the deploy report, since no merge request was merged into the commit.")

		return nil
	}

	for _, mergeRequest := range merged {
		note := config
		note.MergeRequestIid = mergeRequest.IID
		note.Identifier = gitlab.ResolveReportIdentifier(config.Identifier, s.Metadata.JobName, s.Discriminators()...) + suffix

		notes, err := s.Notes(note)
		if err != nil {
			return err
		}

		result, err := gitlab.UpsertMergeRequestReport(ctx, notes, note, body)
		if err != nil {
			return err
		}

		t.Log.Info(fmt.Sprintf("Deploy report note %s: %s!%d -> %d", result.Action(), note.ProjectId, mergeRequest.IID, result.NoteId))
	}

	return nil
}

func (s *Source) job(fallback string) string {
	name := "The " + fallback + " job"
	if s.Metadata.JobName != "" {
		name = "Job " + code(s.Metadata.JobName)
	}

	return link(name, s.Metadata.JobUrl)
}

// The job the note links next to its own, or the pipeline when it was not found.
func (s *Source) other(job *clientgitlab.Job, kind string) string {
	if job == nil {
		return kind + " of " + s.pipeline()
	}

	return link(kind+" "+code(job.Name), job.WebURL)
}

// The stack or state the plan is for, which tells the notes of several apart.
func (s *Source) target() string {
	if s.Metadata.Target == "" {
		return ""
	}

	return " for " + code(s.Metadata.Target)
}

func (s *Source) pipeline() string {
	name := "the pipeline"
	if s.Metadata.PipelineId != "" {
		name = "pipeline #" + s.Metadata.PipelineId
	}

	return link(name, s.Metadata.PipelineUrl)
}

func code(value string) string {
	return "`" + value + "`"
}

func link(text string, url string) string {
	if url == "" {
		return text
	}

	return fmt.Sprintf("[%s](%s)", text, url)
}
