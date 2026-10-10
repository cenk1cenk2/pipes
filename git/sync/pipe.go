package sync

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

const (
	ModePublish = "publish"
	ModeReport  = "report"
)

type (
	Target struct {
		Project string `validate:"required"`
		Branch  string `validate:"required"`
	}

	Identity struct {
		Name  string
		Email string
	}

	Project struct {
		Path              string
		PathSlug          string
		JobSlug           string
		Dir               string
		ServerUrl         string `validate:"required"`
		DefaultBranch     string
		OpenMergeRequests string
		CommitAuthor      string
	}

	Pipe struct {
		Mode               string `validate:"oneof=publish report"`
		Target             Target
		Paths              []git.Path `validate:"min=1"`
		Branch             string     `validate:"required"`
		CommitMessage      string     `validate:"required"`
		AllowEmpty         bool
		Patch              string `validate:"required"`
		Token              string
		Author             Identity
		Committer          Identity
		Project            Project
		Refs               git.Refs
		MergeRequestReport gitlab.MergeRequestReportConfig
		ReportMetadata     terraform.Metadata
	}

	Ctx struct {
		Remote       string
		Repository   git.Repository
		Destinations []string
		Diff         string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

// Dial only once a report or a merge request is actually written, and are swapped out
// by the specs, which have no GitLab to talk to.
var (
	newNotes         gitlab.NotesFactory         = gitlab.NewNotes
	newMergeRequests gitlab.MergeRequestsFactory = gitlab.NewMergeRequests
)

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if !P.MergeRequestReport.Enabled {
				P.MergeRequestReport.MergeRequestIid = 0
			}

			if P.Branch == "" && P.Project.PathSlug != "" && P.Project.JobSlug != "" {
				P.Branch = fmt.Sprintf("sync/%s-%s", P.Project.PathSlug, P.Project.JobSlug)
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			// a publish is what the head-moved guard and the merge request need these for,
			// and a report never pushes anything.
			if P.Mode == ModePublish {
				if P.Token == "" {
					return fmt.Errorf("Publishing needs a GitLab token for the target.")
				}

				if P.Project.Path == "" || P.Project.DefaultBranch == "" || P.ReportMetadata.CommitSha == "" {
					return fmt.Errorf("Publishing needs the project, its default branch and the commit of the pipeline to tell whether it is still the head.")
				}

				// the identity is never guessed, so a consumer that did not set it fails here
				// instead of publishing under a name nobody chose.
				for _, variable := range []struct {
					name  string
					value string
				}{
					{"GIT_PIPES_AUTHOR_NAME", P.Author.Name},
					{"GIT_PIPES_AUTHOR_EMAIL", P.Author.Email},
					{"GIT_PIPES_COMMITTER_NAME", P.Committer.Name},
					{"GIT_PIPES_COMMITTER_EMAIL", P.Committer.Email},
				} {
					if variable.value == "" {
						return fmt.Errorf("Publishing needs the identity of the commit: %s is not set.", variable.name)
					}
				}
			}

			C.Remote = remote(P.Target.Project)

			for _, path := range P.Paths {
				C.Destinations = append(C.Destinations, path.Destination)
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobIfNot(
				selfTarget(tl),
				JobSequence(
					acquire(tl).Job(),
					synchronize(tl).Job(),
					stage(tl).Job(),
					log(tl).Job(),
					report(tl).Job(),
					publish(tl).Job(),
				),
			)
		})
}

func remote(project string) string {
	return fmt.Sprintf("%s/%s.git", strings.TrimSuffix(P.Project.ServerUrl, "/"), project)
}
