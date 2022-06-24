[![GitHub Actions](https://github.com/koki-develop/docker-tags/actions/workflows/main.yml/badge.svg)](https://github.com/koki-develop/docker-tags/actions/workflows/main.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/bb48af807bf1c90e5c05/maintainability)](https://codeclimate.com/github/koki-develop/docker-tags/maintainability)
[![Twitter Follow](https://img.shields.io/twitter/follow/koki_develop?style=social)](https://twitter.com/koki_develop)

# Overview

CLI tool to output list of tags for Docker images.

```
$ docker-tags golang
golang:latest
golang:1
golang:1-alpine
golang:1-alpine3.10
golang:1-alpine3.11
golang:1-alpine3.12
golang:1-alpine3.13
golang:1-alpine3.14
golang:1-alpine3.15
golang:1-alpine3.16
...
```

# Installation

```
go install github.com/koki-develop/docker-tags@latest
```

# Usage

```
Usage:
  docker-tags IMAGE [flags]

Flags:
  -h, --help        help for docker-tags
  -t, --only-tags   show tags only
```

# LICENSE

[MIT](./LICENSE)
