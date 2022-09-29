# pipes-semantic-release

Releases applications through semantic-release library.

`pipes-semantic-release [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |

### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>format(json({ username: string, password: string, registry?: string, useHttps?: boolean }[])) | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | [.npmrc] |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

### Packages

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ADD_APKS` | APK applications to install before running semantic-release. | `StringSlice` | `false` | [] |
| `$ADD_MODULES` | Node packages to install before running semantic-release. | `StringSlice` | `false` | [] |

### Semantic Release

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DRY_RUN` | Node packages to install before running semantic-release. | `Bool` | `false` | false |
| `$RUN_MULTI` | Uses @qiwi/multi-semantic-release package to do a workspace release. | `Bool` | `false` | false |
