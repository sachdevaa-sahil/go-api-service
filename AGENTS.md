# Repository instructions

## Scope and collaboration

This repository is `go-api-service`, a Go learning project backed by MongoDB. The neighboring `react-app` is a separate project; do not apply its Next.js rules here or edit it for a backend-only task.

The user prefers step-by-step explanations and writing code themselves when asking to learn or be guided. For those requests, explain the next concrete step without editing files. When explicitly asked to implement, fix, or change code, make the changes and verify them. Follow the current request over these defaults.

Inspect `git status --short` before editing. Preserve existing work, including uncommitted changes. Commit or push only when requested. Do not add frameworks, abstractions, dependencies, or future features unrelated to the task.

## Read only the context needed

- [README.md](README.md): setup and execution.
- [routes.md](routes.md): implemented API contract and errors.
- [docs/architecture.md](docs/architecture.md): package ownership, data model, and lifecycle.
- [docs/development.md](docs/development.md): checks, change workflow, and limitations.

Source code is authoritative for current behavior. Inspect the relevant implementation before editing and update affected documentation in the same change. Keep this file short; put detailed reference material in the linked documents. Do not record transient session history or assumed test results as permanent facts.

## Implementation conventions

- Use the Go version from `go.mod` and MongoDB driver v2 import paths. Keep standard-library `net/http` routing unless a task requires otherwise.
- Compose dependencies in `cmd/api/main.go`. Reuse the shared MongoDB client; never connect per request or hide connections in package globals.
- Register routes in `internal/http_mod/router.go`. Handlers own HTTP parsing/status codes; MongoDB repositories own queries. The current handler consumes a small repository interface. A service layer is not implemented yet.
- Pass `context.Context` as the first argument for database work. Derive request work from `r.Context()`, apply bounded timeouts, and defer cancellation. Use independent bounded contexts for cleanup.
- Return errors from lower layers. Wrap with `%w` when preserving causes; inspect with `errors.Is`. Keep process termination in `main`, after `run` returns so deferred cleanup runs.
- Use the existing `writeJSON` helper. Preserve empty arrays, ObjectID validation, generic client-facing database errors, and password exclusion via `json:"-"`.
- Keep structs and constants with the package that owns them. Do not create global catch-all structs/constants/utils packages.
- Never print `.env`, connection credentials, full configuration, or password hashes. Use `env.example` for documentation. Do not change live records as a side effect of tests.

## Verification and reporting

Format changed Go files with `gofmt`. For behavior changes run `go test ./...`, `go vet ./...`, and `git diff --check` from the repository root. Use fake repositories for HTTP tests; cover changed success/error behavior rather than mirroring implementation. Documentation-only edits need link/accuracy and diff checks, not new Go tests.

Do not equate fake-repository tests with a verified live MongoDB query. Report what changed, which checks actually ran, and any unverified behavior. Keep endpoint/model/config documentation synchronized when those contracts change.
