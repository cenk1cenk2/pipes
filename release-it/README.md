# pipe-release-it

Releases applications through the release-it library.

`pipe-release-it [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### Environment

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `Bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `String`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `Bool` | `false` | true |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `Bool` | `false` | true |

### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `String`<br/>`json([]struct { username: string, password: string, registry?: string, useHttps?: bool })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | ".npmrc" |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `String` | `false` |  |

### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String`<br/>`enum("npm", "yarn", "pnpm")` | `false` | pnpm |

### release-it

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$RELEASE_IT_DRY_RUN` | Run release-it in dry mode without making changes. | `Bool` | `false` | false |
| `$RELEASE_IT_CONFIG_FILE` | release-it configuration file. | `String` | `false` |  |
