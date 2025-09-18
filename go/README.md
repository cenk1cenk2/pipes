# pipe-go

Build Go applications with the CI pipe.

`pipe-go [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>info</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code>[]</code> |

## Commands

### `pipe-go install`

Vendor go modules.

#### Flags

**Install**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_INSTALL_CWD` | Install CWD for the package manager. | `string` | `false` | <code>.</code> |
| `$GO_INSTALL_VERIFY` | Use the sum file to verify module integrity. | `bool` | `false` | <code>true</code> |
| `$GO_INSTALL_ARGS` | Arguments to append to install command. | `string` | `false` | <code></code> |

**Setup**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_CACHE` | Enable go cache. | `string` | `false` | <code></code> |

### `pipe-go lint`

Lint the repository.

#### Flags

**Lint**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_LINT_CWD` | Lint CWD for the package manager. | `string` | `false` | <code>.</code> |
| `$GO_LINT_ARGS` | Arguments to append to lint command. | `string` | `false` | <code>run --timeout 3m</code> |
| `$GO_LINT_SOURCE` | Source for the linter. | `string`<br/>`enum("tools")` | `false` | <code>tools</code> |
| `$GO_LINT_TOOL` | Binary that provides the linting. | `string` | `false` | <code>golangci-lint</code> |

**Setup**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_CACHE` | Enable go cache. | `string` | `false` | <code></code> |

### `pipe-go build`

Build an application.

#### Flags

**Build**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_BUILD_CWD` | Build CWD for the package manager. | `string` | `false` | <code>.</code> |
| `$GO_BUILD_ARGS` | Arguments to append to build command. | `string` | `false` | <code></code> |
| `$GO_BUILD_OUTPUT` | Output location for the build artifacts. | `string` | `false` | <code>./dist/</code> |
| `$GO_BUILD_BINARY_NAME` | Name of the binary to output during build. | `string` | `false` | <code>bin</code> |
| `$GO_BUILD_BINARY_TEMPLATE` | Binary naming for the build artifact. | `string`<br/>`format(Template(map[string]))` | `false` | <code>{{ .name }}{{ if .os }}-{{ .os }}{{ end }}{{ if .arch }}-{{ .arch }}{{ end }}</code> |
| `$GO_BUILD_LD_FLAGS` | Arguments for the linker during the build process. | `string` | `false` | <code></code> |
| `$GO_BUILD_ENABLE_CGO`<br />`$CGO_ENABLED` | Enable CGO during the build process. | `bool` | `false` | <code>false</code> |
| `$GO_BUILD_TARGETS` | Build targets for the build process. | `string`<br/>`format(yaml([]struct{ os: string?, arch: string? }))` | `false` | <code>[]</code> |

**Setup**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$GO_CACHE` | Enable go cache. | `string` | `false` | <code></code> |
