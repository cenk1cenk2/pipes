# select-env

Selects an set of environment variable prefix depending on the condition.

`select-env [FLAGS]`

## Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | <code>"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"</code> |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | <code>true</code> |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `bool` | `false` | <code>true</code> |
| `$ENVIRONMENT_FILE` | File for writing the environment variables for selected environment. | `string` | `true` | <code>"env.environment"</code> |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |
