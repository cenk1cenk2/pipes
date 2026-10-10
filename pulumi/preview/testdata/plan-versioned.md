## Pulumi preview report

This report carries the plan metadata, the action planned for every resource with the attribute values the plan proposes for it, and the names of the outputs that change. Values the tool marks as sensitive or secret are redacted, values only known after apply are marked as such, and output values are not included.

### Metadata

| Field | Value |
| --- | --- |
| Stack | `dev` |
| Job | `pulumi-preview` |
| Pulumi version | `3.187.0` |
| Plan time | `2026-05-31T12:30:00Z` |
| Plan schema version | `1` |

### Summary

| Action | Resources | Output properties |
| --- | ---: | ---: |
| `-/+ replace` | 1 | 1 |
| `+ create-replacement` | 1 | 1 |
| `- delete-replaced` | 1 | 1 |

Total planned actions: 3.

### Resources

#### `-/+ replace` (1)

<details>
<summary><code>-/+ replace</code> <code>aws:lambda/function:Function/worker</code> (<code>urn:pulumi:stage::example::aws:lambda/function:Function::worker</code>) (create before delete)</summary>

```diff
-/+ replace
+ environment: [secret]
+ runtime: "nodejs22.x"
```

</details>

#### `+ create-replacement` (1)

<details>
<summary><code>+ create-replacement</code> <code>aws:lambda/function:Function/worker</code> (<code>urn:pulumi:stage::example::aws:lambda/function:Function::worker</code>) (create before delete)</summary>

```diff
+ create-replacement
+ environment: [secret]
+ runtime: "nodejs22.x"
```

</details>

#### `- delete-replaced` (1)

<details>
<summary><code>- delete-replaced</code> <code>aws:lambda/function:Function/worker</code> (<code>urn:pulumi:stage::example::aws:lambda/function:Function::worker</code>) (create before delete)</summary>

```diff
- delete-replaced
+ environment: [secret]
+ runtime: "nodejs22.x"
```

</details>

### Output properties

#### `-/+ replace` (1)
- `aws:lambda/function:Function/worker`: `arn`

#### `+ create-replacement` (1)
- `aws:lambda/function:Function/worker`: `arn`

#### `- delete-replaced` (1)
- `aws:lambda/function:Function/worker`: `arn`
