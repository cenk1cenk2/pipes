## Pulumi preview report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| Stack | `default` |
| Job | `pulumi-preview` |
| Pulumi version | `v3.268.0` |
| Plan time | `2026-10-13T11:42:37Z` |

### Summary

| Action | Resources | Output properties |
| --- | ---: | ---: |
| `+ create` | 1 | 0 |
| `= import` | 2 | 0 |
| `~ update` | 1 | 0 |

Total planned actions: 4.

### Resources

#### `+ create` (1)

<details>
<summary><code>+ create</code> <code>pulumi-nodejs:dynamic/gitlab:ResourceGroup/resource-group-12-9544a30e</code> (<code>urn:pulumi:default::gitlab::pulumi-nodejs:dynamic/gitlab:ResourceGroup::resource-group-12-9544a30e</code>)</summary>

```diff
+ create
+ key: "semantic-release"
+ processMode: "oldest_first"
+ project: 12
+ token: [secret]
+ url: "https://gitlab.kilic.dev/api/v4"
```

</details>

#### `= import` (2)

<details>
<summary><code>= import</code> <code>gitlab:index/project:Project/project-11</code> (<code>urn:pulumi:default::gitlab::gitlab:index/project:Project::project-11</code>)</summary>

```diff
= import
+ archiveOnDestroy: true
+ resourceGroupDefaultProcessMode: "oldest_first"
+ squashCommitTemplate: <<-EOT
+   %{title}
+ 
+   %{all_commits}
+ EOT
```

</details>

<details>
<summary><code>= import</code> <code>gitlab:index/project:Project/project-12</code> (<code>urn:pulumi:default::gitlab::gitlab:index/project:Project::project-12</code>)</summary>

```diff
= import
No attribute changes.
```

</details>

#### `~ update` (1)

<details>
<summary><code>~ update</code> <code>gitlab:index/project:Project/project-11</code> (<code>urn:pulumi:default::gitlab::gitlab:index/project:Project::project-11</code>)</summary>

```diff
~ update
+ archiveOnDestroy: true
+ resourceGroupDefaultProcessMode: "oldest_first"
+ squashCommitTemplate: <<-EOT
+   %{title}
+ 
+   %{all_commits}
+ EOT
```

</details>
