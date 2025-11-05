package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           string    ` json:"id" bson:"_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	Name         string    `json:"name" bson:"name"`
	OriginalName string    `json:"original_name" bson:"original_name"`
	Path         string    `json:"path" bson:"path"`
	Size         int64     `json:"size" bson:"size"`
	MimeType     string    `json:"mime_type" bson:"mime_type"`
	Hash         string    `json:"hash" bson:"hash"`
	ParentID     *string   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	IsFolder     bool      `json:"is_folder" bson:"is_folder"`
	IsShared     bool      `json:"is_shared" bson:"is_shared"`
	IsPublic     bool      `json:"is_public" bson:"is_public"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type FileResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	ParentID  *string   `json:"parent_id,omitempty"`
	IsFolder  bool      `json:"is_folder"`
	IsShared  bool      `json:"is_shared"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFile(userID, name, originalName, path string, size int64, mimeType, hash string, parentID *string) *File {
	now := time.Now()

	return &File{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         name,
		OriginalName: originalName,
		Path:         path,
		Size:         size,
		MimeType:     mimeType,
		Hash:         hash,
		ParentID:     parentID,
		IsFolder:     false,
		IsShared:     false,
		IsPublic:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (f *File) ToResponse() FileResponse {
	return FileResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		Name:      f.Name,
		Size:      f.Size,
		MimeType:  f.MimeType,
		ParentID:  f.ParentID,
		IsFolder:  f.IsFolder,
		IsShared:  f.IsShared,
		IsPublic:  f.IsPublic,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

type FolderCreateRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id,omitempty"`
}

func NewFolder(userID, name string, parentID *string) *File {
	now := time.Now()
	return &File{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		ParentID:  parentID,
		IsFolder:  true,
		IsShared:  false,
		IsPublic:  false,
		Size:      0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
