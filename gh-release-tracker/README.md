# pipes-gh-release-tracker

Tracks the given Github repositories latest tags and generates tag file accordingly for further processing.

`pipes-gh-release-tracker [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GH_TOKEN`<br/>`$GITHUB_TOKEN` | Github token for the API requests. | `String` | `false` |  |
| `$GH_REPOSITORY`<br/>`$GH_REPOSITORY` | Target repository to fetch the latest tag. | `String` | `true` |  |
| `$TAGS_FILE` | Read tags from a file. | `String` | `true` |  |

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |
