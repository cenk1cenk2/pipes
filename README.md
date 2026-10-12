# pipes

[![pipeline status](https://gitlab.kilic.dev/devops/pipes/badges/main/pipeline.svg)](https://gitlab.kilic.dev/devops/pipes/-/commits/main)

Automation pipe CLIs used by the shared GitLab CI templates in [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines).

## Consumer CI templates

Consumer-facing GitLab CI templates have moved to [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines). New and migrated repositories should include versioned templates from that repository instead of including files from `devops/pipes`.

## Pipe images

This repository still builds and publishes the pipe images used by `devops/pipelines`:

- [`cenk1cenk2/pipe-buildah`](./buildah/README.md)
- [`cenk1cenk2/pipe-git`](./git/README.md)
- [`cenk1cenk2/pipe-github`](./github/README.md)
- [`cenk1cenk2/pipe-go`](./go/README.md)
- [`cenk1cenk2/pipe-helm`](./helm/README.md)
- [`cenk1cenk2/pipe-kustomize`](./kustomize/README.md)
- [`cenk1cenk2/pipe-mise`](./mise/README.md)
- [`cenk1cenk2/pipe-node`](./node/README.md)
- [`cenk1cenk2/pipe-pulumi`](./pulumi/README.md)
- [`cenk1cenk2/pipe-semantic-release`](./semantic-release/README.md)
- [`cenk1cenk2/pipe-terraform`](./terraform/README.md)
- [`cenk1cenk2/pipe-update-docker-hub-readme`](./update-docker-hub-readme/README.md)

## Methodology

The `_template` directory contains the scaffold for creating a pipe, and every pipe lives in its own directory at the repository root. Pipes use the [plumber](https://gitlab.kilic.dev/libraries/plumber) framework to create a CLI and execute commands in a specific order.

### Releases

Every pipe, and `internal`, is a package of the workspace in the root `package.json` and is versioned on its own by [multi-semantic-release](https://github.com/anolilab/multi-semantic-release), tagged as `<pipe>@x.y.z`. A commit releases the packages whose directory it touches, and a release of `internal` patches every pipe that depends on it. Next to the usual types, `chore(deps)` releases a patch, since a dependency bump changes what the image ships.

Publishing follows the tags. The pipeline of a `<pipe>@x.y.z` tag builds that pipe and publishes the tags of its build matrix, `latest` among them, next to the version: `vx.y.z` for a single image, and `<flavour>-vx.y.z` for a pipe that publishes more than one flavour, like `1.27-alpine-v1.2.0`. A change to a pipe, its `Dockerfile` or its build matrix in `<pipe>/.gitlab-ci.yml` therefore ships through a releasing commit type, and a `ci` or `docs` commit does not.

A failed publish is retried on the released code by running the pipeline of its tag again. Retrying a job of the original pipeline only works while the vendored tree of its `install` job is kept, for a week. A tag on a commit carrying `[skip ci]` gets a skipped pipeline, which "Run pipeline" on the tag recovers. Running the pipeline of an older tag again moves `latest` back to that version, and the tags released before publishing moved to them have no jobs at all.

Merge requests build and test only the pipes they touch, and the default branch builds none, since every change reached it through a merge request. A change outside every pipe directory that no release covers is published by running a pipeline on the default branch by hand, which builds every pipe and publishes `latest` from that commit when the manual `container` job of a pipe is played. The jobs of the parent pipeline run only when what they check changed.

A new pipe needs a `package.json` naming it, a `<pipe>/.gitlab-ci.yml` including `.gitlab-ci.trigger.yml` with its build matrix, that file in the includes of the root `.gitlab-ci.yml`, and its module in `go.work`.

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
