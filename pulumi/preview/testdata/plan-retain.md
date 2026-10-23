## Pulumi preview report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| Stack | `default` |
| Job | `pulumi-preview` |
| Pulumi version | `v3.272.0` |
| Plan time | `2026-10-23T20:34:18Z` |
| Plan schema version | `1` |

### Summary

| Action | Resources | Output properties |
| --- | ---: | ---: |
| `-/+ replace` | 1 | 0 |
| `+ create-replacement` | 1 | 0 |
| `- delete-replaced` | 1 | 0 |
| `- delete` | 3 | 0 |

Total planned actions: 6.

### Resources

#### `-/+ replace` (1)

<details open>
<summary><code>-/+ replace</code> <code>gitlab:index/branchProtection:BranchProtection/branch-protection-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/branchProtection:BranchProtection::branch-protection-265</code>) (delete before create, retain on delete)</summary>

```diff
-/+ replace
+ branch: "release"
```

</details>

#### `+ create-replacement` (1)

<details open>
<summary><code>+ create-replacement</code> <code>gitlab:index/branchProtection:BranchProtection/branch-protection-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/branchProtection:BranchProtection::branch-protection-265</code>) (delete before create, retain on delete)</summary>

```diff
+ create-replacement
+ branch: "release"
```

</details>

#### `- delete-replaced` (1)

<details open>
<summary><code>- delete-replaced</code> <code>gitlab:index/branchProtection:BranchProtection/branch-protection-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/branchProtection:BranchProtection::branch-protection-265</code>) (delete before create, retain on delete)</summary>

```diff
- delete-replaced
+ branch: "release"
```

</details>

#### `- delete` (3)

<details open>
<summary><code>- delete</code> <code>gitlab:index/project:Project/project-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/project:Project::project-265</code>) (retain on delete)</summary>

```diff
- delete
- name: "x"
```

</details>

<details open>
<summary><code>- delete</code> <code>gitlab:index/projectHook:ProjectHook/linear-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/projectHook:ProjectHook::linear-265</code>)</summary>

```diff
- delete
No attribute changes.
```

</details>

<details open>
<summary><code>- delete</code> <code>gitlab:index/tagProtection:TagProtection/tag-protection-265</code> (<code>urn:pulumi:default::gitlab::gitlab:index/tagProtection:TagProtection::tag-protection-265</code>) (retain on delete)</summary>

```diff
- delete
- name: "x"
```

</details>
