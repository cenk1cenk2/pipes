package terraform

import (
	"cmp"
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"text/template"

	. "github.com/cenk1cenk2/plumber/v7"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
)

//go:embed assets/deploy.md.gotmpl
var deployTemplate string

// What the deploy table on top of either note renders, with the job cells already
// written as links.
type deployView struct {
	Title     string
	Status    string
	Branch    string
	Labels    Labels
	Target    string
	PlanJob   string
	DeployJob string
	Pipeline  string
	Commit    string
}

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

			table, err := src.deployTable("ready", "Waiting for the deploy job to run", src.job(), src.other(deploy))
			if err != nil {
				return err
			}

			body, err := RenderReport(report, DiffUnified)
			if err != nil {
				return err
			}

			return src.announce(ctx, t, mergeRequests, table+"\n"+body, ":deploy-planned")
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

		outcome, status := "succeeded", "Applied"
		if applied != nil {
			outcome, status = "failed", "Failed, so the plan may be applied only in part"
		}

		body, err := s.deployTable(outcome, status, s.other(plan), s.job())
		if err != nil {
			return err
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

func (s *Source) deployTable(outcome string, status string, planJob string, deployJob string) (string, error) {
	view := deployView{
		Title:     fmt.Sprintf("%s %s on %s%s", s.Labels.Apply, outcome, code(s.MergeRequest.DefaultBranch), s.target()),
		Status:    status,
		Branch:    s.MergeRequest.DefaultBranch,
		Labels:    s.Labels,
		Target:    s.Metadata.Target,
		PlanJob:   planJob,
		DeployJob: deployJob,
		Commit:    cmp.Or(s.Metadata.CommitShortSha, s.Metadata.CommitSha),
	}
	if s.Metadata.PipelineId != "" || s.Metadata.PipelineUrl != "" {
		view.Pipeline = s.pipeline()
	}

	return InlineTemplate(deployTemplate, view, template.FuncMap{"cell": cell})
}

func (s *Source) job() string {
	name := "This job"
	if s.Metadata.JobName != "" {
		name = code(s.Metadata.JobName)
	}

	return link(name, s.Metadata.JobUrl)
}

// The other job of the pipeline the note links to, or the pipeline when it was not found.
func (s *Source) other(job *clientgitlab.Job) string {
	if job == nil {
		return "In " + s.pipeline()
	}

	return link(code(job.Name), job.WebURL)
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
