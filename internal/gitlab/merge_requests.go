package gitlab

import (
	"fmt"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

// The merge request calls the upsert and the close make, the note the close leaves
// and the user lookup that resolves assignees and reviewers, narrowed from the client
// so the bookkeeping can be driven without a GitLab to talk to. Signatures mirror
// clientgitlab.MergeRequestsService, clientgitlab.NotesService and
// clientgitlab.UsersService exactly.
type MergeRequestsAdapter interface {
	ListProjectMergeRequests(
		pid any,
		opt *clientgitlab.ListProjectMergeRequestsOptions,
		options ...clientgitlab.RequestOptionFunc,
	) ([]*clientgitlab.BasicMergeRequest, *clientgitlab.Response, error)
	CreateMergeRequest(
		pid any,
		opt *clientgitlab.CreateMergeRequestOptions,
		options ...clientgitlab.RequestOptionFunc,
	) (*clientgitlab.MergeRequest, *clientgitlab.Response, error)
	UpdateMergeRequest(
		pid any,
		mergeRequest int64,
		opt *clientgitlab.UpdateMergeRequestOptions,
		options ...clientgitlab.RequestOptionFunc,
	) (*clientgitlab.MergeRequest, *clientgitlab.Response, error)
	CreateMergeRequestNote(
		pid any,
		mergeRequest int64,
		opt *clientgitlab.CreateMergeRequestNoteOptions,
		options ...clientgitlab.RequestOptionFunc,
	) (*clientgitlab.Note, *clientgitlab.Response, error)
	ListUsers(
		opt *clientgitlab.ListUsersOptions,
		options ...clientgitlab.RequestOptionFunc,
	) ([]*clientgitlab.User, *clientgitlab.Response, error)
}

type mergeRequests struct {
	clientgitlab.MergeRequestsServiceInterface
	clientgitlab.NotesServiceInterface
	clientgitlab.UsersServiceInterface
}

var _ MergeRequestsAdapter = (*mergeRequests)(nil)

// Dials only when a merge request is actually going to be written, so a pipe that
// never reaches the publish task never needs a token that parses.
type MergeRequestsFactory func(config MergeRequestConfig) (MergeRequestsAdapter, error)

func NewMergeRequests(config MergeRequestConfig) (MergeRequestsAdapter, error) {
	client, err := clientgitlab.NewClient(
		config.Token,
		clientgitlab.WithBaseURL(config.ApiUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	return &mergeRequests{
		MergeRequestsServiceInterface: client.MergeRequests,
		NotesServiceInterface:         client.Notes,
		UsersServiceInterface:         client.Users,
	}, nil
}
