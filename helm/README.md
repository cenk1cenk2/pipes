# pipe-helm

Helm charts for CI pipelines.

`pipe-helm [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-helm install`

Install Helm chart dependencies.

#### Flags

**Helm**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CWD` | Working directory for Helm commands. | `string` | `false` | <code>"."</code> |

**Helm Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_REGISTRY_URI` | Helm registry url to login to. | `string` | `false` | <code>"docker.io"</code> |
| `$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` | `false` | <code></code> |
| `$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` | `false` | <code></code> |

### `pipe-helm lint`

Lint Helm chart templates.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$KUBERNETES_VERSION` | Kubernetes version to use for linting charts. | `string` | `false` | <code></code> |

**Helm**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CWD` | Working directory for Helm commands. | `string` | `false` | <code>"."</code> |

### `pipe-helm publish`

Publish Helm chart templates.

#### Flags

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |

**Helm**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CWD` | Working directory for Helm commands. | `string` | `false` | <code>"."</code> |

**Helm Chart**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CHART_NAME` | Versions for the helm chart to be published. | `string` | `true` | <code></code> |
| `$HELM_CHART_VERSIONS` | Versions for the helm chart to be published. | `string[]` | `true` | <code></code> |
| `$HELM_CHART_VERSIONS_TEMPLATE` | Modifies every version that matches a certain condition.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>"[]"</code> |
| `$HELM_CHART_SANITIZE_VERSIONS` | Sanitizes the given regex pattern out of version name.<br />    Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>"[\n    { \"match\": \"([^/]*)/(.*)\", \"template\": \"{{ index $ 1 | upper }}_{{ index $ 2 }}\" }\n]"</code> |
| `$HELM_CHART_DESTINATION` | Destination directory for the packaged helm chart. | `string` | `false` | <code>"."</code> |

**Helm Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_REGISTRY_URI` | Helm registry url to login to. | `string` | `false` | <code>"docker.io"</code> |
| `$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` | `false` | <code></code> |
| `$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` | `false` | <code></code> |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `string` | `false` | <code></code> |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` | <code>false</code> |
