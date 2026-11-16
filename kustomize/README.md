# pipe-kustomize

Kustomize operations for CI pipelines.

`pipe-kustomize [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

### `pipe-kustomize build`

Build and validate Kustomize overlays.

`pipe-kustomize build [FLAGS]`

#### Flags

**Container Registry**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$KUSTOMIZE_CONTAINER_LOGIN`<br/>`$CONTAINER_LOGIN` | Container registries to login to. | `string`<br/>`format(yaml([]struct{ uri: string, username: string, password: string }))` |  |

**GitLab Merge Request Report**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` |  |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` |  |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` |  |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `0` |
| `$GITLAB_MR_REPORT_IDENTIFIER` | Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report. | `string` |  |
| `$CI_COMMIT_BRANCH` | Branch of the pipeline, whose merge requests merged into the commit are reported on when it is the default branch. | `string` |  |
| `$CI_DEFAULT_BRANCH` | Default branch of the project, the only branch whose pipelines report on the merge requests merged into the commit. | `string` |  |
| `$GITLAB_MR_REPORT_DEPLOY_JOB` | Name of the job in the pipeline that applies the plan, which the plan report on the merge requests merged into the commit links to run. Links the pipeline when empty or not found. | `string` |  |
| `$GITLAB_MR_REPORT_PLAN_JOB` | Name of the job in the pipeline that planned what is applied, which the deploy report on the merge requests merged into the commit links to. Links the pipeline when empty or not found. | `string` |  |

**GitLab Pipeline**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_JOB_NAME` | GitLab CI job name to include in the plan report metadata. | `string` |  |
| `$CI_JOB_URL` | GitLab CI job URL to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHA` | Git commit sha to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in the plan report metadata. | `string` |  |

**Kustomize**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$KUSTOMIZE_CWD`<br/>`$KUSTOMIZE_ROOT` | Working directory for kustomize commands. | `string` | `"."` |
| `$KUSTOMIZE_PATHS` | Explicit overlay paths to build relative to the working directory. | `string[]` |  |

**Kustomize Build**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$KUSTOMIZE_BUILD_ENABLE_HELM`<br/>`$KUSTOMIZE_ENABLE_HELM` | Enable the Helm chart inflation generator while building overlays. | `bool` | `true` |
| `$KUSTOMIZE_BUILD_HELM_COMMAND`<br/>`$KUSTOMIZE_HELM_COMMAND` | Helm binary to use for the Helm chart inflation generator. | `string` | `"helm"` |
| `$KUSTOMIZE_BUILD_LOAD_RESTRICTOR`<br/>`$KUSTOMIZE_LOAD_RESTRICTOR` | Restricts which files Kustomize may load, where lifting the restriction matches the behaviour of ArgoCD. | `string`<br/>`format(enum("rootOnly", "none"))` | `"none"` |
| `$KUSTOMIZE_BUILD_KUBE_VERSION`<br/>`$KUSTOMIZE_KUBE_VERSION`<br/>`$KUBERNETES_VERSION` | Kubernetes version passed to the Helm chart inflation generator. | `string` |  |
| `$KUSTOMIZE_BUILD_OUTPUT` | Output directory to write the rendered manifests per overlay. Leave empty to skip writing. | `string` |  |
| `$KUSTOMIZE_BUILD_DIFF_BASE`<br/>`$CI_MERGE_REQUEST_TARGET_BRANCH_SHA`<br/>`$CI_MERGE_REQUEST_DIFF_BASE_SHA` | Commit the overlays are rendered on again for the merge request report, which diffs the rendered manifests against it. A merged results pipeline compares against the target branch, any other merge request pipeline against the merge base. | `string` |  |

**Plan Report**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PLAN_REPORT_LOG_ENABLED` | Write the plan report to the job log. | `bool` | `true` |
