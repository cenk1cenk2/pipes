package node_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
)

var _ = Describe("PackageManagers", func() {
	It("knows the package managers the configuration accepts", func() {
		Expect(node.PackageManagers).To(HaveKey("npm"))
		Expect(node.PackageManagers).To(HaveKey("yarn"))
		Expect(node.PackageManagers).To(HaveKey("pnpm"))
		Expect(node.PackageManagers).To(HaveLen(3))
	})

	// The default is what every pipeline that never set the flag runs with.
	It("knows the default package manager", func() {
		Expect(node.PackageManagers).To(HaveKey(node.DEFAULT_PACKAGE_MANAGER))
	})

	// Every command a task builds comes out of this table, so a package manager
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
		Expect(node.PackageManagers["yarn"].RunDelimiter).To(BeEmpty())
		Expect(node.PackageManagers["pnpm"].RunDelimiter).To(BeEmpty())
	})
})

var _ = Describe("NewFlags", func() {
	It("registers the package manager flag", func() {
		cfg := node.Config{}
		flags := node.NewFlags(&cfg)

		Expect(flags).To(HaveLen(1))
		Expect(flags[0].Names()).To(Equal([]string{"node.package_manager"}))
	})

	// A pipe reads the choice back off the same instance it registered, so a flag
	// landing on a copy would leave it on the zero value.
	It("binds the flag onto the given configuration", func() {
		cfg := node.Config{}
		flags := node.NewFlags(&cfg)

		//nolint:errcheck
		Expect(flags[0].(*ucli.StringFlag).Destination).To(BeIdenticalTo(&cfg.PackageManager))
		//nolint:errcheck
		Expect(flags[0].(*ucli.StringFlag).Value).To(Equal(node.DEFAULT_PACKAGE_MANAGER))
	})
})

var _ = Describe("NewLoginFlags", func() {
	names := func(flags []ucli.Flag) []string {
		found := []string{}
		for _, flag := range flags {
			found = append(found, flag.Names()...)
		}

		return found
	}

	It("registers every part of the npmrc the pipe writes", func() {
		cfg := node.Login{}

		Expect(names(node.NewLoginFlags(&cfg))).To(Equal([]string{
			"npm.login",
			"npm.npmrc_file",
			"npm.npmrc",
		}))
	})

	It("unmarshals the credentials onto the given configuration", func() {
		cfg := node.Login{}
		flags := node.NewLoginFlags(&cfg)

		//nolint:errcheck
		Expect(flags[0].(*ucli.StringFlag).Validator(
			`[{ "username": "ci", "token": "npm-token", "registry": "registry.example.com" }]`,
		)).To(Succeed())

		Expect(cfg.Entries).To(Equal([]node.LoginEntry{
			{Username: "ci", Token: "npm-token", Registry: "registry.example.com"},
		}))
	})

	// The credentials are optional, so a pipe that only appends a plain npmrc has
	// to get past validation with nothing set.
	It("leaves the credentials alone for an empty value", func() {
		cfg := node.Login{}

		//nolint:errcheck
		Expect(node.NewLoginFlags(&cfg)[0].(*ucli.StringFlag).Validator("")).To(Succeed())
		Expect(cfg.Entries).To(BeNil())
	})

	// The tasks write to and then read back the files this flag names, so it has to
	// reach the configuration rather than sit unbound on the flag.
	It("binds the npmrc files onto the given configuration", func() {
		cfg := node.Login{}
		flags := node.NewLoginFlags(&cfg)

		//nolint:errcheck
		Expect(flags[1].(*ucli.StringSliceFlag).Destination).To(BeIdenticalTo(&cfg.NpmRcFiles))
		//nolint:errcheck
		Expect(flags[1].(*ucli.StringSliceFlag).Value).To(Equal([]string{".npmrc"}))
		//nolint:errcheck
		Expect(flags[2].(*ucli.StringFlag).Destination).To(BeIdenticalTo(&cfg.NpmRc))
	})
})
