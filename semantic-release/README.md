# pipe-semantic-release

Releases applications through the semantic-release library.

`pipe-semantic-release [FLAGS]`

## Flags

**CI Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME` | Current commit reference that can be branch or tag name of the project.. | `string` | `false` | <code></code> |

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>info</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code>[]</code> |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` | <code>false</code> |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | <code>[<br />    { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />    { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />    { "match" :"^heads/main$", "environment": "develop" },<br />    { "match": "^heads/master$", "environment": "develop" }<br />]</code> |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | <code>true</code> |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `bool` | `false` | <code>true</code> |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`json([]struct { username: string, password: string, registry?: string, useHttps?: bool })` | `false` | <code></code> |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `false` | <code>[.npmrc]</code> |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` | `false` | <code></code> |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `string`<br/>`enum("npm", "yarn", "pnpm")` | `false` | <code>pnpm</code> |

**Semantic Release**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SEMANTIC_RELEASE_DRY_RUN` | Run semantic-release in dry mode without making changes. | `bool` | `false` | <code>false</code> |
| `$SEMANTIC_RELEASE_WORKSPACE` | Use @qiwi/multi-semantic-release package to do a workspace release. | `bool` | `false` | <code>false</code> |
