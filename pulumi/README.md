# pipe-pulumi

Pulumi related tasks in the pipeline.

`pipe-pulumi [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-pulumi preview`](#pipe-pulumi-preview)
- [`pipe-pulumi up`](#pipe-pulumi-up)
- [`pipe-pulumi generate`](#pipe-pulumi-generate)

### `pipe-pulumi preview`

Preview the Pulumi changes.

`pipe-pulumi preview [FLAGS]`

#### Flags

**GitLab Merge Request Report**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` |  |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` |  |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` |  |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `0` |
| `$GITLAB_MR_REPORT_IDENTIFIER` | Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report. | `string` |  |

**GitLab Pipeline**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_JOB_NAME` | GitLab CI job name to include in the plan report metadata. | `string` |  |
| `$CI_JOB_URL` | GitLab CI job URL to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHA` | Git commit sha to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in the plan report metadata. | `string` |  |

**Plan Report**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PLAN_REPORT_LOG_ENABLED` | Write the plan report to the job log. | `bool` | `true` |

**Preview**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_PREVIEW_PLAN`<br/>`$PULUMI_PLAN` | Output file for pulumi plan. | `string` | `"plan.json"` |
| `$PULUMI_PREVIEW_SUMMARY_OUTPUT`<br/>`$PULUMI_SUMMARY_OUTPUT` | Output file for pulumi preview summary. Leave empty to skip summary generation. | `string` | `"pulumi-summary.json"` |

**Pulumi**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_CWD` | Working directory for pulumi commands. | `string` | `"."` |

**Stack**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$PULUMI_STACK`**\* | Stack name to use for pulumi commands. | `string` |  |

\* required

### `pipe-pulumi up`

Apply the Pulumi changes.

`pipe-pulumi up [FLAGS]`

#### Flags

<details>
<summary>Pulumi</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_CWD` | Working directory for pulumi commands. | `string` | `"."` |

</details>

<details>
<summary>Stack</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$PULUMI_STACK`**\* | Stack name to use for pulumi commands. | `string` |  |

\* required

</details>

**Up**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_UP_PLAN`<br/>`$PULUMI_PLAN` | Input file for pulumi plan. | `string` | `"plan.json"` |

### `pipe-pulumi generate`

Generate the Pulumi sources into the working directory. Builds the program, verifies the stack uses the file:// backend of the working directory, removes the local state and the generated paths, then initializes the stack fresh and deploys it. The stack is initialized with an empty passphrase, which only holds when the environment does not set PULUMI_CONFIG_PASSPHRASE, so a generate stack must not carry secrets.

`pipe-pulumi generate [FLAGS]`

#### Flags

**Generate**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_GENERATE_BUILD` | Shell command that builds the Pulumi program before the generation. | `string` |  |
| `$PULUMI_GENERATE_PATHS` | Paths the generation writes, relative to the working directory, removed before every run so the files it no longer generates are gone. | `string`<br/>`format(yaml([]string))` | `"[]"` |
| `$PULUMI_GENERATE_COMMAND` | Shell command that generates the sources against the selected stack. Defaults to deploying the stack. | `string` |  |

<details>
<summary>Pulumi</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PULUMI_CWD` | Working directory for pulumi commands. | `string` | `"."` |

</details>

<details>
<summary>Stack</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$PULUMI_STACK`**\* | Stack name to use for pulumi commands. | `string` |  |

\* required

</details>
