# pipe-template

A template for CLI scaffolding.

`pipe-template [FLAGS]`

## Flags

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PIPE_DEFAULT_FLAG` | Some default flag. | `string` |  |

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |
