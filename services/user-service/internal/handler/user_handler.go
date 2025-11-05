package handler

import (
	"net/http"

	"github.com/cloudbox/services/user-service/internal/service"
	"github.com/cloudbox/shared/middleware"
	"github.com/cloudbox/shared/models"
	"github.com/cloudbox/shared/utils"
	"github.com/gin-gonic/gin"
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