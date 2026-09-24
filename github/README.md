# pipe-github

GitHub App tokens and commit statuses in the pipeline.

`pipe-github [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

### `pipe-github token`

Mint a GitHub App installation token into a dotenv file. Every GitHub App flag is required here.

The token lives for an hour. Expose the file as an artifacts:reports:dotenv report, and every job that needs this one receives the token as a variable, which overrides the job-level variable of the same name.

`pipe-github token [FLAGS]`

#### Flags

**GitHub App**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITHUB_APP_ID` | GitHub App id the token is minted for. | `string` |  |
| `$GITHUB_APP_INSTALLATION_ID` | Installation id of the GitHub App on the owner of the repository. | `string` |  |
| `$GITHUB_APP_PRIVATE_KEY` | Path to the private key of the GitHub App, as a GitLab file variable exports it, or the PEM contents themselves. | `string` |  |
| `$GITHUB_API_URL` | GitHub API URL. | `string` | `"https://api.github.com"` |

**Token**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITHUB_TOKEN_REPOSITORIES` | Repository names, without the owner, to narrow the token down to. Left empty, the token covers every repository of the installation. | `string[]` |  |
| `$GITHUB_TOKEN_PERMISSIONS` | Permissions to narrow the token down to, e.g. {"statuses":"write"}. Left empty, the token carries every permission of the installation. | `string`<br/>`format(json(map[string]string))` |  |
| `$GITHUB_TOKEN_VARIABLE` | Variable name the token is written under in the dotenv file. | `string` | `"GH_TOKEN"` |
| `$GITHUB_TOKEN_GIT_CREDENTIALS_VARIABLE` | Variable name the token is written under as an x-access-token git credential in the dotenv file, for semantic-release to push with. Left empty, no git credential is written. | `string` | `"GIT_CREDENTIALS"` |
| `$GITHUB_TOKEN_FILE` | Dotenv file the token is written into, for the job to expose as a dotenv report. Other variables in an existing file are kept. | `string` | `"github.env"` |

### `pipe-github status`

Post a commit status to GitHub.

Given the GitHub App flags, the status mints its own token, narrowed to the repository and to writing statuses. That needs the private key in the status job as well, so prefer the token from the dotenv file and mint in place only where the pipeline can outlive the token.

A commit GitHub does not have, as on a branch that only exists on GitLab, is skipped with a warning; every other failure fails the job.

`pipe-github status [FLAGS]`

#### Flags

<details>
<summary>GitHub App</summary>

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITHUB_APP_ID` | GitHub App id the token is minted for. | `string` |  |
| `$GITHUB_APP_INSTALLATION_ID` | Installation id of the GitHub App on the owner of the repository. | `string` |  |
| `$GITHUB_APP_PRIVATE_KEY` | Path to the private key of the GitHub App, as a GitLab file variable exports it, or the PEM contents themselves. | `string` |  |
| `$GITHUB_API_URL` | GitHub API URL. | `string` | `"https://api.github.com"` |

</details>

**Status**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITHUB_STATUS_TOKEN`<br/>`$GH_TOKEN` | GitHub token to post the status with. Not needed when the GitHub App flags are given, since the status then mints its own token. | `string` |  |
| **`$GITHUB_STATUS_PROJECT`**\* | GitHub repository to post the status to, as owner/repository. | `string` |  |
| **`$GITHUB_STATUS_STATE`**<br/>**`$GITHUB_STATUS_REPORT`**\* | State of the status. | `string`<br/>`format(enum(pending, success, failure, error))` |  |
| **`$GITHUB_STATUS_SHA`**<br/>**`$CI_COMMIT_SHA`**\* | Commit sha to post the status for. | `string` |  |
| `$GITHUB_STATUS_TARGET_URL`<br/>`$CI_PIPELINE_URL` | URL the status links to. | `string` |  |
| `$GITHUB_STATUS_CONTEXT` | Context that tells this status apart from the others on the commit. | `string` | `"Gitlab CI"` |
| `$GITHUB_STATUS_DESCRIPTION` | Short description of the status. | `string` |  |

\* required
