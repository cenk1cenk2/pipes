# pipes

[![pipeline status](https://gitlab.kilic.dev/devops/pipes/badges/main/pipeline.svg)](https://gitlab.kilic.dev/devops/pipes/-/commits/main)

Automation pipe CLIs used by the shared GitLab CI templates in [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines).

## Consumer CI templates

Consumer-facing GitLab CI templates have moved to [`devops/pipelines`](https://gitlab.kilic.dev/devops/pipelines). New and migrated repositories should include versioned templates from that repository instead of including files from `devops/pipes`.

## Pipe images

This repository still builds and publishes the pipe images used by `devops/pipelines`:

- [`cenk1cenk2/pipe-buildah`](./pipes/buildah/README.md)
- [`cenk1cenk2/pipe-go`](./pipes/go/README.md)
- [`cenk1cenk2/pipe-helm`](./pipes/helm/README.md)
- [`cenk1cenk2/pipe-kustomize`](./pipes/kustomize/README.md)
- [`cenk1cenk2/pipe-node`](./pipes/node/README.md)
- [`cenk1cenk2/pipe-pulumi`](./pipes/pulumi/README.md)
- [`cenk1cenk2/pipe-select-env`](./pipes/select-env/README.md)
- [`cenk1cenk2/pipe-semantic-release`](./pipes/semantic-release/README.md)
- [`cenk1cenk2/pipe-terraform`](./pipes/terraform/README.md)
- [`cenk1cenk2/pipe-update-docker-hub-readme`](./pipes/update-docker-hub-readme/README.md)

## Methodology

The `template` directory contains the scaffold for creating a pipe, and every pipe lives under `pipes/`. Pipes use the [plumber](https://gitlab.kilic.dev/libraries/plumber) framework to create a CLI and execute commands in a specific order.
