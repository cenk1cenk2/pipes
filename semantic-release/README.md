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
| `$SEMANTIC_RELEASE_WORKSPACE` | Use @anolilab/multi-semantic-release package to do a workspace release. | `bool` | `false` |
| `$TAGS_FILE` | File the release writes its tag into, which is handed over to the later jobs as a dotenv report. | `string` | `".tags"` |
| `$SEMANTIC_RELEASE_ENV_FILE` | Dotenv report the released tags are written into for the later jobs. | `string` | `"semantic-release.env"` |
| `$SEMANTIC_RELEASE_ENV_VARIABLE` | Variable the released tag is handed over in, suffixed with the package name of every released package in a workspace. | `string` | `"SEMANTIC_RELEASE_TAG"` |
