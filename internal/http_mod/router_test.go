package http_mod

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"go-api-service/internal/http_mod/handler"
	"go-api-service/internal/user"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeUserRepository struct {
	get  func(context.Context, bson.ObjectID) (user.User, error)
	list func(context.Context) ([]user.User, error)
}

func (f fakeUserRepository) List(ctx context.Context) ([]user.User, error) {
	return f.list(ctx)
}

func (f fakeUserRepository) GetByID(ctx context.Context, id bson.ObjectID) (user.User, error) {
	return f.get(ctx, id)
}

func TestGetUser(t *testing.T) {
	id := bson.NewObjectID()
	for _, tc := range []struct {
		name, path string
		err        error
		status     int
	}{
		{"found", id.Hex(), nil, 200},
		{"invalid", "99", nil, 400},
		{"missing", id.Hex(), user.ErrNotFound, 404},
		{"failure", id.Hex(), errors.New("private database details"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			repo := fakeUserRepository{get: func(ctx context.Context, got bson.ObjectID) (user.User, error) {
				called = true
				if got != id {
					t.Error("incorrect lookup ID")
				}
				if _, ok := ctx.Deadline(); !ok {
					t.Error("missing deadline")
				}
				return user.User{ID: id, Name: "Test user", Password: "private hash"}, tc.err
			}}
			response := httptest.NewRecorder()
			ApplicationRouter(handler.NewUserHandler(repo)).ServeHTTP(response, httptest.NewRequest("GET", "/users/"+tc.path, nil))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			if called != (tc.status != 400) {
				t.Fatal("incorrect repository invocation")
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON")
			}
			if strings.Contains(response.Body.String(), "private") || strings.Contains(response.Body.String(), "password") {
				t.Fatal("private data exposed")
			}
			if tc.status == 200 {
				var result user.User
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				if result.ID != id || result.Name != "Test user" {
					t.Fatal("incorrect user")
				}
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	id := bson.NewObjectID()
	for _, tc := range []struct {
		name   string
		users  []user.User
		err    error
		status int
	}{
		{name: "users", users: []user.User{{ID: id, Name: "Sahil", Email: "sahil@example.com", Password: "private-hash"}}, status: 200},
		{name: "empty", status: 200},
		{name: "failure", err: errors.New("private database error"), status: 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			repo := fakeUserRepository{list: func(ctx context.Context) ([]user.User, error) {
				called = true
				if _, ok := ctx.Deadline(); !ok {
					t.Error("missing query deadline")
				}
				return tc.users, tc.err
			}}
			router := ApplicationRouter(handler.NewUserHandler(repo))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest("GET", "/users", nil))
			if !called {
				t.Fatal("repository not called")
			}
			if response.Code != tc.status {
				t.Fatalf("status = %d", response.Code)
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatal("expected JSON")
			}
			body := response.Body.String()
			if strings.Contains(body, "private") || strings.Contains(body, "password") {
				t.Fatal("private data exposed")
			}
			if tc.err != nil {
				if body != "{\"error\":\"Unable to fetch users\"}\n" {
					t.Fatalf("unexpected error response: %s", body)
				}
				return
			}
			var got []user.User
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tc.users) {
				t.Fatal("incorrect result count")
			}
			if len(got) == 0 && strings.TrimSpace(body) != "[]" {
				t.Fatal("expected empty array")
			}
			if len(got) > 0 && (got[0].ID != id || got[0].Name != "Sahil" || got[0].Email != "sahil@example.com") {
				t.Fatal("incorrect user data")
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	router := ApplicationRouter(handler.NewUserHandler(fakeUserRepository{list: func(context.Context) ([]user.User, error) { return nil, nil }}))
	for _, tc := range []struct {
		name, method, path string
		status             int
		body               string
	}{
		{"health", "GET", "/health", 200, "{\"status\":\"ok\"}\n"},
		{"greeting", "GET", "/hello", 200, "Hello from Go"},
		{"invalid user ID", "GET", "/users/99", 400, "{\"error\":\"Invalid user ID\"}\n"},
		{"unknown route", "GET", "/missing", 404, ""},
		{"unsupported method", "POST", "/users", 405, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, nil))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			if tc.body != "" && response.Body.String() != tc.body {
				t.Fatalf("body = %q, want %q", response.Body.String(), tc.body)
			}
		})
	}
}
