package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/file-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type FileHandler struct {
	fileService *service.FileService
	logger      *utils.Logger
}

func NewFileHandler(fileService *service.FileService, logger *utils.Logger) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("No file provided"))
		return
	}

	var parentID *string
	if pid := c.PostForm("parent_id"); pid != "" {
		parentID = &pid
	}

	response, err := h.fileService.UploadFile(c.Request.Context(), userID, file, parentID)
	if err != nil {
		h.logger.Errorf("Failed to upload file: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("File uploaded successfully: %s", response.ID)
	c.JSON(http.StatusCreated, models.SuccessResponse(response, "File uploaded successfully"))
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	h.logger.Infof("Upload file for userID: %s", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	var parentID *string
	if pid := c.Query("parent_id"); pid != "" {
		parentID = &pid
	}

	files, err := h.fileService.ListFiles(c.Request.Context(), userID, parentID)
	if err != nil {
		h.logger.Errorf("Failed to list files: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(files, "Files retrieved successfully"))
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")

	file, err := h.fileService.DownloadFile(c.Request.Context(), userID, fileID)
	if err != nil {
		h.logger.Errorf("Failed to download file: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	c.FileAttachment(file.Path, file.OriginalName)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")

	if err := h.fileService.DeleteFile(c.Request.Context(), userID, fileID); err != nil {
		h.logger.Errorf("Failed to delete file: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("File deleted successfully: %s", fileID)
	c.JSON(http.StatusOK, models.SuccessResponse(nil, "File deleted successfully"))
}

func (h *FileHandler) CreateFolder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	var req models.FolderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	folder, err := h.fileService.CreateFolder(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Errorf("Failed to create folder: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("Folder created successfully: %s", folder.ID)
	c.JSON(http.StatusCreated, models.SuccessResponse(folder, "Folder created successfully"))
}

func (h *FileHandler) GetFolderContents(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	folderID := c.Param("id")

	files, err := h.fileService.ListFiles(c.Request.Context(), userID, &folderID)
	if err != nil {
		h.logger.Errorf("Failed to get folder contents: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(files, "Folder contents retrieved successfully"))
}

/* Version operations */
func (h *FileHandler) GetFileVersions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")

	versions, err := h.fileService.GetFileVersions(c.Request.Context(), userID, fileID)
	if err != nil {
		h.logger.Errorf("Failed to get file versions: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(versions, "File versions retrieved successfully"))
}

func (h *FileHandler) DownloadFileVersion(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")
	versionStr := c.Param("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid version number"))
		return
	}

	file, versionPath, err := h.fileService.DownloadFileVersion(c.Request.Context(), userID, fileID, version)
	if err != nil {
		h.logger.Errorf("Failed to download file version: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	c.FileAttachment(versionPath, file.OriginalName)
}

func (h *FileHandler) RestoreFileVersion(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")
	versionStr := c.Param("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid version number"))
		return
	}

	response, err := h.fileService.RestoreFileVersion(c.Request.Context(), userID, fileID, version)
	if err != nil {
		h.logger.Errorf("Failed to restore file version: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("File version restored successfully: %s v%d", fileID, version)
	c.JSON(http.StatusOK, models.SuccessResponse(response, "File version restored successfully"))
}

func (h *FileHandler) DeleteFileVersion(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	fileID := c.Param("id")
	versionStr := c.Param("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid version number"))
		return
	}

	if err := h.fileService.DeleteFileVersion(c.Request.Context(), userID, fileID, version); err != nil {
		h.logger.Errorf("Failed to delete file version: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("File version deleted successfully: %s v%d", fileID, version)
	c.JSON(http.StatusOK, models.SuccessResponse(nil, "File version deleted successfully"))
}
