package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/user-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/middleware"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type UserHandler struct {
	userService *service.UserService
	logger      *utils.Logger
}

func NewUserHandler(userService *service.UserService, logger *utils.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(user, "User retrieved successfully"))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(user, "User retrieved successfully"))
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Errorf("Failed to update user: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("User updated successfully: %s", userID)
	c.JSON(http.StatusOK, models.SuccessResponse(user, "User updated successfully"))
}

func (h *UserHandler) UpdateStorageUsed(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	var req struct {
		Increment int64 `json:"increment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	if err := h.userService.UpdateStorageUsed(c.Request.Context(), userID, req.Increment); err != nil {
		h.logger.Errorf("Failed to update storage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to update storage"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil, "Storage updated successfully"))
}
