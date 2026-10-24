# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
While the major version is 0, minor releases may contain breaking changes;
they are always listed under a **Breaking Changes** heading.

## [Unreleased]

### Breaking Changes

- `Response.Body` is now an `io.ReadCloser` streamed from the server instead of
  a `[]byte`, like `http.Response.Body`. Callers of methods that return a
  `*Response` (`File.DownloadFile`, `Request.Send`, `BatchRequest.ExecuteBatch`)
  must close the body, or call the new `Response.ReadBody()` to read and close
  it in one step. Large downloads no longer need to fit in memory.
- `ApiOtherError` and `ApiStatusError` no longer implement `fmt.Formatter` /
  `xerrors.Formatter`: the `Format(fmt.State, rune)` and
  `FormatError(xerrors.Printer) error` methods are removed. `%v`, `%+v` and `%s`
  now all print the same text as `Error()`; the caller frame that `%+v` used to
  show is gone. `Error()`, `Unwrap()`, `errors.Is` and `errors.As` behave as before.

### Changed

- Dropped the dependency on `golang.org/x/xerrors`. Internal error wrapping uses
  `errors.New` / `fmt.Errorf` with `%w`.
- `%v` on an `ApiOtherError` no longer repeats the wrapped cause twice.
- When response body logging is enabled, the JSON body is logged in a separate
  `ResponseBody` entry once it has been read, after the status and headers.

### Added

- `Response.ReadBody()` reads the whole body, closes it and caches the bytes;
  repeated calls return the same result.
- Tests for error formatting, unwrapping, `NewApiStatusError`, `Response.ReadBody`
  and streamed downloads.

### Fixed

- The retry loop did not close the body of a 429 / 5xx response before retrying,
  leaking the connection until the transport reclaimed it.

## [0.1.0] - 2026-10-24

First tagged release.

### Breaking Changes

- Requires **Go 1.26 or later** (the latest `golang.org/x/*` modules require it).

### Changed

- `APIConn.MaxRequestAttempts` is now honored by the retry loop. Retries were
  previously fixed at 5 attempts regardless of the field. The default stays 5
  and values below 1 are treated as the default. (#10)
- `ParseResource` selects the resource type from the top-level `type` field
  only. Previously a nested object's `type` appearing earlier in the JSON could
  be picked up instead. (#11)
- Resource methods share common send / status check / decode helpers; the
  caller frame recorded in `ApiStatusError` now points at the helper. (#8)
- `interface{}` is written as `any`; `math/rand` replaced by `math/rand/v2`. (#8)

### Fixed

- `SharedLink.String()` printed `DownloadCount` / `PreviewCount` as a single
  rune instead of a decimal number. (#9)
- Batch requests (`ExecuteBatch`) produced invalid JSON when a header value
  contained `"` or `\`; headers are now encoded with `json.Marshal`. (#9)
- `ExecuteBatch` ignored errors while reading the response body and parsed a
  possibly truncated body; it now returns an `ApiOtherError`. (#9)
- The error returned when a JWT assertion could not be created did not wrap
  the underlying error (`%w` was missing). (#7)

### Dependencies

- `github.com/golang-jwt/jwt` v4 → v5. (#7)
- `github.com/google/go-cmp` 0.6.0 → 0.7.0 (tests only). (#7)
- `golang.org/x/crypto`, `golang.org/x/sys`, `golang.org/x/text` and the viper
  transitive dependencies updated to their latest versions. (#7)
- Archived `github.com/mitchellh/go-homedir` replaced by `os.UserHomeDir`. (#7)
- Removed the deprecated `rand.Seed` calls. (#7)

### Documentation

- README: GitHub Actions / pkg.go.dev badges, OAuth 2.0 and JWT usage, Go
  version requirement, `MaxRequestAttempts` configuration. (#8, #10)
- CI matrix updated to Go 1.26 and 1.27. (#7)

## [0.0.1]

Initial development version (never tagged). Files, folders, collaborations,
users, groups, memberships, events, file locks, batch requests, OAuth 2.0 and
JWT authentication.

[Unreleased]: https://github.com/jparound30/goboxer/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/jparound30/goboxer/releases/tag/v0.1.0
[0.0.1]: https://github.com/jparound30/goboxer/commit/0273273
