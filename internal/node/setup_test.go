package node_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("PackageManagers", func() {
	// the default is what every pipeline that never set the flag runs with.
	It("knows the default package manager", func() {
		Expect(node.PackageManagers).To(HaveKey(node.DEFAULT_PACKAGE_MANAGER))
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
		Expect(flag.Value).To(Equal(node.DEFAULT_PACKAGE_MANAGER))
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
