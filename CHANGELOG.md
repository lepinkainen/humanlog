# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Configurable message width and JSON output support
- Project infrastructure (.gitignore, .golangci.yml, Taskfile)
- Improved logging attributes handling

### Changed
- Gemini refactor

## [0.1.0] - 2025-07-18

### Added
- Initial release
- Human-readable formatter for slog
- Fixed-width message formatting (40 characters)
- Colored level output
- Structured attribute support with grouping
- Time formatting in RFC3339
- Attribute quoting for strings with spaces
- Basic documentation and examples

[Unreleased]: https://github.com/lepinkainen/humanlog/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/lepinkainen/humanlog/releases/tag/v0.1.0
