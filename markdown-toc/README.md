# pipe-markdown-toc

Finds the markdown files and adds TOC to them.

`pipe-markdown-toc [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$MARKDOWN_TOC_PATTERNS` | Pattern for markdown. | `string[]` | `false` | [README.md] |
| `$MARKDOWN_TOC_START_DEPTH` | Start depth for the elements in the given document. | `int` | `false` | 1 |
| `$MARKDOWN_TOC_END_DEPTH` | End depth for the elements in the given document. | `int` | `false` | 5 |
| `$MARKDOWN_TOC_INDENTATION` | Indentation for each element. | `int` | `false` | 2 |
| `$MARKDOWN_TOC_LIST_IDENTIFIER` | Identifier for each list element. | `string` | `false` | - |

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | [] |
