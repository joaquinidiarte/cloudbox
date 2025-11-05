package handler

import (
	"net/http"

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
