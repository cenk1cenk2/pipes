# pipe-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipe-node [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-node login`

Login to the given NPM registries.

#### Flags

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))` | `false` | <code></code> |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `false` | <code>".npmrc"</code> |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` | `false` | <code></code> |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `false` | <code>"pnpm"</code> |

### `pipe-node install`

Install node.js dependencies with the given package manager.

#### Flags

**Install**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_INSTALL_CWD` | Working directory for the install operation. | `string` | `false` | <code>"."</code> |
| `$NODE_INSTALL_USE_LOCK_FILE` | Use the lockfile while installing the packages. | `bool` | `false` | <code>true</code> |
| `$NODE_INSTALL_ARGS` | Arguments to append to the install command. | `string` | `false` | <code></code> |
| `$NODE_INSTALL_CACHE`<br />`$NODE_INSTALL_CACHE_ENABLE` | Enable caching for the package manager. | `bool` | `false` | <code>true</code> |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))` | `false` | <code></code> |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `false` | <code>".npmrc"</code> |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` | `false` | <code></code> |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `false` | <code>"pnpm"</code> |

### `pipe-node add`

Install node packages with the given package manager.

#### Flags

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `false` | <code>"pnpm"</code> |

**Packages**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_ADD_PACKAGES`<br />`$PACKAGES_NODE` | Install node packages before performing operations. | `string[]` | `true` | <code></code> |
| `$NODE_ADD_GLOBAL`<br />`$PACKAGES_NODE_GLOBAL` | Install node packages globally. | `bool` | `false` | <code>true</code> |
| `$NODE_ADD_SCRIPT_ARGS`<br />`$PACKAGES_NODE_SCRIPT_ARGS` | Script arguments to append to the install command. | `string` | `false` | <code></code> |
| `$NODE_ADD_CWD`<br />`$PACKAGES_NODE_CWD` | Working directory for the add operation. | `string` | `false` | <code>"."</code> |

### `pipe-node build`

#### Flags

**Build**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_BUILD_SCRIPT` | package.json script for the build operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` | <code>"build"</code> |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for the build operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` | <code></code> |
| `$NODE_BUILD_CWD` | Working directory for the build operation. | `string` | `false` | <code>"."</code> |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` | <code>false</code> |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(json([]struct{ match: RegExp, environment: string }))` | `false` | <code>"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"</code> |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | <code>true</code> |
| `$ENVIRONMENT_STRICT` | Fail when no environment is selected. | `bool` | `false` | <code>true</code> |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `false` | <code>"pnpm"</code> |

### `pipe-node run`

#### Flags

**Command**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_RUN_SCRIPT`<br />`$NODE_COMMAND_SCRIPT` | package.json script for the given command operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `false` | <code></code> |
| `$NODE_RUN_CWD`<br />`$NODE_COMMAND_CWD` | Working directory for the given command operation. | `string` | `false` | <code>"."</code> |

**Environment**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` | <code>false</code> |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br />Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(json([]struct{ match: RegExp, environment: string }))` | `false` | <code>"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"</code> |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `false` | <code>true</code> |
| `$ENVIRONMENT_STRICT` | Fail when no environment is selected. | `bool` | `false` | <code>true</code> |

**GIT**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br />`$BITBUCKET_BRANCH` | Source control branch. | `string` | `false` | <code></code> |
| `$CI_COMMIT_TAG`<br />`$BITBUCKET_TAG` | Source control tag. | `string` | `false` | <code></code> |

**Package Manager**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `false` | <code>"pnpm"</code> |
