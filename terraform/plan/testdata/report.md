## Terraform plan report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| State | `production` |
| Job | `tf-plan` |
| Terraform version | `1.9.8` |
| Plan time | `2026-05-31T12:00:00Z` |
| Plan schema version | `1.2` |

### Summary

| Action | Resources | Outputs |
| --- | ---: | ---: |
| `+ create` | 1 | 1 |
| `~ update` | 1 | 1 |
| `- delete` | 0 | 1 |
| `-/+ replace` | 1 | 0 |
| `> move` | 1 | 0 |
| `<= read` | 1 | 0 |

Total planned actions: 5.

### Resources

#### `+ create` (1)

<details>
<summary><code>+ create</code> <code>aws_s3_bucket.logs</code></summary>

```diff
+ create
+ arn: (known after apply)
+ bucket: (sensitive value)
+ cors_rule:
+   [0]:
+     allowed_methods: ["GET"]
+     id: (known after apply)
+     max_age_seconds: 300
+ force_destroy: false
+ id: (known after apply)
+ lifecycle_rule: []
+ tags:
+   env: "dev"
+   team: "platform"
```

</details>

#### `~ update` (1)

<details>
<summary><code>~ update</code> <code>aws_instance.web</code></summary>

```diff
~ update
- instance_type: "t3.micro"
+ instance_type: "t3.small"
- user_data: (sensitive value)
+ user_data: (sensitive value)
```

</details>

#### `-/+ replace` (1)

<details>
<summary><code>-/+ replace</code> <code>aws_db_instance.main</code> (destroy before create)</summary>

```diff
-/+ replace
- engine_version: "15"
+ engine_version: "16" # forces replacement
- password: (sensitive value)
+ password: (sensitive value)
```

</details>

#### `> move` (1)

<details>
<summary><code>&gt; move</code> <code>aws_sqs_queue.jobs</code> (moved from <code>aws_sqs_queue.legacy_jobs</code>)</summary>

```diff
> move
No attribute changes.
```

</details>

#### `<= read` (1)

<details>
<summary><code>&lt;= read</code> <code>data.aws_caller_identity.current</code></summary>

```diff
<= read
+ account_id: (sensitive value)
```

</details>

### Outputs

#### `+ create` (1)
- `bucket_name`

#### `~ update` (1)
- `endpoint`

#### `- delete` (1)
- `password`
