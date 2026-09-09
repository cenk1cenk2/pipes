# pipe-semantic-release

Releases applications through the semantic-release library.

`pipe-semantic-release [FLAGS]`

## Flags

**CI Variables**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$SEMANTIC_RELEASE_CI_COMMIT_REFERENCE`<br/>`$CI_COMMIT_REF_NAME` | Current commit reference, either the branch or the tag name of the project. | `string` |  |

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

**Environment**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$ENVIRONMENT_ENABLE` | Enable environment injection. | `bool` | `false` |
| `$ENVIRONMENT_CONDITIONS` | Regex pattern to select an environment.<br/>Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. | `string`<br/>`format(json([]struct{ match: RegExp, environment: string }))` | `"[\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+$\", \"environment\": \"production\" },\n    { \"match\": \"^tags/v?\\\\d+.\\\\d+.\\\\d+-.*\\\\.\\\\d+$\", \"environment\": \"stage\" },\n    { \"match\" :\"^heads/main$\", \"environment\": \"develop\" },\n    { \"match\": \"^heads/master$\", \"environment\": \"develop\" }\n]"` |
| `$ENVIRONMENT_FAIL_ON_NO_REFERENCE` | Fail on missing environment references. | `bool` | `true` |
| `$ENVIRONMENT_STRICT` | Fail when no environment is selected. | `bool` | `true` |

**GIT**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` |  |

**Login**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NPM_LOGIN` | NPM registries to login. | `string`<br/>`format(json([]struct{ username: string, token: string, registry?: string, useHttps?: bool }))` |  |
| `$NPM_NPMRC_FILE` | .npmrc file to use. | `string[]` | `".npmrc"` |
| `$NPM_NPMRC` | Direct contents of .npmrc file. | `string` |  |

**Package Manager**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$NODE_PACKAGE_MANAGER` | Preferred package manager for nodejs. | `string`<br/>`format(enum("npm", "yarn", "pnpm"))` | `"pnpm"` |

**Semantic Release**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$SEMANTIC_RELEASE_DRY_RUN` | Run semantic-release in dry mode without making changes. | `bool` | `false` |
| `$SEMANTIC_RELEASE_WORKSPACE` | Use @qiwi/multi-semantic-release package to do a workspace release. | `bool` | `false` |
