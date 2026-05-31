# pipe-pulumi

Pulumi actions for CI pipelines.

`pipe-pulumi [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-pulumi preview`

Preview the Pulumi changes.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name for the pulumi to be used in the commands. | `string` | `true` | <code></code> |
| `$PULUMI_PLAN` | Output file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |

**Gitlab Merge Request Report**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` | <code>false</code> |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `false` | <code>0</code> |
| `$GITLAB_MR_REPORT_IDENTIFIER`<br />`$CI_JOB_NAME` | Hidden marker identifier for merge request report notes. | `string` | `false` | <code></code> |

**Gitlab Pipeline**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_JOB_NAME` | GitLab CI job name to include in merge request report metadata. | `string` | `false` | <code></code> |
| `$CI_JOB_URL` | GitLab CI job URL to include in merge request report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in merge request report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in merge request report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHA` | Git commit sha to include in merge request report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in merge request report metadata. | `string` | `false` | <code></code> |

**pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Path to the Pulumi working directory. | `string` | `false` | <code>"."</code> |

### `pipe-pulumi up`

Apply the Pulumi changes.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name for the pulumi to be used in the commands. | `string` | `true` | <code></code> |
| `$PULUMI_PLAN` | Input file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |

**pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Path to the Pulumi working directory. | `string` | `false` | <code>"."</code> |
