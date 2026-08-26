package gitlab

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func mergeRequestReportMarker(identifier string) string {
	return fmt.Sprintf("<!-- %s:%s -->", "gitlab-pipes:mr-report", identifier)
}

const (
	CATEGORY_GITLAB_MERGE_REQUEST_REPORT = "Gitlab Merge Request Report"
)

type MergeRequestReportConfig struct {
	Enabled        bool
	Token          string `validate:"required_with=MergeRequestId"`
	ApiUrl         string `validate:"required_with=MergeRequestId"`
	ProjectId      string `validate:"required_with=MergeRequestId"`
	MergeRequestId int64  `validate:"omitempty,gt=0"`
	Identifier     string `validate:"omitempty,printascii,excludes=-->"`
	// Markers of earlier identifier schemes, so a note already posted under one of
	// them is adopted instead of orphaned next to a duplicate.
	LegacyIdentifiers []string
}

type MergeRequestReportResult struct {
	NoteId     int64
	Identifier string
	Created    bool
	// Set when the note was found under an earlier identifier scheme rather than the
	// current one, which is the only signal that a marker migration took effect.
	AdoptedLegacy bool
}

// Reads as the past tense of what happened to the note, for the job log.
func (r MergeRequestReportResult) Action() string {
	switch {
	case r.Created:
		return "created"
	case r.AdoptedLegacy:
		return "adopted"
	}

	return "updated"
}

func NewMergeRequestReportFlags(config *MergeRequestReportConfig) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.enabled",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_ENABLED"),
			),
			Usage:       "Enable GitLab merge request report note on the given merge request.",
			Required:    false,
			Value:       false,
			Destination: &config.Enabled,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.token",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GL_PIPES_TOKEN"),
			),
			Usage:       "GitLab API token for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.Token,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.api-url",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_API_V4_URL"),
			),
			Usage:       "GitLab API URL for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.ApiUrl,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.project-id",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_PROJECT_ID"),
			),
			Usage:       "GitLab project id for merge request report notes.",
			Required:    false,
			Value:       "",
			Destination: &config.ProjectId,
		},

		&cli.Int64Flag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.merge-request-iid",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_MERGE_REQUEST_IID"),
			),
			Usage:       "GitLab merge request iid for merge request report notes.",
			Required:    false,
			Value:       0,
			Destination: &config.MergeRequestId,
		},

		&cli.StringFlag{
			Category: CATEGORY_GITLAB_MERGE_REQUEST_REPORT,
			Name:     "gitlab-mr-report.identifier",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("GITLAB_MR_REPORT_IDENTIFIER"),
			),
			Usage:       "Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report.",
			Required:    false,
			Value:       "",
			Destination: &config.Identifier,
		},
	}
}

// Builds the marker identifier that keeps concurrent report jobs on the same merge
// request from overwriting each other's note. Discriminators name what the job is
// reporting on -- a stack, a state, a working directory -- and are what makes the
// identifier unique when a matrix or child pipeline runs the same job name several
// times. They must be stable across pipeline runs, so never derive one from a job or
// pipeline id: that would post a new note per push instead of updating the old one.
func ResolveReportIdentifier(override string, jobName string, discriminators ...string) string {
	if override != "" {
		return override
	}

	parts := []string{}
	for _, part := range append([]string{jobName}, discriminators...) {
		if part = sanitizeReportIdentifier(part); part != "" {
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, ":")
}

// Stack and state names reach the identifier from user configuration, and the
// identifier ends up inside an HTML comment marker.
func sanitizeReportIdentifier(value string) string {
	value = strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return -1
		}

		return r
	}, value)

	// One pass splices a fresh terminator out of what it leaves behind, so "--->->"
	// would come back out as "-->" and cut the marker comment short.
	for strings.Contains(value, "-->") {
		value = strings.ReplaceAll(value, "-->", "")
	}

	return strings.TrimSpace(value)
}

// The note carrying the current marker always wins, wherever it sits in the
// listing. Falling back to a legacy marker before exhausting the current one would
// strand the note this job already owns and leave a stale plan on the merge request.
func selectMergeRequestReportNote(notes []*clientgitlab.Note, marker string, legacy []string) *clientgitlab.Note {
	for _, candidate := range append([]string{marker}, legacy...) {
		if index := slices.IndexFunc(notes, func(note *clientgitlab.Note) bool {
			return note != nil && strings.Contains(note.Body, candidate)
		}); index >= 0 {
			return notes[index]
		}
	}

	return nil
}

// The three calls the upsert makes against a merge request's notes, narrowed from
// the client so the note bookkeeping below can be driven without a GitLab to talk to.
// Signatures mirror clientgitlab.NotesService exactly.
type mergeRequestNotes interface {
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

func UpsertMergeRequestReport(
	ctx context.Context,
	config MergeRequestReportConfig,
	body string,
) (*MergeRequestReportResult, error) {
	client, err := clientgitlab.NewClient(
		config.Token,
		clientgitlab.WithBaseURL(config.ApiUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	return upsertMergeRequestReport(ctx, client.Notes, config, body)
}

func upsertMergeRequestReport(
	ctx context.Context,
	notes mergeRequestNotes,
	config MergeRequestReportConfig,
	body string,
) (*MergeRequestReportResult, error) {
	if config.Identifier == "" {
		return nil, fmt.Errorf("merge request report identifier can not be empty")
	}

	marker := mergeRequestReportMarker(config.Identifier)
	if !strings.Contains(body, marker) {
		body = fmt.Sprintf("%s\n\n%s", strings.TrimRight(body, "\n"), marker)
	}

	legacyMarkers := []string{}
	for _, legacy := range config.LegacyIdentifiers {
		if legacy == "" || legacy == config.Identifier {
			continue
		}

		legacyMarkers = append(legacyMarkers, mergeRequestReportMarker(legacy))
	}

	listed := []*clientgitlab.Note{}
	page := int64(1)

	for {
		existing, response, err := notes.ListMergeRequestNotes(
			config.ProjectId,
			config.MergeRequestId,
			&clientgitlab.ListMergeRequestNotesOptions{
				Page:    page,
				PerPage: 100,
			},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("list GitLab merge request notes: %w", err)
		}

		listed = append(listed, existing...)

		if response == nil || response.NextPage == 0 {
			break
		}

		page = response.NextPage
	}

	note := selectMergeRequestReportNote(listed, marker, legacyMarkers)

	if note != nil {
		updated, _, err := notes.UpdateMergeRequestNote(
			config.ProjectId,
			config.MergeRequestId,
			note.ID,
			&clientgitlab.UpdateMergeRequestNoteOptions{
				Body: new(body),
			},
			clientgitlab.WithContext(ctx),
		)
		if err != nil {
			return nil, fmt.Errorf("update GitLab merge request report note: %w", err)
		}

		return &MergeRequestReportResult{
			NoteId:        updated.ID,
			Identifier:    config.Identifier,
			AdoptedLegacy: !strings.Contains(note.Body, marker),
		}, nil
	}

	created, _, err := notes.CreateMergeRequestNote(
		config.ProjectId,
		config.MergeRequestId,
		&clientgitlab.CreateMergeRequestNoteOptions{
			Body: new(body),
		},
		clientgitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab merge request report note: %w", err)
	}

	return &MergeRequestReportResult{
		NoteId:     created.ID,
		Identifier: config.Identifier,
		Created:    true,
	}, nil
}
