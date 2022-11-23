# pipe-gitlab-artifacts

Downloads gitlab artifacts from the API for creating downstream pipelines.

`pipe-gitlab-artifacts [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

### Gitlab

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_API_V4_URL` | Gitlab API URL of the instance. | `String` | `true` |  |
| `$GL_TOKEN`<br/>`$GITLAB_TOKEN` | Token for Gitlab API authentication. | `String` | `true` |  |
| `$CI_JOB_TOKEN` | Job token coming from the build job. | `String` | `false` |  |

### Gitlab Pipeline

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_PROJECT_ID` | Parent project id. | `String` | `true` |  |
| `$PARENT_PIPELINE_ID` | Pipeline id of the parent pipeline. | `String` | `true` |  |
| `$PARENT_DOWNLOAD_ARTIFACTS` | Names of the jobs that yield artifacts from the parent job. | `String`<br/>`multiple("|")` | `true` |  |
