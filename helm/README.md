# pipe-helm

Helm chart toolkit for pipelines.

`pipe-helm [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-helm install`](#pipe-helm-install)
- [`pipe-helm lint`](#pipe-helm-lint)
- [`pipe-helm publish`](#pipe-helm-publish)

### `pipe-helm install`

Install Helm chart dependencies.

`pipe-helm install [FLAGS]`

#### Flags

**Helm**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_CWD`<br/>`$HELM_ROOT` | Working directory for helm commands. | `string` | `"."` |

**Helm Registry**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_LOGIN_REGISTRY_URI`<br/>`$HELM_REGISTRY_URI` | Helm registry URL to login to. | `string` | `"docker.io"` |
| `$HELM_LOGIN_REGISTRY_USERNAME`<br/>`$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` |  |
| `$HELM_LOGIN_REGISTRY_PASSWORD`<br/>`$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` |  |

### `pipe-helm lint`

Lint Helm chart templates.

`pipe-helm lint [FLAGS]`

#### Flags

<details>
<summary>Helm</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_CWD`<br/>`$HELM_ROOT` | Working directory for helm commands. | `string` | `"."` |

</details>

**Helm Lint**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_LINT_KUBERNETES_VERSION`<br/>`$KUBERNETES_VERSION` | Kubernetes version to use for linting charts. | `string` |  |
| `$HELM_LINT_SHOULD_TEMPLATE` | Template the chart while linting. | `bool` | `true` |

### `pipe-helm publish`

Publish Helm chart templates.

`pipe-helm publish [FLAGS]`

#### Flags

<details>
<summary>Helm</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_CWD`<br/>`$HELM_ROOT` | Working directory for helm commands. | `string` | `"."` |

</details>

**Helm Chart**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$HELM_PUBLISH_CHART_TARGET`**<br/>**`$HELM_CHART_TARGET`**\* | Helm chart repository target to publish to. | `string` |  |
| `$HELM_PUBLISH_CHART_VERSIONS`<br/>`$HELM_CHART_VERSIONS` | Versions for the helm chart to be published. | `string[]` |  |
| `$HELM_PUBLISH_CHART_VERSIONS_TEMPLATE`<br/>`$HELM_CHART_VERSIONS_TEMPLATE` | Modifies every version that matches a certain condition.<br/>Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `"[]"` |
| `$HELM_PUBLISH_CHART_VERSIONS_SANITIZE`<br/>`$HELM_CHART_SANITIZE_VERSIONS` | Sanitizes the given regex pattern out of version name.<br/>Template is interpolated with the given matches in the regular expression. | `string`<br/>`format(yaml([]struct{ match: RegExp, template: Template(match) }))` | `"[\n    { \"match\": \"([^/]*)/(.*)\", \"template\": \"{{ index $ 1 \| upper }}_{{ index $ 2 }}\" }\n]"` |
| `$HELM_PUBLISH_CHART_DESTINATION`<br/>`$HELM_CHART_DESTINATION` | Destination directory for the packaged helm chart. | `string` | `"./dist/"` |
| `$HELM_PUBLISH_CHART_APP_VERSION`<br/>`$HELM_CHART_APP_VERSION` | Application version for the packaged helm chart. | `string` |  |

\* required

<details>
<summary>Helm Registry</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HELM_LOGIN_REGISTRY_URI`<br/>`$HELM_REGISTRY_URI` | Helm registry URL to login to. | `string` | `"docker.io"` |
| `$HELM_LOGIN_REGISTRY_USERNAME`<br/>`$HELM_REGISTRY_USERNAME` | Helm registry username for the given registry. | `string` |  |
| `$HELM_LOGIN_REGISTRY_PASSWORD`<br/>`$HELM_REGISTRY_PASSWORD` | Helm registry password for the given registry. | `string` |  |

</details>

**Tags File**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$TAGS_FILE` | Read tags from a comma separated file. | `string` |  |
| `$TAGS_FILE_STRICT` | Fail on missing tags file. | `bool` | `false` |
