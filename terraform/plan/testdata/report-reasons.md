## Terraform plan report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| State | `production` |
| Job | `tf-plan` |
| Terraform version | `1.14.0` |
| Plan time | `2026-10-23T12:00:00Z` |
| Plan schema version | `1.2` |

### Summary

| Action | Resources | Outputs |
| --- | ---: | ---: |
| `~ update` | 1 | 0 |
| `- delete` | 3 | 0 |
| `-/+ replace` without property changes | 1 | 0 |
| `> move` | 1 | 0 |

Total planned actions: 6.

### Resources

#### `~ update` (1)

<details open>
<summary><code>~ update</code> <code>aws_instance.web</code></summary>

```diff
~ update
- instance_type: "t3.micro"
+ instance_type: "t3.small"
```

</details>

#### `- delete` (3)

<details open>
<summary><code>- delete</code> <code>aws_eip.ip[2]</code> (because index [2] is out of range for count)</summary>

```diff
- delete
- domain: "vpc"
```

</details>

<details open>
<summary><code>- delete</code> <code>aws_instance.web</code> (deposed object 00000001)</summary>

```diff
- delete
- instance_type: "t3.nano"
```

</details>

<details open>
<summary><code>- delete</code> <code>aws_sqs_queue.jobs</code> (moved from <code>aws_sqs_queue.legacy_jobs</code>) (because aws_sqs_queue.legacy_jobs was moved to aws_sqs_queue.jobs, which is not in configuration)</summary>

```diff
- delete
- name: "jobs"
```

</details>

#### `-/+ replace` without property changes (1)

<details open>
<summary><code>-/+ replace</code> <code>aws_db_instance.main</code> (create before destroy, tainted)</summary>

```diff
-/+ replace
No attribute changes.
```

</details>

#### `> move` (1)

<details open>
<summary><code>&gt; move</code> <code>aws_sqs_queue.jobs</code> (moved from <code>aws_sqs_queue.legacy_jobs</code>) (because aws_sqs_queue.legacy_jobs was moved to aws_sqs_queue.jobs, which is not in configuration)</summary>

```diff
> move
- name: "jobs"
```

</details>
