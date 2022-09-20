# pipes-node

Pipe for installing node.js dependencies and building node.js applications on CI/CD.

`pipes-node [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `login` 

Login to the given NPM registries.

### Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --node.package_manager, $NODE_PACKAGE_MANAGER | Preferred Package manager for nodejs. |  String  | false | &#34;yarn&#34; |
| --npm.login, $NPM_LOGIN | npm registries to login to. (format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])) |  String  | false |  |
| --npm.npmrc_file, $NPM_NPMRC_FILE | .npmrc file to use. |  StringSlice  | false | [.npmrc] |
| --npm.npmrc, $NPM_NPMRC | Pass direct contents of the NPMRC file. |  String  | false |  |

## `install` 

Install node.js dependencies with the given package manager.

`pipes-node install`

### Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --node.package_manager, $NODE_PACKAGE_MANAGER | Preferred Package manager for nodejs. |  String  | false | &#34;yarn&#34; |
| --npm.login, $NPM_LOGIN | npm registries to login to. (format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])) |  String  | false |  |
| --npm.npmrc_file, $NPM_NPMRC_FILE | .npmrc file to use. |  StringSlice  | false | [.npmrc] |
| --npm.npmrc, $NPM_NPMRC | Pass direct contents of the NPMRC file. |  String  | false |  |
| --node.install_cwd, $NODE_INSTALL_CWD | Install CWD for nodejs. |  String  | false | &#34;.&#34; |
| --node.use_lock_file, $NODE_INSTALL_USE_LOCK_FILE | Whether to use lock file or not. |  Bool  | false | true |
| --node.install_args, $NODE_INSTALL_ARGS | Arguments for appending to installation. |  String  | false |  |

## `build` 

### Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --node.package_manager, $NODE_PACKAGE_MANAGER | Preferred Package manager for nodejs. |  String  | false | &#34;yarn&#34; |
| --git.branch, $CI_COMMIT_REF_NAME | Source control management branch. |  String  | false |  |
| --git.tag, $CI_COMMIT_TAG | Source control management tag. |  String  | false |  |
| --node.build_script, $NODE_BUILD_SCRIPT | package.json script for building operation. |  String  | false | &#34;build&#34; |
| --node.build_script_args, $NODE_BUILD_SCRIPT_ARGS | package.json script arguments for building operation. |  String  | false |  |
| --node.build_cwd, $NODE_BUILD_CWD | Working directory for build operation. |  String  | false | &#34;.&#34; |
| --node.build_environment_files, $NODE_BUILD_ENVIRONMENT_FILES | Yaml files to inject to build. |  StringSlice  | false | [] |
| --node.build_environment_conditions, $NODE_BUILD_ENVIRONMENT_CONDITIONS | Tagging regex patterns to match. json({ [name: string]: RegExp }) |  String  | false | &#34;{ \&#34;production\&#34;: \&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*$\&#34;, \&#34;stage\&#34;: \&#34;^v\\\\d*\\\\.\\\\d*\\\\.\\\\d*-.*$\&#34; }&#34; |
| --node.build_environment_fallback, $NODE_BUILD_ENVIRONMENT_FALLBACK | Fallback, if it does not match any conditions. Defaults to current branch name. |  String  | false | &#34;develop&#34; |

## `run` 

### Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --node.package_manager, $NODE_PACKAGE_MANAGER | Preferred Package manager for nodejs. |  String  | false | &#34;yarn&#34; |
| --node.command_script, $NODE_COMMAND_SCRIPT | package.json script for given command operation. |  String  | false |  |
| --node.command_script_args, $NODE_COMMAND_SCRIPT_ARGS | package.json script arguments for given command operation. |  String  | false |  |
| --node.command_cwd, $NODE_COMMAND_CWD | Working directory for the given command operation. |  String  | false | &#34;.&#34; |

## `help` , `h`

`Shows a list of commands or help for one command`
