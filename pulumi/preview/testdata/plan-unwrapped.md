## Pulumi preview report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| Stack | `dev` |
| Job | `pulumi-preview` |
| Pulumi version | `3.187.0` |
| Plan time | `2026-05-31T12:00:00Z` |

### Summary

| Action | Resources | Output properties |
| --- | ---: | ---: |
| `+ create` | 1 | 1 |
| `~ update` | 1 | 1 |

Total planned actions: 2.

### Resources

#### `+ create` (1)

<details>
<summary><code>+ create</code> <code>aws:s3/bucket:Bucket/logs</code> (<code>urn:pulumi:dev::example::aws:s3/bucket:Bucket::logs</code>)</summary>

```diff
+ create
+ acl: "private"
+ region: [unknown]
+ rules:
+   [0]:
+     days: 30
+     name: "expire"
+ tags:
+   env: "dev"
+   owner: [secret]
+ token: [secret]
```

</details>

#### `~ update` (1)

<details>
<summary><code>~ update</code> <code>kubernetes:core/v1:ConfigMap/settings</code> (<code>urn:pulumi:dev::example::kubernetes:core/v1:ConfigMap::settings</code>)</summary>

```diff
~ update
  data:
+   level: "debug"
+   password: [secret]
- immutable
```

</details>

### Output properties

#### `+ create` (1)
- `aws:s3/bucket:Bucket/logs`: `bucketName`

#### `~ update` (1)
- `kubernetes:core/v1:ConfigMap/settings`: `metadata`
