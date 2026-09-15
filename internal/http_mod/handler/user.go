package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"go-api-service/internal/user"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	List(context.Context) ([]user.User, error)
	GetByID(context.Context, bson.ObjectID) (user.User, error)
	Create(context.Context, user.CreateUserInput) (user.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := h.service.List(ctx)
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
	result, err := h.service.GetByID(ctx, id)
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input *user.CreateUserInput
	err := decoder.Decode(&input)
	if err == nil {
		var extra any
		if nextErr := decoder.Decode(&extra); nextErr != io.EOF {
			if nextErr == nil {
				nextErr = errors.New("multiple JSON values")
			}
			err = nextErr
		}
	}
	if err != nil || input == nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "Request body too large"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Expected a single JSON object with name, email, and password fields"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	created, err := h.service.Create(ctx, *input)
	if errors.Is(err, user.ErrInvalidInput) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Unable to create user"})
		return
	}
	w.Header().Set("Location", "/users/"+created.ID.Hex())
	writeJSON(w, http.StatusCreated, created)
}
