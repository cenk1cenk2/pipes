# pipes-markdown-toc

Finds the markdown files and adds TOC to them.

`pipes-markdown-toc [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$MARKDOWN_TOC_PATTERNS`<br/>`$PLUGIN_MARKDOWN_TOC_PATTERNS` | Pattern for markdown. Supports multiple patterns with comma-separated values. | `StringSlice` | `false` | &#34;README.md&#34; |
| `$MARKDOWN_TOC_ARGUMENTS`<br/>`$PLUGIN_MARKDOWN_TOC_ARGUMENTS` | Pass in the arguments for markdown-toc. | `String` | `false` |  |

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |
