package http_mod

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-service/internal/http_mod/handler"
)

func TestUserEndpointsAgree(t *testing.T) {
	router := ApplicationRouter()
	list := httptest.NewRecorder()
	router.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/users", nil))
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d", list.Code)
	}
	if got := list.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q", got)
	}
	var users []handler.User
	if err := json.Unmarshal(list.Body.Bytes(), &users); err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatalf("got %d users, want 4", len(users))
	}
	for _, want := range users {
		t.Run(want.ID, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/users/"+want.ID, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d", response.Code)
			}
			var got handler.User
			if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("got %+v, want %+v", got, want)
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	router := ApplicationRouter()
	for _, tc := range []struct {
		name, method, path string
		status             int
		body               string
	}{
		{"health", "GET", "/health", 200, "{\"status\":\"ok\"}\n"},
		{"greeting", "GET", "/hello", 200, "Hello from Go"},
		{"missing user", "GET", "/users/99", 404, "{\"error\":\"User not found\"}\n"},
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
