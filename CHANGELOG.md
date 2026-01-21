# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0-beta](https://github.com/pixelfactory-go/observability-trace/compare/v0.4.2-beta...v0.5.0-beta) (2026-01-21)


### Features

* add SAST workflows (CodeQL and Scorecard) ([b9ae8e2](https://github.com/pixelfactory-go/observability-trace/commit/b9ae8e2e5473da74cf4e7b7ca98f97cb3d2c3648))
* add security policy and enhance dependency updates ([13770db](https://github.com/pixelfactory-go/observability-trace/commit/13770db86c35644c45fad02e96da1c3d399c4455))
* enhance CI with security improvements and govulncheck ([5bfd9d7](https://github.com/pixelfactory-go/observability-trace/commit/5bfd9d71a13fca52be0cef7574302e7cc33f422a))
* update CI and workflows for improved security and dependency management ([3f00c9e](https://github.com/pixelfactory-go/observability-trace/commit/3f00c9e81ddb67438de377f3bbcaaba52e668027))


### Bug Fixes

* add names to CI jobs for clarity ([bb821f5](https://github.com/pixelfactory-go/observability-trace/commit/bb821f5f5a42fcda4f626269e88fd90c901dc9fd))
* add names to CI jobs for clarity ([c0d01be](https://github.com/pixelfactory-go/observability-trace/commit/c0d01be09415e4b4c66be9f81e75d9d726276059))

## [0.4.2-beta](https://github.com/pixelfactory-go/observability-trace/compare/v0.4.1-beta...v0.4.2-beta) (2026-01-20)


### Bug Fixes

* add missing goreleaser config ([96a1c23](https://github.com/pixelfactory-go/observability-trace/commit/96a1c237f6ac92ee63d2a1a061450819fc93309f))
* add missing goreleaser config ([8cfb8c4](https://github.com/pixelfactory-go/observability-trace/commit/8cfb8c465add573a4c5123965ef92aab7e57c2d3))

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
