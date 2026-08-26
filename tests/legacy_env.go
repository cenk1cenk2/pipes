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
		"buildah.login.registry.password": {"CONTAINER_REGISTRY_PASSWORD", "BUILDAH_LOGIN_REGISTRY_PASSWORD"},
		"buildah.login.registry.uri":      {"CONTAINER_REGISTRY_URI", "BUILDAH_LOGIN_REGISTRY_URI"},
		"buildah.login.registry.username": {"CONTAINER_REGISTRY_USERNAME", "BUILDAH_LOGIN_REGISTRY_USERNAME"},
		"container-image.storage-driver":  {"CONTAINER_IMAGE_STORAGE_DRIVER", "BUILDAH_STORAGE_DRIVER"},
		"git.branch":                      {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                         {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"go": {
		"go.build.enable-cgo": {"GO_BUILD_ENABLE_CGO", "CGO_ENABLED"},
	},

	"helm": {
		"git.branch":                   {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                      {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		"helm.cwd":                     {"HELM_ROOT", "HELM_CWD"},
		"helm.login.registry.password": {"HELM_REGISTRY_PASSWORD", "HELM_LOGIN_REGISTRY_PASSWORD"},
		"helm.login.registry.uri":      {"HELM_REGISTRY_URI", "HELM_LOGIN_REGISTRY_URI"},
		"helm.login.registry.username": {"HELM_REGISTRY_USERNAME", "HELM_LOGIN_REGISTRY_USERNAME"},
	},

	"kustomize": {
		"kustomize.cwd": {"KUSTOMIZE_ROOT", "KUSTOMIZE_CWD"},
	},

	"node": {
		"git.branch": {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":    {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"select-env": {
		"git.branch": {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":    {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"semantic-release": {
		"git.branch": {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":    {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},

	"terraform": {
		"terraform-apply.out":                       {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		"terraform-config.log-level":                {"TF_LOG_LEVEL", "TF_LOG"},
		"terraform-module.cwd":                      {"TF_MODULE_CWD", "TF_ROOT"},
		"terraform-module.name":                     {"TF_MODULE_NAME", "CI_PROJECT_NAME"},
		"terraform-plan.out":                        {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT"},
		"terraform-state.gitlab-http.http-address":  {"TF_HTTP_ADDRESS", "TF_ADDRESS"},
		"terraform-state.gitlab-http.http-password": {"TF_HTTP_PASSWORD", "TF_PASSWORD", "CI_JOB_TOKEN"},
		"terraform-state.gitlab-http.http-username": {"TF_HTTP_USERNAME", "TF_USERNAME"},
		"terraform-var.api-url":                     {"TF_VAR_CI_API_V4_URL", "CI_API_V4_URL"},
		"terraform-var.project-id":                  {"TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID"},
		"terraform.cwd":                             {"TF_ROOT", "TERRAFORM_CWD"},
	},

	"update-docker-hub-readme": {
		"readme.repository": {"DOCKER_IMAGE_NAME", "CONTAINER_IMAGE_NAME", "README_REPOSITORY"},
	},
}

// uncategorizedFlags are the visible flags that predate the categories the
// generated documentation files flags under. The list only ever shrinks: a new
// flag without a category fails, and so does an entry that has since been given
// one.
var uncategorizedFlags = map[string][]string{
	"helm": {
		"helm-lint.should-template",
		"kubernetes.version",
	},

	"pulumi": {
		"pulumi.preview.plan",
		"pulumi.preview.summary-output",
		"pulumi.stack",
		"pulumi.up.plan",
	},

	"terraform": {
		"terraform-apply.args",
		"terraform-apply.out",
		"terraform-install.args",
		"terraform-install.reconfigure",
		"terraform-install.use-lockfile",
		"terraform-lint.fmt-check.args",
		"terraform-lint.fmt-check.enable",
		"terraform-lint.validate.args",
		"terraform-lint.validate.enable",
		"terraform-plan.args",
		"terraform-plan.out",
		"terraform-plan.pipeline-source",
		"terraform-plan.preview-for-mrs",
		"terraform-plan.retry-delay",
		"terraform-plan.retry-tries",
		"terraform-plan.summary-output",
		"terraform-registry.credentials",
	},
}
