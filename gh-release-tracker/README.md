# pipes-gh-release-tracker

Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

`pipes-gh-release-tracker [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --gh.token, $GH_TOKEN, $GITHUB_TOKEN | Github token for the API requests. |  String  | false |  |
| --gh.repository, $GH_REPOSITORY, $GH_REPOSITORY | Target repository to fetch the latest tag. |  String  | true |  |
| --docker_image.tags_file, $TAGS_FILE | Read tags from a file. |  String  | true |  |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
