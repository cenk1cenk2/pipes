# pipes-gitlab-artifacts

Downloads gitlab artifacts from the API for creating downstream pipelines.

`pipes-gitlab-artifacts [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` |  |

### Gitlab

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_API_V4_URL` | Gitlab API URL of the instance. | `String` | `true` |  |
| `$GL_TOKEN` | Token for Gitlab API authentication. | `String` | `true` |  |
| `$CI_JOB_TOKEN` | Job token coming from the build job. | `String` | `false` |  |

### Gitlab Pipeline

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_PROJECT_ID` | Parent project id. | `String` | `true` |  |
| `$PARENT_PIPELINE_ID` | Pipeline id of the parent pipeline. | `String` | `true` |  |
| `$PARENT_DOWNLOAD_ARTIFACTS` | Names of the jobs that yield artifacts from the parent job. | `String` | `true` |  |
