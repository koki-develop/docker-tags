# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Testing
- `go test ./... -race -coverprofile=coverage.out -covermode=atomic` - Run tests with race detection and coverage
- `go build` - Build the CLI tool
- `make build-cli-plugin` - Build and install as Docker CLI plugin (installs to ~/.docker/cli-plugins/)
- `golangci-lint run ./...` - Run linting (requires golangci-lint from mise.toml)
- `goreleaser check` - Validate goreleaser configuration  
- `goreleaser release --snapshot --clean` - Build release artifacts locally

### Tool Management
- `mise install` - Install tools specified in mise.toml (Go 1.24.3, golangci-lint, goreleaser)

## Architecture Overview

This is a CLI tool that fetches Docker image tags from various container registries. The tool can run standalone or as a Docker CLI plugin.

### Core Components

**Registry Abstraction** (`internal/registry/`):
- Common `Registry` interface with `ListTags(name string) ([]string, error)` method
- Registry selection based on domain detection in `registry.New()`
- Supported registries: Docker Hub, Amazon ECR/ECR Public, Google Container Registry, Google Artifact Registry
- Each registry implementation handles authentication and API communication differently

**Output Formatting** (`internal/printers/`):
- Plugin-style printer system using `Register()` pattern
- Common `Printer` interface with `Print(w io.Writer, params *PrintParameters) error`
- Formats: text, json, yaml (registered at package init)
- `WithName` option prefixes tags with image name

**CLI Structure** (`cmd/`):
- Single root command using Cobra
- Dual-mode operation: standalone CLI or Docker plugin (controlled by build-time `cliPlugin` flag)
- Image parsing using Docker's reference library to extract domain and path

### Registry Domain Detection Logic
- `docker.io` → Docker Hub
- `public.ecr.aws` → ECR Public  
- `*.amazonaws.com` → Private ECR
- `gcr.io` → Google Container Registry
- `*-docker.pkg.dev` → Google Artifact Registry

### AWS Integration
Registry implementations for ECR services use AWS SDK with configurable profile support via `--aws-profile` flag.

### Docker CLI Plugin Mode
When built with `-ldflags "-X github.com/koki-develop/docker-tags/cmd.cliPlugin=true"`, the tool integrates as `docker tags` command using Docker CLI plugin framework.