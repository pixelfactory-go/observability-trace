# Contributing to Trace

Thank you for your interest in contributing to this project! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Be respectful, inclusive, and considerate in all interactions.

## How to Contribute

### Reporting Issues

If you find a bug or have a feature request:

1. **Search existing issues** to avoid duplicates
2. **Use issue templates** if available
3. **Provide detailed information**:
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant code snippets or logs

### Submitting Pull Requests

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Ensure all tests pass** locally
6. **Submit a pull request** with a clear description

## Development Setup

### Prerequisites

- **Go**: Version 1.20 or higher
- **Git**: For version control
- **Make**: Optional but recommended for using Makefile targets
- **golangci-lint**: For code quality checks

### Clone and Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/observability-trace.git
cd observability-trace

# Add upstream remote
git remote add upstream https://github.com/pixelfactory-go/observability-trace.git

# Install dependencies
go mod download
```

### Local Development Commands

```bash
# Format code
make fmt
# or
go fmt ./...

# Run linter
make lint
# or
golangci-lint run

# Run tests
make test
# or
go test ./...

# Run tests with coverage
make test-coverage
# or
go test -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage in browser
go tool cover -html=coverage.out

# Build the project
make build
# or
go build ./...
```

## Coding Standards

### Go Guidelines

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html) guidelines:

- Use `gofmt` for formatting (handled by `make fmt`)
- Follow standard Go naming conventions
- Keep functions small and focused
- Prefer clear code over clever code
- Handle errors explicitly

### Documentation

- **Exported functions/types** must have godoc comments
- Comments should explain **why**, not just what
- Keep comments up-to-date with code changes
- Include examples for complex functionality

Example:

```go
// NewProvider creates and initializes an OpenTelemetry trace provider.
// It configures the provider using the provided options and environment variables.
// The returned provider must be shut down via Shutdown() when the application exits.
func NewProvider(opts ...Option) (*Provider, error) {
    // implementation
}
```

### Testing

- Write tests for all new functionality
- Aim for high test coverage (>80%)
- Use table-driven tests where appropriate
- Test edge cases and error conditions

Example table-driven test:

```go
func TestWithServiceName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"valid name", "my-service", "my-service"},
        {"empty name", "", ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var c Config
            WithServiceName(tt.input)(&c)
            if c.ServiceName != tt.expected {
                t.Errorf("got %q, want %q", c.ServiceName, tt.expected)
            }
        })
    }
}
```

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear and semantic commit messages:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **chore**: Maintenance tasks
- **ci**: CI/CD changes
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **perf**: Performance improvements

### Examples

```
feat(http): add timeout configuration for HTTP client

Add support for configuring request timeouts in the HTTP client
wrapper to prevent hanging requests.

Closes #123
```

```
fix(provider): correct propagator initialization order

The propagator configuration was being applied before validation,
causing panics on invalid configurations.

Fixes #456
```

```
docs: update README with custom exporter examples

Add examples showing how to use custom trace exporters beyond
the default OTLP exporter.
```

### Scope

Use scope to indicate the affected component:
- `http` - HTTP instrumentation
- `provider` - Trace provider
- `span` - Span utilities
- `config` - Configuration
- `deps` - Dependencies

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass locally (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits format
- [ ] Branch is up-to-date with `main`

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing performed

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented if unavoidable)
```

### Review Process

1. **Automated checks** must pass (CI, linting, tests)
2. **Code review** by at least one maintainer
3. **Address feedback** promptly
4. **Maintainer approval** required before merge

## Release Process

Releases are automated using Release Please:

1. Commits following conventional commits automatically update CHANGELOG
2. Release Please creates release PRs
3. Maintainers merge release PRs to publish new versions
4. Semantic versioning is automatically determined from commit types

## Need Help?

- **Questions**: Open a discussion or issue
- **Stuck**: Ask for help in your PR
- **Ideas**: Open an issue to discuss before implementing

## License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers this project.

## Recognition

All contributors will be recognized in the project. Thank you for helping improve this library!
