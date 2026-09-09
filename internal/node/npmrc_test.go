package node_test

import (
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("LoginTaskList", func() {
	var dir string

	BeforeEach(func() {
		dir = GinkgoT().TempDir()
	})

	npmrc := func(name string) string {
		return filepath.Join(dir, name)
	}

	run := func(runner *tests.TestingCommandRunner, cfg node.Login) error {
		GinkgoHelper()

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return node.LoginTaskList(p, &cfg)
				},
			},
		}).Run()
	}

	entry := func() node.LoginEntry {
		return node.LoginEntry{
			Username: "ci",
			Token:    "npm-token",
			Registry: "registry.example.com",
			UseHttps: true,
		}
	}

	// a workspace hands the pipe one npmrc per package, and every one of them has to
	// carry every registry the pipeline authenticated against.
	It("writes an auth line per entry into every npmrc file", func() {
		Expect(os.Mkdir(filepath.Join(dir, "project"), 0700)).To(Succeed())

		other := entry()
		other.Registry = "npm.example.com"

		Expect(run(fixtures.Runner(), node.Login{
			Entries:    []node.LoginEntry{entry(), other},
			NpmRcFiles: []string{npmrc(".npmrc"), npmrc("project/.npmrc")},
		})).To(Succeed())

		for _, file := range []string{npmrc(".npmrc"), npmrc("project/.npmrc")} {
			Expect(os.ReadFile(file)).To(SatisfyAll(
				ContainSubstring("//registry.example.com/:_authToken=npm-token"),
				ContainSubstring("//npm.example.com/:_authToken=npm-token"),
			), file)
		}
	})

	// a pipeline may hand the file over whole instead of, or alongside, the
	// credentials the pipe formats for it.
	It("appends the given npmrc behind the credentials", func() {
		Expect(run(fixtures.Runner(), node.Login{
			Entries:    []node.LoginEntry{entry()},
			NpmRc:      "always-auth=true",
			NpmRcFiles: []string{npmrc(".npmrc")},
		})).To(Succeed())

		Expect(os.ReadFile(npmrc(".npmrc"))).To(SatisfyAll(
			ContainSubstring("//registry.example.com/:_authToken=npm-token"),
			ContainSubstring("always-auth=true"),
		))
	})

	// a pipe composes this list unconditionally, so a pipeline that configured
	// neither must not have an npmrc appeared underneath it.
	It("writes nothing without credentials and without an npmrc", func() {
		Expect(run(fixtures.Runner(), node.Login{NpmRcFiles: []string{npmrc(".npmrc")}})).To(Succeed())

		Expect(npmrc(".npmrc")).NotTo(BeAnExistingFile())
	})

	It("checks every registry against the file the credentials landed in", func() {
		runner := fixtures.Runner()

		Expect(run(runner, node.Login{
			Entries:    []node.LoginEntry{entry()},
			NpmRcFiles: []string{npmrc(".npmrc")},
		})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("npm"))
		Expect(invocation.Args).To(Equal([]string{
			"whoami", "--configfile", npmrc(".npmrc"), "--registry", "https://registry.example.com",
		}))
	})

	// without an npmrc file there is nothing holding the credentials for npm to
	// read back, so there is nothing to verify either.
	It("verifies nothing when no npmrc file was named", func() {
		runner := fixtures.Runner()

		Expect(run(runner, node.Login{Entries: []node.LoginEntry{entry()}})).To(Succeed())

		Expect(runner.Invocations()).To(BeEmpty())
	})
})
