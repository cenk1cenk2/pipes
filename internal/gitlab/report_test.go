package gitlab

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	validator "github.com/go-playground/validator/v10"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

var _ = Describe("Merge request report config", func() {
	It("allows an enabled report outside a merge request", func() {
		Expect(validator.New().Struct(MergeRequestReportConfig{
			Enabled: true,
		})).To(Succeed())
	})

	It("rejects an invalid merge request id", func() {
		Expect(validator.New().Struct(MergeRequestReportConfig{
			MergeRequestIid: -1,
		})).NotTo(Succeed())
	})

	It("requires the GitLab context alongside a merge request id", func() {
		Expect(validator.New().Struct(MergeRequestReportConfig{
			Enabled:         true,
			MergeRequestIid: 1,
		})).NotTo(Succeed())
	})
})

var _ = Describe("Report identifier", func() {
	DescribeTable(
		"resolves",
		func(override string, jobName string, discriminators []string, expected string) {
			Expect(ResolveReportIdentifier(override, jobName, discriminators...)).To(Equal(expected))
		},
		Entry("an explicit override ahead of everything else", "custom", "plan", []string{"cache"}, "custom"),
		Entry("the job name alone when nothing else is set", "", "tf-plan", nil, "tf-plan"),
		Entry("the target alongside the job name", "", "plan", []string{"cache"}, "plan:cache"),
		Entry("every target in order", "", "plan", []string{"cache", ".deploy/sun"}, "plan:cache:.deploy/sun"),
		Entry("without the empty parts", "", "", []string{"", "cache"}, "cache"),
		Entry("without the marker terminator", "", "plan", []string{"ca-->che"}, "plan:cache"),
		Entry("without non ascii runes", "", "plan", []string{"caché"}, "plan:cach"),
		Entry("without a terminator spliced from its own leftovers", "", "plan", []string{"--->->"}, "plan"),
	)
})

var _ = Describe("Report note selection", func() {
	var (
		current = mergeRequestReportMarker("plan:cache")
		legacy  = mergeRequestReportMarker("plan")
	)

	notes := func(bodies ...string) []*clientgitlab.Note {
		listed := []*clientgitlab.Note{}
		for index, body := range bodies {
			listed = append(listed, &clientgitlab.Note{ID: int64(index + 1), Body: body})
		}

		return listed
	}

	It("prefers the current marker over an earlier legacy note", func() {
		note := selectMergeRequestReportNote(notes("stale\n"+legacy, "fresh\n"+current), current, []string{legacy})
		Expect(note).NotTo(BeNil())
		Expect(note.ID).To(Equal(int64(2)))
	})

	It("adopts a legacy note when the current marker is absent", func() {
		note := selectMergeRequestReportNote(notes("unrelated", "stale\n"+legacy), current, []string{legacy})
		Expect(note).NotTo(BeNil())
		Expect(note.ID).To(Equal(int64(2)))
	})

	It("returns nothing when no marker matches", func() {
		Expect(selectMergeRequestReportNote(notes("unrelated"), current, []string{legacy})).To(BeNil())
	})
})
