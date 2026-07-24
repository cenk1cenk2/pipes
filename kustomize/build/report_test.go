package build

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kustomize build merge request report", func() {
	It("summarizes overlay build results and renders the report", func() {
		results := []OverlayResult{
			{Overlay: "overlays/production", DocCount: 12},
			{Overlay: "overlays/staging", Err: errors.New("accumulating resources: broken")},
		}

		report := newMergeRequestReport(results)

		Expect(report.Summary).To(Equal(mergeRequestReportSummary{
			Total:     2,
			Succeeded: 1,
			Failed:    1,
		}))
		Expect(report.Overlays).To(ConsistOf(
			mergeRequestReportOverlay{Overlay: "overlays/production", Status: "success", DocCount: 12},
			mergeRequestReportOverlay{Overlay: "overlays/staging", Status: "failed", Error: "accumulating resources: broken"},
		))

		body, err := renderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("Kustomize build report"))
		Expect(body).To(ContainSubstring("`overlays/production`"))
		Expect(body).To(ContainSubstring("1/2"))
		Expect(body).To(ContainSubstring("accumulating resources: broken"))
	})

	It("reports when no overlays were discovered", func() {
		report := newMergeRequestReport(nil)

		body, err := renderMergeRequestReport(report)
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(ContainSubstring("No Kustomize overlays were discovered."))
	})
})
