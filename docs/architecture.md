# Architecture

## Startup and request flow

`cmd/api/main.go` loads configuration, creates and pings one MongoDB client, selects the configured database, constructs the repository and handler, builds the router, and starts the HTTP server.

```text
HTTP request → ServeMux → UserHandler → UserRepository → MongoDB users collection
```

The handler's repository interface makes HTTP tests independent of a running database. There is currently no service layer between handler and repository.

## File map

| Location | Responsibility |
| --- | --- |
| `cmd/api/main.go` | Dependency wiring, startup errors, deferred disconnect |
| `internal/config/config.go` | Optional dotenv loading and required-variable validation |
| `internal/http_mod/server.go` | Server construction and HTTP timeouts |
| `internal/http_mod/router.go` | Method/path registration |
| `internal/http_mod/handler/user.go` | Repository interface, request deadlines, ID validation, HTTP errors |
| `internal/http_mod/handler/response.go` | JSON encoding before headers, response writes |
| `internal/http_mod/handler/health.go` | Process health response |
| `internal/user/model.go` | User struct with BSON and JSON mappings |
| `internal/user/repository.go` | Shared `ErrNotFound` sentinel |
| `internal/storage/mongodb/db.go` | Client creation, bounded ping, failed-connect cleanup |
| `internal/storage/mongodb/user.go` | List and `_id` lookup queries; maps missing documents to `ErrNotFound` |
| `internal/user/service.go` | Placeholder, not an implemented business layer |
| `internal/http_mod/middleware/` | Placeholder files; no active auth or logging middleware |

Paths above are relative to the repository root. Keep the existing `http_mod` package naming unless renaming is part of the requested work.

## Stored user fields

| MongoDB field | Go type | JSON field |
| --- | --- | --- |
| `_id` | `bson.ObjectID` | `id` |
| `name` | `string` | `name` |
| `email` | `string` | `email` |
| `password` | `string` | Excluded |
| `createdAt` | `time.Time` | `createdAt` |
| `updatedAt` | `time.Time` | `updatedAt` |

Dates must be BSON Date values, not nested objects containing a literal `$date` key. Password hashing, input validation for creation, and automatic timestamps are not implemented. BSON tags describe mappings; they do not enforce a database schema or uniqueness. Repositories currently retrieve the password field, but JSON serialization excludes it.

## Context and lifecycle

User handlers derive a five-second context from the request and pass it through to repository operations. The list query and cursor decoding share that deadline. `cursor.All` closes its cursor. The MongoDB client is shared across requests and disconnected with a fresh bounded context when `run()` returns.

`main()` exits with status 1 for an error returned by `run()`. No signal handler calls `server.Shutdown` yet, so Ctrl+C is not a verified graceful cleanup path. Server timeout values are defined in `internal/http_mod/server.go`; query and connection timeouts live at their operation boundaries.

See [routes.md](../routes.md) for external behavior and [development.md](development.md) for verification.
