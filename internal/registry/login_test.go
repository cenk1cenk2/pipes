package registry_test

import (
	"bytes"

	"github.com/cenk1cenk2/plumber/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

var buildah = registry.Spec{
	Category:  "Container Registry",
	Label:     "Container registry",
	Prefix:    "buildah",
	Command:   "login",
	LegacyEnv: "CONTAINER",
}

var _ = Describe("NewFlags", func() {
	var (
		creds *registry.Credentials
		flags []ucli.Flag
	)

	BeforeEach(func() {
		creds = &registry.Credentials{}
		flags = registry.NewFlags(buildah, creds)
	})

	flag := func(index int) *ucli.StringFlag {
		return flags[index].(*ucli.StringFlag)
	}

	It("declares the uri, the username and the password", func() {
		Expect(flags).To(HaveLen(3))
		Expect(flag(0).Name).To(Equal("buildah.login.registry.uri"))
		Expect(flag(1).Name).To(Equal("buildah.login.registry.username"))
		Expect(flag(2).Name).To(Equal("buildah.login.registry.password"))
	})

	// The legacy name is what every pipeline consuming this image sets today, so
	// it has to keep winning over the name that replaced it.
	It("puts the legacy environment name ahead of the canonical one", func() {
		Expect(flag(0).Sources.Chain).To(HaveLen(2))
		Expect(flag(0).Sources.Chain[0].String()).To(ContainSubstring("CONTAINER_REGISTRY_URI"))
		Expect(flag(0).Sources.Chain[1].String()).To(ContainSubstring("BUILDAH_LOGIN_REGISTRY_URI"))
	})

	It("uppercases the canonical environment name", func() {
		Expect(flag(2).Sources.Chain[1].String()).To(ContainSubstring("BUILDAH_LOGIN_REGISTRY_PASSWORD"))
	})

	It("reads the label into the usage of every flag", func() {
		Expect(flag(0).Usage).To(Equal("Container registry url to login to."))
		Expect(flag(1).Usage).To(Equal("Container registry username for the given registry."))
		Expect(flag(2).Usage).To(Equal("Container registry password for the given registry."))
	})

	It("files every flag under the given category", func() {
		for i := range flags {
			Expect(flag(i).Category).To(Equal("Container Registry"))
		}
	})

	It("defaults the uri and leaves the credentials empty", func() {
		Expect(flag(0).Value).To(Equal(registry.DEFAULT_URI))
		Expect(flag(1).Value).To(BeEmpty())
		Expect(flag(2).Value).To(BeEmpty())
	})

	It("writes into the given credentials", func() {
		Expect(flag(0).Destination).To(BeIdenticalTo(&creds.Uri))
		Expect(flag(1).Destination).To(BeIdenticalTo(&creds.Username))
		Expect(flag(2).Destination).To(BeIdenticalTo(&creds.Password))
	})

	// The password is the one value that must never reach a log line, and marking
	// it here is what keeps a new pipe from having to remember to.
	It("marks the password as a secret", func() {
		p := plumber.NewPlumber(func(_ *plumber.Plumber) *ucli.Command {
			return &ucli.Command{Name: "test"}
		})
		output := &bytes.Buffer{}
		p.Log.SetOutput(output)

		creds.Password = "not-in-the-log"
		Expect(cli.Validated(p, creds)).To(Succeed())

		p.Log.Infof("logging in with %s", creds.Password)
		Expect(output.String()).NotTo(ContainSubstring("not-in-the-log"))
	})
})

var _ = Describe("LoginTask", func() {
	var (
		p     *plumber.Plumber
		tl    *plumber.TaskList
		creds *registry.Credentials
	)

	BeforeEach(func() {
		p = plumber.NewPlumber(func(_ *plumber.Plumber) *ucli.Command {
			return &ucli.Command{Name: "test"}
		})
		p.Log.SetOutput(GinkgoWriter)

		tl = &plumber.TaskList{}
		tl.New(p)

		creds = &registry.Credentials{Uri: "docker.io", Username: "user", Password: "secret"}
	})

	// Half a credential pair is not something to guess at, and a pipeline with
	// neither is relying on an ambient login the pipe must not clobber.
	DescribeTable(
		"disables itself without a complete credential pair",
		func(username, password string, disabled bool) {
			creds.Username = username
			creds.Password = password

			Expect(registry.LoginTask(tl, creds, "buildah", "login").IsDisabled()).To(Equal(disabled))
		},
		Entry("both set", "user", "secret", false),
		Entry("no password", "user", "", true),
		Entry("no username", "", "secret", true),
		Entry("neither", "", "", true),
	)
})

var _ = Describe("LoginArgs", func() {
	creds := &registry.Credentials{Uri: "docker.io", Username: "user", Password: "secret"}

	It("puts the given arguments ahead of the credentials", func() {
		Expect(registry.LoginArgs(creds, "registry", "login")).To(Equal([]string{
			"registry", "login", "docker.io", "--username", "user", "--password-stdin",
		}))
	})

	It("works without any leading arguments", func() {
		Expect(registry.LoginArgs(creds)).To(Equal([]string{
			"docker.io", "--username", "user", "--password-stdin",
		}))
	})

	// An argument list shows up in a process listing and in the command trace log,
	// which is exactly what the password has to stay out of.
	It("never carries the password", func() {
		Expect(registry.LoginArgs(creds, "login")).NotTo(ContainElement("secret"))
	})

	// The arguments arrive as a variadic, so appending in place would grow into
	// whatever array backs the caller's slice.
	It("leaves the given arguments alone", func() {
		args := []string{"registry", "login"}

		registry.LoginArgs(creds, args...)

		Expect(args).To(Equal([]string{"registry", "login"}))
	})
})
