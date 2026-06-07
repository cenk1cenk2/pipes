# pipe-pulumi

Pulumi actions for CI pipelines.

`pipe-pulumi [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

## Commands

### `pipe-pulumi preview`

Preview the Pulumi changes.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name for the pulumi to be used in the commands. | `string` | `true` | <code></code> |
| `$PULUMI_PLAN` | Output file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |

**Report**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_REPORT_ENABLED` | Generate Pulumi report artifact. | `bool` | `false` | <code>true</code> |
| `$PULUMI_REPORT_OUTPUT` | Output file for Pulumi report artifact. | `string` | `false` | <code>"pulumi-report.html"</code> |

**pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Path to the Pulumi working directory. | `string` | `false` | <code>"."</code> |

### `pipe-pulumi up`

Apply the Pulumi changes.

#### Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_STACK` | Stack name for the pulumi to be used in the commands. | `string` | `true` | <code></code> |
| `$PULUMI_PLAN` | Input file for pulumi plan. | `string` | `false` | <code>"plan.json"</code> |

**pulumi**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$PULUMI_CWD` | Path to the Pulumi working directory. | `string` | `false` | <code>"."</code> |
