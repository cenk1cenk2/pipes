package terraform

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
	// the report model of this package shares its names with the reporting dsl of
	// ginkgo, which a dot import would shadow it with.
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = ginkgo.Describe("Plan report source", func() {
	// one of the pipes reaches for the plan by running a command, so a report that every
	// task read for itself would run it once per task.
	ginkgo.It("reaches for the plan once for every task that reports on it", func() {
		reads := 0
		src := &Source{
			Read: func(_ context.Context, _ *Task) (Report, error) {
				reads++

				return Report{Title: "Example report"}, nil
			},
		}

		task := fixtures.Task("read")

		for range 3 {
			report, err := src.read(context.Background(), task)
			Expect(err).NotTo(HaveOccurred())
			Expect(report.Title).To(Equal("Example report"))
		}

		Expect(reads).To(Equal(1))
	})

	ginkgo.It("holds nothing back when the plan cannot be read", func() {
		src := &Source{
			Read: func(_ context.Context, _ *Task) (Report, error) {
				return Report{}, fmt.Errorf("no plan")
			},
		}

		_, err := src.read(context.Background(), fixtures.Task("read"))
		Expect(err).To(MatchError("no plan"))
	})
})
