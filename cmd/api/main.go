package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-api-service/internal/config"
	apihttp "go-api-service/internal/http_mod"
	"go-api-service/internal/http_mod/handler"
	"go-api-service/internal/storage/mongodb"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	client, err := mongodb.Connect(context.Background(), cfg.MongoURI)
	if err != nil {
		return fmt.Errorf("connect to MongoDB: %w", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Disconnect MongoDB: %v", err)
		}
	}()

	log.Println("Connected to MongoDB")

	db := client.Database(cfg.MongoDatabase)
	userRepository := mongodb.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepository)
	router := apihttp.ApplicationRouter(userHandler)
	server := apihttp.NewServer(":8000", router)
	log.Println("Starting server on port 8000")

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	return nil
}
