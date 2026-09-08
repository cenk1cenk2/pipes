package gitlab

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func mergeRequestReportMarker(identifier string) string {
	return fmt.Sprintf("<!-- %s:%s -->", "gitlab-pipes:mr-report", identifier)
}

type MergeRequestReportConfig struct {
	Enabled         bool
	Token           string `validate:"required_with=MergeRequestIid"`
	ApiUrl          string `validate:"required_with=MergeRequestIid"`
	ProjectId       string `validate:"required_with=MergeRequestIid"`
	MergeRequestIid int64  `validate:"omitempty,gt=0"`
	Identifier      string `validate:"omitempty,printascii,excludes=-->"`
	// markers of earlier identifier schemes, so a note already posted under one of
	// them is adopted instead of orphaned next to a duplicate.
	LegacyIdentifiers []string
}

type MergeRequestReportResult struct {
	NoteId     int64
	Identifier string
	Created    bool
	// set when the note was found under an earlier identifier scheme rather than the
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

// Builds the marker identifier that keeps concurrent report jobs on one merge request
// from overwriting each other's note. Discriminators name what the job reports on and
// must be stable across pipeline runs: one derived from a job or pipeline id posts a
// new note per push instead of updating the old one.
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

	// one pass splices a fresh terminator out of what it leaves behind, so "--->->"
	// would come back out as "-->" and cut the marker comment short.
	for strings.Contains(value, "-->") {
		value = strings.ReplaceAll(value, "-->", "")
	}

	return strings.TrimSpace(value)
}

// The note carrying the current marker always wins, wherever it sits in the listing:
// falling back to a legacy marker early would strand the note this job already owns.
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

func UpsertMergeRequestReport(
	ctx context.Context,
	notes Notes,
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
			config.MergeRequestIid,
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
			config.MergeRequestIid,
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
		config.MergeRequestIid,
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
