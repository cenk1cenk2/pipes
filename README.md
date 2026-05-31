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
- [`cenk1cenk2/pipe-node`](./node/README.md)
- [`cenk1cenk2/pipe-pulumi`](./pulumi/README.md)
- [`cenk1cenk2/pipe-select-env`](./select-env/README.md)
- [`cenk1cenk2/pipe-semantic-release`](./semantic-release/README.md)
- [`cenk1cenk2/pipe-terraform`](./terraform/README.md)
- [`cenk1cenk2/pipe-update-docker-hub-readme`](./update-docker-hub-readme/README.md)

## Methodology

The `_template` directory contains the scaffold for creating a pipe. Pipes use the [plumber](https://gitlab.kilic.dev/libraries/plumber) framework to create a CLI and execute commands in a specific order.
