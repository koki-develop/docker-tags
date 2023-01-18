# docker-tags

Command line tool to get a list of tags for docker images.
It can also be used as a docker cli plugin.

# Installation

## Homebrew

```sh
$ brew install koki-develop/tap/docker-tags
```

## go install

```sh
$ go install github.com/koki-develop/docker-tags@latest
```

## docker cli plugin

```sh
$ git clone https://github.com/koki-develop/docker-tags.git
$ cd docker-tags
$ make
$ mkdir -p $HOME/.docker/cli-plugins/
$ mv ./dist/docker-tags $HOME/.docker/cli-plugins/
```

# Usage

```sh
Command line tool to get a list of tags for docker images.

Usage:
  docker-tags [IMAGE] [flags]

Flags:
      --aws-profile string   aws profile
  -h, --help                 help for docker-tags
  -v, --version              version for docker-tags
```

```sh
$ docker-tags alpine
latest
edge
3.9.6
3.9.5
3.9.4
3.9.3
3.9.2
3.9
3.8.5
3.8.4
3.8
...
```

# Supported Registry

- [Docker Hub](https://hub.docker.com/)
- [Amazon ECR](https://aws.amazon.com/ecr/)
- [Amazon ECR Public](https://docs.aws.amazon.com/AmazonECR/latest/public/index.html)
- [Google Container Registry](https://cloud.google.com/container-registry)
- [Google Artifact Registry](https://cloud.google.com/artifact-registry)

# LICENSE

[LICENSE](./LICENSE)
