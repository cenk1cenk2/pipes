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

**Kustomize**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$KUSTOMIZE_ROOT`<br />`$KUSTOMIZE_CWD` | Working directory for kustomize commands. | `string` | `false` | <code>"."</code> |
| `$KUSTOMIZE_PATHS` | Explicit overlay paths to build relative to the working directory. | `string[]` | `false` | <code></code> |

**Kustomize Build**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$KUSTOMIZE_ENABLE_HELM`<br />`$KUSTOMIZE_BUILD_ENABLE_HELM` | Enable the Helm chart inflation generator while building overlays. | `bool` | `false` | <code>true</code> |
| `$KUSTOMIZE_HELM_COMMAND`<br />`$KUSTOMIZE_BUILD_HELM_COMMAND` | Helm binary to use for the Helm chart inflation generator. | `string` | `false` | <code>"helm"</code> |
| `$KUSTOMIZE_LOAD_RESTRICTOR`<br />`$KUSTOMIZE_BUILD_LOAD_RESTRICTOR` | Load restrictor for Kustomize file access. "rootOnly" restricts loads to the overlay root, "none" allows loading files outside the overlay directory (matches ArgoCD). | `string`<br/>`format(enum("rootOnly", "none"))` | `false` | <code>"none"</code> |
| `$KUSTOMIZE_KUBE_VERSION`<br />`$KUSTOMIZE_BUILD_KUBE_VERSION` | Kubernetes version passed to the Helm chart inflation generator. | `string` | `false` | <code></code> |
| `$KUSTOMIZE_BUILD_OUTPUT` | Output directory to write the rendered manifests per overlay. Leave empty to skip writing. | `string` | `false` | <code></code> |
