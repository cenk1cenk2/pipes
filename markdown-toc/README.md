# pipe-markdown-toc

Finds the markdown files and adds TOC to them.

`pipe-markdown-toc [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$MARKDOWN_TOC_PATTERNS`<br/>`$PLUGIN_MARKDOWN_TOC_PATTERNS` | Pattern for markdown. Supports multiple patterns with comma-separated values. | `StringSlice` | `false` | "README.md" |
| `$MARKDOWN_TOC_ARGUMENTS`<br/>`$PLUGIN_MARKDOWN_TOC_ARGUMENTS` | Pass in the arguments for markdown-toc. | `String` | `false` | --bullets='-' |

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |
