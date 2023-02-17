# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [Go's conventions for module version numbering](https://go.dev/doc/modules/version-numbers).

## [Unreleased]
### Changed
- Response data types now embed `*http.Response` instead of type aliasing it, for a simplicity and extensibility.

## [0.0.1] - 2023-02-17
### Added
- Implementation of Barentswatch HTTP API as [documented](https://live.ais.barentswatch.net/index.html#/), both streams and queries.
- Go-native data types for API schemas
- Country code, ship type and model format data types, with additional facility methods for human-readableness
- Licence, changelog and readme files.
 
###  Fixed
- Discrepancies between published swagger documentation and actual API behavior adjusted for, and verified in private communication with API developers.

[unreleased]: https://github.com/ilder-as/go-barentswatch-ais/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/ilder-as/go-barentswatch-ais/releases/tag/v0.0.1