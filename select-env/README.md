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
| `$ENVIRONMENT_CONDITIONS` | Regular expression patterns to match for selecting the environment.  | `String`<br/>json({ tags: map[string]RegExp, branches: map[string]RegExp }) | `false` |  |
