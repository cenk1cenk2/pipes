# pipe-terraform

Terraform actions for pipelines.

`pipe-terraform [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-terraform install`

Install terraform project.

#### Flags

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOG_LEVEL`<br />`$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`format(enum("trace", "debug", "info", "warn", "error"))` | `false` | <code></code> |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CI_API_URL`<br />`$TF_VAR_CI_API_V4_URL`<br />`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` | <code></code> |
| `$TERRAFORM_CI_PROJECT_ID`<br />`$TF_VAR_CI_PROJECT_ID`<br />`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` | <code></code> |

**Install**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_INSTALL_RECONFIGURE`<br />`$TF_INSTALL_RECONFIGURE` | Reconfigure flag for terraform init. | `bool` | `false` | <code>false</code> |
| `$TERRAFORM_INSTALL_USE_LOCKFILE`<br />`$TF_INSTALL_USE_LOCKFILE` | Use lockfile for terraform init. | `bool` | `false` | <code>false</code> |
| `$TERRAFORM_INSTALL_ARGS`<br />`$TF_INSTALL_ARGS` | Additional arguments for terraform init. | `string` | `false` | <code></code> |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOGIN_REGISTRY_CREDENTIALS`<br />`$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`format(json([]struct{ registry: string, token: string }))` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CWD`<br />`$TF_ROOT` | Working directory for terraform commands. | `string` | `false` | <code>"."</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_STATE_TYPE`<br />`$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`format(enum("gitlab-http"))` | `false` | <code></code> |
| `$TERRAFORM_STATE_NAME`<br />`$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>"default"</code> |
| `$TERRAFORM_STATE_STRICT`<br />`$TF_STATE_STRICT` | Fail when no Terraform state type is configured. | `bool` | `false` | <code>false</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS`<br />`$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | HTTP address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS`<br />`$TF_HTTP_LOCK_ADDRESS` | HTTP lock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD`<br />`$TF_HTTP_LOCK_METHOD` | HTTP lock method for the GitLab HTTP state backend. | `string` | `false` | <code>"POST"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS`<br />`$TF_HTTP_UNLOCK_ADDRESS` | HTTP unlock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD`<br />`$TF_HTTP_UNLOCK_METHOD` | HTTP unlock method for the GitLab HTTP state backend. | `string` | `false` | <code>"DELETE"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME`<br />`$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | HTTP username for the GitLab HTTP state backend. | `string` | `false` | <code>"gitlab-ci-token"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD`<br />`$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | HTTP password for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN`<br />`$TF_HTTP_RETRY_WAIT_MIN` | Minimum time to wait between HTTP retries for the GitLab HTTP state backend, in seconds. | `string` | `false` | <code>"5"</code> |

### `pipe-terraform lint`

Lint terraform project with terraform.

#### Flags

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOG_LEVEL`<br />`$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`format(enum("trace", "debug", "info", "warn", "error"))` | `false` | <code></code> |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CI_API_URL`<br />`$TF_VAR_CI_API_V4_URL`<br />`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` | <code></code> |
| `$TERRAFORM_CI_PROJECT_ID`<br />`$TF_VAR_CI_PROJECT_ID`<br />`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` | <code></code> |

**Lint**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LINT_FORMAT_CHECK_ENABLE`<br />`$TF_LINT_FMT_CHECK_ENABLE` | Enable terraform fmt. | `bool` | `false` | <code>true</code> |
| `$TERRAFORM_LINT_FORMAT_CHECK_ARGS`<br />`$TF_LINT_FMT_CHECK_ARGS` | Additional arguments for terraform fmt. | `string` | `false` | <code></code> |
| `$TERRAFORM_LINT_VALIDATE_ENABLE`<br />`$TF_LINT_VALIDATE_ENABLE` | Enable terraform validate. | `bool` | `false` | <code>true</code> |
| `$TERRAFORM_LINT_VALIDATE_ARGS`<br />`$TF_LINT_VALIDATE_ARGS` | Additional arguments for terraform validate. | `string` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CWD`<br />`$TF_ROOT` | Working directory for terraform commands. | `string` | `false` | <code>"."</code> |

### `pipe-terraform plan`

Plan terraform project.

#### Flags

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOG_LEVEL`<br />`$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`format(enum("trace", "debug", "info", "warn", "error"))` | `false` | <code></code> |

**GitLab Merge Request Report**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` | <code>false</code> |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` | `false` | <code></code> |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `false` | <code>0</code> |
| `$GITLAB_MR_REPORT_IDENTIFIER` | Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report. | `string` | `false` | <code></code> |

**GitLab Pipeline**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_JOB_NAME` | GitLab CI job name to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_JOB_URL` | GitLab CI job URL to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHA` | Git commit sha to include in the plan report metadata. | `string` | `false` | <code></code> |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in the plan report metadata. | `string` | `false` | <code></code> |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CI_API_URL`<br />`$TF_VAR_CI_API_V4_URL`<br />`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` | <code></code> |
| `$TERRAFORM_CI_PROJECT_ID`<br />`$TF_VAR_CI_PROJECT_ID`<br />`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` | <code></code> |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOGIN_REGISTRY_CREDENTIALS`<br />`$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`format(json([]struct{ registry: string, token: string }))` | `false` | <code></code> |

**Plan**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_PLAN_OUTPUT`<br />`$TF_PLAN_CACHE`<br />`$TF_APPLY_OUTPUT`<br />`$TF_PLAN_OUTPUT` | Output file for terraform plan. | `string` | `false` | <code>"plan"</code> |
| `$TERRAFORM_PLAN_ARGS`<br />`$TF_PLAN_ARGS` | Additional arguments for terraform plan. | `string` | `false` | <code></code> |
| `$TERRAFORM_PLAN_PREVIEW_FOR_MERGE_REQUESTS`<br />`$TF_PLAN_PREVIEW_FOR_MRS` | Run merge request terraform plans as previews without state locking. | `bool` | `false` | <code>true</code> |
| `$TERRAFORM_PLAN_PIPELINE_SOURCE`<br />`$CI_PIPELINE_SOURCE` | GitLab CI pipeline source used to detect merge request pipelines. | `string` | `false` | <code></code> |
| `$TERRAFORM_PLAN_RETRY_TRIES`<br />`$TF_PLAN_RETRY_TRIES` | Number of retries for terraform plan command. | `uint` | `false` | <code>5</code> |
| `$TERRAFORM_PLAN_RETRY_DELAY`<br />`$TF_PLAN_RETRY_DELAY` | Delay between retries for terraform plan command. | `duration` | `false` | <code>1m0s</code> |
| `$TERRAFORM_PLAN_SUMMARY_OUTPUT`<br />`$TERRAFORM_SUMMARY_OUTPUT` | Output file for terraform plan summary. Leave empty to skip summary generation. | `string` | `false` | <code>"terraform-summary.json"</code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CWD`<br />`$TF_ROOT` | Working directory for terraform commands. | `string` | `false` | <code>"."</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_STATE_TYPE`<br />`$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`format(enum("gitlab-http"))` | `false` | <code></code> |
| `$TERRAFORM_STATE_NAME`<br />`$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>"default"</code> |
| `$TERRAFORM_STATE_STRICT`<br />`$TF_STATE_STRICT` | Fail when no Terraform state type is configured. | `bool` | `false` | <code>false</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS`<br />`$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | HTTP address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS`<br />`$TF_HTTP_LOCK_ADDRESS` | HTTP lock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD`<br />`$TF_HTTP_LOCK_METHOD` | HTTP lock method for the GitLab HTTP state backend. | `string` | `false` | <code>"POST"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS`<br />`$TF_HTTP_UNLOCK_ADDRESS` | HTTP unlock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD`<br />`$TF_HTTP_UNLOCK_METHOD` | HTTP unlock method for the GitLab HTTP state backend. | `string` | `false` | <code>"DELETE"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME`<br />`$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | HTTP username for the GitLab HTTP state backend. | `string` | `false` | <code>"gitlab-ci-token"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD`<br />`$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | HTTP password for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN`<br />`$TF_HTTP_RETRY_WAIT_MIN` | Minimum time to wait between HTTP retries for the GitLab HTTP state backend, in seconds. | `string` | `false` | <code>"5"</code> |

### `pipe-terraform apply`

Apply terraform project.

#### Flags

**Apply**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_APPLY_OUTPUT`<br />`$TF_PLAN_CACHE`<br />`$TF_APPLY_OUTPUT`<br />`$TF_PLAN_OUTPUT` | Output file for terraform apply. | `string` | `false` | <code>"plan"</code> |
| `$TERRAFORM_APPLY_ARGS`<br />`$TF_APPLY_ARGS` | Additional arguments for terraform apply. | `string` | `false` | <code></code> |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOG_LEVEL`<br />`$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`format(enum("trace", "debug", "info", "warn", "error"))` | `false` | <code></code> |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CI_API_URL`<br />`$TF_VAR_CI_API_V4_URL`<br />`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` | <code></code> |
| `$TERRAFORM_CI_PROJECT_ID`<br />`$TF_VAR_CI_PROJECT_ID`<br />`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` | <code></code> |

**Login**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_LOGIN_REGISTRY_CREDENTIALS`<br />`$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`format(json([]struct{ registry: string, token: string }))` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_CWD`<br />`$TF_ROOT` | Working directory for terraform commands. | `string` | `false` | <code>"."</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_STATE_TYPE`<br />`$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`format(enum("gitlab-http"))` | `false` | <code></code> |
| `$TERRAFORM_STATE_NAME`<br />`$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>"default"</code> |
| `$TERRAFORM_STATE_STRICT`<br />`$TF_STATE_STRICT` | Fail when no Terraform state type is configured. | `bool` | `false` | <code>false</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS`<br />`$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | HTTP address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS`<br />`$TF_HTTP_LOCK_ADDRESS` | HTTP lock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD`<br />`$TF_HTTP_LOCK_METHOD` | HTTP lock method for the GitLab HTTP state backend. | `string` | `false` | <code>"POST"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS`<br />`$TF_HTTP_UNLOCK_ADDRESS` | HTTP unlock address for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD`<br />`$TF_HTTP_UNLOCK_METHOD` | HTTP unlock method for the GitLab HTTP state backend. | `string` | `false` | <code>"DELETE"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME`<br />`$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | HTTP username for the GitLab HTTP state backend. | `string` | `false` | <code>"gitlab-ci-token"</code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD`<br />`$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | HTTP password for the GitLab HTTP state backend. | `string` | `false` | <code></code> |
| `$TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN`<br />`$TF_HTTP_RETRY_WAIT_MIN` | Minimum time to wait between HTTP retries for the GitLab HTTP state backend, in seconds. | `string` | `false` | <code>"5"</code> |

### `pipe-terraform publish`

Publish terraform project.

#### Flags

**Module**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_PUBLISH_MODULE_NAME`<br />`$TF_MODULE_NAME`<br />`$CI_PROJECT_NAME` | Name for the module that will be published. | `string` | `true` | <code></code> |
| `$TERRAFORM_PUBLISH_MODULE_CWD`<br />`$TF_MODULE_CWD`<br />`$TF_ROOT` | Directory for the module that will be published. | `string` | `false` | <code>"."</code> |
| `$TERRAFORM_PUBLISH_MODULE_SYSTEM`<br />`$TF_MODULE_SYSTEM` | Module system for the module that will be published. | `string` | `false` | <code>"local"</code> |

**Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_PUBLISH_REGISTRY_NAME`<br />`$TF_MODULE_REGISTRY` | Registry of the module that will be published. | `string`<br/>`format(enum("gitlab"))` | `false` | <code>"gitlab"</code> |

**Registry - GitLab**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL`<br />`$CI_API_V4_URL` | GitLab API URL for the publish call. | `string` | `false` | <code></code> |
| `$TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID`<br />`$CI_PROJECT_ID` | GitLab project id for the publish call. | `string` | `false` | <code></code> |
| `$TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN`<br />`$CI_JOB_TOKEN` | GitLab API token for the publish call. | `string` | `false` | <code></code> |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a comma separated file. | `string` | `false` | <code>".tags"</code> |
