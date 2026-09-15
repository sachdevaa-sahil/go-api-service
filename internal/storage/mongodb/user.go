package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go-api-service/internal/user"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (repo *UserRepository) List(ctx context.Context) ([]user.User, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(100)

	cursor, err := repo.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}

	users := make([]user.User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}

	return users, nil
}

func (repo *UserRepository) GetByID(ctx context.Context, id bson.ObjectID) (user.User, error) {
	var result user.User
	err := repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, user.ErrNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("find user: %w", err)
	}
	return result, nil
}
