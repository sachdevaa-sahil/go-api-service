package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"go-api-service/internal/user"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	List(context.Context) ([]user.User, error)
	GetByID(context.Context, bson.ObjectID) (user.User, error)
}

type UserHandler struct {
	repository UserRepository
}

func NewUserHandler(repository UserRepository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := h.repository.List(ctx)
	if err != nil {
		log.Printf("List users: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to fetch users"})
		return
	}
	if users == nil {
		users = make([]user.User, 0)
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := h.repository.GetByID(ctx, id)
	if errors.Is(err, user.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	if err != nil {
		log.Printf("Get user: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to fetch user"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}
