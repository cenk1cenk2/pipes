package build

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

var _ = Describe("Container image tags", func() {
	format := func(uri, name, tag string) string {
		GinkgoHelper()

		P.Image.Name = name

		return ContainerImageTags(Deps{Registry: &registry.Credentials{Uri: uri}}).Format(tag)
	}

	It("prefixes the image with the registry the login step authenticated against", func() {
		Expect(format("registry.example.com", "group/image", "v1.0.0")).
			To(Equal("registry.example.com/group/image:v1.0.0"))
	})

	// A pipeline with no registry configured builds the image locally, where a
	// prefix would name something that is never pushed.
	It("leaves the image unprefixed without a registry", func() {
		Expect(format("", "group/image", "v1.0.0")).To(Equal("group/image:v1.0.0"))
	})

	It("keeps a registry that already carries a port and a path", func() {
		Expect(format("registry.example.com:5000/mirror", "image", "latest")).
			To(Equal("registry.example.com:5000/mirror/image:latest"))
	})
})
