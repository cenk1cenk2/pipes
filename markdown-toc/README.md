# pipe-markdown-toc

Finds the markdown files and adds TOC to them.

`pipe-markdown-toc [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$MARKDOWN_TOC_PATTERNS` | Pattern for markdown. | `StringSlice` | `false` | "README.md" |
| `$MARKDOWN_TOC_START_DEPTH` | Start depth for the elements in the given document. | `Int` | `false` | 0 |
| `$MARKDOWN_TOC_END_DEPTH` | End depth for the elements in the given document. | `Int` | `false` | 0 |
| `$MARKDOWN_TOC_INDENTATION` | Indentation for each element. | `Int` | `false` | 0 |
| `$MARKDOWN_TOC_LIST_IDENTIFIER` | Identifier for each list element. | `String` | `false` | - |

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |
