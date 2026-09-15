package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		name, uri, database, file, wantError string
	}{
		{name: "environment", uri: "mongodb://localhost:27017", database: "test"},
		{name: "missing URI", database: "test", wantError: "MONGODB_URI is required"},
		{name: "missing database", uri: "mongodb://localhost:27017", wantError: "MONGODB_DATABASE is required"},
		{name: "invalid file", file: "MONGODB_URI='secret", wantError: "could not load .env"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("MONGODB_URI", tc.uri)
			t.Setenv("MONGODB_DATABASE", tc.database)
			if tc.file != "" {
				if err := os.WriteFile(".env", []byte(tc.file), 0600); err != nil {
					t.Fatal(err)
				}
			}
			cfg, err := Load()
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatal("expected configuration error")
				}
				if strings.Contains(err.Error(), "secret") {
					t.Fatal("error exposed file contents")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if cfg.MongoURI != tc.uri || cfg.MongoDatabase != tc.database {
				t.Fatal("configuration differs from environment")
			}
		})
	}
}

func TestLoadDotEnvAndPreserveEnvironment(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("MONGODB_URI", "")
	if err := os.Unsetenv("MONGODB_URI"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MONGODB_DATABASE", "existing")
	if err := os.WriteFile(".env", []byte("MONGODB_URI=mongodb://localhost:27017\nMONGODB_DATABASE=from_file\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MongoURI != "mongodb://localhost:27017" || cfg.MongoDatabase != "existing" {
		t.Fatal("dotenv loading or environment precedence failed")
	}
}
