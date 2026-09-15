package http_mod

import (
	"net/http"

	"go-api-service/internal/http_mod/handler"
)

func ApplicationRouter(userHandler *handler.UserHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", handler.HealthCheck)
	router.HandleFunc("GET /users", userHandler.GetUsers)
	router.HandleFunc("POST /users", userHandler.CreateUser)
	router.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello from Go"))
	})

	return router
}
