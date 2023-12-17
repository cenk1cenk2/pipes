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
| `$TF_INSTALL_ARGS` | Additional arguments for terraform init. | `String` | `false` |  |

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

### `plan`

Plan terraform project.

`pipe-terraform plan [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `String`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform plan. | `String` | `false` | plan |
| `$TF_PLAN_ARGS` | Additional arguments for terraform plan. | `String` | `false` |  |

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

### `apply`

Apply terraform project.

`pipe-terraform apply [GLOBAL FLAGS] [FLAGS]`

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$TF_REGISTRY_CREDENTIALS` | Terraform registry credentials. | `String`<br/>`json([]struct { registry: string, token: string })` | `false` |  |
| `$TF_PLAN_CACHE`<br/>`$TF_APPLY_OUTPUT`<br/>`$TF_PLAN_OUTPUT` | Output file for terraform apply. | `String` | `false` | plan |
| `$TF_APPLY_ARGS` | Additional arguments for terraform apply. | `String` | `false` |  |

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
