# pipes

[![pipeline status](https://gitlab.kilic.dev/devops/pipes/badges/master/pipeline.svg)](https://gitlab.kilic.dev/devops/pipes/-/commits/master)

## Description

Automation tools for various pipelines.

<!-- toc -->

- [Description](#description)
- [Pipes](#pipes-1)
  - [cenk1cenk2/pipe-docker](#cenk1cenk2pipe-docker)
  - [cenk1cenk2/pipe-gl-artifacts](#cenk1cenk2pipe-gl-artifacts)
  - [cenk1cenk2/pipe-markdown-toc](#cenk1cenk2pipe-markdown-toc)
  - [cenk1cenk2/pipe-node](#cenk1cenk2pipe-node)
  - [cenk1cenk2/pipe-s3-upload](#cenk1cenk2pipe-s3-upload)
  - [cenk1cenk2/select-env](#cenk1cenk2select-env)
  - [cenk1cenk2/pipe-semantic-release](#cenk1cenk2pipe-semantic-release)
  - [cenk1cenk2/update-docker-hub-readme](#cenk1cenk2update-docker-hub-readme)
- [Templates](#templates)
- [Methodology](#methodology)

<!-- tocstop -->

--

## Pipes

### cenk1cenk2/pipe-docker

[Read more...](./docker/README.md)

### cenk1cenk2/pipe-gl-artifacts

[Read more...](./gl-artifacts/README.md)

### cenk1cenk2/pipe-markdown-toc

[Read more...](./markdown-toc/README.md)

### cenk1cenk2/pipe-node

[Read more...](./node/README.md)

### cenk1cenk2/pipe-s3-upload

[Read more...](./s3-upload/README.md)

### cenk1cenk2/select-env

[Read more...](./select-env/README.md)

### cenk1cenk2/pipe-semantic-release

[Read more...](./semantic-release/README.md)

### cenk1cenk2/update-docker-hub-readme

[Read more...](./update-docker-hub-readme/README.md)

## Templates

Templates contains gitlab-ci templates to extend basic tasks from.

[Read more...](./templates)

## Methodology

`_template` folder contains the basic setup for creating a pipe, these pipes use the [plumber](https://gitlab.kilic.dev/libraries/plumber) framework to create a cli and execute the commands in a specific order.
