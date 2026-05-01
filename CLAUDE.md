# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About

`goboxer` is a Go client library for the [Box API](https://developer.box.com/). Module: `github.com/jparound30/goboxer`. It is under active development (v0.0.1) and the API may change destructively.

## Commands

```bash
# Build
go build -v ./...

# Run all tests with coverage
go test -cover ./...

# Run a single test
go test -run TestFunctionName ./...

# Run tests for a specific file (e.g., files_test.go)
go test -run TestFile ./...

# Run with verbose output
go test -v -run TestFunctionName ./...
```

There is no Makefile, linter config, or lint step in CI. CI runs `go test -cover` via GitHub Actions (`.github/workflows/go.yml`) on push and pull requests.

## Architecture

This is a flat-package library — all core code lives in the root `github.com/jparound30/goboxer` package. There are no sub-packages for the library itself; the `cmd/` directory contains standalone CLI tools built on top of the library.

### Core request pipeline

`APIConn` (`apiconn.go`) is the central struct that holds credentials and HTTP state. All API operations are methods on resource structs (e.g., `File`, `Folder`, `User`) but take an `*APIConn` to execute. Actual HTTP execution flows through `Request` (`request.go`), which handles:
- Bearer token injection
- User impersonation via the `As-User` header
- Auto-retry on HTTP 429 (rate limit) and 5xx with exponential backoff
- Redirect following (max 3)

Token lifecycle is managed in `APIConn`: access tokens auto-refresh 60 seconds before expiry using either OAuth2 refresh tokens or JWT-signed assertions. Token refresh events are broadcast to registered listeners via a notification channel. All token reads/writes are protected by `sync.RWMutex`.

### Authentication modes

1. **OAuth2** — standard user-delegated auth with access + refresh tokens. `cmd/goboxertokengen` runs a local HTTPS server to handle the OAuth callback and print the resulting tokens.
2. **JWT** — service account or user-level auth. Config loaded from a Box-generated JSON file via `jwtconfigloader.go`. JWT signing uses PKCS#8 keys (`golang.org/x/crypto`, `youmark/pkcs8`). The `cmd/sandbox` directory shows a usage example.

### Resource structs

Each Box resource type has a dedicated file:

| File | Resource | Notable |
|---|---|---|
| `files.go` | `File` | Upload, download, lock/unlock, versioning, shared links |
| `folders.go` | `Folder` | CRUD, item listing, upload email, collaborations |
| `collaborations.go` | `Collaboration` | Role-based access, pending/accepted/rejected states |
| `users.go` | `User` | App users, enterprise enumeration, impersonation |
| `groups.go` | `Group` | Enterprise groups CRUD |
| `membership.go` | `Membership` | Links users to groups |
| `events.go` | `Event` | User and enterprise event streams, polling |
| `filelocks.go` | `FileLock` | File lock/unlock |

### Test fixtures

Tests use JSON fixtures under `testdata/` (one subdirectory per resource type). Tests mock HTTP responses by intercepting at the `http.Transport` level rather than spinning up a real server.

### CLI tools (`cmd/`)

| Directory | Purpose |
|---|---|
| `cmd/goboxer` | General-purpose Box CLI (Cobra-based) |
| `cmd/goboxertokengen` | OAuth2 token generator (requires self-signed certs in `cert/`) |
| `cmd/parallel` | Batch/parallel API operation utilities |
| `cmd/sandbox` | JWT auth example |

### Error handling

`goboxererror.go` defines the custom error type returned from API calls. Callers should type-assert to inspect Box-specific error codes.
