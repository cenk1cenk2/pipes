# select-env

Selects an set of environment variable prefix depending on the condition.

`select-env [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |

### Environment

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment. Use either &#34;heads/&#34; for narrowing the search to branches or &#34;tags/&#34; for narrowing the search to tags.  | `String`<br/>json([]struct{ match: RegExp, environment: string }) | `false` |  |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Whether to fail on missing environment references. | `Bool` | `false` | false |
| `$ENVIRONMENT_FILE` | File for writing the environment variables for selected environment. | `String` | `true` |  |

### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |
