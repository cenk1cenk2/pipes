package sync

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/markdown"
)

// A job writing into the branch its own pipeline runs on would trigger itself again
// with every commit it publishes, so it stops before touching anything instead.
func selfTarget(tl *TaskList) Predicate {
	return func() bool {
		if P.Target.Project != P.Project.Path || P.Target.Branch != P.Refs.Branch {
			return false
		}

		tl.Log.Info(fmt.Sprintf("Skipping the sync, since the pipeline runs on the target branch itself: %s@%s", P.Target.Project, P.Target.Branch))

		return true
	}
}

// The pipeline already has the project checked out, so the target branch of the same
// project only costs a shallow fetch next to it, while another project is cloned with
// nothing but the destinations.
func acquire(tl *TaskList) *Task {
	return tl.CreateTask("acquire").
		Set(func(ctx context.Context, t *Task) error {
			dir, err := os.MkdirTemp("", "git-sync-")
			if err != nil {
				return fmt.Errorf("Can not create a directory for the target: %w", err)
			}

			if P.Target.Project == P.Project.Path {
				t.Log.Info(fmt.Sprintf("Checking the target branch out next to the checkout of the pipeline: %s", P.Target.Branch))

				C.Repository, err = git.Repository{Dir: P.Project.Dir, Token: P.Token}.Worktree(ctx, t, C.Remote, P.Target.Branch, dir)

				return err
			}

			t.Log.Info(fmt.Sprintf("Cloning the target: %s@%s", P.Target.Project, P.Target.Branch))

			C.Repository = git.Repository{Dir: filepath.Join(dir, "target"), Token: P.Token}

			return C.Repository.Clone(ctx, t, C.Remote, P.Target.Branch, C.Destinations)
		})
}

func synchronize(tl *TaskList) *Task {
	return tl.CreateTask("sync").
		Set(func(_ context.Context, t *Task) error {
			return git.Sync(t, C.Repository.Dir, P.Paths, P.AllowEmpty)
		})
}

func stage(tl *TaskList) *Task {
	return tl.CreateTask("stage").
		Set(func(ctx context.Context, t *Task) error {
			if err := C.Repository.Stage(ctx, t, C.Destinations); err != nil {
				return err
			}

			if err := C.Repository.Patch(ctx, t, P.Patch); err != nil {
				return err
			}

			diff, err := os.ReadFile(P.Patch)
			if err != nil {
				return fmt.Errorf("Can not read the staged changes: %s -> %w", P.Patch, err)
			}

			C.Diff = string(diff)

			if C.Diff == "" {
				t.Log.Info("Nothing changed in the target.")
			}

			return nil
		})
}

func log(tl *TaskList) *Task {
	return tl.CreateTask("log").
		Set(func(_ context.Context, t *Task) error {
			return markdown.Log(t.Log, render())
		})
}

func report(tl *TaskList) *Task {
	return tl.CreateTask("report").
		ShouldDisable(func(_ *Task) bool {
			return !P.MergeRequestReport.Enabled
		}).
		Set(func(ctx context.Context, t *Task) error {
			targets, err := reportTargets(ctx)
			if err != nil {
				return err
			}

			if len(targets) == 0 {
				t.Log.Debug("Skipping the merge request report, since no merge request is open for the pipeline.")

				return nil
			}

			body := render()

			for _, target := range targets {
				config := P.MergeRequestReport
				config.ProjectId = target.ProjectId
				config.MergeRequestIid = target.MergeRequestIid
				config.Identifier = gitlab.ResolveReportIdentifier(config.Identifier, P.ReportMetadata.JobName, P.Target.Project, P.Target.Branch)

				notes, err := newNotes(config)
				if err != nil {
					return err
				}

				result, err := gitlab.UpsertMergeRequestReport(ctx, notes, config, body)
				if err != nil {
					return err
				}

				t.Log.Info(fmt.Sprintf("Merge request report note %s: %s!%d -> %d", result.Action(), target.ProjectId, target.MergeRequestIid, result.NoteId))
			}

			return nil
		})
}

func publish(tl *TaskList) *Task {
	return tl.CreateTask("publish").
		ShouldDisable(func(t *Task) bool {
			if P.Mode != ModePublish {
				return true
			}

			if C.Diff == "" {
				t.Log.Info("Skipping the publish, since there is nothing to commit.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			head, err := git.Repository{Dir: P.Project.Dir, Token: P.Token}.Lease(ctx, t, remote(P.Project.Path), P.Project.DefaultBranch)
			if err != nil {
				return err
			}

			// the default branch always exists, so an empty head is a lookup that failed and
			// not a branch that moved.
			if head == "" {
				return fmt.Errorf("Can not find the head of the default branch: %s", P.Project.DefaultBranch)
			}

			// a newer pipeline on the default branch publishes the newer result, and this
			// one would only race it onto the branch with an older one.
			if head != P.ReportMetadata.CommitSha {
				return fmt.Errorf(
					"Commit of the pipeline is no longer the head of %s: %s -> %s; a newer pipeline publishes the newer result",
					P.Project.DefaultBranch,
					P.ReportMetadata.CommitSha,
					head,
				)
			}

			if err := commit(t).Run(ctx); err != nil {
				return fmt.Errorf("Can not commit the changes: %w", err)
			}

			lease, err := C.Repository.Lease(ctx, t, C.Remote, P.Branch)
			if err != nil {
				return err
			}

			if err := C.Repository.Push(ctx, t, C.Remote, P.Branch, lease); err != nil {
				return err
			}

			result, err := upsertMergeRequest(ctx, t)
			if err != nil {
				return err
			}

			t.Log.Info(fmt.Sprintf("Merge request %s: %s", result.Action(), result.WebUrl))

			return nil
		})
}

// The job carries an identity of its own in the environment, which git would take
// over anything passed to the commit, so the inherited one is dropped for this call.
func commit(t *Task) *Command {
	environment := []string{}
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "GIT_AUTHOR_") && !strings.HasPrefix(variable, "GIT_COMMITTER_") {
			environment = append(environment, variable)
		}
	}

	args := []string{"commit", "--quiet", "--message", P.CommitMessage}

	// the one who triggered the generation is credited without being the author, which
	// stays the identity that opens the merge request.
	if P.Project.CommitAuthor != "" {
		args = append(args, "--trailer", fmt.Sprintf("Co-authored-by: %s", P.Project.CommitAuthor))
	}

	return t.CreateCommand("git", args...).
		SetDir(C.Repository.Dir).
		SetMaskOsEnvironment().
		AppendDirectEnvironment(environment...).
		AppendEnvironment(map[string]string{
			"GIT_AUTHOR_NAME":     P.Author.Name,
			"GIT_AUTHOR_EMAIL":    P.Author.Email,
			"GIT_COMMITTER_NAME":  P.Committer.Name,
			"GIT_COMMITTER_EMAIL": P.Committer.Email,
		})
}

// Another run creating the merge request between the listing and the creation makes
// GitLab refuse the second one, which then only has to update what the other opened.
func upsertMergeRequest(ctx context.Context, t *Task) (*gitlab.MergeRequestResult, error) {
	subject, _, _ := strings.Cut(P.CommitMessage, "\n")

	config := gitlab.MergeRequestConfig{
		Token:        P.Token,
		ApiUrl:       P.MergeRequestReport.ApiUrl,
		ProjectId:    P.Target.Project,
		SourceBranch: P.Branch,
		TargetBranch: P.Target.Branch,
		Title:        strings.TrimSpace(subject),
		Assignees:    P.Assignees,
		Reviewers:    P.Reviewers,
	}

	mergeRequests, err := newMergeRequests(config)
	if err != nil {
		return nil, err
	}

	body := render()

	result, err := gitlab.UpsertMergeRequest(ctx, mergeRequests, config, body)
	if clientgitlab.HasStatusCode(err, http.StatusConflict) {
		t.Log.Warn("Merge request was opened by another run in the meantime, updating it instead.")

		return gitlab.UpsertMergeRequest(ctx, mergeRequests, config, body)
	}

	return result, err
}
