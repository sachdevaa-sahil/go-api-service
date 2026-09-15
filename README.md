# Go API service

A Go HTTP API using the standard library. Both user endpoints read from MongoDB; no write endpoints are implemented.

## Run

Use the Go version specified in `go.mod`. For initial setup, copy `env.example` to `.env` and replace the placeholders with your MongoDB settings (keep an existing `.env`). The application loads `.env` from the working directory; existing environment variables take precedence. Both `MONGODB_URI` and `MONGODB_DATABASE` are required. The application verifies MongoDB connectivity before starting the HTTP server.

Run from the project root:

```sh
go run ./cmd/api
```

The server listens on port `8000`. Stop it with Ctrl+C. If that port is occupied, stop the existing process or change the address in `cmd/api/main.go`.

## Endpoints

| Method | Path | Response |
| --- | --- | --- |
| GET | `/health` | Health status |
| GET | `/hello` | Plain-text greeting |
| GET | `/users` | Up to 100 MongoDB users sorted by `_id`, without password hashes |
| GET | `/users/{id}` | MongoDB user object, or JSON error |

```sh
curl -i http://localhost:8000/health
curl -i http://localhost:8000/users
curl -i http://localhost:8000/users/REPLACE_WITH_OBJECT_ID
```

`GET /users/{id}` accepts a 24-character hexadecimal MongoDB ObjectID. Invalid IDs return 400, missing users return 404, and database failures return 500. An empty MongoDB collection returns `[]` from `/users`. Unknown routes return 404 and unsupported methods return 405 using the standard router's plain-text responses.

## Structure

- `cmd/api`: application entry point.
- `internal/http_mod`: server configuration, routes, and handlers.
- `internal/config`: environment configuration.
- `internal/user`: MongoDB user model, with a shared not-found error and a service placeholder.
- `internal/storage/mongodb`: shared client setup and user queries.
- Middleware files: placeholders for later work.

## Checks

```sh
go test ./...
go vet ./...
```

## Project reference

- [AGENTS.md](AGENTS.md): instructions for coding assistants and collaboration preferences.
- [routes.md](routes.md): endpoint contracts, response fields, and status codes.
- [Architecture](docs/architecture.md): dependency flow, package ownership, and data model.
- [Development](docs/development.md): configuration, verification, and current limitations.
