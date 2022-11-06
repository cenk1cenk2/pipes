# pipes-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipes-node [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |

## Commands

### `login`

Login to the given NPM registries.

`pipes-node login [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>json(slice({ username: string, password: string, registry?: string, useHttps?: boolean })) | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | &#34;.npmrc&#34; |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs.  | `String`<br/>enum(&#34;npm&#34;, &#34;yarn&#34;) | `false` |  |

### `install`

Install node.js dependencies with the given package manager.

`pipes-node install [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Install

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_INSTALL_CWD` | Install CWD for nodejs. | `String` | `false` |  |
| `$NODE_INSTALL_USE_LOCK_FILE` | Whether to use lock file or not. | `Bool` | `false` | false |
| `$NODE_INSTALL_ARGS` | Arguments for appending to installation. | `String` | `false` |  |

##### Login

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NPM_LOGIN` | npm registries to login to.  | `String`<br/>json(slice({ username: string, password: string, registry?: string, useHttps?: boolean })) | `false` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `StringSlice` | `false` | &#34;.npmrc&#34; |
| `$NPM_NPMRC` | Pass direct contents of the NPMRC file. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs.  | `String`<br/>enum(&#34;npm&#34;, &#34;yarn&#34;) | `false` |  |

### `build`

`pipes-node build [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Build

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_BUILD_SCRIPT` | package.json script for building operation.  | `String`<br/>format(Template(struct{ Environment: string, EnvVars: map[string]string })) | `false` |  |
| `$NODE_BUILD_SCRIPT_ARGS` | package.json script arguments for building operation.  | `String`<br/>format(Template(struct{ Environment: string, EnvVars: map[string]string })) | `false` |  |
| `$NODE_BUILD_CWD` | Working directory for build operation. | `String` | `false` |  |

##### GIT

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `String` | `false` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs.  | `String`<br/>enum(&#34;npm&#34;, &#34;yarn&#34;) | `false` |  |

### `run`

`pipes-node run [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Command

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_COMMAND_SCRIPT` | package.json script for given command operation.  | `String`<br/>format(Template(struct{ Environment: string, EnvVars: map[string]string })) | `false` |  |
| `$NODE_COMMAND_SCRIPT_ARGS` | package.json script arguments for given command operation.  | `String`<br/>format(Template(struct{ Environment: string, EnvVars: map[string]string })) | `false` |  |
| `$NODE_COMMAND_CWD` | Working directory for the given command operation. | `String` | `false` |  |

##### Package Manager

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NODE_PACKAGE_MANAGER` | Preferred Package manager for nodejs.  | `String`<br/>enum(&#34;npm&#34;, &#34;yarn&#34;) | `false` |  |
