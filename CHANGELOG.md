# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.1-beta](https://github.com/pixelfactory-go/observability-trace/compare/v0.4.0-beta...v0.4.1-beta) (2026-01-20)


### Bug Fixes

* resolve all linter errors and fix Windows CI ([1ae9191](https://github.com/pixelfactory-go/observability-trace/commit/1ae9191f6c1151965dea9e81dda3e9dc0362f528))

## [Unreleased]

### Added
- Comprehensive README with usage examples and configuration guide
- CONTRIBUTING.md with development guidelines
- CI/CD workflows for automated testing and releases
- golangci-lint configuration for code quality
- Makefile for common development tasks

### Changed
- Updated to latest stable Go version and dependencies
- Enhanced documentation across all files

## [0.4.0] - 2024-03-15

### Added
- Span creation helper functions (`NewSpan`, `SpanFromContext`)
- Span utility functions (`AddSpanTags`, `AddSpanEvents`, `AddSpanError`, `FailSpan`)
- SpanCustomiser interface for custom span options

### Changed
- Refactored example to include HTTP client usage

## [0.3.0] - 2024-03-10

### Added
- HTTP instrumentation helpers (`HTTPHandler`, `HTTPHandlerFunc`, `HTTPClientTransporter`)
- HTTP client transport wrapper for tracing
- WithTraceExporter option for custom exporters

### Fixed
- WithTraceExporter comment documentation

## [0.2.0] - 2024-03-05

### Added
- Configuration refactoring for better usability
- Environment variable support for all configuration options
- Multiple propagator support (B3, TraceContext, Baggage, OT)

## [0.1.0] - 2024-03-01

### Added
- Initial release
- Basic OpenTelemetry trace provider setup
- OTLP gRPC exporter support
- Configuration via options pattern
- Resource attributes and service name/version support

[Unreleased]: https://github.com/pixelfactory-go/observability-trace/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/pixelfactory-go/observability-trace/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/pixelfactory-go/observability-trace/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pixelfactory-go/observability-trace/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pixelfactory-go/observability-trace/releases/tag/v0.1.0
