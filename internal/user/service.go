package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidInput = errors.New("invalid user input")
var ErrNotFound = errors.New("user not found")

type Service struct {
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{collection: db.Collection("users")}
}

func (s *Service) Create(ctx context.Context, input CreateUserInput) (User, error) {
	newUser, err := prepareUser(input)
	if err != nil {
		return User{}, err
	}
	if _, err := s.collection.InsertOne(ctx, newUser); err != nil {
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return newUser, nil
}

func prepareUser(input CreateUserInput) (User, error) {
	if err := input.Validate(); err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()

	newUser := User{
		ID:           bson.NewObjectID(),
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.TrimSpace(input.Email),
		Introduction: strings.TrimSpace(input.Introduction),
		Password:     string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return newUser, nil
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(100)

	cursor, err := s.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}

	users := make([]User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}

	return users, nil
}

func (s *Service) GetByID(ctx context.Context, id bson.ObjectID) (User, error) {
	var result User
	err := s.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user: %w", err)
	}
	return result, nil
}

func (input CreateUserInput) Validate() error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("name is required")
	}

	email := strings.TrimSpace(input.Email)

	address, err := mail.ParseAddress(email)

	if err != nil || address.Address != email {
		return errors.New("a valid email is required")
	}

	if len(input.Password) < 15 {
		return errors.New("password must be at least 15 bytes")
	}

	if len(input.Password) > 72 {
		return errors.New("password must be at most 72 bytes")
	}

	return nil
}
