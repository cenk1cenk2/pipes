package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// envChains is the environment source chain a flag is already known to resolve
// through, keyed by the canonical name the chain starts with. The names a pipeline
// already sets are kept forever behind the canonical one, so the order is the
// precedence and reordering an entry silently changes which value a running
// pipeline picks up.
//
// Every chain here has to stay the head of the chain the pipe actually documents.
// A new flag needs no entry, and a fallback appended behind the ones below changes
// no precedence, so only reordering, removing or inserting ahead fails the spec.
var envChains = map[string]map[string][]string{
	"buildah": {
		"BUILDAH_BUILD_FILE_CONTEXT":         {"CONTAINER_FILE_CONTEXT"},
		"BUILDAH_BUILD_FILE_NAME":            {"CONTAINER_FILE_NAME"},
		"BUILDAH_BUILD_IMAGE_BUILD_ARGS":     {"CONTAINER_IMAGE_BUILD_ARGS"},
		"BUILDAH_BUILD_IMAGE_CACHE":          {"CONTAINER_IMAGE_CACHE"},
		"BUILDAH_BUILD_IMAGE_FORMAT":         {"CONTAINER_IMAGE_FORMAT"},
		"BUILDAH_BUILD_IMAGE_LATEST_TAG":     {"CONTAINER_IMAGE_LATEST_TAG"},
		"BUILDAH_BUILD_IMAGE_NAME":           {"CONTAINER_IMAGE_NAME"},
		"BUILDAH_BUILD_IMAGE_PLATFORMS":      {"CONTAINER_IMAGE_PLATFORMS"},
		"BUILDAH_BUILD_IMAGE_PULL":           {"CONTAINER_IMAGE_PULL"},
		"BUILDAH_BUILD_IMAGE_PUSH":           {"CONTAINER_IMAGE_PUSH"},
		"BUILDAH_BUILD_IMAGE_STORAGE_DRIVER": {"CONTAINER_IMAGE_STORAGE_DRIVER", "BUILDAH_STORAGE_DRIVER"},
		"BUILDAH_BUILD_IMAGE_TAGS":           {"CONTAINER_IMAGE_TAGS"},
		"BUILDAH_BUILD_IMAGE_TAGS_SANITIZE":  {"CONTAINER_IMAGE_SANITIZE_TAGS"},
		"BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE":  {"CONTAINER_IMAGE_TAGS_TEMPLATE"},
		"BUILDAH_BUILD_IMAGE_TAG_AS_LATEST":  {"CONTAINER_IMAGE_TAGS_AS_LATEST"},
		"BUILDAH_BUILD_MANIFEST_FILE":        {"CONTAINER_MANIFEST_FILE"},
		"BUILDAH_BUILD_MANIFEST_TARGET":      {"CONTAINER_MANIFEST_TARGET"},
		"BUILDAH_LOGIN_REGISTRY_PASSWORD":    {"CONTAINER_REGISTRY_PASSWORD"},
		"BUILDAH_LOGIN_REGISTRY_URI":         {"CONTAINER_REGISTRY_URI"},
		"BUILDAH_LOGIN_REGISTRY_USERNAME":    {"CONTAINER_REGISTRY_USERNAME"},
		"BUILDAH_MANIFEST_FILES":             {"CONTAINER_MANIFEST_FILES"},
		"BUILDAH_MANIFEST_IMAGES":            {"CONTAINER_MANIFEST_IMAGES"},
		"BUILDAH_MANIFEST_MATRIX":            {"CONTAINER_MANIFEST_MATRIX"},
		"BUILDAH_MANIFEST_TARGET":            {"CONTAINER_MANIFEST_TARGET"},
		"CI_COMMIT_REF_NAME":                 {"BITBUCKET_BRANCH"},
		"CI_COMMIT_TAG":                      {"BITBUCKET_TAG"},
	},

	"go": {
		"GO_BUILD_ENABLE_CGO":   {"CGO_ENABLED"},
		"GO_BUILD_LINKER_FLAGS": {"GO_BUILD_LINKER"},
		"GO_WORKSPACE":          {"GO_LINT_WORKSPACE"},
	},

	"helm": {
		"CI_COMMIT_REF_NAME":                   {"BITBUCKET_BRANCH"},
		"CI_COMMIT_TAG":                        {"BITBUCKET_TAG"},
		"HELM_CWD":                             {"HELM_ROOT"},
		"HELM_LINT_KUBERNETES_VERSION":         {"KUBERNETES_VERSION"},
		"HELM_LOGIN_REGISTRY_PASSWORD":         {"HELM_REGISTRY_PASSWORD"},
		"HELM_LOGIN_REGISTRY_URI":              {"HELM_REGISTRY_URI"},
		"HELM_LOGIN_REGISTRY_USERNAME":         {"HELM_REGISTRY_USERNAME"},
		"HELM_PUBLISH_CHART_APP_VERSION":       {"HELM_CHART_APP_VERSION"},
		"HELM_PUBLISH_CHART_DESTINATION":       {"HELM_CHART_DESTINATION"},
		"HELM_PUBLISH_CHART_TARGET":            {"HELM_CHART_TARGET"},
		"HELM_PUBLISH_CHART_VERSIONS":          {"HELM_CHART_VERSIONS"},
		"HELM_PUBLISH_CHART_VERSIONS_SANITIZE": {"HELM_CHART_SANITIZE_VERSIONS"},
		"HELM_PUBLISH_CHART_VERSIONS_TEMPLATE": {"HELM_CHART_VERSIONS_TEMPLATE"},
	},

	"kustomize": {
		"KUSTOMIZE_BUILD_ENABLE_HELM":     {"KUSTOMIZE_ENABLE_HELM"},
		"KUSTOMIZE_BUILD_HELM_COMMAND":    {"KUSTOMIZE_HELM_COMMAND"},
		"KUSTOMIZE_BUILD_KUBE_VERSION":    {"KUSTOMIZE_KUBE_VERSION", "KUBERNETES_VERSION"},
		"KUSTOMIZE_BUILD_LOAD_RESTRICTOR": {"KUSTOMIZE_LOAD_RESTRICTOR"},
		"KUSTOMIZE_CWD":                   {"KUSTOMIZE_ROOT"},
	},

	"node": {
		"CI_COMMIT_REF_NAME":   {"BITBUCKET_BRANCH"},
		"CI_COMMIT_TAG":        {"BITBUCKET_TAG"},
		"NODE_ADD_CWD":         {"PACKAGES_NODE_CWD"},
		"NODE_ADD_GLOBAL":      {"PACKAGES_NODE_GLOBAL"},
		"NODE_ADD_PACKAGES":    {"PACKAGES_NODE"},
		"NODE_ADD_SCRIPT_ARGS": {"PACKAGES_NODE_SCRIPT_ARGS"},
		"NODE_INSTALL_CACHE":   {"NODE_INSTALL_CACHE_ENABLE"},
		"NODE_RUN_CWD":         {"NODE_COMMAND_CWD"},
		"NODE_RUN_SCRIPT":      {"NODE_COMMAND_SCRIPT"},
	},

	"pulumi": {
		"PULUMI_PREVIEW_PLAN":           {"PULUMI_PLAN"},
		"PULUMI_PREVIEW_SUMMARY_OUTPUT": {"PULUMI_SUMMARY_OUTPUT"},
		"PULUMI_UP_PLAN":                {"PULUMI_PLAN"},
	},

	"select-env": {
		"CI_COMMIT_REF_NAME": {"BITBUCKET_BRANCH"},
		"CI_COMMIT_TAG":      {"BITBUCKET_TAG"},
	},

	"semantic-release": {
		"CI_COMMIT_REF_NAME":                   {"BITBUCKET_BRANCH"},
		"CI_COMMIT_TAG":                        {"BITBUCKET_TAG"},
		"SEMANTIC_RELEASE_CI_COMMIT_REFERENCE": {"CI_COMMIT_REF_NAME"},
	},

	"terraform": {
		"TERRAFORM_APPLY_ARGS":                            {"TF_APPLY_ARGS"},
		"TERRAFORM_APPLY_OUTPUT":                          {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		"TERRAFORM_CI_API_URL":                            {"TF_VAR_CI_API_V4_URL", "CI_API_V4_URL"},
		"TERRAFORM_CI_PROJECT_ID":                         {"TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID"},
		"TERRAFORM_CWD":                                   {"TF_ROOT"},
		"TERRAFORM_INSTALL_ARGS":                          {"TF_INSTALL_ARGS"},
		"TERRAFORM_INSTALL_RECONFIGURE":                   {"TF_INSTALL_RECONFIGURE"},
		"TERRAFORM_INSTALL_USE_LOCKFILE":                  {"TF_INSTALL_USE_LOCKFILE"},
		"TERRAFORM_LINT_FORMAT_CHECK_ARGS":                {"TF_LINT_FMT_CHECK_ARGS"},
		"TERRAFORM_LINT_FORMAT_CHECK_ENABLE":              {"TF_LINT_FMT_CHECK_ENABLE"},
		"TERRAFORM_LINT_VALIDATE_ARGS":                    {"TF_LINT_VALIDATE_ARGS"},
		"TERRAFORM_LINT_VALIDATE_ENABLE":                  {"TF_LINT_VALIDATE_ENABLE"},
		"TERRAFORM_LOGIN_REGISTRY_CREDENTIALS":            {"TF_REGISTRY_CREDENTIALS"},
		"TERRAFORM_LOG_LEVEL":                             {"TF_LOG_LEVEL", "TF_LOG"},
		"TERRAFORM_PLAN_ARGS":                             {"TF_PLAN_ARGS"},
		"TERRAFORM_PLAN_OUTPUT":                           {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		"TERRAFORM_PLAN_PIPELINE_SOURCE":                  {"CI_PIPELINE_SOURCE"},
		"TERRAFORM_PLAN_PREVIEW_FOR_MERGE_REQUESTS":       {"TF_PLAN_PREVIEW_FOR_MRS"},
		"TERRAFORM_PLAN_RETRY_DELAY":                      {"TF_PLAN_RETRY_DELAY"},
		"TERRAFORM_PLAN_RETRY_TRIES":                      {"TF_PLAN_RETRY_TRIES"},
		"TERRAFORM_PLAN_SUMMARY_OUTPUT":                   {"TERRAFORM_SUMMARY_OUTPUT"},
		"TERRAFORM_PUBLISH_MODULE_CWD":                    {"TF_MODULE_CWD", "TF_ROOT"},
		"TERRAFORM_PUBLISH_MODULE_NAME":                   {"TF_MODULE_NAME", "CI_PROJECT_NAME"},
		"TERRAFORM_PUBLISH_MODULE_SYSTEM":                 {"TF_MODULE_SYSTEM"},
		"TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL":       {"CI_API_V4_URL"},
		"TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID":    {"CI_PROJECT_ID"},
		"TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN":         {"CI_JOB_TOKEN"},
		"TERRAFORM_PUBLISH_REGISTRY_NAME":                 {"TF_MODULE_REGISTRY"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS":        {"TF_HTTP_ADDRESS", "TF_ADDRESS"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS":   {"TF_HTTP_LOCK_ADDRESS"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD":    {"TF_HTTP_LOCK_METHOD"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD":       {"TF_HTTP_PASSWORD", "TF_PASSWORD", "CI_JOB_TOKEN"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN": {"TF_HTTP_RETRY_WAIT_MIN"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS": {"TF_HTTP_UNLOCK_ADDRESS"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD":  {"TF_HTTP_UNLOCK_METHOD"},
		"TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME":       {"TF_HTTP_USERNAME", "TF_USERNAME"},
		"TERRAFORM_STATE_NAME":                            {"TF_STATE_NAME"},
		"TERRAFORM_STATE_STRICT":                          {"TF_STATE_STRICT"},
		"TERRAFORM_STATE_TYPE":                            {"TF_STATE_TYPE"},
	},

	"update-docker-hub-readme": {
		"DOCKER_HUB_PASSWORD":           {"DOCKER_PASSWORD"},
		"DOCKER_HUB_README_DESCRIPTION": {"README_SHORT_DESCRIPTION"},
		"DOCKER_HUB_README_FILE":        {"README_FILE"},
		"DOCKER_HUB_README_MATRIX":      {"README_MATRIX"},
		"DOCKER_HUB_README_REPOSITORY":  {"DOCKER_IMAGE_NAME", "CONTAINER_IMAGE_NAME", "README_REPOSITORY"},
		"DOCKER_HUB_USERNAME":           {"DOCKER_USERNAME"},
	},
}

var _ = Describe("Environment sources", func() {
	for _, dir := range Pipes() {
		It(dir, func() {
			chains, err := ReadEnvChains(dir)
			Expect(err).NotTo(HaveOccurred())

			for canonical, recorded := range envChains[dir] {
				Expect(chains).To(
					HaveKey(canonical),
					"%s no longer falls back to %v, which changes the value a running pipeline resolves",
					canonical, recorded,
				)

				actual := chains[canonical]
				Expect(len(actual)).To(
					BeNumerically(">=", len(recorded)),
					"%s dropped a source it used to fall back to: %v", canonical, recorded,
				)
				Expect(actual[:len(recorded)]).To(
					Equal(recorded),
					"the source chain of %s moved, which changes the value a running pipeline resolves",
					canonical,
				)
			}
		})
	}
})
