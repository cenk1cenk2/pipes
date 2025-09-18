# pipe-template

template-cli

`pipe-template [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PIPE_DEFAULT_FLAG` | Some default flag. | `string` | `false` | `` |

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | `info` |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | `[]` |
