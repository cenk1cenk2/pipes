# pipe-mise

Pipe for installing tools with mise.

`pipe-mise [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

### `pipe-mise install`

Install the tools of the mise configuration into the data directory.

The data directory carries its own copy of the mise binary and the shims link to that copy, so a job that restores the directory runs the tools without mise in its image.

`pipe-mise install [FLAGS]`

#### Flags

**Install**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$MISE_INSTALL_PRUNE` | Remove the tool versions the configuration no longer asks for, so the data directory does not grow with every bump. | `bool` | `true` |
| `$MISE_INSTALL_ARGS` | Arguments to append to the install command. | `string` |  |

**Setup**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$MISE_CWD` | Working directory for mise commands. | `string` | `"."` |
| `$MISE_DATA_DIR` | Data directory for mise, which holds the binary, the installs and the shims. Every job that restores it has to use the same absolute path, since the shims link into it. | `string` | `"./.mise/"` |
