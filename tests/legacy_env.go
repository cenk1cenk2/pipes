package tests

// legacyEnvAliases is the exact, ordered environment source chain of every
// visible flag that answers to more than one name. The names a pipeline already
// sets are kept forever and listed ahead of the canonical one, so the order is
// the precedence and dropping or reordering an entry silently changes which
// value a running pipeline picks up.
//
// The table is closed in both directions: a flag listed here has to keep this
// chain, and a flag that grows a second source has to be added here.
var legacyEnvAliases = map[string]map[string][]string{
	"buildah": {
		"buildah.build.file.context":         {"CONTAINER_FILE_CONTEXT", "BUILDAH_BUILD_FILE_CONTEXT"},
		"buildah.build.file.name":            {"CONTAINER_FILE_NAME", "BUILDAH_BUILD_FILE_NAME"},
		"buildah.build.image.build-args":     {"CONTAINER_IMAGE_BUILD_ARGS", "BUILDAH_BUILD_IMAGE_BUILD_ARGS"},
		"buildah.build.image.cache":          {"CONTAINER_IMAGE_CACHE", "BUILDAH_BUILD_IMAGE_CACHE"},
		"buildah.build.image.format":         {"CONTAINER_IMAGE_FORMAT", "BUILDAH_BUILD_IMAGE_FORMAT"},
		"buildah.build.image.latest-tag":     {"CONTAINER_IMAGE_LATEST_TAG", "BUILDAH_BUILD_IMAGE_LATEST_TAG"},
		"buildah.build.image.name":           {"CONTAINER_IMAGE_NAME", "BUILDAH_BUILD_IMAGE_NAME"},
		"buildah.build.image.platforms":      {"CONTAINER_IMAGE_PLATFORMS", "BUILDAH_BUILD_IMAGE_PLATFORMS"},
		"buildah.build.image.pull":           {"CONTAINER_IMAGE_PULL", "BUILDAH_BUILD_IMAGE_PULL"},
		"buildah.build.image.push":           {"CONTAINER_IMAGE_PUSH", "BUILDAH_BUILD_IMAGE_PUSH"},
		"buildah.build.image.storage-driver": {"CONTAINER_IMAGE_STORAGE_DRIVER", "BUILDAH_STORAGE_DRIVER", "BUILDAH_BUILD_IMAGE_STORAGE_DRIVER"},
		"buildah.build.image.tag-as-latest":  {"CONTAINER_IMAGE_TAGS_AS_LATEST", "BUILDAH_BUILD_IMAGE_TAG_AS_LATEST"},
		"buildah.build.image.tags":           {"CONTAINER_IMAGE_TAGS", "BUILDAH_BUILD_IMAGE_TAGS"},
		"buildah.build.image.tags-sanitize":  {"CONTAINER_IMAGE_SANITIZE_TAGS", "BUILDAH_BUILD_IMAGE_TAGS_SANITIZE"},
		"buildah.build.image.tags-template":  {"CONTAINER_IMAGE_TAGS_TEMPLATE", "BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE"},
		"buildah.build.manifest.file":        {"CONTAINER_MANIFEST_FILE", "BUILDAH_BUILD_MANIFEST_FILE"},
		"buildah.build.manifest.target":      {"CONTAINER_MANIFEST_TARGET", "BUILDAH_BUILD_MANIFEST_TARGET"},
		"buildah.login.registry.password":    {"CONTAINER_REGISTRY_PASSWORD", "BUILDAH_LOGIN_REGISTRY_PASSWORD"},
		"buildah.login.registry.uri":         {"CONTAINER_REGISTRY_URI", "BUILDAH_LOGIN_REGISTRY_URI"},
		"buildah.login.registry.username":    {"CONTAINER_REGISTRY_USERNAME", "BUILDAH_LOGIN_REGISTRY_USERNAME"},
		"buildah.manifest.files":             {"CONTAINER_MANIFEST_FILES", "BUILDAH_MANIFEST_FILES"},
		"buildah.manifest.images":            {"CONTAINER_MANIFEST_IMAGES", "BUILDAH_MANIFEST_IMAGES"},
		"buildah.manifest.matrix":            {"CONTAINER_MANIFEST_MATRIX", "BUILDAH_MANIFEST_MATRIX"},
		"buildah.manifest.target":            {"CONTAINER_MANIFEST_TARGET", "BUILDAH_MANIFEST_TARGET"},
		"git.branch":                         {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                            {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"go": {
		"go.build.enable-cgo":   {"GO_BUILD_ENABLE_CGO", "CGO_ENABLED"},
		"go.build.linker-flags": {"GO_BUILD_LINKER", "GO_BUILD_LINKER_FLAGS"},
	},

	"helm": {
		"git.branch":                           {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                              {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		"helm.cwd":                             {"HELM_ROOT", "HELM_CWD"},
		"helm.lint.kubernetes.version":         {"KUBERNETES_VERSION", "HELM_LINT_KUBERNETES_VERSION"},
		"helm.login.registry.password":         {"HELM_REGISTRY_PASSWORD", "HELM_LOGIN_REGISTRY_PASSWORD"},
		"helm.login.registry.uri":              {"HELM_REGISTRY_URI", "HELM_LOGIN_REGISTRY_URI"},
		"helm.login.registry.username":         {"HELM_REGISTRY_USERNAME", "HELM_LOGIN_REGISTRY_USERNAME"},
		"helm.publish.chart.app-version":       {"HELM_CHART_APP_VERSION", "HELM_PUBLISH_CHART_APP_VERSION"},
		"helm.publish.chart.destination":       {"HELM_CHART_DESTINATION", "HELM_PUBLISH_CHART_DESTINATION"},
		"helm.publish.chart.target":            {"HELM_CHART_TARGET", "HELM_PUBLISH_CHART_TARGET"},
		"helm.publish.chart.versions":          {"HELM_CHART_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS"},
		"helm.publish.chart.versions-sanitize": {"HELM_CHART_SANITIZE_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS_SANITIZE"},
		"helm.publish.chart.versions-template": {"HELM_CHART_VERSIONS_TEMPLATE", "HELM_PUBLISH_CHART_VERSIONS_TEMPLATE"},
	},

	"kustomize": {
		"kustomize.build.enable-helm":     {"KUSTOMIZE_ENABLE_HELM", "KUSTOMIZE_BUILD_ENABLE_HELM"},
		"kustomize.build.helm-command":    {"KUSTOMIZE_HELM_COMMAND", "KUSTOMIZE_BUILD_HELM_COMMAND"},
		"kustomize.build.kube-version":    {"KUSTOMIZE_KUBE_VERSION", "KUSTOMIZE_BUILD_KUBE_VERSION"},
		"kustomize.build.load-restrictor": {"KUSTOMIZE_LOAD_RESTRICTOR", "KUSTOMIZE_BUILD_LOAD_RESTRICTOR"},
		"kustomize.cwd":                   {"KUSTOMIZE_ROOT", "KUSTOMIZE_CWD"},
	},

	"node": {
		"git.branch":         {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":            {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		"node.install.cache": {"NODE_INSTALL_CACHE_ENABLE", "NODE_INSTALL_CACHE"},
		"node.run.cwd":       {"NODE_COMMAND_CWD", "NODE_RUN_CWD"},
		"node.run.script":    {"NODE_COMMAND_SCRIPT", "NODE_RUN_SCRIPT"},
	},

	// Both pulumi commands have always read $PULUMI_PLAN, so a job that runs
	// preview and up in turn hands them the same file. The canonical names are
	// per command, which is the way out of that without breaking the jobs that
	// rely on it.
	"pulumi": {
		"pulumi.preview.plan":           {"PULUMI_PLAN", "PULUMI_PREVIEW_PLAN"},
		"pulumi.preview.summary.output": {"PULUMI_SUMMARY_OUTPUT", "PULUMI_PREVIEW_SUMMARY_OUTPUT"},
		"pulumi.up.plan":                {"PULUMI_PLAN", "PULUMI_UP_PLAN"},
	},

	"select-env": {
		"git.branch": {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":    {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"semantic-release": {
		"git.branch":                           {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                              {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		"semantic-release.ci.commit-reference": {"CI_COMMIT_REF_NAME", "SEMANTIC_RELEASE_CI_COMMIT_REFERENCE"},
	},

	"terraform": {
		"terraform.apply.args":                            {"TF_APPLY_ARGS", "TERRAFORM_APPLY_ARGS"},
		"terraform.apply.output":                          {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT", "TERRAFORM_APPLY_OUTPUT"},
		"terraform.ci.api-url":                            {"TF_VAR_CI_API_V4_URL", "CI_API_V4_URL", "TERRAFORM_CI_API_URL"},
		"terraform.ci.project-id":                         {"TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID", "TERRAFORM_CI_PROJECT_ID"},
		"terraform.cwd":                                   {"TF_ROOT", "TERRAFORM_CWD"},
		"terraform.install.args":                          {"TF_INSTALL_ARGS", "TERRAFORM_INSTALL_ARGS"},
		"terraform.install.reconfigure":                   {"TF_INSTALL_RECONFIGURE", "TERRAFORM_INSTALL_RECONFIGURE"},
		"terraform.install.use-lockfile":                  {"TF_INSTALL_USE_LOCKFILE", "TERRAFORM_INSTALL_USE_LOCKFILE"},
		"terraform.lint.format-check.args":                {"TF_LINT_FMT_CHECK_ARGS", "TERRAFORM_LINT_FORMAT_CHECK_ARGS"},
		"terraform.lint.format-check.enable":              {"TF_LINT_FMT_CHECK_ENABLE", "TERRAFORM_LINT_FORMAT_CHECK_ENABLE"},
		"terraform.lint.validate.args":                    {"TF_LINT_VALIDATE_ARGS", "TERRAFORM_LINT_VALIDATE_ARGS"},
		"terraform.lint.validate.enable":                  {"TF_LINT_VALIDATE_ENABLE", "TERRAFORM_LINT_VALIDATE_ENABLE"},
		"terraform.log-level":                             {"TF_LOG_LEVEL", "TF_LOG", "TERRAFORM_LOG_LEVEL"},
		"terraform.login.registry.credentials":            {"TF_REGISTRY_CREDENTIALS", "TERRAFORM_LOGIN_REGISTRY_CREDENTIALS"},
		"terraform.plan.args":                             {"TF_PLAN_ARGS", "TERRAFORM_PLAN_ARGS"},
		"terraform.plan.output":                           {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT", "TERRAFORM_PLAN_OUTPUT"},
		"terraform.plan.pipeline-source":                  {"CI_PIPELINE_SOURCE", "TERRAFORM_PLAN_PIPELINE_SOURCE"},
		"terraform.plan.preview-for-merge-requests":       {"TF_PLAN_PREVIEW_FOR_MRS", "TERRAFORM_PLAN_PREVIEW_FOR_MERGE_REQUESTS"},
		"terraform.plan.retry-delay":                      {"TF_PLAN_RETRY_DELAY", "TERRAFORM_PLAN_RETRY_DELAY"},
		"terraform.plan.retry-tries":                      {"TF_PLAN_RETRY_TRIES", "TERRAFORM_PLAN_RETRY_TRIES"},
		"terraform.plan.summary.output":                   {"TERRAFORM_SUMMARY_OUTPUT", "TERRAFORM_PLAN_SUMMARY_OUTPUT"},
		"terraform.publish.module.cwd":                    {"TF_MODULE_CWD", "TF_ROOT", "TERRAFORM_PUBLISH_MODULE_CWD"},
		"terraform.publish.module.name":                   {"TF_MODULE_NAME", "CI_PROJECT_NAME", "TERRAFORM_PUBLISH_MODULE_NAME"},
		"terraform.publish.module.system":                 {"TF_MODULE_SYSTEM", "TERRAFORM_PUBLISH_MODULE_SYSTEM"},
		"terraform.publish.registry.gitlab.api-url":       {"CI_API_V4_URL", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL"},
		"terraform.publish.registry.gitlab.project-id":    {"CI_PROJECT_ID", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID"},
		"terraform.publish.registry.gitlab.token":         {"CI_JOB_TOKEN", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN"},
		"terraform.publish.registry.name":                 {"TF_MODULE_REGISTRY", "TERRAFORM_PUBLISH_REGISTRY_NAME"},
		"terraform.state.gitlab-http.http-address":        {"TF_HTTP_ADDRESS", "TF_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS"},
		"terraform.state.gitlab-http.http-lock-address":   {"TF_HTTP_LOCK_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS"},
		"terraform.state.gitlab-http.http-lock-method":    {"TF_HTTP_LOCK_METHOD", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD"},
		"terraform.state.gitlab-http.http-password":       {"TF_HTTP_PASSWORD", "TF_PASSWORD", "CI_JOB_TOKEN", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD"},
		"terraform.state.gitlab-http.http-retry-wait-min": {"TF_HTTP_RETRY_WAIT_MIN", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN"},
		"terraform.state.gitlab-http.http-unlock-address": {"TF_HTTP_UNLOCK_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS"},
		"terraform.state.gitlab-http.http-unlock-method":  {"TF_HTTP_UNLOCK_METHOD", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD"},
		"terraform.state.gitlab-http.http-username":       {"TF_HTTP_USERNAME", "TF_USERNAME", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME"},
		"terraform.state.name":                            {"TF_STATE_NAME", "TERRAFORM_STATE_NAME"},
		"terraform.state.strict":                          {"TF_STATE_STRICT", "TERRAFORM_STATE_STRICT"},
		"terraform.state.type":                            {"TF_STATE_TYPE", "TERRAFORM_STATE_TYPE"},
	},

	"update-docker-hub-readme": {
		"docker-hub.password":           {"DOCKER_PASSWORD", "DOCKER_HUB_PASSWORD"},
		"docker-hub.readme.description": {"README_SHORT_DESCRIPTION", "DOCKER_HUB_README_DESCRIPTION"},
		"docker-hub.readme.file":        {"README_FILE", "DOCKER_HUB_README_FILE"},
		"docker-hub.readme.matrix":      {"README_MATRIX", "DOCKER_HUB_README_MATRIX"},
		"docker-hub.readme.repository":  {"DOCKER_IMAGE_NAME", "CONTAINER_IMAGE_NAME", "README_REPOSITORY", "DOCKER_HUB_README_REPOSITORY"},
		"docker-hub.username":           {"DOCKER_USERNAME", "DOCKER_HUB_USERNAME"},
	},
}

// uncategorizedFlags are the visible flags that predate the categories the
// generated documentation files flags under. The list only ever shrinks: a new
// flag without a category fails, and so does an entry that has since been given
// one.
var uncategorizedFlags = map[string][]string{
	"helm": {
		"helm.lint.kubernetes.version",
		"helm.lint.should-template",
	},

	"pulumi": {
		"pulumi.preview.plan",
		"pulumi.preview.summary.output",
		"pulumi.stack",
		"pulumi.up.plan",
	},

	"terraform": {
		"terraform.apply.args",
		"terraform.apply.output",
		"terraform.install.args",
		"terraform.install.reconfigure",
		"terraform.install.use-lockfile",
		"terraform.lint.format-check.args",
		"terraform.lint.format-check.enable",
		"terraform.lint.validate.args",
		"terraform.lint.validate.enable",
		"terraform.login.registry.credentials",
		"terraform.plan.args",
		"terraform.plan.output",
		"terraform.plan.pipeline-source",
		"terraform.plan.preview-for-merge-requests",
		"terraform.plan.retry-delay",
		"terraform.plan.retry-tries",
		"terraform.plan.summary.output",
	},
}
