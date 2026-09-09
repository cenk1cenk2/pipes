# pipe-pulumi

Pulumi related tasks in the pipeline.

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

**GitLab Merge Request Report**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` | <code>false</code> |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `false` | <code>0</code> |
| `$GITLAB_MR_REPORT_IDENTIFIER` | Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report. | `string` | `false` | <code></code> |

**GitLab Pipeline**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_JOB_NAME` | GitLab CI job name to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_JOB_URL` | GitLab CI job URL to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHA` | Git commit sha to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in the plan report metadata. | `string` | `false` | <code></code> |

**Preview**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_PREVIEW_PLAN`<br />`$PULUMI_PLAN` | Output file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |
| `$PULUMI_PREVIEW_SUMMARY_OUTPUT`<br />`$PULUMI_SUMMARY_OUTPUT` | Output file for pulumi preview summary. Leave empty to skip summary generation. | `string` | `false` | <code>"pulumi-summary.json"</code> |

**Pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Working directory for pulumi commands. | `string` | `false` | <code>"."</code> |

**Stack**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name to use for pulumi commands. | `string` | `true` | <code></code> |

### `pipe-pulumi up`

Apply the Pulumi changes.

#### Flags

**Pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Working directory for pulumi commands. | `string` | `false` | <code>"."</code> |

**Stack**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name to use for pulumi commands. | `string` | `true` | <code></code> |

**Up**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_UP_PLAN`<br />`$PULUMI_PLAN` | Input file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |
