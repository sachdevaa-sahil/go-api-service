# Architecture

## Startup and request flow

`cmd/api/main.go` loads configuration, creates and pings one MongoDB client, selects the configured database, constructs the user service and handler, builds the router, and starts the HTTP server.

```text
HTTP request → ServeMux → UserHandler → user.Service → MongoDB users collection
```

All user endpoints call one service. `user.Service` owns validation, password hashing, timestamps, and MongoDB reads/writes. The handler handles HTTP parsing and responses. A single `UserService` interface at the handler lets HTTP tests use a fake service without a database.

## File map

| Location | Responsibility |
| --- | --- |
| `cmd/api/main.go` | Dependency wiring, startup errors, deferred disconnect |
| `internal/config/config.go` | Optional dotenv loading and required-variable validation |
| `internal/http_mod/server.go` | Server construction and HTTP timeouts |
| `internal/http_mod/router.go` | Method/path registration |
| `internal/http_mod/handler/user.go` | User service interface, request deadlines, ID validation, HTTP errors |
| `internal/http_mod/handler/response.go` | JSON encoding before headers, response writes |
| `internal/http_mod/handler/health.go` | Process health response |
| `internal/user/model.go` | User and creation-input structs with BSON/JSON mappings |
| `internal/storage/mongodb/db.go` | Client creation, bounded ping, failed-connect cleanup |
| `internal/user/service.go` | List, lookup, creation, validation, bcrypt hashing, IDs, timestamps, and shared errors |
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

Dates must be BSON Date values, not nested objects containing a literal `$date` key. The creation service validates input, hashes passwords with bcrypt, and assigns UTC creation/update timestamps. These operations run for `POST /users`. BSON tags describe mappings; they do not enforce a database schema or uniqueness. Read queries currently retrieve the password field, but JSON serialization excludes it.

## Context and lifecycle

User handlers derive a five-second context from the request and pass it through to service database operations. Creation applies the deadline to the service call; bcrypt hashing itself is not interruptible by context cancellation. The list query and cursor decoding share that deadline. `cursor.All` closes its cursor. The MongoDB client is shared across requests and disconnected with a fresh bounded context when `run()` returns.

`main()` exits with status 1 for an error returned by `run()`. No signal handler calls `server.Shutdown` yet, so Ctrl+C is not a verified graceful cleanup path. Server timeout values are defined in `internal/http_mod/server.go`; query and connection timeouts live at their operation boundaries.

See [routes.md](../routes.md) for external behavior and [development.md](development.md) for verification.
