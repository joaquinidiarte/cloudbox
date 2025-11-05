package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaquinidiarte/cloudbox/services/auth-service/internal/service"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *utils.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	response, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorf("Registration failed: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("User registered successfully: %s", req.Email)
	c.JSON(http.StatusCreated, models.SuccessResponse(response, "User registered successfully"))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorf("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(err.Error()))
		return
	}

	h.logger.Infof("User logged in successfully: %s", req.Email)
	c.JSON(http.StatusOK, models.SuccessResponse(response, "Login successful"))
}

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	claims, err := h.authService.VerifyToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(claims, "Token is valid"))
}
