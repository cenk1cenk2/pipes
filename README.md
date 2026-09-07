# pipes

[![pipeline status](https://gitlab.kilic.dev/devops/pipes/badges/main/pipeline.svg)](https://gitlab.kilic.dev/devops/pipes/-/commits/main)

Automation pipe CLIs used by the shared GitLab CI templates in [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines).

## Consumer CI templates

Consumer-facing GitLab CI templates have moved to [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines). New and migrated repositories should include versioned templates from that repository instead of including files from `devops/pipes`.

## Pipe images

This repository still builds and publishes the pipe images used by `devops/pipelines`:

- [`cenk1cenk2/pipe-buildah`](./buildah/README.md)
- [`cenk1cenk2/pipe-go`](./go/README.md)
- [`cenk1cenk2/pipe-helm`](./helm/README.md)
- [`cenk1cenk2/pipe-kustomize`](./kustomize/README.md)
- [`cenk1cenk2/pipe-node`](./node/README.md)
- [`cenk1cenk2/pipe-pulumi`](./pulumi/README.md)
- [`cenk1cenk2/pipe-select-env`](./select-env/README.md)
- [`cenk1cenk2/pipe-semantic-release`](./semantic-release/README.md)
- [`cenk1cenk2/pipe-terraform`](./terraform/README.md)
- [`cenk1cenk2/pipe-update-docker-hub-readme`](./update-docker-hub-readme/README.md)

## Methodology

The `_template` directory contains the scaffold for creating a pipe, and every pipe lives in its own directory at the repository root. Pipes use the [plumber](https://gitlab.kilic.dev/libraries/plumber) framework to create a CLI and execute commands in a specific order.

### Flag descriptions

The `Usage` of a flag is one or more plain sentences. Whenever the Go type of the value does not describe the payload, a single annotation closes the description and documents what the validator accepts:

- `format(yaml(<shape>))` and `format(json(<shape>))` for the values parsed by the unmarshal wrappers in `internal/flags`.
- `format(Template(<context>))` for the values rendered as a Go template, where the context is what the template is interpolated against.
- `format(enum(<values>))` for the values a `oneof` validation tag restricts, since validation tags never reach the generated documentation.
- `format(glob)` and `format(RegExp)` for the remaining opaque strings.

Shapes are written in a TypeScript flavour: `string`, `[]T`, `map[string]T`, `struct{ field: T, optional?: T }`.

```go
Usage: "Read published tags from a file. format(glob)",
```

Descriptions that span multiple lines carry no indentation of their own and close with the annotation on the last line, after a blank one:

```go
Usage: strings.TrimSpace(`
Modifies every tag that matches a certain condition.
Template is interpolated with the given matches in the regular expression.

format(yaml([]struct{ match: RegExp, template: Template(match) }))
`),
```

The documentation generator lifts the annotation out of the description into the type column, matching from the first `format(`, `json(`, `yaml(`, `Template(`, `RegExp(`, `enum(` or `multiple(` up to the last `)` in the whole string. Nothing may follow the annotation, and the prose itself must never place one of those words directly in front of a parenthesis, or the description gets cut from that word onwards.
