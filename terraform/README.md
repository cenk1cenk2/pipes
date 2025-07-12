# pipe-terraform

Running terraform inside the pipelines.

`pipe-terraform [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | [] |

## Commands

### `pipe-terraform install`

Install terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_INSTALL_RECONFIGURE` | Reconfigure flag for terraform init. | `bool` | `false` | false |
| `$TF_INSTALL_USE_LOCKFILE` | Use lockfile for terraform init. | `bool` | `false` | false |
| `$TF_INSTALL_ARGS` | Additional arguments for terraform init. | `string` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` |  |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | . |
| `$TF_WORKSPACES` | Workspaces that this command will be executed on. | `string[]` | `false` | [] |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` |  |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | default |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | false |
| `$TF_HTTP_ADDRESS`<br/>`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | POST |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` |  |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | DELETE |
| `$TF_HTTP_USERNAME`<br/>`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | gitlab-ci-token |
| `$TF_HTTP_PASSWORD`<br/>`$TF_PASSWORD`<br/>`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` |  |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | 5 |

### `pipe-terraform lint`

Lint terraform project with terraform.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LINT_FMT_CHECK_ENABLE` | Enable terraform fmt. | `bool` | `false` | true |
| `$TF_LINT_FMT_CHECK_ARGS` | Additional arguments for terraform fmt. | `string` | `false` |  |
| `$TF_LINT_VALIDATE_ENABLE` | Enable terraform validate. | `bool` | `false` | true |
| `$TF_LINT_VALIDATE_ARGS` | Additional arguments for terraform validate. | `string` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` |  |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | . |
| `$TF_WORKSPACES` | Workspaces that this command will be executed on. | `string[]` | `false` | [] |

### `pipe-terraform plan`

Plan terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform plan. | `string` | `false` | plan |
| `$TF_PLAN_ARGS` | Additional arguments for terraform plan. | `string` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` |  |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | . |
| `$TF_WORKSPACES` | Workspaces that this command will be executed on. | `string[]` | `false` | [] |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` |  |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | default |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | false |
| `$TF_HTTP_ADDRESS`<br/>`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | POST |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` |  |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | DELETE |
| `$TF_HTTP_USERNAME`<br/>`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | gitlab-ci-token |
| `$TF_HTTP_PASSWORD`<br/>`$TF_PASSWORD`<br/>`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` |  |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | 5 |

### `pipe-terraform apply`

Apply terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform apply. | `string` | `false` | plan |
| `$TF_APPLY_ARGS` | Additional arguments for terraform apply. | `string` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br/>`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` |  |

**Injected Variables**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_VAR_CI_JOB_ID`<br/>`$CI_JOB_ID` | Injected CI job-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_COMMIT_SHA`<br/>`$CI_COMMIT_SHA` | Injected CI commit-sha variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_JOB_STAGE`<br/>`$CI_JOB_STAGE` | Injected CI job-stage variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_ID`<br/>`$CI_PROJECT_ID` | Injected CI project-id variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAME`<br/>`$CI_PROJECT_NAME` | Injected CI project-name variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_NAMESPACE`<br/>`$CI_PROJECT_NAMESPACE` | Injected CI project-namespace variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_PATH`<br/>`$CI_PROJECT_PATH` | Injected CI project-path variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_PROJECT_URL`<br/>`$CI_PROJECT_URL` | Injected CI project-url variable to the deployment. | `string` | `false` |  |
| `$TF_VAR_CI_API_V4_URL`<br/>`$CI_API_V4_URL` | Injected CI api-url variable to the deployment. | `string` | `false` |  |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | . |
| `$TF_WORKSPACES` | Workspaces that this command will be executed on. | `string[]` | `false` | [] |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` |  |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | default |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | false |
| `$TF_HTTP_ADDRESS`<br/>`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` |  |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | POST |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` |  |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | DELETE |
| `$TF_HTTP_USERNAME`<br/>`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | gitlab-ci-token |
| `$TF_HTTP_PASSWORD`<br/>`$TF_PASSWORD`<br/>`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` |  |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | 5 |

### `pipe-terraform publish`

Publish terraform project.

#### Flags

**Module**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_NAME`<br/>`$CI_PROJECT_NAME` | Name for the module that will be published. | `string` | `true` |  |
| `$TF_MODULE_CWD`<br/>`$TF_ROOT` | Directory for the module that will be published. | `string` | `false` | . |
| `$TF_MODULE_SYSTEM` | Module system for the module that will be published. | `string` | `false` | local |

**Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_REGISTRY` | Registry of the module that will be published. | `string` | `false` | gitlab |

**Registry - Gitlab**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_API_V4_URL` | Gitlab API URL for publish call. | `string` | `false` |  |
| `$CI_PROJECT_ID` | Gitlab project id for publish call. | `string` | `false` |  |
| `$CI_JOB_TOKEN` | Gitlab API token for publish call. | `string` | `false` |  |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `string` | `false` | .tags |
