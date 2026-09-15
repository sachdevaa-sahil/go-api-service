# Development workflow

## Setup

Run commands from the repository root. Use the toolchain declared in `go.mod`.

For a new checkout, copy `env.example` to `.env` only if `.env` does not already exist, then supply your local MongoDB settings. Do not overwrite credentials during setup.

| Variable | Required | Purpose |
| --- | --- | --- |
| `MONGODB_URI` | Yes | Connection URI, including deployment-specific authentication |
| `MONGODB_DATABASE` | Yes | Database containing the `users` collection |

Dotenv reads from the working directory. Existing process environment variables take precedence. A missing `.env` is permitted when variables are supplied externally. The HTTP port is currently fixed at `:8000` in `cmd/api/main.go`, not configured by an environment variable.

```sh
go run ./cmd/api
```

The MongoDB ping must succeed before HTTP starts. If a port is occupied, identify the process rather than killing an unknown server. Never include credentials or raw user records in diagnostic output.

## Debug in VS Code

Open `go-api-service` as the workspace folder, choose **Debug Go Project** in Run and Debug, and press F5. The configuration in [.vscode/launch.json](../.vscode/launch.json) launches `cmd/api` and keeps the working directory at the project root so the application can load `.env`. It requires the Go extension and Delve. The API uses port 8000 in debug mode too; stop an existing instance before launching another.

## Verify a change

```sh
go test ./...
go vet ./...
git diff --check
```

Format only changed Go files with `gofmt -w <files>`. If a sandbox blocks the default Go build cache, set `GOCACHE` to a writable temporary directory for the command; do not change project configuration for a machine-specific restriction.

- `internal/config/config_test.go`: required variables, dotenv loading, precedence, sanitized parse errors.
- `internal/http_mod/router_test.go`: fake-service HTTP behavior, invalid/missing IDs, query failures, password exclusion, and JSON results.
- `internal/http_mod/create_user_test.go`: HTTP creation with a fake service, including JSON parsing, body limits, error mapping, and password exclusion.
- `internal/user/service_test.go`: real input validation, password hashing, normalization, generated IDs, and timestamps.
- These tests do not establish that MongoDB is reachable or that real stored documents decode successfully. Report live checks separately. Do not seed, update, or delete live data just to validate a read endpoint.

## Change checklist

1. Inspect current source and Git changes; preserve unrelated work.
2. For API changes, update the model/service/handler/router wiring only where needed.
3. Add behavior tests using the existing fake service; include errors and response shape.
4. Update [routes.md](../routes.md) for API changes, [architecture.md](architecture.md) for structural changes, and `env.example` plus [README.md](../README.md) for configuration changes.
5. Run relevant checks and state their results. Do not claim live verification from mocked tests.

## Work not implemented

Update/delete endpoints, authentication, pagination, queues, and signal-driven graceful shutdown are future work. Creation input validation, bcrypt password hashing, service logic, and MongoDB insertion exist, and are exposed through `POST /users`. Email uniqueness is not enforced by application logic. Implement them only when requested; empty placeholder files do not imply functioning features.
