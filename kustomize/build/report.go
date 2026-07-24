package build

import (
	"bytes"
	"cmp"
	_ "embed"
	"fmt"
	"slices"
	"text/template"
)

//go:embed assets/mr-report.md.gotmpl
var mergeRequestReportTemplate string

type (
	mergeRequestReport struct {
		Metadata []mergeRequestReportMetadata
		Summary  mergeRequestReportSummary
		Overlays []mergeRequestReportOverlay
	}

	mergeRequestReportMetadata struct {
		Name  string
		Value string
	}

	mergeRequestReportSummary struct {
		Total     int
		Succeeded int
		Failed    int
	}

	mergeRequestReportOverlay struct {
		Overlay  string
		Status   string
		DocCount int
		Error    string
	}
)

func newMergeRequestReport(results []OverlayResult) *mergeRequestReport {
	report := &mergeRequestReport{
		Summary: mergeRequestReportSummary{Total: len(results)},
	}

	overlays := make([]mergeRequestReportOverlay, 0, len(results))
	for _, result := range results {
		overlay := mergeRequestReportOverlay{
			Overlay:  result.Overlay,
			DocCount: result.DocCount,
		}

		if result.Err != nil {
			overlay.Status = "failed"
			overlay.Error = result.Err.Error()
			report.Summary.Failed++
		} else {
			overlay.Status = "success"
			report.Summary.Succeeded++
		}

		overlays = append(overlays, overlay)
	}

	slices.SortFunc(overlays, func(left, right mergeRequestReportOverlay) int {
		return cmp.Compare(left.Overlay, right.Overlay)
	})

	report.Overlays = overlays

	return report
}

func renderMergeRequestReport(report *mergeRequestReport) (string, error) {
	tmpl, err := template.New("mr-report.md.gotmpl").
		Parse(mergeRequestReportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse GitLab merge request report template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, report); err != nil {
		return "", fmt.Errorf("render GitLab merge request report template: %w", err)
	}

	return body.String(), nil
}
