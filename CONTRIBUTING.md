# Contributing to humanlog

Thank you for your interest in contributing! This document provides guidelines for contributing to the humanlog project.

## Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automatic version bumping and changelog generation.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat:` - New feature (triggers MINOR version bump)
- `fix:` - Bug fix (triggers PATCH version bump)
- `docs:` - Documentation only
- `refactor:` - Code refactoring (no functional changes)
- `perf:` - Performance improvement
- `test:` - Adding or updating tests
- `chore:` - Build process, tooling changes

### Breaking Changes

Add `!` after type or include `BREAKING CHANGE:` in footer (triggers MAJOR version bump):

```
feat!: change Handler interface signature

BREAKING CHANGE: Handler.Handle now requires context.Context
```

### Examples

```
feat: add support for custom time formats
fix: correct attribute quoting for edge cases
docs: update README with installation instructions
refactor: simplify attribute formatting logic
perf: optimize string concatenation in formatter
test: add tests for grouped attributes
chore: update dependencies
```

## Development Workflow

### Setup

1. Clone the repository
2. Install dependencies:
   ```sh
   go mod download
   ```
3. Install required tools:
   ```sh
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Making Changes

1. Create a new branch for your changes
2. Make your changes following the project conventions
3. Write tests for new features
4. Run quality checks:
   ```sh
   task build  # Runs tests and linting
   ```
5. Commit with conventional commit messages
6. Push and create a pull request

### Code Quality

- All code must pass `task build` before submission
- Run `task lint` to check code formatting
- Run `task test` to run tests
- Follow existing code style and patterns

## Release Process

Releases are fully automated. Developers only need to:

1. Ensure commits follow conventional format
2. Run `task build` to verify quality
3. Create and push a tag:

```sh
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

GitHub Actions will:
- Run tests and linting
- Generate CHANGELOG.md
- Create GitHub release with notes
- Update pkg.go.dev automatically

### Version Guidelines

- `v0.x.x` - Pre-release, API may change
- `v1.0.0` - First stable release, semver guarantees
- `v2.0.0+` - Major versions require `/vN` module path suffix

### Semantic Versioning Rules

- **PATCH** (v0.1.0 → v0.1.1): Bug fixes only
- **MINOR** (v0.1.0 → v0.2.0): New features, backward compatible
- **MAJOR** (v0.x.x → v1.0.0 or v1.x.x → v2.0.0): Breaking changes

## Testing Locally

### Test Release Process

```sh
# Test changelog generation
task changelog

# Test GoReleaser without publishing
task release-snapshot

# Pre-release validation
task release-check
```

### Test with Beta Release

```sh
# Create and push beta tag
git tag v0.2.0-beta.1
git push origin v0.2.0-beta.1

# Monitor GitHub Actions
# Delete if successful
git tag -d v0.2.0-beta.1
git push origin :refs/tags/v0.2.0-beta.1
```

## Project Structure

```
humanlog/
├── .github/workflows/  # CI/CD workflows
├── .chglog/           # Changelog configuration
├── example/           # Usage examples
├── handler.go         # Main handler implementation
├── options.go         # Configuration options
├── humanlog.go        # Package entry point
├── handler_test.go    # Tests
├── Taskfile.yml       # Task automation
└── README.md          # Documentation
```

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues and pull requests first
- Be respectful and constructive in discussions

Thank you for contributing!
