package login

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
)

var _ = Describe("Flags", func() {
	flag := func(index int) *cli.StringFlag {
		GinkgoHelper()

		return Flags[index].(*cli.StringFlag)
	}

	It("declares the uri, the username and the password", func() {
		Expect(Flags).To(HaveLen(3))
		Expect(flag(0).Name).To(Equal("helm.login.registry.uri"))
		Expect(flag(1).Name).To(Equal("helm.login.registry.username"))
		Expect(flag(2).Name).To(Equal("helm.login.registry.password"))
	})

	It("reads the registry out of the environment of the pipeline", func() {
		Expect(flag(0).Sources.Chain[0].String()).To(ContainSubstring("HELM_LOGIN_REGISTRY_URI"))
		Expect(flag(1).Sources.Chain[0].String()).To(ContainSubstring("HELM_LOGIN_REGISTRY_USERNAME"))
		Expect(flag(2).Sources.Chain[0].String()).To(ContainSubstring("HELM_LOGIN_REGISTRY_PASSWORD"))
	})

	It("files every flag under the helm registry category", func() {
		for i := range Flags {
			Expect(flag(i).Category).To(Equal(CATEGORY_HELM_REGISTRY))
		}
	})

	It("defaults the uri and leaves the credentials empty", func() {
		Expect(flag(0).Value).To(Equal("docker.io"))
		Expect(flag(1).Value).To(BeEmpty())
		Expect(flag(2).Value).To(BeEmpty())
	})

	It("writes into the registry the rest of the pipe reads back", func() {
		Expect(flag(0).Destination).To(BeIdenticalTo(&P.Uri))
		Expect(flag(1).Destination).To(BeIdenticalTo(&P.Username))
		Expect(flag(2).Destination).To(BeIdenticalTo(&P.Password))
	})
})
