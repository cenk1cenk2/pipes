package install

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terraform plugin cache prune", func() {
	var (
		root      string
		cache     string
		providers string
	)

	const null = "registry.terraform.io/hashicorp/null"

	write := func(path string, content string) {
		GinkgoHelper()

		Expect(os.MkdirAll(filepath.Dir(path), 0o755)).To(Succeed())
		Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())
	}

	BeforeEach(func() {
		root = GinkgoT().TempDir()
		cache = filepath.Join(root, "cache")
		providers = filepath.Join(root, ".terraform", "providers")

		write(filepath.Join(cache, null, "3.2.1", "linux_amd64", "terraform-provider-null"), "old")
		write(filepath.Join(cache, null, "3.2.3", "linux_amd64", "terraform-provider-null"), "new")
		write(filepath.Join(cache, "registry.terraform.io/hashicorp/random/3.6.0/linux_amd64/terraform-provider-random"), "other")
		write(filepath.Join(cache, "registry.terraform.io/hashicorp/aws/5.0.0/linux_amd64/terraform-provider-aws"), "gone")

		// init links the selected versions into the cache.
		Expect(os.MkdirAll(filepath.Join(providers, null, "3.2.3"), 0o755)).To(Succeed())
		Expect(os.Symlink(filepath.Join(cache, null, "3.2.3", "linux_amd64"), filepath.Join(providers, null, "3.2.3", "linux_amd64"))).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(providers, "registry.terraform.io/hashicorp/random/3.6.0"), 0o755)).To(Succeed())
	})

	It("removes the versions init did not select and keeps the selected ones", func() {
		removed, freed, err := pruneCache(cache, providers)

		Expect(err).NotTo(HaveOccurred())
		Expect(removed).To(ConsistOf(
			filepath.Join(null, "3.2.1"),
			"registry.terraform.io/hashicorp/aws/5.0.0",
		))
		Expect(freed).To(Equal(int64(len("old") + len("gone"))))

		Expect(filepath.Join(cache, null, "3.2.1")).NotTo(BeADirectory())
		Expect(filepath.Join(cache, null, "3.2.3", "linux_amd64", "terraform-provider-null")).To(BeARegularFile())
		Expect(filepath.Join(cache, "registry.terraform.io/hashicorp/random/3.6.0")).To(BeADirectory())
	})

	It("removes the parents a pruned provider leaves empty", func() {
		_, _, err := pruneCache(cache, providers)

		Expect(err).NotTo(HaveOccurred())
		Expect(filepath.Join(cache, "registry.terraform.io/hashicorp/aws")).NotTo(BeADirectory())
		Expect(filepath.Join(cache, null)).To(BeADirectory())
		Expect(cache).To(BeADirectory())
	})

	It("leaves the cache alone when the root selected nothing", func() {
		Expect(os.RemoveAll(filepath.Join(root, ".terraform"))).To(Succeed())

		removed, _, err := pruneCache(cache, providers)

		Expect(err).NotTo(HaveOccurred())
		Expect(removed).To(BeEmpty())
		Expect(filepath.Join(cache, null, "3.2.1")).To(BeADirectory())
	})

	It("treats a missing cache as nothing to prune", func() {
		removed, freed, err := pruneCache(filepath.Join(root, "missing"), providers)

		Expect(err).NotTo(HaveOccurred())
		Expect(removed).To(BeEmpty())
		Expect(freed).To(BeZero())
	})
})
