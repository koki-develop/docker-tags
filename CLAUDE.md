# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Testing
- `go test ./... -race -coverprofile=coverage.out -covermode=atomic` - Run tests with race detection and coverage
- `go test ./internal/printers/... -v` - Run specific package tests with verbose output
- `go tool cover -html=coverage.out -o coverage.html` - Generate HTML coverage report
- `go build` - Build the CLI tool
- `make build-cli-plugin` - Build and install as Docker CLI plugin (installs to ~/.docker/cli-plugins/)
- `golangci-lint run ./...` - Run linting (requires golangci-lint from mise.toml)
- `goreleaser check` - Validate goreleaser configuration  
- `goreleaser release --snapshot --clean` - Build release artifacts locally (not used in CI)
- `go run . <image>` - Test CLI with specific Docker image (e.g., `go run . alpine`)
- `go run . <image> --output json` - Test with different output formats

### Tool Management
- `mise install` - Install tools specified in mise.toml (Go 1.24.3, golangci-lint, goreleaser)

### Code Formatting and Modernization
- `goimports -w .` - Format Go code and organize imports
- `modernize -test -fix ./...` - Apply modern Go patterns (e.g., interface{} → any)

## Architecture Overview

This is a CLI tool that fetches Docker image tags from various container registries. The tool can run standalone or as a Docker CLI plugin.

### Core Components

**Registry Abstraction** (`internal/registry/`):
- Common `Registry` interface with `ListTags(name string) ([]string, error)` method
- Registry selection based on domain detection in `registry.New()`
- Supported registries: Docker Hub, Amazon ECR/ECR Public, Google Container Registry, GitHub Container Registry, Google Artifact Registry
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
- `ghcr.io` → GitHub Container Registry
- `*-docker.pkg.dev` → Google Artifact Registry

### Code Consistency Patterns
- Import organization: standard library → third-party packages → local packages
- Comments should explain "why" not "what" (e.g., "// Reverse tags to show most recent first")
- Use modern Go patterns: `any` instead of `interface{}` for Go 1.18+
- Test naming conventions: `Test<StructName>_<MethodName>` for methods, `Test_<functionName>` for functions

### AWS Integration
Registry implementations for ECR services use AWS SDK with configurable profile support via `--aws-profile` flag.

### Docker CLI Plugin Mode
When built with `-ldflags "-X github.com/koki-develop/docker-tags/cmd.cliPlugin=true"`, the tool integrates as `docker tags` command using Docker CLI plugin framework.

## Registry Implementation Guidelines

### Adding New Registry Support
When adding support for a new container registry:

1. **Create registry package** in `internal/registry/{name}/` with `Registry` struct implementing the `Registry` interface
2. **Update domain detection** in `internal/registry/registry.go` to include new domain pattern
3. **Follow error handling pattern** - use `io.ReadAll(resp.Body)` and return response content as error for non-OK HTTP responses (see `dockerutil.Client.do()`)
4. **Authentication approaches**:
   - **Anonymous/Public**: Direct token requests (GHCR, Docker Hub pattern)
   - **Cloud Provider**: Use provider-specific clients then Docker registry tokens (GCR, ECR pattern)
   - **dockerutil.Client**: Only when registry follows exact Docker Hub token format with `service` parameter
5. **Variable naming**: Use consistent `token` naming for authentication tokens, avoid abbreviations like `tkn`
6. **Error handling**: Follow existing patterns - use `io.ReadAll(resp.Body)` for HTTP error responses

## Testing Guidelines

### Test Coverage Strategy
- Focus on unit tests for core business logic in `internal/` packages
- Current coverage: ~27% overall, with high coverage on critical components:
  - `internal/printers/`: 91.2% coverage (output formatting)
  - `internal/registry/`: 100% coverage (domain detection logic)
  - `internal/util/dockerutil/`: 87.5% coverage (HTTP client utilities)
- Individual registry implementations are not unit tested (require complex HTTP mocking)
- Uses `github.com/stretchr/testify` for assertions and test utilities

### Test Implementation Patterns
- Use table-driven tests for multiple scenarios
- Test both success and error cases for critical functions
- Focus on testing interfaces and public APIs rather than implementation details
- Use `assert.IsType()` for concrete type verification in factory patterns
- Use `assert.Implements()` to verify interface compliance

### Conventional Commits
This project uses conventional commit format: `type: description`
- `feat:` for new features
- `fix:` for bug fixes  
- `docs:` for documentation
- `ci:` for CI/CD changes
- `refactor:` for code refactoring
- `test:` for adding or updating tests