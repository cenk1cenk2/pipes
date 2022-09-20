# pipes-semantic-release

Releases applications through semantic-release library.

`pipes-semantic-release [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --npm.login, $NPM_LOGIN | npm registries to login to. (format: json({username: string, password: string, registry?: string, useHttps?: boolean}[])) |  String  | false |  |
| --npm.npmrc_file, $NPM_NPMRC_FILE | .npmrc file to use. |  StringSlice  | false | [.npmrc] |
| --npm.npmrc, $NPM_NPMRC | Pass direct contents of the NPMRC file. |  String  | false |  |
| --packages.apk, $ADD_APKS | APK applications to install before running semantic-release. |  StringSlice  | false | [] |
| --packages.node, $ADD_MODULES | Node packages to install before running semantic-release. |  StringSlice  | false | [] |
| --semantic_release.dry_run, $DRY_RUN | Node packages to install before running semantic-release. |  Bool  | false | false |
| --semantic_release.run_multi, $RUN_MULTI | Uses @qiwi/multi-semantic-release package to do a workspace release. |  Bool  | false | false |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
