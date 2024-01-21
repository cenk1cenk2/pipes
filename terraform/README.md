# pipe-terraform

Running terraform inside the pipelines.

`pipe-terraform [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

## Commands

### `install`

Install terraform project.

`pipe-terraform install [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `String`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_INSTALL_RECONFIGURE` | Reconfigure flag for terraform init. | `Bool` | `false` | false |
| `$TF_INSTALL_USE_LOCKFILE` | Use lockfile for terraform init. | `Bool` | `false` | false |
| `$TF_INSTALL_ARGS` | Additional arguments for terraform init. | `String` | `false` |  |

##### Config

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `String`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

##### Injected Variables

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `String` | `false` |  |

##### Project

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `String` | `false` | . |

### `lint`

Lint terraform project with terraform.

`pipe-terraform lint [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LINT_FMT_CHECK_ENABLE` | Enable terraform fmt. | `Bool` | `false` | true |
| `$TF_LINT_FMT_CHECK_ARGS` | Additional arguments for terraform fmt. | `String` | `false` |  |
| `$TF_LINT_VALIDATE_ENABLE` | Enable terraform validate. | `Bool` | `false` | true |
| `$TF_LINT_VALIDATE_ARGS` | Additional arguments for terraform validate. | `String` | `false` |  |

##### Config

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `String`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

##### Injected Variables

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `String` | `false` |  |

##### Project

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `String` | `false` | . |

### `plan`

Plan terraform project.

`pipe-terraform plan [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `String`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform plan. | `String` | `false` | plan |
| `$TF_PLAN_ARGS` | Additional arguments for terraform plan. | `String` | `false` |  |

##### Config

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `String`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

##### Injected Variables

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `String` | `false` |  |

##### Project

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `String` | `false` | . |

##### State

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `String`<br/>`enum("gitlab-http")` | `true` |  |
| `$TF_STATE_NAME` | Terraform state name. | `String` | `true` |  |
| `$TF_HTTP_ADDRESS`<br/>`$TF_ADDRESS` | State configuration for terraform: http-address | `String` | `false` |  |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `String` | `false` |  |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `String` | `false` | POST |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `String` | `false` |  |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `String` | `false` | DELETE |
| `$TF_HTTP_USERNAME`<br/>`$TF_USERNAME` | State configuration for terraform: http-username | `String` | `false` | gitlab-ci-token |
| `$TF_HTTP_PASSWORD`<br/>`$TF_PASSWORD`<br/>`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `String` | `false` |  |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `String` | `false` | 5 |

### `apply`

Apply terraform project.

`pipe-terraform apply [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `String`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform apply. | `String` | `false` | plan |
| `$TF_APPLY_ARGS` | Additional arguments for terraform apply. | `String` | `false` |  |

##### Config

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `String`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

##### Injected Variables

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `String` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `String` | `false` |  |

##### Project

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `String` | `false` | . |

##### State

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `String`<br/>`enum("gitlab-http")` | `true` |  |
| `$TF_STATE_NAME` | Terraform state name. | `String` | `true` |  |
| `$TF_HTTP_ADDRESS`<br/>`$TF_ADDRESS` | State configuration for terraform: http-address | `String` | `false` |  |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `String` | `false` |  |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `String` | `false` | POST |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `String` | `false` |  |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `String` | `false` | DELETE |
| `$TF_HTTP_USERNAME`<br/>`$TF_USERNAME` | State configuration for terraform: http-username | `String` | `false` | gitlab-ci-token |
| `$TF_HTTP_PASSWORD`<br/>`$TF_PASSWORD`<br/>`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `String` | `false` |  |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `String` | `false` | 5 |

### `publish`

Publish terraform project.

`pipe-terraform publish [GLOBAL FLAGS] [FLAGS]`

#### Flags

##### Module

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_NAME`<br/>`$CI_PROJECT_NAME` | Name for the module that will be published. | `String` | `true` |  |
| `$TF_MODULE_CWD`<br/>`$TF_ROOT` | Directory for the module that will be published. | `String` | `false` | . |
| `$TF_MODULE_SYSTEM` | Module system for the module that will be published. | `String` | `false` | local |

##### Registry

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_REGISTRY` | Registry of the module that will be published. | `String` | `false` | gitlab |

##### Registry - Gitlab

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_API_V4_URL` | Gitlab API URL for publish call. | `String` | `false` |  |
| `$CI_PROJECT_ID` | Gitlab project id for publish call. | `String` | `false` |  |
| `$CI_JOB_TOKEN` | Gitlab API token for publish call. | `String` | `false` |  |

##### Tags File

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `String` | `true` |  |
