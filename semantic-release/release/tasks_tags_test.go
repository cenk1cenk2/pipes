package release

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Released tags dotenv report", func() {
	const envFile = "semantic-release.env"

	BeforeEach(func() {
		GinkgoT().Chdir(GinkgoT().TempDir())
	})

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe
		*C = Ctx{}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-semantic-release",
			CommandName: "release",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(
								release(tl).Job(),
								tags(tl).Job(),
							)
						})
				},
			},
		}).Run()
	}

	single := func() Pipe {
		return Pipe{SemanticRelease: SemanticRelease{
			TagsFile:    ".tags",
			EnvFile:     envFile,
			EnvVariable: "SEMANTIC_RELEASE_TAG",
		}}
	}

	workspace := func() Pipe {
		p := single()
		p.SemanticRelease.Workspace = true

		return p
	}

	write := func(path string, content string) {
		GinkgoHelper()

		Expect(os.MkdirAll(filepath.Dir(path), 0o700)).To(Succeed())
		Expect(os.WriteFile(path, []byte(content), 0o600)).To(Succeed())
	}

	report := func() string {
		GinkgoHelper()

		content, err := os.ReadFile(envFile)
		Expect(err).NotTo(HaveOccurred())

		return string(content)
	}

	It("hands the tag of a single package over", func() {
		write(".tags", "v1.2.3")

		Expect(run(fixtures.Runner(), single())).To(Succeed())

		Expect(report()).To(Equal("SEMANTIC_RELEASE_TAG=v1.2.3\n"))
	})

	// the preset writes the tag with a trailing newline and exec writers echo it.
	It("trims the whitespace around the tag", func() {
		write(".tags", "v1.2.3\n")

		Expect(run(fixtures.Runner(), single())).To(Succeed())

		Expect(report()).To(Equal("SEMANTIC_RELEASE_TAG=v1.2.3\n"))
	})

	It("writes nothing when nothing was released", func() {
		Expect(run(fixtures.Runner(), single())).To(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	It("writes nothing on a dry run", func() {
		write(".tags", "v1.2.3")

		p := single()
		p.SemanticRelease.DryRun = true

		Expect(run(fixtures.Runner(), p)).To(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	It("writes nothing when the release failed", func() {
		write(".tags", "v1.2.3")

		runner := fixtures.Runner(tests.TestingCommandResponse{Name: SemanticReleaseExe, Err: errors.New("release failed")})

		Expect(run(runner, single())).NotTo(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	It("writes the report and names the variable as the pipeline asks", func() {
		write(".tags", "v1.2.3")

		p := single()
		p.SemanticRelease.EnvFile = "release.env"
		p.SemanticRelease.EnvVariable = "RELEASED"

		Expect(run(fixtures.Runner(), p)).To(Succeed())

		content, err := os.ReadFile("release.env")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(Equal("RELEASED=v1.2.3\n"))
	})

	It("reads an absolute tags file directly", func() {
		path := filepath.Join(GinkgoT().TempDir(), "tags")
		write(path, "v1.2.3")

		p := single()
		p.SemanticRelease.TagsFile = path

		Expect(run(fixtures.Runner(), p)).To(Succeed())

		Expect(report()).To(Equal("SEMANTIC_RELEASE_TAG=v1.2.3\n"))
	})

	It("writes nothing for a tag that spans more than one line", func() {
		write(".tags", "v1.2.3\nv1.2.4")

		Expect(run(fixtures.Runner(), single())).To(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	Describe("workspace", func() {
		// the shape of operator-lb: object form workspaces with a glob, where only the
		// released charts get a tags file.
		BeforeEach(func() {
			write("package.json", `{"name":"operator-lb","workspaces":{"packages":["charts/*"]}}`)
			write("charts/operator/package.json", `{"name":"@cenk1cenk2/operator-lb-chart"}`)
			write("charts/crds/package.json", `{"name":"crds"}`)
			write("charts/unreleased/package.json", `{"name":"unreleased"}`)
		})

		It("keys every released package by its name", func() {
			write("charts/operator/.tags", "1.2.3\n")
			write("charts/crds/.tags", "0.4.0")

			Expect(run(fixtures.Runner(), workspace())).To(Succeed())

			Expect(report()).To(Equal(
				"SEMANTIC_RELEASE_TAG_CENK1CENK2_OPERATOR_LB_CHART=1.2.3\n" +
					"SEMANTIC_RELEASE_TAG_CRDS=0.4.0\n",
			))
		})

		It("skips the installed dependencies", func() {
			write("charts/operator/.tags", "1.2.3")
			write("node_modules/dependency/package.json", `{"name":"dependency"}`)
			write("node_modules/dependency/.tags", "9.9.9")

			Expect(run(fixtures.Runner(), workspace())).To(Succeed())

			Expect(report()).To(Equal("SEMANTIC_RELEASE_TAG_CENK1CENK2_OPERATOR_LB_CHART=1.2.3\n"))
		})

		It("writes nothing when two package names reduce to the same variable", func() {
			write("charts/operator/.tags", "1.2.3")
			write("charts/clash/package.json", `{"name":"@cenk1cenk2/operator.lb.chart"}`)
			write("charts/clash/.tags", "1.2.4")

			Expect(run(fixtures.Runner(), workspace())).To(Succeed())

			Expect(envFile).NotTo(BeAnExistingFile())
		})
	})
})
