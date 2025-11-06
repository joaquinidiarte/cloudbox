package repository

import (
	"context"
	"errors"
	"time"

	"github.com/joaquinidiarte/cloudbox/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, update *models.UserUpdateRequest) error {
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if update.FirstName != "" {
		updateDoc["$set"].(bson.M)["first_name"] = update.FirstName
	}
	if update.LastName != "" {
		updateDoc["$set"].(bson.M)["last_name"] = update.LastName
	}
	if update.Email != "" {
		updateDoc["$set"].(bson.M)["email"] = update.Email
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) UpdateStorageUsed(ctx context.Context, userID string, delta int64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"storage_used": delta}},
	)
	return err
}