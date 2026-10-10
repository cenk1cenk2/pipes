# pipe-git

Git related tasks in the pipeline.

`pipe-git [GLOBAL FLAGS] [COMMAND] [FLAGS]`

## Global Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

## Commands

### `pipe-git sync`

Sync the paths a job generated into a branch of a GitLab project. Replaces every destination on the target branch with its source, writes the staged diff to a patch, logs it and reports it on the merge requests of the pipeline. In publish mode, the changes are committed onto the tip of the target branch, pushed to the sync branch and opened as a merge request against the target branch. A pipeline running on the target branch itself skips the sync, and a publish fails once the commit of the pipeline is no longer the head of the default branch.

`pipe-git sync [FLAGS]`

#### Flags

**GIT**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_REF_NAME`<br/>`$BITBUCKET_BRANCH` | Source control branch. | `string` |  |
| `$CI_COMMIT_TAG`<br/>`$BITBUCKET_TAG` | Source control tag. | `string` |  |

**GitLab Merge Request Report**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GITLAB_MR_REPORT_ENABLED` | Enable GitLab merge request report note on the given merge request. | `bool` | `false` |
| `$GL_PIPES_TOKEN` | GitLab API token for merge request report notes. | `string` |  |
| `$CI_API_V4_URL` | GitLab API URL for merge request report notes. | `string` |  |
| `$CI_PROJECT_ID` | GitLab project id for merge request report notes. | `string` |  |
| `$CI_MERGE_REQUEST_IID` | GitLab merge request iid for merge request report notes. | `int` | `0` |
| `$GITLAB_MR_REPORT_IDENTIFIER` | Hidden marker identifier for merge request report notes. Defaults to the job name combined with the stack or state under report. | `string` |  |

**GitLab Pipeline**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_JOB_NAME` | GitLab CI job name to include in the plan report metadata. | `string` |  |
| `$CI_JOB_URL` | GitLab CI job URL to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_ID` | GitLab CI pipeline id to include in the plan report metadata. | `string` |  |
| `$CI_PIPELINE_URL` | GitLab CI pipeline URL to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHA` | Git commit sha to include in the plan report metadata. | `string` |  |
| `$CI_COMMIT_SHORT_SHA` | Short git commit sha to include in the plan report metadata. | `string` |  |

**GitLab Project**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$CI_COMMIT_AUTHOR` | Author of the commit of the pipeline, credited as a co-author of the sync commit. | `string` |  |
| `$CI_PROJECT_PATH` | Path of the GitLab project the pipeline runs for. | `string` |  |
| `$CI_PROJECT_PATH_SLUG` | Slug of the path of the GitLab project the pipeline runs for. | `string` |  |
| `$CI_JOB_NAME_SLUG` | Slug of the name of the job. | `string` |  |
| `$CI_PROJECT_DIR` | Checkout of the project the target branch is checked out next to when the project is the target. | `string` | `"."` |
| `$CI_SERVER_URL` | URL of the GitLab instance the projects are fetched from and pushed to. | `string` |  |
| `$CI_DEFAULT_BRANCH` | Default branch of the project, which the commit of the pipeline has to still be the head of for a publish. | `string` |  |
| `$CI_OPEN_MERGE_REQUESTS` | Merge requests open from the branch of the pipeline, which the report is posted on outside of a merge request pipeline. | `string` |  |

**Identity**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GIT_PIPES_AUTHOR_NAME` | Name of the author of the commit. Required in publish mode. | `string` |  |
| `$GIT_PIPES_AUTHOR_EMAIL` | Email of the author of the commit. Required in publish mode. | `string` |  |
| `$GIT_PIPES_COMMITTER_NAME` | Name of the committer of the commit. Required in publish mode. | `string` |  |
| `$GIT_PIPES_COMMITTER_EMAIL` | Email of the committer of the commit. Required in publish mode. | `string` |  |

**Merge Request**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GIT_PIPES_ASSIGNEES` | GitLab usernames assigned to the merge request, replacing its assignees on every publish. Leaves the assignees of the merge request as they are when empty. | `string[]` |  |
| `$GIT_PIPES_REVIEWERS` | GitLab usernames requested to review the merge request, replacing its reviewers on every publish. Leaves the reviewers of the merge request as they are when empty. | `string[]` |  |

**Sync**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$GIT_SYNC_MODE` | Whether the changes are only reported or committed and opened as a merge request on the target as well. | `string`<br/>`format(enum(publish, report))` | `"report"` |
| `$GIT_SYNC_TARGET_PROJECT`<br/>`$CI_PROJECT_PATH` | Path of the GitLab project the paths are synced into. Defaults to the project of the pipeline. | `string` |  |
| `$GIT_SYNC_TARGET_BRANCH` | Branch of the target the paths are synced onto, and the merge request targets. | `string` | `"next"` |
| **`$GIT_SYNC_PATHS`**\* | Sources of the job synced onto destinations in the target, which are replaced completely. Include and exclude are globs anchored at the source. | `string`<br/>`format(yaml([]{ source: string, destination: string, include?: []string, exclude?: []string }))` |  |
| `$GIT_SYNC_BRANCH` | Branch the changes are pushed to and the merge request is opened from. Defaults to sync/$CI_PROJECT_PATH_SLUG-$CI_JOB_NAME_SLUG. | `string` |  |
| `$GIT_SYNC_COMMIT_MESSAGE` | Message of the commit, whose subject is the title of the merge request as well. | `string` | `"chore: sync generated files"` |
| `$GIT_SYNC_ALLOW_EMPTY` | Sync a missing or empty source as an empty destination instead of failing. | `bool` | `false` |
| `$GIT_SYNC_PATCH` | File the full staged diff is written to, for the job artifacts. A relative path resolves against the working directory. | `string` | `"git-sync.patch"` |
| `$GIT_SYNC_GL_TOKEN` | GitLab token that fetches and pushes the target and opens the merge request on it. A publish also reads the head of the default branch of the project of the pipeline with it, so a target in another project needs read access to the project of the pipeline as well. | `string` |  |

\* required
