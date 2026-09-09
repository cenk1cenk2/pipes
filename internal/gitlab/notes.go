package gitlab

import (
	"fmt"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

// The three note calls the report upsert makes, narrowed from the client so the
// bookkeeping can be driven without a GitLab to talk to. Signatures mirror
// clientgitlab.NotesService exactly.
type NotesAdapter interface {
	ListMergeRequestNotes(
		pid any,
		mergeRequest int64,
		opt *clientgitlab.ListMergeRequestNotesOptions,
		options ...clientgitlab.RequestOptionFunc,
	) ([]*clientgitlab.Note, *clientgitlab.Response, error)
	CreateMergeRequestNote(
		pid any,
		mergeRequest int64,
		opt *clientgitlab.CreateMergeRequestNoteOptions,
		options ...clientgitlab.RequestOptionFunc,
	) (*clientgitlab.Note, *clientgitlab.Response, error)
	UpdateMergeRequestNote(
		pid any,
		mergeRequest, note int64,
		opt *clientgitlab.UpdateMergeRequestNoteOptions,
		options ...clientgitlab.RequestOptionFunc,
	) (*clientgitlab.Note, *clientgitlab.Response, error)
}

var _ NotesAdapter = (*clientgitlab.NotesService)(nil)

// Dials only when a report is actually going to be written, so a pipe that never
// reaches the report task never needs a token that parses.
type NotesFactory func(config MergeRequestReportConfig) (NotesAdapter, error)

func NewNotes(config MergeRequestReportConfig) (NotesAdapter, error) {
	client, err := clientgitlab.NewClient(
		config.Token,
		clientgitlab.WithBaseURL(config.ApiUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	return client.Notes, nil
}
