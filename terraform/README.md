# pipe-terraform

Running terraform inside the pipelines.

`pipe-terraform [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>info</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code>[]</code> |

## Commands

### `pipe-terraform install`

Install terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` | <code></code> |
| `$TF_INSTALL_RECONFIGURE` | Reconfigure flag for terraform init. | `bool` | `false` | <code>false</code> |
| `$TF_INSTALL_USE_LOCKFILE` | Use lockfile for terraform init. | `bool` | `false` | <code>false</code> |
| `$TF_INSTALL_ARGS` | Additional arguments for terraform init. | `string` | `false` | <code></code> |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | <code>.</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` | <code></code> |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>default</code> |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | <code>false</code> |
| `$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | <code>POST</code> |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | <code>DELETE</code> |
| `$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | <code>gitlab-ci-token</code> |
| `$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` | <code></code> |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | <code>5</code> |

### `pipe-terraform lint`

Lint terraform project with terraform.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LINT_FMT_CHECK_ENABLE` | Enable terraform fmt. | `bool` | `false` | <code>true</code> |
| `$TF_LINT_FMT_CHECK_ARGS` | Additional arguments for terraform fmt. | `string` | `false` | <code></code> |
| `$TF_LINT_VALIDATE_ENABLE` | Enable terraform validate. | `bool` | `false` | <code>true</code> |
| `$TF_LINT_VALIDATE_ARGS` | Additional arguments for terraform validate. | `string` | `false` | <code></code> |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | <code>.</code> |

### `pipe-terraform plan`

Plan terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` | <code></code> |
| `$TF_PLAN_CACHE`<br />`$TF_APPLY_OUTPUT`<br />`$TF_PLAN_OUTPUT` | Output file for terraform plan. | `string` | `false` | <code>plan</code> |
| `$TF_PLAN_ARGS` | Additional arguments for terraform plan. | `string` | `false` | <code></code> |
| `$TF_PLAN_RETRY_TRIES` | Number of retries for terraform plan command. | `uint` | `false` | <code>5</code> |
| `$TF_PLAN_RETRY_DELAY` | Delay between retries for terraform plan command. | `duration` | `false` | <code>1m0s</code> |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | <code>.</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` | <code></code> |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>default</code> |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | <code>false</code> |
| `$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | <code>POST</code> |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | <code>DELETE</code> |
| `$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | <code>gitlab-ci-token</code> |
| `$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` | <code></code> |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | <code>5</code> |

### `pipe-terraform apply`

Apply terraform project.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `string`<br/>`json([]struct { registry: string, token: string })` | `false` | <code></code> |
| `$TF_PLAN_CACHE`<br />`$TF_APPLY_OUTPUT`<br />`$TF_PLAN_OUTPUT` | Output file for terraform apply. | `string` | `false` | <code>plan</code> |
| `$TF_APPLY_ARGS` | Additional arguments for terraform apply. | `string` | `false` | <code></code> |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_LOG_LEVEL`<br />`$TF_LOG` | Terraform log level. | `string`<br/>`enum("trace", "debug", "info", "warn", "error")` | `false` | <code></code> |

**Project**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_ROOT` | Terraform project working directory | `string` | `false` | <code>.</code> |

**State**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_STATE_TYPE` | Terraform state type. | `string`<br/>`enum("gitlab-http")` | `false` | <code></code> |
| `$TF_STATE_NAME` | Terraform state name. | `string` | `false` | <code>default</code> |
| `$TF_STATE_STRICT` | Terraform state strict. | `bool` | `false` | <code>false</code> |
| `$TF_HTTP_ADDRESS`<br />`$TF_ADDRESS` | State configuration for terraform: http-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_ADDRESS` | State configuration for terraform: http-lock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_LOCK_METHOD` | State configuration for terraform: http-lock-method | `string` | `false` | <code>POST</code> |
| `$TF_HTTP_UNLOCK_ADDRESS` | State configuration for terraform: http-unlock-address | `string` | `false` | <code></code> |
| `$TF_HTTP_UNLOCK_METHOD` | State configuration for terraform: http-unlock-method | `string` | `false` | <code>DELETE</code> |
| `$TF_HTTP_USERNAME`<br />`$TF_USERNAME` | State configuration for terraform: http-username | `string` | `false` | <code>gitlab-ci-token</code> |
| `$TF_HTTP_PASSWORD`<br />`$TF_PASSWORD`<br />`$CI_JOB_TOKEN` | State configuration for terraform: http-password | `string` | `false` | <code></code> |
| `$TF_HTTP_RETRY_WAIT_MIN` | State configuration for terraform: http-retry-wait-min | `string` | `false` | <code>5</code> |

### `pipe-terraform publish`

Publish terraform project.

#### Flags

**Module**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_NAME`<br />`$CI_PROJECT_NAME` | Name for the module that will be published. | `string` | `true` | <code></code> |
| `$TF_MODULE_CWD`<br />`$TF_ROOT` | Directory for the module that will be published. | `string` | `false` | <code>.</code> |
| `$TF_MODULE_SYSTEM` | Module system for the module that will be published. | `string` | `false` | <code>local</code> |

**Registry**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_MODULE_REGISTRY` | Registry of the module that will be published. | `string` | `false` | <code>gitlab</code> |

**Registry - Gitlab**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$CI_API_V4_URL` | Gitlab API URL for publish call. | `string` | `false` | <code></code> |
| `$CI_PROJECT_ID` | Gitlab project id for publish call. | `string` | `false` | <code></code> |
| `$CI_JOB_TOKEN` | Gitlab API token for publish call. | `string` | `false` | <code></code> |

**Tags File**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TAGS_FILE` | Read tags from a file. | `string` | `false` | <code>.tags</code> |
