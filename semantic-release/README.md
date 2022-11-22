# pipe-semantic-release

Releases applications through semantic-release library.

`pipe-semantic-release [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warning", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### Environment

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `Bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `Bool` | `false` | false |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `Bool` | `false` | false |

### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to. | `String`<br/>`json([]struct{ username: string, password: string, registry?: string, useHttps?: boolean })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | ".npmrc" |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn")` | `false` | yarn |

### Packages

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PACKAGES_NODE` | Install node packages before performing operations. | `StringSlice` | `false` |  |
| `$PACKAGES_NODE_GLOBAL` | Install node packages globally. | `Bool` | `false` | false |
| `$PACKAGES_NODE_SCRIPT_ARGS` | package.json script arguments for building operation. | `String`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` |  |
| `$PACKAGES_NODE_CWD` | Working directory for build operation. | `String` | `false` | . |
| `$ADD_APKS` | APK applications to install before running semantic-release. | `StringSlice` | `false` |  |

### Semantic Release

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SEMANTIC_RELEASE_DRY_RUN` | Node packages to install before running semantic-release. | `Bool` | `false` | false |
| `$SEMANTIC_RELEASE_RUN_MULTI` | Use @qiwi/multi-semantic-release package to do a workspace release. | `Bool` | `false` | false |
