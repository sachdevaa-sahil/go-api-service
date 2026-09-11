# Go API service

A Go HTTP API using the standard library. User data is currently dummy data; no database is connected and no writes are persisted.

## Run

Use the Go version specified in `go.mod`, then run from the project root:

```sh
go run ./cmd/api
```

The server listens on port `8000`. Stop it with Ctrl+C. If that port is occupied, stop the existing process or change the address in `cmd/api/main.go`.

## Endpoints

| Method | Path | Response |
| --- | --- | --- |
| GET | `/health` | Health status |
| GET | `/hello` | Plain-text greeting |
| GET | `/users` | Array of user objects |
| GET | `/users/{id}` | User object, or JSON error with status 404 |

```sh
curl -i http://localhost:8000/health
curl -i http://localhost:8000/users
curl -i http://localhost:8000/users/2
```

User IDs are `1` through `4`. Unknown routes return 404 and unsupported methods return 405 using the standard router's plain-text responses.

## Structure

- `cmd/api`: application entry point.
- `internal/http_mod`: server configuration, routes, and handlers.
- `internal/config`, `internal/user`, `internal/storage/mongodb`, and middleware files: placeholders for later work.

## Checks

```sh
go test ./...
go vet ./...
```
