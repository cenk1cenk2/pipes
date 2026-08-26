package gitlab

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"
	mockgitlab "gitlab.kilic.dev/devops/pipes/internal/test/mocks/gitlab"
)

var _ = Describe("Merge request report upsert", func() {
	var (
		config = MergeRequestReportConfig{
			ProjectId:         "3",
			MergeRequestIid:   452,
			Identifier:        "plan:cache",
			LegacyIdentifiers: []string{"plan"},
		}
		marker = mergeRequestReportMarker("plan:cache")
		legacy = mergeRequestReportMarker("plan")
	)

	var (
		notes       *mockgitlab.MockNotes
		listedPages []int64
		createdBody string
		updatedNote int64
		updatedBody string
	)

	BeforeEach(func() {
		notes = mockgitlab.NewMockNotes(GinkgoT())
		listedPages = nil
		createdBody = ""
		updatedNote = 0
		updatedBody = ""
	})

	note := func(id int64, body string) *clientgitlab.Note {
		return &clientgitlab.Note{ID: id, Body: body}
	}

	// Serves the notes a page at a time the way the API does, and records which
	// pages the upsert actually asked for.
	listing := func(pages ...[]*clientgitlab.Note) {
		notes.EXPECT().
			ListMergeRequestNotes(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(
				_ any,
				_ int64,
				opt *clientgitlab.ListMergeRequestNotesOptions,
				_ ...clientgitlab.RequestOptionFunc,
			) ([]*clientgitlab.Note, *clientgitlab.Response, error) {
				listedPages = append(listedPages, opt.Page)

				index := int(opt.Page) - 1
				if index < 0 || index >= len(pages) {
					return nil, &clientgitlab.Response{}, nil
				}

				next := int64(0)
				if index+1 < len(pages) {
					next = opt.Page + 1
				}

				return pages[index], &clientgitlab.Response{NextPage: next}, nil
			})
	}

	listingFails := func(err error) {
		notes.EXPECT().
			ListMergeRequestNotes(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil, err)
	}

	creating := func() {
		notes.EXPECT().
			CreateMergeRequestNote(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(
				_ any,
				_ int64,
				opt *clientgitlab.CreateMergeRequestNoteOptions,
				_ ...clientgitlab.RequestOptionFunc,
			) (*clientgitlab.Note, *clientgitlab.Response, error) {
				createdBody = *opt.Body

				return &clientgitlab.Note{ID: 99}, nil, nil
			})
	}

	creatingFails := func(err error) {
		notes.EXPECT().
			CreateMergeRequestNote(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil, err)
	}

	updating := func() {
		notes.EXPECT().
			UpdateMergeRequestNote(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(
				_ any,
				_ int64,
				note int64,
				opt *clientgitlab.UpdateMergeRequestNoteOptions,
				_ ...clientgitlab.RequestOptionFunc,
			) (*clientgitlab.Note, *clientgitlab.Response, error) {
				updatedNote = note
				updatedBody = *opt.Body

				return &clientgitlab.Note{ID: note}, nil, nil
			})
	}

	upsert := func() (*MergeRequestReportResult, error) {
		return UpsertMergeRequestReport(context.Background(), notes, config, "report body")
	}

	It("refuses to write without an identifier", func() {
		_, err := UpsertMergeRequestReport(context.Background(), notes, MergeRequestReportConfig{}, "report body")
		Expect(err).To(MatchError(ContainSubstring("identifier can not be empty")))
		Expect(listedPages).To(BeEmpty())
	})

	It("creates a note carrying the marker when none exists", func() {
		listing([]*clientgitlab.Note{note(1, "unrelated")})
		creating()

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(99)))
		Expect(result.Created).To(BeTrue())
		Expect(result.AdoptedLegacy).To(BeFalse())
		Expect(result.Action()).To(Equal("created"))
		Expect(createdBody).To(HaveSuffix(marker))
		Expect(createdBody).To(HavePrefix("report body"))
	})

	It("updates the note already carrying the marker", func() {
		listing([]*clientgitlab.Note{note(7, "stale\n"+marker)})
		updating()

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(7)))
		Expect(result.Created).To(BeFalse())
		Expect(result.AdoptedLegacy).To(BeFalse())
		Expect(result.Action()).To(Equal("updated"))
		Expect(updatedNote).To(Equal(int64(7)))
	})

	It("adopts a note left under a legacy marker and rewrites it to the current one", func() {
		listing([]*clientgitlab.Note{note(4, "stale\n"+legacy)})
		updating()

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(4)))
		Expect(result.AdoptedLegacy).To(BeTrue())
		Expect(result.Action()).To(Equal("adopted"))
		Expect(updatedBody).To(HaveSuffix(marker))
		Expect(updatedBody).NotTo(ContainSubstring(legacy))
	})

	It("keeps paging until it runs out of notes", func() {
		listing(
			[]*clientgitlab.Note{note(1, "unrelated")},
			[]*clientgitlab.Note{note(2, "also unrelated")},
			[]*clientgitlab.Note{note(3, "stale\n"+marker)},
		)
		updating()

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(3)))
		Expect(listedPages).To(Equal([]int64{1, 2, 3}))
	})

	It("prefers the current marker over a legacy note on an earlier page", func() {
		listing(
			[]*clientgitlab.Note{note(1, "stale\n"+legacy)},
			[]*clientgitlab.Note{note(2, "fresh\n"+marker)},
		)
		updating()

		result, err := upsert()
		Expect(err).NotTo(HaveOccurred())
		Expect(result.NoteId).To(Equal(int64(2)))
		Expect(result.AdoptedLegacy).To(BeFalse())
	})

	It("appends the marker only when the body lacks it", func() {
		listing()
		creating()

		_, err := UpsertMergeRequestReport(
			context.Background(),
			notes,
			config,
			"report body\n\n"+marker,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Count(createdBody, marker)).To(Equal(1))
	})

	It("surfaces a listing failure", func() {
		listingFails(fmt.Errorf("boom"))

		_, err := upsert()
		Expect(err).To(MatchError(ContainSubstring("list GitLab merge request notes")))
	})

	It("surfaces a write failure", func() {
		listing()
		creatingFails(fmt.Errorf("boom"))

		_, err := upsert()
		Expect(err).To(MatchError(ContainSubstring("create GitLab merge request report note")))
	})
})
