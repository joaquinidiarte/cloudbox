package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/shared/models"
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

	// Check if file with same name exists (for versioning)
	existingFile, err := s.fileRepo.FindByOriginalName(ctx, userID, fileHeader.Filename, parentID)
	if err != nil {
		return nil, err
	}

	if existingFile != nil {
		// File with same name exists - create new version
		return s.addNewVersion(ctx, existingFile, fileHeader, src)
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
		parentID,
	)

	if err := s.fileRepo.Create(ctx, file); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	response := file.ToResponse()
	return &response, nil
}

func (c *FileService) ListFiles(ctx context.Context, userID string, parentID *string) ([]*models.FileResponse, error) {
	files, err := c.fileRepo.FindByUserID(ctx, userID, parentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.FileResponse, len(files))
	for i, file := range files {
		response := file.ToResponse()
		responses[i] = &response
	}
	return responses, nil
}

func (s *FileService) DownloadFile(ctx context.Context, userID, fileID string) (*models.File, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if file.UserID != userID && !file.IsPublic {
		return nil, errors.New("unauthorized access to file")
	}

	if file.IsFolder {
		return nil, errors.New("cannot download a folder")
	}

	return file, nil
}

func (s *FileService) DeleteFile(ctx context.Context, userID, fileID string) error {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return err
	}

	if file.UserID != userID {
		return errors.New("unauthorized access to file")
	}

	// Delete file from disk
	if !file.IsFolder && file.Path != "" {
		os.Remove(file.Path)
	}

	return s.fileRepo.Delete(ctx, fileID)
}

/* Folder operations */
func (s *FileService) CreateFolder(ctx context.Context, userID string, req *models.FolderCreateRequest) (*models.FileResponse, error) {
	folder := models.NewFolder(userID, req.Name, req.ParentID)

	if err := s.fileRepo.Create(ctx, folder); err != nil {
		return nil, err
	}

	response := folder.ToResponse()
	return &response, nil
}

/* Version operations */
func (s *FileService) addNewVersion(ctx context.Context, existingFile *models.File, fileHeader *multipart.FileHeader, src multipart.File) (*models.FileResponse, error) {
	// Reset file pointer
	src.Seek(0, 0)

	// Create user directory
	userDir := filepath.Join(s.storagePath, existingFile.UserID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename for new version
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

	// Create new version
	newVersionNumber := existingFile.CurrentVersion + 1
	newVersion := models.FileVersion{
		Version:    newVersionNumber,
		Size:       fileHeader.Size,
		Path:       filePath,
		MimeType:   fileHeader.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	// Add version to database
	if err := s.fileRepo.AddVersion(ctx, existingFile.ID, newVersion, newVersionNumber, filePath, newVersion.MimeType, fileHeader.Size); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	// Get updated file
	updatedFile, err := s.fileRepo.FindByID(ctx, existingFile.ID)
	if err != nil {
		return nil, err
	}

	response := updatedFile.ToResponse()
	return &response, nil
}

func (s *FileService) GetFileVersions(ctx context.Context, userID, fileID string) ([]models.FileVersionResponse, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if file.UserID != userID {
		return nil, errors.New("unauthorized access to file")
	}

	if file.IsFolder {
		return nil, errors.New("folders do not have versions")
	}

	return file.GetVersionResponses(), nil
}

func (s *FileService) DownloadFileVersion(ctx context.Context, userID, fileID string, version int) (*models.File, string, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, "", err
	}

	if file.UserID != userID && !file.IsPublic {
		return nil, "", errors.New("unauthorized access to file")
	}

	if file.IsFolder {
		return nil, "", errors.New("cannot download a folder")
	}

	// Find the requested version
	var versionPath string
	for _, v := range file.Versions {
		if v.Version == version {
			versionPath = v.Path
			break
		}
	}

	if versionPath == "" {
		return nil, "", errors.New("version not found")
	}

	return file, versionPath, nil
}
