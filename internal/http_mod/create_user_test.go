package http_mod

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-api-service/internal/http_mod/handler"
	"go-api-service/internal/user"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCreateUser(t *testing.T) {
	const password = "a-long-example-password"
	valid := `{"name":" Sahil ","email":" sahil@example.com ","password":"` + password + `"}`
	for _, tc := range []struct {
		name, body string
		status     int
		dbError    bool
	}{
		{"created", valid, 201, false},
		{"database failure", valid, 500, true},
		{"empty", "", 400, false},
		{"malformed", `{`, 400, false},
		{"null", `null`, 400, false},
		{"array", `[]`, 400, false},
		{"wrong field type", `{"name":42}`, 400, false},
		{"unknown field", `{"admin":true}`, 400, false},
		{"multiple values", valid + ` {}`, 400, false},
		{"trailing garbage", valid + ` !`, 400, false},
		{"invalid input", `{}`, 400, false},
		{"oversized body", `{"name":"` + strings.Repeat("a", 1<<20) + `"}`, 413, false},
		{"oversized trailing whitespace", valid + strings.Repeat(" ", 1<<20), 413, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			id := bson.NewObjectID()
			service := fakeUserService{create: func(ctx context.Context, input user.CreateUserInput) (user.User, error) {
				called = true
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) > 5*time.Second {
					t.Error("expected bounded context")
				}
				if tc.name == "invalid input" {
					return user.User{}, user.ErrInvalidInput
				}
				if input.Name != " Sahil " || input.Email != " sahil@example.com " || input.Password != password {
					t.Error("incorrect decoded input")
				}
				if tc.dbError {
					return user.User{}, errors.New("private database details")
				}
				return user.User{ID: id, Name: "Sahil", Email: "sahil@example.com", Password: "private hash"}, nil
			}}
			router := ApplicationRouter(handler.NewUserHandler(service))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest("POST", "/users", strings.NewReader(tc.body)))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if called != (tc.status == 201 || tc.status == 500 || tc.name == "invalid input") {
				t.Error("unexpected service invocation")
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Error("expected JSON")
			}
			body := response.Body.String()
			if strings.Contains(body, password) || strings.Contains(body, "private") || strings.Contains(body, "$2a$") {
				t.Error("sensitive data exposed")
			}
			if tc.status == 201 {
				var result map[string]any
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				if _, exists := result["password"]; exists {
					t.Error("password field exposed")
				}
				if result["id"] != id.Hex() || result["name"] != "Sahil" {
					t.Error("incorrect created user")
				}
				if response.Header().Get("Location") != "/users/"+id.Hex() {
					t.Error("incorrect location")
				}
			} else if response.Header().Get("Location") != "" {
				t.Error("location on failed creation")
			}
			if tc.status == 500 && body != "{\"error\":\"Unable to create user\"}\n" {
				t.Error("unexpected database error response")
			}
		})
	}
}
