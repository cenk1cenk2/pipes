# pipes-gh-release-tracker

Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

`pipes-gh-release-tracker [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |

### Github

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GH_TOKEN`<br/>`$GITHUB_TOKEN` | Github token for the API requests. | `String` | `false` |  |
| `$GH_REPOSITORY`<br/>`$GH_REPOSITORY` | Target repository to fetch the latest tag. | `String` | `true` |  |

### Tags File

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `String` | `true` | .tags |
