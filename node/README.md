# pipe-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipe-node [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | [] |

## Commands

### `pipe-node login`

Login to the given NPM registries.

#### Flags

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`json([]struct { username: string, password: string, registry?: string, useHttps?: bool })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `false` | [.npmrc] |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` | `false` |  |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `string`<br/>`enum("npm", "yarn", "pnpm")` | `false` | pnpm |

### `pipe-node install`

Install node.js dependencies with the given package manager.

#### Flags

**Install**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_INSTALL_CWD` | Install CWD for the package manager. | `string` | `false` | . |
| `$NODE_INSTALL_USE_LOCK_FILE` | Use the lockfile while installing the packages. | `bool` | `false` | true |
| `$NODE_INSTALL_ARGS` | Arguments to append to install command. | `string` | `false` |  |
| `$NODE_INSTALL_CACHE_ENABLE` | Enable caching for the package manager. | `bool` | `false` | true |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`json([]struct { username: string, password: string, registry?: string, useHttps?: bool })` | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `false` | [.npmrc] |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` | `false` |  |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `string`<br/>`enum("npm", "yarn", "pnpm")` | `false` | pnpm |

### `pipe-node build`

#### Flags

**Build**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_BUILD_SCRIPT` | package.json script for building operation. | `string`<br/>`Template(struct { Environment: string, EnvVars: map[string]string })` | `false` | build |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for building operation. | `string`<br/>`Template(struct { Environment: string, EnvVars: map[string]string })` | `false` |  |
| `$NODE_BUILD_CWD` | Working directory for build operation. | `string` | `false` | . |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | true |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `bool` | `false` | true |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` | `false` |  |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `string`<br/>`enum("npm", "yarn", "pnpm")` | `false` | pnpm |

### `pipe-node run`

#### Flags

**Command**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_COMMAND_SCRIPT` | package.json script for given command operation. | `string`<br/>`Template(struct { Environment: string, EnvVars: map[string]string })` | `false` |  |
| `$NODE_COMMAND_CWD` | Working directory for the given command operation. | `string` | `false` | . |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` | false |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`json([]struct{ match: RegExp, environment: string })` | `false` | [<br />  { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },<br />  { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },<br />  { "match" :"^heads/main$", "environment": "develop" },<br />  { "match": "^heads/master$", "environment": "develop" }<br />] |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | true |
| `$ENVIRONMENT_STRICT` | Fail on no environment selected. | `bool` | `false` | true |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` | `false` |  |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `string`<br/>`enum("npm", "yarn", "pnpm")` | `false` | pnpm |
