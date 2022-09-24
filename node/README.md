# pipes-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipes-node [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DEBUG` | Enable debugging for the application. | `Bool` | `false` | false |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |

## Commands

### `login` 

Login to the given NPM registries.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String` | `false` | &#34;yarn&#34; |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>format(json({ username: string, password: string, registry?: string, useHttps?: boolean }[])) | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | [.npmrc] |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

### `install` 

Install node.js dependencies with the given package manager.

`pipes-node install`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String` | `false` | &#34;yarn&#34; |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>format(json({ username: string, password: string, registry?: string, useHttps?: boolean }[])) | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | [.npmrc] |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |
| `$NODE_INSTALL_CWD` | Install CWD for nodejs. | `String` | `false` | &#34;.&#34; |
| `$NODE_INSTALL_USE_LOCK_FILE` | Whether to use lock file or not. | `Bool` | `false` | true |
| `$NODE_INSTALL_ARGS` | Arguments for appending to installation. | `String` | `false` |  |

### `build` 

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String` | `false` | &#34;yarn&#34; |
| `$CI_COMMIT_REF_NAME` | Source control management branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG` | Source control management tag. | `String` | `false` |  |
| `$NODE_BUILD_SCRIPT` | package.json script for building operation. | `String` | `false` | &#34;build&#34; |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for building operation. | `String` | `false` |  |
| `$NODE_BUILD_CWD` | Working directory for build operation. | `String` | `false` | &#34;.&#34; |
| `$NODE_BUILD_ENVIRONMENT_FILES` | Yaml files to inject to build. | `StringSlice` | `false` | [] |
| `$NODE_BUILD_ENVIRONMENT_CONDITIONS` | Tagging regex patterns to match.  | `String`<br/>format(json({ [name: string]: RegExp })) | `false` | &#34;{ \&#34;production\&#34;: \&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*$\&#34;, \&#34;stage\&#34;: \&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*-.*$\&#34; }&#34; |
| `$NODE_BUILD_ENVIRONMENT_FALLBACK` | Fallback, if it does not match any conditions. Defaults to current branch name. | `String` | `false` | &#34;develop&#34; |

### `run` 

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs. | `String` | `false` | &#34;yarn&#34; |
| `$NODE_COMMAND_SCRIPT` | package.json script for given command operation. | `String` | `false` |  |
| `$NODE_COMMAND_SCRIPT_ARGS` | package.json script arguments for given command operation. | `String` | `false` |  |
| `$NODE_COMMAND_CWD` | Working directory for the given command operation. | `String` | `false` | &#34;.&#34; |

### `help` , `h`

`Shows a list of commands or help for one command`
