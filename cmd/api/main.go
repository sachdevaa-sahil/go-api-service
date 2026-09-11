package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	apihttp "go-api-service/internal/http_mod"
)

func main() {
	server := apihttp.NewServer(":8000")
	log.Println("Starting server on port 8000")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server failed: %v", err)
		os.Exit(1)
	}
}
