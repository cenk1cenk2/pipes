# pipe-go

Pipe for Go builds.

`pipe-go [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-go install`](#pipe-go-install)
- [`pipe-go build`](#pipe-go-build)
- [`pipe-go lint`](#pipe-go-lint)
- [`pipe-go test`](#pipe-go-test)
- [`pipe-go tool`](#pipe-go-tool)

### `pipe-go install`

Vendor go modules.

`pipe-go install [FLAGS]`

#### Flags

**Install**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_INSTALL_VERIFY` | Use the sum file to verify module integrity. | `bool` | `true` |
| `$GO_INSTALL_ARGS` | Arguments to append to the install command. | `string` |  |

**Setup**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_CWD` | Working directory for go commands. | `string` | `"."` |
| `$GO_CACHE` | Cache directory for go commands. Leave empty to use the environment defaults. | `string` | `"./.go/"` |

### `pipe-go build`

Build an application.

`pipe-go build [FLAGS]`

#### Flags

**Build**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_BUILD_ARGS` | Arguments to append to the build command. | `string` |  |
| `$GO_BUILD_OUTPUT` | Output location for the build artifacts. | `string` | `"./dist/"` |
| `$GO_BUILD_BINARY_NAME` | Name of the binary to output during build. | `string` | `"bin"` |
| `$GO_BUILD_BINARY_TEMPLATE` | Binary naming for the build artifact. | `string`<br/>`format(Template(map[string]string))` | `"{{ .name }}{{ if .os }}-{{ .os }}{{ end }}{{ if .arch }}-{{ .arch }}{{ end }}"` |
| `$GO_BUILD_LINKER_FLAGS`<br/>`$GO_BUILD_LINKER` | Arguments for the linker during the build process. | `string`<br/>`format(Template())` |  |
| `$GO_BUILD_ENABLE_CGO`<br/>`$CGO_ENABLED` | Enable CGO during the build process. | `bool` | `false` |
| `$GO_BUILD_TARGETS` | Build targets for the build process. | `string`<br/>`format(yaml([]struct{ os?: string, arch?: string }))` | `"[]"` |
| `$GO_BUILD_TAGS` | Build tags for the build process. | `string[]` |  |
| `$GO_BUILD_VARIABLES` | Build variables for the build process. | `string`<br/>`format(yaml(map[string]string))` | `"{}"` |

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_CWD` | Working directory for go commands. | `string` | `"."` |
| `$GO_CACHE` | Cache directory for go commands. Leave empty to use the environment defaults. | `string` | `"./.go/"` |

</details>

### `pipe-go lint`

Run golangci-lint on the project.

`pipe-go lint [FLAGS]`

#### Flags

**Lint**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_LINT_ARGS` | Arguments to append to the lint command. | `string` |  |
| `$GO_LINT_TIMEOUT` | Timeout for the lint command. | `duration` | `5m0s` |

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_CWD` | Working directory for go commands. | `string` | `"."` |
| `$GO_CACHE` | Cache directory for go commands. Leave empty to use the environment defaults. | `string` | `"./.go/"` |

</details>

### `pipe-go test`

Run ginkgo on the project.

`pipe-go test [FLAGS]`

#### Flags

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_CWD` | Working directory for go commands. | `string` | `"."` |
| `$GO_CACHE` | Cache directory for go commands. Leave empty to use the environment defaults. | `string` | `"./.go/"` |

</details>

**Test**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_TEST_RUNNER` | Test runner for the test command. | `string`<br/>`format(enum("ginkgo", "go"))` | `"ginkgo"` |
| `$GO_TEST_ARGS` | Arguments to append to the test command. | `string` |  |
| `$GO_TEST_COVERAGE` | Coverage profile to write. Leave empty to skip coverage. | `string` | `"coverage.out"` |

### `pipe-go tool`

Run a specified go tool.

`pipe-go tool [FLAGS] Tool to run.`

#### Arguments

| Argument | Usage | Type | Values |
| --- | --- | --- | --- |
| `arg` | `Tool to run.` | `string[]` | `0..*` |

#### Flags

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_CWD` | Working directory for go commands. | `string` | `"."` |
| `$GO_CACHE` | Cache directory for go commands. Leave empty to use the environment defaults. | `string` | `"./.go/"` |

</details>

**Tool**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GO_TOOL` | Binary that provides the tooling. | `string` |  |
| `$GO_TOOL_ARGS` | Arguments to append to the tool command. | `string` |  |
