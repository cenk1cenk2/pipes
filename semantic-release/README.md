# pipe-semantic-release

Releases applications through semantic-release library.

`pipe-semantic-release [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>`json([]struct{ username: string, password: string, registry?: string, useHttps?: boolean })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | ".npmrc" |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

### Packages

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ADD_APKS` | APK applications to install before running semantic-release. | `StringSlice` | `false` |  |
| `$ADD_MODULES` | Node packages to install before running semantic-release. | `StringSlice` | `false` |  |

### Semantic Release

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DRY_RUN` | Node packages to install before running semantic-release. | `Bool` | `false` | false |
| `$RUN_MULTI` | Uses @qiwi/multi-semantic-release package to do a workspace release. | `Bool` | `false` | false |
