# docker-tags

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/koki-develop/docker-tags)](https://github.com/koki-develop/docker-tags/releases/latest)
[![GitHub](https://img.shields.io/github/license/koki-develop/docker-tags)](./LICENSE)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/koki-develop/docker-tags/ci.yml?logo=github)](https://github.com/koki-develop/docker-tags/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/koki-develop/docker-tags)](https://goreportcard.com/report/github.com/koki-develop/docker-tags)

Command line tool to get a list of tags for docker images.
It can also be used as a docker cli plugin.

# Supported Registry

> [!NOTE]
> For the [Amazon ECR](https://aws.amazon.com/ecr/) and [ECR Public](https://docs.aws.amazon.com/AmazonECR/latest/public/index.html), an AWS Profile must be configured.  
> See [documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) for details.

- [Docker Hub](https://hub.docker.com/)
- [Amazon ECR](https://aws.amazon.com/ecr/)
- [Amazon ECR Public](https://docs.aws.amazon.com/AmazonECR/latest/public/index.html)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Google Artifact Registry](https://cloud.google.com/artifact-registry)

# Installation

## Homebrew

```sh
$ brew install koki-develop/tap/docker-tags
```

## go install

```sh
$ go install github.com/koki-develop/docker-tags@latest
```

## Docker CLI Plugin

```sh
$ git clone https://github.com/koki-develop/docker-tags.git
$ cd docker-tags
$ make
$ docker tags --help
```

## Release

Download the binary from the [releases page](https://github.com/koki-develop/docker-tags/releases/latest).

# Usage

```sh
Command line tool to get a list of tags for docker images.

Usage:
  docker-tags [IMAGE] [flags]

Flags:
      --aws-profile string   aws profile
  -h, --help                 help for docker-tags
  -o, --output string        output format (text|json) (default "text")
  -v, --version              version for docker-tags
  -n, --with-name            print with image name
```

```sh
$ docker-tags alpine
latest
edge
3.9.6
3.9.5
...
```

```sh
# json format
$ docker-tags alpine -o json
[
  "latest",
  "edge",
  "3.9.6",
  "3.9.5",
  ...
]
```

# LICENSE

[LICENSE](./LICENSE)
