package user

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestPrepareUser(t *testing.T) {
	input := CreateUserInput{Name: " Sahil ", Email: " sahil@example.com ", Password: "  a-long-example-password  "}
	started := time.Now()
	u, err := prepareUser(input)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID.IsZero() || u.Name != "Sahil" || u.Email != "sahil@example.com" {
		t.Fatal("incorrect user fields")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		t.Fatal("password not correctly hashed")
	}
	if u.CreatedAt.Before(started) || u.CreatedAt.After(time.Now()) || !u.CreatedAt.Equal(u.UpdatedAt) || u.CreatedAt.Location() != time.UTC {
		t.Fatal("incorrect timestamps")
	}
}

func TestCreateRejectsInvalidInputBeforeDatabase(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input CreateUserInput
	}{
		{"missing fields", CreateUserInput{}},
		{"blank name", CreateUserInput{Name: " ", Email: "user@example.com", Password: strings.Repeat("a", 15)}},
		{"invalid email", CreateUserInput{Name: "Sahil", Email: "invalid", Password: strings.Repeat("a", 15)}},
		{"display name", CreateUserInput{Name: "Sahil", Email: "Sahil <user@example.com>", Password: strings.Repeat("a", 15)}},
		{"short password", CreateUserInput{Name: "Sahil", Email: "user@example.com", Password: strings.Repeat("a", 14)}},
		{"long password", CreateUserInput{Name: "Sahil", Email: "user@example.com", Password: strings.Repeat("a", 73)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// A nil collection makes accidental database access fail this test.
			_, err := (&Service{}).Create(context.Background(), tc.input)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("expected invalid input, got %v", err)
			}
		})
	}
}

func TestPasswordLengthBoundaries(t *testing.T) {
	for _, length := range []int{15, 72} {
		input := CreateUserInput{Name: "Sahil", Email: "user@example.com", Password: strings.Repeat("a", length)}
		if err := input.Validate(); err != nil {
			t.Errorf("length %d: %v", length, err)
		}
	}
}
