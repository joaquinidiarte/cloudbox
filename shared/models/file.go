package models

import (
	"time"

	"github.com/google/uuid"
)

// FileVersion represents a specific version of a file
type FileVersion struct {
	Version    int       `json:"version" bson:"version"`
	Size       int64     `json:"size" bson:"size"`
	Path       string    `json:"path" bson:"path"`
	MimeType   string    `json:"mime_type" bson:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
	Comment    string    `json:"comment,omitempty" bson:"comment,omitempty"`
}

type File struct {
	ID             string        ` json:"id" bson:"_id"`
	UserID         string        `json:"user_id" bson:"user_id"`
	Name           string        `json:"name" bson:"name"`
	OriginalName   string        `json:"original_name" bson:"original_name"`
	Path           string        `json:"path" bson:"path"`
	Size           int64         `json:"size" bson:"size"`
	MimeType       string        `json:"mime_type" bson:"mime_type"`
	ParentID       *string       `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	IsFolder       bool          `json:"is_folder" bson:"is_folder"`
	IsShared       bool          `json:"is_shared" bson:"is_shared"`
	IsPublic       bool          `json:"is_public" bson:"is_public"`
	CurrentVersion int           `json:"current_version" bson:"current_version"`       // Current version number
	Versions       []FileVersion `json:"versions,omitempty" bson:"versions,omitempty"` // Version history
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
}

type FileResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	OriginalName   string    `json:"original_name"`
	Size           int64     `json:"size"`
	MimeType       string    `json:"mime_type"`
	ParentID       *string   `json:"parent_id,omitempty"`
	IsFolder       bool      `json:"is_folder"`
	IsShared       bool      `json:"is_shared"`
	IsPublic       bool      `json:"is_public"`
	CurrentVersion int       `json:"current_version"`
	VersionCount   int       `json:"version_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewFile(userID, name, originalName, path string, size int64, mimeType string, parentID *string) *File {
	now := time.Now()

	firstVersion := FileVersion{
		Version:    1,
		Size:       size,
		Path:       path,
		MimeType:   mimeType,
		UploadedAt: now,
	}

	return &File{
		ID:             uuid.New().String(),
		UserID:         userID,
		Name:           name,
		OriginalName:   originalName,
		Path:           path,
		Size:           size,
		MimeType:       mimeType,
		ParentID:       parentID,
		IsFolder:       false,
		IsShared:       false,
		IsPublic:       false,
		CurrentVersion: 1,
		Versions:       []FileVersion{firstVersion},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (f *File) ToResponse() FileResponse {
	return FileResponse{
		ID:             f.ID,
		UserID:         f.UserID,
		Name:           f.Name,
		OriginalName:   f.OriginalName,
		Size:           f.Size,
		MimeType:       f.MimeType,
		ParentID:       f.ParentID,
		IsFolder:       f.IsFolder,
		IsShared:       f.IsShared,
		IsPublic:       f.IsPublic,
		CurrentVersion: f.CurrentVersion,
		VersionCount:   len(f.Versions),
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
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

type FileVersionResponse struct {
	Version    int       `json:"version"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
	Comment    string    `json:"comment,omitempty"`
	IsCurrent  bool      `json:"is_current"`
}

func (f *File) GetVersionResponses() []FileVersionResponse {
	responses := make([]FileVersionResponse, len(f.Versions))
	for i, v := range f.Versions {
		responses[i] = FileVersionResponse{
			Version:    v.Version,
			Size:       v.Size,
			MimeType:   v.MimeType,
			UploadedAt: v.UploadedAt,
			Comment:    v.Comment,
			IsCurrent:  v.Version == f.CurrentVersion,
		}
	}
	return responses
}
