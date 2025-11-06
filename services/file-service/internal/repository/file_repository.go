package repository

import (
	"context"
	"errors"

	"github.com/joaquinidiarte/cloudbox/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FileRepository defines the interface for file data access
type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	FindByUserID(ctx context.Context, userID string, parentID *string) ([]*models.File, error)
	FindByID(ctx context.Context, id string) (*models.File, error)
	FindByOriginalName(ctx context.Context, userID, originalName string, parentID *string) (*models.File, error)
	Delete(ctx context.Context, id string) error
	AddVersion(ctx context.Context, id string, version models.FileVersion, currentVersion int, path, mimeType string, size int64) error
	UpdateCurrentVersion(ctx context.Context, id string, version int, path, mimeType string, size int64) error
	DeleteVersion(ctx context.Context, id string, version int) error
}

// MongoDBFileRepository is the MongoDB implementation of FileRepository
type MongoDBFileRepository struct {
	collection *mongo.Collection
}

// NewFileRepository creates a new MongoDB file repository
func NewFileRepository(db *mongo.Database) FileRepository {
	return &MongoDBFileRepository{
		collection: db.Collection("files"),
	}
}

func (r *MongoDBFileRepository) Create(ctx context.Context, file *models.File) error {
	_, err := r.collection.InsertOne(ctx, file)
	return err
}

func (r *MongoDBFileRepository) FindByUserID(ctx context.Context, userID string, parentID *string) ([]*models.File, error) {
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

func (r *MongoDBFileRepository) FindByID(ctx context.Context, id string) (*models.File, error) {
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

// FindByOriginalName finds a file by user ID and original name
func (r *MongoDBFileRepository) FindByOriginalName(ctx context.Context, userID, originalName string, parentID *string) (*models.File, error) {
	filter := bson.M{
		"user_id":       userID,
		"original_name": originalName,
		"is_folder":     false,
	}
	if parentID != nil {
		filter["parent_id"] = *parentID
	} else {
		filter["parent_id"] = bson.M{"$exists": false}
	}

	var file models.File
	err := r.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *MongoDBFileRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}

func (r *MongoDBFileRepository) AddVersion(ctx context.Context, id string, version models.FileVersion, currentVersion int, path, mimeType string, size int64) error {
	update := bson.M{
		"$push": bson.M{"versions": version},
		"$set": bson.M{
			"current_version": currentVersion,
			"path":            path,
			"mime_type":       mimeType,
			"size":            size,
			"updated_at":      version.UploadedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}

func (r *MongoDBFileRepository) UpdateCurrentVersion(ctx context.Context, id string, version int, path, mimeType string, size int64) error {
	update := bson.M{
		"$set": bson.M{
			"current_version": version,
			"path":            path,
			"mime_type":       mimeType,
			"size":            size,
		},
		"$currentDate": bson.M{
			"updated_at": true,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}

func (r *MongoDBFileRepository) DeleteVersion(ctx context.Context, id string, version int) error {
	update := bson.M{
		"$pull": bson.M{"versions": bson.M{"version": version}},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}
