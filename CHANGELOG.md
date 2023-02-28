# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [Go's conventions for module version numbering](https://go.dev/doc/modules/version-numbers).

## [Unreleased]
### Added 
- EOF error which signifies end of stream for `StreamResponse` when using `UnmarshalStream`, as well as an `IsEOF` function to check whether an error is an EOF error.
- Tests
- CONTRIBUTING.md
 
### Changed
- The `SSEStreamResponse[T]` data type is merged into and replaced by the `StreamResponse[T]` data type.
- The `Combined` data type is replaced by `CombinedSimpleJson` since these are identical.
- Response data types now embed `*http.Response` instead of type aliasing it, for a simplicity and extensibility.
- Signatures of `UnmarshalStream` on `StreamResponse` and `SSEStreamResponse` have changed substantially. They no longer supply an error channel, errors are instead provided when the data channel closes, and can be checked with the `Error()` method. They also do not return cancellation functions anymore, cancellation is instead achieved by the user cancelling the request's context. 
- `NewClient` now optionally takes a set of URLs for the HTTP endpoints.
- Query options are moved to new options package.

### Fixed
- Incorrect query parameters generated in HTTP request in GetLatestAIS/GetLatestAISContext and GetLatestCombined/GetLatestCombinedContext

## [0.0.1] - 2023-02-17
### Added
- Implementation of Barentswatch HTTP API as [documented](https://live.ais.barentswatch.net/index.html#/), both streams and queries.
- Go-native data types for API schemas
- Country code, ship type and model format data types, with additional facility methods for human-readableness
- License, changelog and readme files.
 
###  Fixed
- Discrepancies between published swagger documentation and actual API behavior adjusted for, and verified in private communication with API developers.

[unreleased]: https://github.com/ilder-as/go-barentswatch-ais/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/ilder-as/go-barentswatch-ais/releases/tag/v0.0.1
