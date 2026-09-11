package http_mod

import (
	"net/http"

	"go-api-service/internal/http_mod/handler"
)

func ApplicationRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", handler.HealthCheck)
	router.HandleFunc("GET /users", handler.GetUsers)
	router.HandleFunc("GET /users/{id}", handler.GetUserByID)
	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello from Go"))
	})

	return router
}
