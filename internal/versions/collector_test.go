package versions_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

// The default the buildah and helm flags ship, which is what almost every
// pipeline actually runs with.
var defaultSanitize = []versions.Match{
	{Match: "([^/]*)/(.*)", Template: "{{ index $ 1 | upper }}_{{ index $ 2 }}"},
}

func testLog() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(GinkgoWriter)
	logger.SetLevel(logrus.TraceLevel)

	return logrus.NewEntry(logger)
}

// The expectations here are the output of the buildah and helm implementations
// this collector replaced, captured before they were deleted.
var _ = Describe("Process", func() {
	var log *logrus.Entry

	BeforeEach(func() {
		log = testLog()
	})

	Describe("as buildah publishes container image tags", func() {
		image := func() *versions.Collector {
			return &versions.Collector{
				Sanitize: defaultSanitize,
				Format: func(tag string) string {
					return fmt.Sprintf("%s/%s:%s", "docker.io", "x/y", tag)
				},
			}
		}

		DescribeTable(
			"prefixes the registry and the image name",
			func(tag, expected string) {
				Expect(image().Process(log, tag)).To(Equal(expected))
			},
			Entry("a semver tag", "v1.2.3", "docker.io/x/y:v1.2.3"),
			Entry("a branch name", "main", "docker.io/x/y:main"),
			Entry("the latest tag", "latest", "docker.io/x/y:latest"),
			// A slash is not legal in a tag, which is the whole reason the default
			// sanitizer exists.
			Entry("a slashed branch name", "feature/foo", "docker.io/x/y:FEATURE_foo"),
			// Only the first slash is consumed, so a deeper branch keeps the rest.
			Entry("a twice slashed branch name", "renovate/deps/bump", "docker.io/x/y:RENOVATE_deps/bump"),
		)

		It("omits the registry when the pipe formats without one", func() {
			collector := &versions.Collector{
				Sanitize: defaultSanitize,
				Format: func(tag string) string {
					return fmt.Sprintf("%s:%s", "x/y", tag)
				},
			}

			Expect(collector.Process(log, "v1.2.3")).To(Equal("x/y:v1.2.3"))
		})
	})

	Describe("as helm publishes chart versions", func() {
		chart := func() *versions.Collector {
			return &versions.Collector{Sanitize: defaultSanitize}
		}

		DescribeTable(
			"leaves the value as it is without a format",
			func(version, expected string) {
				Expect(chart().Process(log, version)).To(Equal(expected))
			},
			Entry("a semver version", "v1.2.3", "v1.2.3"),
			Entry("a branch name", "main", "main"),
			Entry("a slashed branch name", "feature/foo", "FEATURE_foo"),
			Entry("a twice slashed branch name", "renovate/deps/bump", "RENOVATE_deps/bump"),
		)
	})

	// The template runs first, so a value it rewrote is what the sanitizer sees.
	Describe("with both a template and a sanitizer", func() {
		collector := func() *versions.Collector {
			return &versions.Collector{
				Templates: []versions.Match{{Match: "^heads/(.*)$", Template: "branch-{{ index $ 1 }}"}},
				Sanitize:  defaultSanitize,
			}
		}

		It("applies the template and leaves the sanitizer nothing to match", func() {
			Expect(collector().Process(log, "heads/main")).To(Equal("branch-main"))
		})

		It("falls through to the sanitizer when no template matches", func() {
			Expect(collector().Process(log, "tags/v1.0.0")).To(Equal("TAGS_v1.0.0"))
		})
	})

	// This is what lets a pipeline pass a template in as the tag itself rather
	// than writing a condition for it.
	It("renders a value that matches nothing as a template of its own", func() {
		collector := &versions.Collector{}

		Expect(collector.Process(log, "{{ 1 | add 1 }}")).To(Equal("2"))
	})

	It("refuses a value that a template sanitized away entirely", func() {
		collector := &versions.Collector{Sanitize: []versions.Match{{Match: "^(.*)$", Template: ""}}}

		_, err := collector.Process(log, "anything")
		Expect(err).To(MatchError("Can not add empty tag to list."))
	})

	It("reports a sanitizer that does not compile", func() {
		collector := &versions.Collector{Sanitize: []versions.Match{{Match: "([", Template: "x"}}}

		_, err := collector.Process(log, "anything")
		Expect(err).To(MatchError(ContainSubstring("Can not compile sanitize regular expression")))
	})

	It("reports a template that does not compile", func() {
		collector := &versions.Collector{Templates: []versions.Match{{Match: "([", Template: "x"}}}

		_, err := collector.Process(log, "anything")
		Expect(err).To(MatchError(ContainSubstring("Can not compile tag template regular expression")))
	})
})

var _ = Describe("Tasks", func() {
	var (
		p   *plumber.Plumber
		tl  *plumber.TaskList
		dir string
	)

	BeforeEach(func() {
		p = plumber.NewPlumber(func(_ *plumber.Plumber) *ucli.Command {
			return &ucli.Command{Name: "test"}
		})
		p.Log.SetOutput(GinkgoWriter)

		tl = &plumber.TaskList{}
		tl.New(p)

		dir = GinkgoT().TempDir()
	})

	collect := func(collector *versions.Collector) []string {
		out := []string{}
		Expect(p.RunJobs(collector.Tasks(tl, &out).Job())).To(Succeed())

		return out
	}

	It("collects what the user gave, in order, without the duplicates", func() {
		out := collect(&versions.Collector{
			Name:     "tags",
			Label:    "Image tags",
			FromUser: []string{"v1.2.3", "v1.2.3", "main"},
			Sanitize: defaultSanitize,
		})

		Expect(out).To(Equal([]string{"v1.2.3", "main"}))
	})

	// The user tags are a package level flag destination shared with the rest of
	// the pipe, so collecting must not truncate them.
	It("leaves the user values alone", func() {
		user := []string{"v1.2.3", "v1.2.3", "main"}

		collect(&versions.Collector{Name: "tags", FromUser: user})

		Expect(user).To(Equal([]string{"v1.2.3", "v1.2.3", "main"}))
	})

	It("collects the comma separated values out of the file", func() {
		Expect(os.WriteFile(filepath.Join(dir, "tags"), []byte("from-file-a,from-file-b\n"), 0600)).To(Succeed())

		out := collect(&versions.Collector{Name: "tags", File: "tags", FileDir: dir})

		Expect(out).To(ConsistOf("from-file-a", "from-file-b"))
	})

	// The user and the file source run in parallel, so which of them lands first
	// is not something a pipe may depend on.
	It("collects both sources at once", func() {
		Expect(os.WriteFile(filepath.Join(dir, "tags"), []byte("from-file"), 0600)).To(Succeed())

		out := collect(&versions.Collector{
			Name:     "tags",
			FromUser: []string{"v1.2.3"},
			File:     "tags",
			FileDir:  dir,
		})

		Expect(out).To(ConsistOf("v1.2.3", "from-file"))
	})

	It("adds the latest value when a reference matches", func() {
		out := collect(&versions.Collector{
			Name:        "tags",
			FromUser:    []string{"v1.2.3"},
			LatestWhen:  []string{`^tags/v?\d+.\d+.\d+$`},
			LatestValue: "latest",
			References:  []string{"tags/v1.2.3", "heads/main"},
		})

		Expect(out).To(Equal([]string{"v1.2.3", "latest"}))
	})

	It("adds nothing when no reference matches", func() {
		out := collect(&versions.Collector{
			Name:        "tags",
			FromUser:    []string{"main"},
			LatestWhen:  []string{`^tags/v?\d+.\d+.\d+$`},
			LatestValue: "latest",
			References:  []string{"heads/main"},
		})

		Expect(out).To(Equal([]string{"main"}))
	})

	// A pipe with no notion of a latest version, such as helm, leaves the patterns
	// nil rather than passing an empty slice.
	It("leaves the latest task out without any patterns", func() {
		out := collect(&versions.Collector{
			Name:        "versions",
			FromUser:    []string{"v1.2.3"},
			LatestValue: "latest",
			References:  []string{"tags/v1.2.3"},
		})

		Expect(out).To(Equal([]string{"v1.2.3"}))
	})

	It("formats every collected value", func() {
		out := collect(&versions.Collector{
			Name:     "tags",
			FromUser: []string{"v1.2.3", "feature/foo"},
			Sanitize: defaultSanitize,
			Format: func(tag string) string {
				return fmt.Sprintf("docker.io/x/y:%s", tag)
			},
		})

		Expect(out).To(Equal([]string{"docker.io/x/y:v1.2.3", "docker.io/x/y:FEATURE_foo"}))
	})

	It("collects nothing when the pipe gave it no source", func() {
		Expect(collect(&versions.Collector{Name: "tags"})).To(BeEmpty())
	})
})
