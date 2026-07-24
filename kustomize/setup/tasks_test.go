package setup

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kustomize overlay discovery", func() {
	It("keeps only overlays not nested under another discovered overlay", func() {
		dirs := []string{
			".deploy/base",
			".deploy/base/blackbox-exporter",
			".deploy/base/opentelemetry-collector/logs",
			".deploy/moon",
			".deploy/moon/grafana-alloy",
			".deploy/rubik",
		}

		Expect(rootOverlays(dirs)).To(Equal([]string{
			".deploy/base",
			".deploy/moon",
			".deploy/rubik",
		}))
	})

	It("keeps sibling roots and the working directory root", func() {
		dirs := []string{".", "app", "app/overlays/prod"}

		Expect(rootOverlays(dirs)).To(Equal([]string{".", "app"}))
	})
})
