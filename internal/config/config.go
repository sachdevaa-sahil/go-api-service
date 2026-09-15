package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
}

func Load() (Config, error) {
	// A .env file is optional when variables are supplied by the environment.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		// Parser errors can contain file contents, so do not include credentials.
		return Config{}, errors.New("could not load .env; check its permissions and syntax")
	}

	cfg := Config{
		MongoURI:      os.Getenv("MONGODB_URI"),
		MongoDatabase: os.Getenv("MONGODB_DATABASE"),
	}
	if cfg.MongoURI == "" {
		return Config{}, fmt.Errorf("MONGODB_URI is required")
	}

	if cfg.MongoDatabase == "" {
		return Config{}, fmt.Errorf("MONGODB_DATABASE is required")
	}
	return cfg, nil
}
