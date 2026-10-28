package node_test

import (
	"log/slog"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("SetupTaskList", func() {
	run := func(probe tests.TestingCommandResponse, metadata tests.TestingCommandResponse) *tests.LogCapture {
		GinkgoHelper()

		probe.Name = "pnpm"
		probe.Args = []string{"--version"}
		metadata.Name = "cat"
		metadata.Args = []string{"/usr/local/share/pipe-node-pnpm-version"}
		logger, logs := tests.NewCaptureLogger()
		runner := fixtures.Runner(
			tests.TestingCommandResponse{Name: "node", Args: []string{"--version"}, Stdout: "v24.0.0\n"},
			probe,
			metadata,
		)

		Expect(fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					p.Log = logger

					return node.SetupTaskList(p, &node.Config{PackageManager: "pnpm"}, &node.Ctx{})
				},
			},
		}).Run()).To(Succeed())

		return logs
	}

	DescribeTable("reports the probe without mixing stderr into the version", func(stderr string, warnings int) {
		logs := run(
			tests.TestingCommandResponse{Stdout: "12.14.0\n", Stderr: stderr},
			tests.TestingCommandResponse{Stdout: "12.14.0\n"},
		)
		var info []string
		var warning []string
		var debug []string
		for _, record := range logs.Records() {
			switch record.Level {
			case slog.LevelInfo:
				if strings.HasPrefix(record.Message, "pnpm ") {
					info = append(info, record.Message)
				}
			case slog.LevelWarn:
				warning = append(warning, record.Message)
			case slog.LevelDebug:
				debug = append(debug, record.Message)
			}
		}
		Expect(info).To(Equal([]string{"pnpm available in image: v12.14.0", "pnpm version: v12.14.0"}))
		Expect(warning).To(HaveLen(warnings))
		if warnings != 0 {
			Expect(warning).To(ConsistOf("Downloaded pnpm v12.14.0 during version check; image provides v12.14.0."))
		}
		for _, line := range strings.Split(strings.TrimSpace(stderr), "\n") {
			if line != "" {
				Expect(debug).To(ContainElement(ContainSubstring(line)))
			}
		}
		Expect(logs.Messages()).To(ContainElement("node.js version: v24.0.0"))
	},
		Entry("warm cache", "", 0),
		Entry("warm cache with unrelated stderr", "pnpm diagnostic\n", 0),
		Entry("corepack package download", "! Corepack is about to download https://registry.npmjs.org/pnpm/-/pnpm-12.14.0.tgz\n", 1),
		Entry("pnpm native binary download", "Downloading the pnpm 12.14.0 binary for linux-x64...\n", 1),
		Entry("both downloads", "! Corepack is about to download pnpm\nDownloading the pnpm 12.14.0 binary for linux-x64...\n", 1),
	)

	DescribeTable("reports unknown image metadata", func(stderr string) {
		failure := tests.TestingCommandFailure(1)
		logs := run(
			tests.TestingCommandResponse{Stdout: "12.14.0\n"},
			tests.TestingCommandResponse{Stderr: stderr, Result: &failure},
		)

		Expect(logs.Records()).To(ContainElement(tests.LogRecord{Level: slog.LevelInfo, Message: "pnpm available in image: unknown"}))
		Expect(logs.Messages()).To(ContainElement("pnpm version: v12.14.0"))
	},
		Entry("missing file", "cat: /usr/local/share/pipe-node-pnpm-version: No such file or directory\n"),
		Entry("unreadable file", "cat: /usr/local/share/pipe-node-pnpm-version: Permission denied\n"),
	)

	DescribeTable("does not claim a version or download for a failed probe", func(stdout string, exitCode int) {
		result := tests.TestingCommandSuccess()
		if exitCode != 0 {
			result = tests.TestingCommandFailure(exitCode)
		}
		logs := run(
			tests.TestingCommandResponse{Stdout: stdout, Stderr: "Downloading pnpm failed\n", Result: &result},
			tests.TestingCommandResponse{Stdout: "12.14.0\n"},
		)

		Expect(logs.Records()).To(ContainElement(tests.LogRecord{Level: slog.LevelDebug, Message: "Can not fetch package manager version."}))
		for _, record := range logs.Records() {
			Expect(record.Level).NotTo(Equal(slog.LevelWarn))
			Expect(record.Message).NotTo(HavePrefix("pnpm version:"))
		}
	},
		Entry("empty stdout", "", 0),
		Entry("nonzero exit", "", 1),
		Entry("nonzero exit with stdout", "12.14.0\n", 1),
	)
})

var _ = Describe("PackageManagers", func() {
	// the default is what every pipeline that never set the flag runs with.
	It("knows the default package manager", func() {
		Expect(node.PackageManagers).To(HaveKey(node.DefaultPackageManager))
	})

	// every command a task builds comes out of this table, so a package manager
	// missing one of them would produce a command with a hole in it.
	It("spells every operation for each of them", func() {
		for name, commands := range node.PackageManagers {
			Expect(commands.Install).NotTo(BeEmpty(), name)
			Expect(commands.InstallWithLock).NotTo(BeEmpty(), name)
			Expect(commands.Run).NotTo(BeEmpty(), name)
			Expect(commands.Add).NotTo(BeEmpty(), name)
			Expect(commands.Global).NotTo(BeEmpty(), name)
			Expect(commands.Cache).NotTo(BeEmpty(), name)
			Expect(commands.Version).NotTo(BeEmpty(), name)
		}
	})

	// npm is the only one that needs the arguments of a script separated from the
	// arguments of the package manager itself.
	It("only delimits the run arguments for npm", func() {
		Expect(node.PackageManagers["npm"].RunDelimiter).To(Equal([]string{"--"}))

		for name, commands := range node.PackageManagers {
			if name == "npm" {
				continue
			}

			Expect(commands.RunDelimiter).To(BeEmpty(), name)
		}
	})
})

var _ = Describe("NewFlags", func() {
	// a pipe reads the choice back off the same instance it registered, so a flag
	// landing on a copy would leave it on the zero value.
	It("binds the package manager flag onto the given configuration", func() {
		cfg := node.Config{}
		flag := fixtures.Flag[*cli.StringFlag](node.NewFlags(node.Options{Destination: &cfg}), "node.package-manager")

		Expect(flag.Destination).To(BeIdenticalTo(&cfg.PackageManager))
		Expect(flag.Value).To(Equal(node.DefaultPackageManager))
	})
})

var _ = Describe("NewLoginFlags", func() {
	flags := func(cfg *node.Login) []cli.Flag {
		return node.NewLoginFlags(node.LoginOptions{Destination: cfg})
	}

	It("unmarshals the credentials onto the given configuration", func() {
		cfg := node.Login{}

		Expect(fixtures.Flag[*cli.StringFlag](flags(&cfg), "npm.login").Validator(
			`[{ "username": "ci", "token": "npm-token", "registry": "registry.example.com" }]`,
		)).To(Succeed())

		Expect(cfg.Entries).To(Equal([]node.LoginEntry{
			{Username: "ci", Token: "npm-token", Registry: "registry.example.com"},
		}))
	})

	// the credentials are optional, so a pipe that only appends a plain npmrc has
	// to get past validation with nothing set.
	It("leaves the credentials alone for an empty value", func() {
		cfg := node.Login{}

		Expect(fixtures.Flag[*cli.StringFlag](flags(&cfg), "npm.login").Validator("")).To(Succeed())
		Expect(cfg.Entries).To(BeNil())
	})

	// the tasks write to and then read back the files this flag names, so it has to
	// reach the configuration rather than sit unbound on the flag.
	It("binds the npmrc files onto the given configuration", func() {
		cfg := node.Login{}
		registered := flags(&cfg)

		Expect(fixtures.Flag[*cli.StringSliceFlag](registered, "npm.npmrc-file").Destination).
			To(BeIdenticalTo(&cfg.NpmRcFiles))
		Expect(fixtures.Flag[*cli.StringSliceFlag](registered, "npm.npmrc-file").Value).
			To(Equal([]string{".npmrc"}))
		Expect(fixtures.Flag[*cli.StringFlag](registered, "npm.npmrc").Destination).
			To(BeIdenticalTo(&cfg.NpmRc))
	})
})
