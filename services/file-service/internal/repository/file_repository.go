package repository

import (
	"context"
	"errors"

	"github.com/joaquinidiarte/cloudbox/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) *FileRepository {
	return &FileRepository{
		collection: db.Collection("files"),
	}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	_, err := r.collection.InsertOne(ctx, file)
	return err
}

func (r *FileRepository) FindByUserID(ctx context.Context, userID string, parentID *string) ([]*models.File, error) {
	filter := bson.M{"user_id": userID}
	if parentID != nil {
		filter["parent_id"] = *parentID
	} else {
		filter["parent_id"] = bson.M{"$exists": false}
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*models.File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) FindByID(ctx context.Context, id string) (*models.File, error) {
	var file models.File
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}
