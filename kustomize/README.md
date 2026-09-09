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
