package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type FileService struct {
	fileRepo    *repository.FileRepository
	storagePath string
	maxFileSize int64
}

func NewFileService(fileRepo *repository.FileRepository, storagePath string, maxFileSize int64) *FileService {
	return &FileService{
		fileRepo:    fileRepo,
		storagePath: storagePath,
		maxFileSize: maxFileSize,
	}
}

func (s *FileService) UploadFile(ctx context.Context, userID string, fileHeader *multipart.FileHeader, parentID *string) (*models.FileResponse, error) {
	// Check file size
	if fileHeader.Size > s.maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxFileSize)
	}

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Calculate hash
	hash, err := utils.HashFile(src)
	if err != nil {
		return nil, err
	}
	// Reset file pointer
	src.Seek(0, 0)

	// Create user directory
	userDir := filepath.Join(s.storagePath, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := filepath.Join(userDir, filename)

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	// Create file record
	file := models.NewFile(
		userID,
		filename,
		fileHeader.Filename,
		filePath,
		fileHeader.Size,
		fileHeader.Header.Get("Content-Type"),
		hash,
		parentID,
	)

	if err := s.fileRepo.Create(ctx, file); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	response := file.ToResponse()
	return &response, nil
}
