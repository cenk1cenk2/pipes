package gitlab

import (
	"context"
	"fmt"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

type MergeRequestConfig struct {
	Token        string `validate:"required"`
	ApiUrl       string `validate:"required"`
	ProjectId    string `validate:"required"`
	SourceBranch string `validate:"required"`
	TargetBranch string `validate:"required"`
	Title        string `validate:"required"`
}

type MergeRequestResult struct {
	MergeRequestIid int64
	WebUrl          string
	Created         bool
}

// Reads as the past tense of what happened to the merge request, for the job log.
func (r MergeRequestResult) Action() string {
	if r.Created {
		return "created"
	}

	return "updated"
}

// GitLab allows a single open merge request per source and target branch pair, so
// the pair alone identifies the merge request to update. A merged or closed one is
// never listed and gets a fresh merge request instead of an update.
func UpsertMergeRequest(
	ctx context.Context,
	mergeRequests MergeRequestsAdapter,
	config MergeRequestConfig,
	description string,
) (*MergeRequestResult, error) {
	existing, _, err := mergeRequests.ListProjectMergeRequests(
		config.ProjectId,
		&clientgitlab.ListProjectMergeRequestsOptions{
			State:        new("opened"),
			SourceBranch: new(config.SourceBranch),
			TargetBranch: new(config.TargetBranch),
		},
		clientgitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("list GitLab merge requests: %w", err)
	}

	if len(existing) > 0 {
		updated, _, err := mergeRequests.UpdateMergeRequest(
			config.ProjectId,
			existing[0].IID,
			&clientgitlab.UpdateMergeRequestOptions{
				Description: new(description),
			},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("update GitLab merge request: %w", err)
		}

		return &MergeRequestResult{
			MergeRequestIid: updated.IID,
			WebUrl:          updated.WebURL,
		}, nil
	}

	created, _, err := mergeRequests.CreateMergeRequest(
		config.ProjectId,
		&clientgitlab.CreateMergeRequestOptions{
			Title:        new(config.Title),
			Description:  new(description),
			SourceBranch: new(config.SourceBranch),
			TargetBranch: new(config.TargetBranch),
		},
		clientgitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab merge request: %w", err)
	}

	return &MergeRequestResult{
		MergeRequestIid: created.IID,
		WebUrl:          created.WebURL,
		Created:         true,
	}, nil
}
