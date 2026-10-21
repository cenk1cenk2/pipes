# pipe-python

Pipe for installing, checking and building Python projects with uv on CI/CD.

`pipe-python [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

- [`pipe-python install`](#pipe-python-install)
- [`pipe-python lint`](#pipe-python-lint)
- [`pipe-python test`](#pipe-python-test)
- [`pipe-python build`](#pipe-python-build)
- [`pipe-python run`](#pipe-python-run)

### `pipe-python install`

Install the locked dependencies of the project.

`pipe-python install [FLAGS]`

#### Flags

**Install**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_INSTALL_LOCKED` | Fail when the lockfile is out of date with the project. Otherwise the lockfile is installed as it is. | `bool` | `true` |
| `$PYTHON_INSTALL_ARGS` | Arguments to append to the install command. | `string` |  |

**Setup**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_CWD` | Working directory of the uv project. | `string` | `"."` |
| `$PYTHON_CACHE` | Cache directory for uv. Leave empty to use the environment defaults. | `string` | `"./.uv/"` |

### `pipe-python lint`

Run the lint commands inside the installed environment.

`pipe-python lint [FLAGS]`

#### Flags

**Lint**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_LINT_COMMANDS` | Commands to lint with, each run in order inside the installed environment. | `string`<br/>`format(yaml([]string))` | `"[\"ruff format --check .\", \"ruff check .\", \"ty check .\"]"` |

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_CWD` | Working directory of the uv project. | `string` | `"."` |
| `$PYTHON_CACHE` | Cache directory for uv. Leave empty to use the environment defaults. | `string` | `"./.uv/"` |

</details>

### `pipe-python test`

Run the test commands inside the installed environment.

`pipe-python test [FLAGS]`

#### Flags

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_CWD` | Working directory of the uv project. | `string` | `"."` |
| `$PYTHON_CACHE` | Cache directory for uv. Leave empty to use the environment defaults. | `string` | `"./.uv/"` |

</details>

**Test**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_TEST_COMMANDS` | Commands to test with, each run in order inside the installed environment. | `string`<br/>`format(yaml([]string))` | `"[\"pytest --junitxml=junit.xml\"]"` |

### `pipe-python build`

Build the source distribution and the wheel of the project.

`pipe-python build [FLAGS]`

#### Flags

**Build**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_BUILD_OUTPUT` | Directory the distributions are written to, relative to the working directory. | `string` | `"dist"` |
| `$PYTHON_BUILD_ARGS` | Arguments to append to the build command. | `string` |  |

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_CWD` | Working directory of the uv project. | `string` | `"."` |
| `$PYTHON_CACHE` | Cache directory for uv. Leave empty to use the environment defaults. | `string` | `"./.uv/"` |

</details>

### `pipe-python run`

Run a command inside the installed environment.

`pipe-python run [FLAGS] Command to run inside the installed environment.`

#### Arguments

| Argument | Usage | Type | Values |
| --- | --- | --- | --- |
| **`arg`**\* | `Command to run inside the installed environment.` | `string[]` | `1..*` |

\* required

#### Flags

<details>
<summary>Setup</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$PYTHON_CWD` | Working directory of the uv project. | `string` | `"."` |
| `$PYTHON_CACHE` | Cache directory for uv. Leave empty to use the environment defaults. | `string` | `"./.uv/"` |

</details>
