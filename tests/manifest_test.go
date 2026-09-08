package tests

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// published is DOCKER_HUB_README_MATRIX written out verbatim, in the order the
// job carries it. Changing the job and changing this list is the same commit.
var published = []ReadmeMatrixEntry{
	{Repository: "cenk1cenk2/pipe-buildah", File: "./buildah/README.md", Description: "Builds and publishes a container image with given conditions."},
	{Repository: "cenk1cenk2/pipe-node", File: "./node/README.md", Description: "Node.JS operations for pipelines."},
	{Repository: "cenk1cenk2/pipe-pulumi", File: "./pulumi/README.md", Description: "Pulumi operations for pipelines."},
	{Repository: "cenk1cenk2/pipe-go", File: "./go/README.md", Description: "Golang operations for pipelines."},
	{Repository: "cenk1cenk2/pipe-helm", File: "./helm/README.md", Description: "Helm operations for pipelines."},
	{Repository: "cenk1cenk2/pipe-kustomize", File: "./kustomize/README.md", Description: "Kustomize operations for pipelines."},
	{Repository: "cenk1cenk2/pipe-select-env", File: "./select-env/README.md", Description: "Selects an environment given on the conditions."},
	{Repository: "cenk1cenk2/pipe-semantic-release", File: "./semantic-release/README.md", Description: "semantic-release embedded inside a container for CI jobs."},
	{Repository: "cenk1cenk2/pipe-terraform", File: "./terraform/README.md", Description: "Terraform helper pipe."},
	{Repository: "cenk1cenk2/pipe-update-docker-hub-readme", File: "./update-docker-hub-readme/README.md", Description: "Updates the README on DockerHub for given repository."},
}

var _ = Describe("update-docker-hub-readme", func() {
	It("publishes exactly this matrix", func() {
		matrix, err := ReadReadmeMatrix()
		Expect(err).NotTo(HaveOccurred())

		Expect(matrix).To(Equal(published))
	})

	It("points every entry at a readme that is there", func() {
		for _, entry := range published {
			Expect(filepath.Join(Root(), entry.File)).To(BeAnExistingFile())
		}
	})
})
