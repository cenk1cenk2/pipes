# pipe-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipe-node [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-node login`](#pipe-node-login)
- [`pipe-node install`](#pipe-node-install)
- [`pipe-node add`](#pipe-node-add)
- [`pipe-node build`](#pipe-node-build)
- [`pipe-node run`](#pipe-node-run)

### `pipe-node login`

Login to the given NPM registries.

`pipe-node login [FLAGS]`

#### Flags

**Login**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `".npmrc"` |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` |  |

**Package Manager**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

### `pipe-node install`

Install node.js dependencies with the given package manager.

`pipe-node install [FLAGS]`

#### Flags

**Install**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_INSTALL_CWD` | Working directory for the install operation. | `string` | `"."` |
| `$NODE_INSTALL_USE_LOCK_FILE` | Use the lockfile while installing the packages. | `bool` | `true` |
| `$NODE_INSTALL_ARGS` | Arguments to append to the install command. | `string` |  |
| `$NODE_INSTALL_CACHE` | Enable caching for the package manager. | `bool` | `true` |

<details>
<summary>Login</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `".npmrc"` |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` |  |

</details>

<details>
<summary>Package Manager</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

</details>

### `pipe-node add`

Install node packages with the given package manager.

`pipe-node add [FLAGS]`

#### Flags

<details>
<summary>Package Manager</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

</details>

**Packages**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$NODE_ADD_PACKAGES`**\* | Install node packages before performing operations. | `string[]` |  |
| `$NODE_ADD_GLOBAL` | Install node packages globally. | `bool` | `true` |
| `$NODE_ADD_SCRIPT_ARGS` | Script arguments to append to the install command. | `string` |  |
| `$NODE_ADD_CWD` | Working directory for the add operation. | `string` | `"."` |

\* required

### `pipe-node build`

`pipe-node build [FLAGS]`

#### Flags

**Build**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_BUILD_SCRIPT` | package.json script for the build operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` | `"build"` |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for the build operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` |  |
| `$NODE_BUILD_CWD` | Working directory for the build operation. | `string` | `"."` |

**Environment**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br/>Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(json([]struct{ match: RegExp, environment: string }))` | `"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"` |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `true` |
| `$ENVIRONMENT_STRICT` | Fail when no environment is selected. | `bool` | `true` |

**GIT**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` |  |

<details>
<summary>Package Manager</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

</details>

### `pipe-node run`

`pipe-node run [FLAGS] Arguments appended to the script.`

#### Arguments

| Argument | Usage | Type | Values |
| --- | --- | --- | --- |
| `arg` | `Arguments appended to the script.` | `string[]` | `0..*` |

#### Flags

**Command**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_RUN_SCRIPT`<br/>`$NODE_COMMAND_SCRIPT` | package.json script for the given command operation. | `string`<br/>`format(Template(struct{ Environment: string, EnvVars: map[string]string }))` |  |
| `$NODE_RUN_CWD`<br/>`$NODE_COMMAND_CWD` | Working directory for the given command operation. | `string` | `"."` |

<details>
<summary>Environment</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br/>Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(json([]struct{ match: RegExp, environment: string }))` | `"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"` |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `true` |
| `$ENVIRONMENT_STRICT` | Fail when no environment is selected. | `bool` | `true` |

</details>

<details>
<summary>GIT</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` |  |

</details>

<details>
<summary>Package Manager</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

</details>
