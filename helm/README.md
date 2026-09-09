# pipe-helm

Helm chart toolkit for pipelines.

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
| `$HELM_CWD`<br />`$HELM_ROOT` | Working directory for helm commands. | `string` | `false` | <code>"."</code> |

**Helm Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_LOGIN_REGISTRY_URI`<br />`$HELM_REGISTRY_URI` | Helm registry URL to login to. | `string` | `false` | <code>"docker.io"</code> |
| `$HELM_LOGIN_REGISTRY_USERNAME`<br />`$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` | `false` | <code></code> |
| `$HELM_LOGIN_REGISTRY_PASSWORD`<br />`$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` | `false` | <code></code> |

### `pipe-helm lint`

Lint Helm chart templates.

#### Flags

**Helm**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CWD`<br />`$HELM_ROOT` | Working directory for helm commands. | `string` | `false` | <code>"."</code> |

**Helm Lint**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_LINT_KUBERNETES_VERSION`<br />`$KUBERNETES_VERSION` | Kubernetes version to use for linting charts. | `string` | `false` | <code></code> |
| `$HELM_LINT_SHOULD_TEMPLATE` | Template the chart while linting. | `bool` | `false` | <code>true</code> |

### `pipe-helm publish`

Publish Helm chart templates.

#### Flags

**Helm**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_CWD`<br />`$HELM_ROOT` | Working directory for helm commands. | `string` | `false` | <code>"."</code> |

**Helm Chart**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_PUBLISH_CHART_TARGET`<br />`$HELM_CHART_TARGET` | Helm chart repository target to publish to. | `string` | `true` | <code></code> |
| `$HELM_PUBLISH_CHART_VERSIONS`<br />`$HELM_CHART_VERSIONS` | Versions for the helm chart to be published. | `string[]` | `false` | <code></code> |
| `$HELM_PUBLISH_CHART_VERSIONS_TEMPLATE`<br />`$HELM_CHART_VERSIONS_TEMPLATE` | Modifies every version that matches a certain condition.<br />Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>"[]"</code> |
| `$HELM_PUBLISH_CHART_VERSIONS_SANITIZE`<br />`$HELM_CHART_SANITIZE_VERSIONS` | Sanitizes the given regex pattern out of version name.<br />Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `false` | <code>"[\n    { \"match\": \"([^/]*)/(.*)\", \"template\": \"{{ index $ 1 | upper }}_{{ index $ 2 }}\" }\n]"</code> |
| `$HELM_PUBLISH_CHART_DESTINATION`<br />`$HELM_CHART_DESTINATION` | Destination directory for the packaged helm chart. | `string` | `false` | <code>"./dist/"</code> |
| `$HELM_PUBLISH_CHART_APP_VERSION`<br />`$HELM_CHART_APP_VERSION` | Application version for the packaged helm chart. | `string` | `false` | <code></code> |

**Helm Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HELM_LOGIN_REGISTRY_URI`<br />`$HELM_REGISTRY_URI` | Helm registry URL to login to. | `string` | `false` | <code>"docker.io"</code> |
| `$HELM_LOGIN_REGISTRY_USERNAME`<br />`$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` | `false` | <code></code> |
| `$HELM_LOGIN_REGISTRY_PASSWORD`<br />`$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` | `false` | <code></code> |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a comma separated file. | `string` | `false` | <code></code> |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` | <code>false</code> |
