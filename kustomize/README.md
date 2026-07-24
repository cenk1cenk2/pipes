# pipe-kustomize

Kustomize operations for CI pipelines.

`pipe-kustomize [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-kustomize build`

Build and validate Kustomize overlays.

#### Flags

**Gitlab Merge Request Report**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` | <code>false</code> |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `false` | <code>0</code> |
| `$GITLAB_MR_REPORT_IDENTIFIER`<br />`$CI_JOB_NAME` | Hidden marker identifier for merge request report notes. | `string` | `false` | <code></code> |

**Kustomize**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$KUSTOMIZE_ROOT` | Working directory for Kustomize commands. | `string` | `false` | <code>"."</code> |
| `$KUSTOMIZE_PATHS` | Explicit overlay paths to build relative to the working directory. When set, overlay discovery is skipped. | `string[]` | `false` | <code></code> |
| `$KUSTOMIZE_DISCOVERY_PATTERN` | Glob patterns to discover Kustomize overlays under the working directory. | `string[]`<br/>`format(glob)` | `false` | <code>"**/kustomization.yaml", "**/kustomization.yml", "**/Kustomization"</code> |
| `$KUSTOMIZE_DISCOVERY_STRATEGY` | How to filter discovered overlays. "roots" keeps only overlays that are not nested under another discovered overlay, mirroring what ArgoCD renders. "all" keeps every discovered kustomization, including nested ones. | `string`<br/>`format(enum("roots", "all"))` | `false` | <code>"roots"</code> |

**Kustomize Build**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$KUSTOMIZE_ENABLE_HELM` | Enable the Helm chart inflation generator while building overlays. | `bool` | `false` | <code>true</code> |
| `$KUSTOMIZE_HELM_COMMAND` | Helm binary to use for the Helm chart inflation generator. | `string` | `false` | <code>"helm"</code> |
| `$KUSTOMIZE_LOAD_RESTRICTOR` | Load restrictor for Kustomize file access. "rootOnly" restricts loads to the overlay root, "none" allows loading files outside the overlay directory (matches ArgoCD). | `string`<br/>`format(enum("rootOnly", "none"))` | `false` | <code>"none"</code> |
| `$KUSTOMIZE_KUBE_VERSION` | Kubernetes version passed to the Helm chart inflation generator. | `string` | `false` | <code></code> |
| `$KUSTOMIZE_BUILD_OUTPUT` | Output directory to write the rendered manifests per overlay. Leave empty to skip writing. | `string` | `false` | <code></code> |
| `$KUSTOMIZE_BUILD_FAIL_FAST` | Fail on the first overlay that can not be built instead of collecting all results. | `bool` | `false` | <code>false</code> |
