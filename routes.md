# API routes

Default local base URL: `http://localhost:8000`.
Registration: [router.go](internal/http_mod/router.go). User handlers: [user.go](internal/http_mod/handler/user.go).

## Implemented endpoints

| Method | Path | Success | Failure |
| --- | --- | --- | --- |
| GET | `/health` | 200 JSON: `{"status":"ok"}` | No database check in handler |
| GET | `/hello` | 200 text: `Hello from Go` | — |
| GET | `/users` | 200 JSON array, at most 100 users | 500 JSON: `{"error":"Unable to fetch users"}` |
| GET | `/users/{id}` | 200 JSON user object | 400 invalid ID; 404 missing user; 500 query failure |
| POST | `/users` | 201 JSON user object; Location header | 400 invalid input; 413 oversized body; 500 creation failure |

User endpoints use the `users` collection in `MONGODB_DATABASE`. The list is sorted by `_id` ascending and has a fixed limit of 100. No pagination parameters, filtering, updates, deletes, or authentication are implemented. An empty list is `[]`, not `null`.

The standard router also permits HEAD for GET patterns. Unknown paths return 404 and unsupported methods return 405 with standard-library plain-text responses. There is no uniform JSON middleware for router errors.

## User response

Illustrative data only:

```json
{
  "id": "68c28d101234567890abcdef",
  "name": "Example User",
  "email": "user@example.com",
  "createdAt": "2026-09-11T09:00:00Z",
  "updatedAt": "2026-09-11T09:00:00Z"
}
```

The list endpoint wraps these objects in an array. IDs are hexadecimal MongoDB ObjectIDs. Date fields serialize from `time.Time`; password hashes are never included in JSON. See [model.go](internal/user/model.go).

## Single-user errors

| Status | JSON body | Trigger |
| --- | --- | --- |
| 400 | `{"error":"Invalid user ID"}` | ID is not a valid 24-character hexadecimal ObjectID |
| 404 | `{"error":"User not found"}` | Valid ID has no matching document |
| 500 | `{"error":"Unable to fetch user"}` | Database lookup or decoding fails |

Each user handler allows five seconds for database work, bounded by request cancellation. Query timeout failures currently use the same 500 response as other database failures; 504 is not implemented.

## Try locally

Start the app from the repository root using `go run ./cmd/api`, then:

```sh
curl -i http://localhost:8000/health
curl -i http://localhost:8000/users
curl -i http://localhost:8000/users/REPLACE_WITH_ID_FROM_LIST
```

Use an actual ID returned by the list endpoint for the final request. `/health` confirms the HTTP process can respond, not that MongoDB is still available.

## Create a user

`POST /users` accepts one JSON object, limited to 1 MiB:

```json
{"name":"Example User","email":"user@example.com","password":"a-long-example-password"}
```

Name must be nonblank, email must be a plain valid email address, and password must be 15–72 bytes. Name and email are trimmed; password whitespace is preserved. Empty/malformed bodies, null, unknown fields, wrong field types, and additional JSON values return 400. Validation errors return JSON with an `error` message. Oversized bodies return 413 with `{"error":"Request body too large"}`.

The service hashes the password using bcrypt and assigns UTC timestamps. The same service generates an ObjectID and inserts into `users`. Success returns 201 with the user response above and `Location: /users/{id}`. Failures during hashing or insertion return 500 with `{"error":"Unable to create user"}`. Passwords and hashes are excluded from responses. Email uniqueness is not enforced by application logic; database constraint failures currently return the generic 500 response.
