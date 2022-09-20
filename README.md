# pipes

[![pipeline status](https://gitlab.kilic.dev/devops/pipes/badges/master/pipeline.svg)](https://gitlab.kilic.dev/devops/pipes/-/commits/master)

## Description

Automation tools for various pipelines.

<!-- toc -->

- [Containers](#containers)
  - [gitlab-pipes-docker](#gitlab-pipes-docker)
  - [cenk1cenk2/pipe-gh-release-tracker](#cenk1cenk2pipe-gh-release-tracker)
  - [gitlab-pipes-gl-artifacts](#gitlab-pipes-gl-artifacts)
  - [cenk1cenk2/pipe-markdown-toc](#cenk1cenk2pipe-markdown-toc)
  - [gitlab-pipes-node](#gitlab-pipes-node)
  - [cenk1cenk2/pipe-s3-upload](#cenk1cenk2pipe-s3-upload)
  - [cenk1cenk2/select-env](#cenk1cenk2select-env)
  - [cenk1cenk2/pipe-semantic-release](#cenk1cenk2pipe-semantic-release)
  - [cenk1cenk2/update-docker-hub-readme](#cenk1cenk2update-docker-hub-readme)
- [Templates](#templates)
- [Methodology](#methodology)

<!-- tocstop -->

--

## Containers

### gitlab-pipes-docker

[Read more...](./docker/README.md)

### cenk1cenk2/pipe-gh-release-tracker

[Read more...](./gh-release-tracker/README.md)

### gitlab-pipes-gl-artifacts

[Read more...](./gl-artifacts/README.md)

### cenk1cenk2/pipe-markdown-toc

[Read more...](./markdown-toc/README.md)

### gitlab-pipes-node

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
