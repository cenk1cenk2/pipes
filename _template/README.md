# pipes-template

template-cli

`pipes-template [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --default.flag, $DEFAULT_FLAG | Some default flag. |  String  | false |  |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
