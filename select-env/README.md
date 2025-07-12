# select-env

Selects an set of environment variable prefix depending on the condition.

`select-env [FLAGS]`

## Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | [] |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | true |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `bool` | `false` | true |
| `$ENVIRONMENT_FILE` | File for writing the environment variables for selected environment. | `string` | `true` | env.environment |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` | `false` |  |
