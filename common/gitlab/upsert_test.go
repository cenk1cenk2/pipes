package gitlab

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

// Serves notes a page at a time and records what the upsert wrote back.
type fakeNotes struct {
	pages    [][]*clientgitlab.Note
	listErr  error
	writeErr error

	listedPages []int64
	updatedNote int64
	createdBody string
	updatedBody string
}

func (f *fakeNotes) ListMergeRequestNotes(
	_ any,
	_ int64,
	opt *clientgitlab.ListMergeRequestNotesOptions,
	_ ...clientgitlab.RequestOptionFunc,
) ([]*clientgitlab.Note, *clientgitlab.Response, error) {
	if f.listErr != nil {
		return nil, nil, f.listErr
	}

	f.listedPages = append(f.listedPages, opt.Page)

	index := int(opt.Page) - 1
	if index < 0 || index >= len(f.pages) {
		return nil, &clientgitlab.Response{}, nil
	}

	next := int64(0)
	if index+1 < len(f.pages) {
		next = opt.Page + 1
	}

	return f.pages[index], &clientgitlab.Response{NextPage: next}, nil
}

func (f *fakeNotes) CreateMergeRequestNote(
	_ any,
	_ int64,
	opt *clientgitlab.CreateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	if f.writeErr != nil {
		return nil, nil, f.writeErr
	}

	f.createdBody = *opt.Body

	return &clientgitlab.Note{ID: 99}, nil, nil
}

func (f *fakeNotes) UpdateMergeRequestNote(
	_ any,
	_ int64,
	note int64,
	opt *clientgitlab.UpdateMergeRequestNoteOptions,
	_ ...clientgitlab.RequestOptionFunc,
) (*clientgitlab.Note, *clientgitlab.Response, error) {
	if f.writeErr != nil {
		return nil, nil, f.writeErr
	}

	f.updatedNote = note
	f.updatedBody = *opt.Body

	return &clientgitlab.Note{ID: note}, nil, nil
}

var _ = Describe("Merge request report upsert", func() {
	var (
		config = MergeRequestReportConfig{
			ProjectId:         "3",
			MergeRequestId:    452,
			Identifier:        "plan:cache",
			LegacyIdentifiers: []string{"plan"},
		}
		marker = mergeRequestReportMarker("plan:cache")
		legacy = mergeRequestReportMarker("plan")
	)

	note := func(id int64, body string) *clientgitlab.Note {
		return &clientgitlab.Note{ID: id, Body: body}
	}

	upsert := func(notes *fakeNotes) (*MergeRequestReportResult, error) {
		return upsertMergeRequestReport(context.Background(), notes, config, "report body")
	}

	It("refuses to write without an identifier", func() {
		notes := &fakeNotes{}

		_, err := upsertMergeRequestReport(context.Background(), notes, MergeRequestReportConfig{}, "report body")
		Expect(err).To(MatchError(ContainSubstring("identifier can not be empty")))
		Expect(notes.listedPages).To(BeEmpty())
	})

	It("creates a note carrying the marker when none exists", func() {
		notes := &fakeNotes{pages: [][]*clientgitlab.Note{{note(1, "unrelated")}}}

		result, err := upsert(notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(99)))
		Expect(result.Created).To(BeTrue())
		Expect(result.AdoptedLegacy).To(BeFalse())
		Expect(result.Action()).To(Equal("created"))
		Expect(notes.createdBody).To(HaveSuffix(marker))
		Expect(notes.createdBody).To(HavePrefix("report body"))
	})

	It("updates the note already carrying the marker", func() {
		notes := &fakeNotes{pages: [][]*clientgitlab.Note{{note(7, "stale\n"+marker)}}}

		result, err := upsert(notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(7)))
		Expect(result.Created).To(BeFalse())
		Expect(result.AdoptedLegacy).To(BeFalse())
		Expect(result.Action()).To(Equal("updated"))
		Expect(notes.updatedNote).To(Equal(int64(7)))
	})

	It("adopts a note left under a legacy marker and rewrites it to the current one", func() {
		notes := &fakeNotes{pages: [][]*clientgitlab.Note{{note(4, "stale\n"+legacy)}}}

		result, err := upsert(notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(4)))
		Expect(result.AdoptedLegacy).To(BeTrue())
		Expect(result.Action()).To(Equal("adopted"))
		Expect(notes.updatedBody).To(HaveSuffix(marker))
		Expect(notes.updatedBody).NotTo(ContainSubstring(legacy))
	})

	It("keeps paging until it runs out of notes", func() {
		notes := &fakeNotes{pages: [][]*clientgitlab.Note{
			{note(1, "unrelated")},
			{note(2, "also unrelated")},
			{note(3, "stale\n"+marker)},
		}}

		result, err := upsert(notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(3)))
		Expect(notes.listedPages).To(Equal([]int64{1, 2, 3}))
	})

	It("prefers the current marker over a legacy note on an earlier page", func() {
		notes := &fakeNotes{pages: [][]*clientgitlab.Note{
			{note(1, "stale\n"+legacy)},
			{note(2, "fresh\n"+marker)},
		}}

		result, err := upsert(notes)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(2)))
		Expect(result.AdoptedLegacy).To(BeFalse())
	})

	It("appends the marker only when the body lacks it", func() {
		notes := &fakeNotes{}

		_, err := upsertMergeRequestReport(
			context.Background(),
			notes,
			config,
			"report body\n\n"+marker,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Count(notes.createdBody, marker)).To(Equal(1))
	})

	It("surfaces a listing failure", func() {
		notes := &fakeNotes{listErr: fmt.Errorf("boom")}

		_, err := upsert(notes)
		Expect(err).To(MatchError(ContainSubstring("list GitLab merge request notes")))
	})

	It("surfaces a write failure", func() {
		notes := &fakeNotes{writeErr: fmt.Errorf("boom")}

		_, err := upsert(notes)
		Expect(err).To(MatchError(ContainSubstring("create GitLab merge request report note")))
	})
})
