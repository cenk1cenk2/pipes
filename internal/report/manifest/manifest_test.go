package manifest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"gitlab.kilic.dev/devops/pipes/internal/report/manifest"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Diff", func() {
	parse := func(documents ...string) []*yaml.RNode {
		GinkgoHelper()

		nodes := []*yaml.RNode{}
		for _, document := range documents {
			node, err := yaml.Parse(document)
			Expect(err).NotTo(HaveOccurred())

			nodes = append(nodes, node)
		}

		return nodes
	}

	configMap := func(name string, value string) string {
		return "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: " + name + "\n  namespace: app\ndata:\n  key: " + value + "\n"
	}

	It("adds every resource of a group the base does not have", func(ctx SpecContext) {
		result, err := manifest.Diff(ctx, fixtures.Task(), []manifest.Group{{
			Name: "overlay",
			Head: parse(
				configMap("settings", "value"),
				"apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole\nmetadata:\n  name: reader\nrules: []\n",
			),
		}})
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Patch).To(ContainSubstring("+++ b/overlay/app/ConfigMap.core/settings.yaml"))
		Expect(result.Patch).To(ContainSubstring("+++ b/overlay/_cluster/ClusterRole.rbac.authorization.k8s.io/reader.yaml"))
		Expect(result.Groups).To(Equal(map[string]manifest.Counts{"overlay": {Added: 2}}))
	})

	It("keeps the same resource of two groups apart", func(ctx SpecContext) {
		result, err := manifest.Diff(ctx, fixtures.Task(), []manifest.Group{
			{Name: "staging", Base: parse(configMap("settings", "old")), Head: parse(configMap("settings", "staging"))},
			{Name: "production", Base: parse(configMap("settings", "old")), Head: parse(configMap("settings", "production"))},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Patch).To(ContainSubstring("+++ b/staging/app/ConfigMap.core/settings.yaml"))
		Expect(result.Patch).To(ContainSubstring("+++ b/production/app/ConfigMap.core/settings.yaml"))
		Expect(result.Patch).To(ContainSubstring("+  key: staging"))
		Expect(result.Patch).To(ContainSubstring("+  key: production"))
		Expect(result.Groups).To(Equal(map[string]manifest.Counts{"staging": {Changed: 1}, "production": {Changed: 1}}))
	})

	It("removes a resource only the base has", func(ctx SpecContext) {
		result, err := manifest.Diff(ctx, fixtures.Task(), []manifest.Group{{
			Name: "overlay",
			Base: parse(configMap("settings", "value"), configMap("removed", "value")),
			Head: parse(configMap("settings", "value")),
		}})
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Patch).To(ContainSubstring("--- a/overlay/app/ConfigMap.core/removed.yaml"))
		Expect(result.Patch).To(ContainSubstring("deleted file mode"))
		Expect(result.Groups).To(Equal(map[string]manifest.Counts{"overlay": {Removed: 1}}))
	})
})
