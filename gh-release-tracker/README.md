# pipe-gh-release-tracker

Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

`pipe-gh-release-tracker [FLAGS]`

## Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

**Github**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GH_TOKEN`<br/>`$GITHUB_TOKEN` | Github token for the API requests. | `String` | `false` |  |
| `$GH_REPOSITORY`<br/>`$GITHUB_REPOSITORY` | Target repository to fetch the latest tag. | `String` | `true` |  |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `String` | `true` | .tags |
