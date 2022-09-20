# pipes-gitlab-artifacts

Downloads gitlab artifacts from the API for creating downstream pipelines.

`pipes-gitlab-artifacts [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags
| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| --debug, $$DEBUG | Enable debugging for the application. |  Bool  | false | false |
| --log-level, $$LOG_LEVEL | Define the log level for the application. (format: enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;)) |  String  | false | &#34;info&#34; |
| --gl.api_url, $$CI_API_V4_URL | Gitlab API URL of the instance. |  String  | true |  |
| --gl.token, $$GL_TOKEN | Token for gitlab api authentication. |  String  | true |  |
| --gl.job_token, $$CI_JOB_TOKEN | Job token coming from the build job. |  String  | false |  |
| --gl_pipeline.project_id, $$CI_PROJECT_ID | Parent project id. |  String  | true |  |
| --gl_pipeline.parent_pipeline_id, $$PARENT_PIPELINE_ID | Pipeline id of the parent pipeline. |  String  | true |  |
| --gl_pipeline.download_artifacts, $$PARENT_DOWNLOAD_ARTIFACTS | Names of the jobs that yield artifacts from the parent job. |  String  | true |  |
| --help, -h | show help |  Bool  | false | false |
| --version, -v | print the version |  Bool  | false | false |

# Commands

## `help` , `h`

`Shows a list of commands or help for one command`
